package rebus

import (
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"

	"github.com/insolar/assured-ledger/ledger-core/v2/conveyor"
	"github.com/insolar/assured-ledger/ledger-core/v2/conveyor/smachine"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar/bus"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar/bus/meta"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar/jet"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar/payload"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar/utils"
	"github.com/insolar/assured-ledger/ledger-core/v2/instrumentation/instracer"
	"github.com/insolar/assured-ledger/ledger-core/v2/vanilla/injector"
)

type SharedResponseBucket struct {
	replies []*message.Message
}

type ResultCallback func(message []*message.Message, err error) smachine.StateFunc

type MessageSend struct {
	jetCoordinator jet.Coordinator
	pulseSlot      *conveyor.PulseSlot
	publisher      message.Publisher
	timeout        time.Duration
	catalog        MessageCatalog

	// input:
	// if Origin is nil:
	// * if Role is undefined - send to target
	// * if Role is defined - send by role
	// if Origin is not nil:
	// * this is reply
	// panic otherwise
	Message        *message.Message
	Role           insolar.DynamicRole
	ObjectOrTarget insolar.Reference
	Origin         payload.Meta

	resultCallback ResultCallback

	// field to pass between steps
	waitStarted time.Time
	target      insolar.Reference
	isStarted   bool
	key         MessageKey

	// fields for reply
	shared SharedResponseBucket
}

func (s *MessageSend) checkConstraints() {
	if s.isStarted {
		panic("already started")
	}
}

func (s *MessageSend) Prepare(ctx smachine.ExecutionContext, cb ResultCallback) MessageSendFabric {
	return MessageSendFabric{s: s, ctx: ctx, cb: cb}
}

type MessageSendFabric struct {
	s   *MessageSend
	ctx smachine.ExecutionContext
	cb  ResultCallback
}

func (f MessageSendFabric) SendRole(msg *message.Message, object insolar.Reference, role insolar.DynamicRole) smachine.StateUpdate {
	f.s.checkConstraints()

	f.s.Message = msg
	f.s.Role = role
	f.s.ObjectOrTarget = object

	return f.s.JumpSelectTarget(f.ctx, f.cb)
}

func (f MessageSendFabric) SendTarget(msg *message.Message, target insolar.Reference) smachine.StateUpdate {
	f.s.checkConstraints()

	f.s.Message = msg
	f.s.ObjectOrTarget = target

	return f.s.JumpSelectTarget(f.ctx, f.cb)

}

func (f MessageSendFabric) SendResponse(msg *message.Message, origin payload.Meta) smachine.StateUpdate {
	f.s.checkConstraints()

	f.s.Message = msg
	f.s.Origin = origin

	return f.s.JumpSelectTarget(f.ctx, f.cb)
}

func (s *MessageSend) InjectDependencies(injector *injector.DependencyInjector) {
	injector.MustInject(&s.jetCoordinator)
	injector.MustInject(&s.pulseSlot)
	injector.MustInject(&s.publisher)
}

func (s *MessageSend) JumpSelectTarget(ctx smachine.ExecutionContext, cb ResultCallback) smachine.StateUpdate {
	s.isStarted = true

	var (
		goCtx       = ctx.GetContext()
		pulseNumber = s.pulseSlot.PulseData().PulseNumber
	)

	if s.Role == insolar.DynamicRoleUndefined {
		s.target = s.ObjectOrTarget
	} else {
		targetReferences, err := s.jetCoordinator.QueryRole(goCtx, s.Role, *s.ObjectOrTarget.GetLocal(), pulseNumber)
		if err != nil {
			return s.error(ctx, errors.Wrap(err, "failed to get target node"))
		}

		s.target = targetReferences[0]
	}

	s.resultCallback = cb
	return ctx.Jump(s.stepSendMessage)
}

func (s *MessageSend) JumpWaitMoreResponses(ctx smachine.ExecutionContext, cb ResultCallback) smachine.StateUpdate {
	s.resultCallback = cb
	return s.stepWaitMoreResponses(ctx)
}

// this method modifies message in-place
func (s *MessageSend) internalWrapMessageMeta(
	msg *message.Message,
	receiver insolar.Reference,
	originHash payload.MessageHash,
	pulseNumber insolar.PulseNumber,
) (payload.Meta, error) {
	payloadMeta := payload.Meta{
		Polymorph:  uint32(payload.TypeMeta),
		Payload:    msg.Payload,
		Sender:     s.jetCoordinator.Me(),
		Receiver:   receiver,
		Pulse:      pulseNumber,
		ID:         []byte(msg.UUID),
		OriginHash: originHash,
	}

	buf, err := payloadMeta.Marshal()
	if err != nil {
		return payload.Meta{}, errors.Wrap(err, "failed to wrap message")
	}

	msg.Payload = buf
	msg.Metadata.Set(meta.Receiver, receiver.String())

	return payloadMeta, nil
}

func (s *MessageSend) stepSendMessage(ctx smachine.ExecutionContext) smachine.StateUpdate {
	var (
		goCtx       = ctx.GetContext()
		traceID     = utils.TraceID(goCtx)
		pulseNumber = s.pulseSlot.PulseData().PulseNumber
	)

	msg := s.Message.Copy()
	msg.Metadata.Set(meta.TraceID, traceID)

	if spanDataSerialized, err := instracer.Serialize(goCtx); err != nil {
		ctx.Log().Error("failed to serialize context", err)
	} else {
		msg.Metadata.Set(meta.SpanData, string(spanDataSerialized))
	}

	msg.SetContext(goCtx)
	wrapped, err := s.internalWrapMessageMeta(msg, s.target, payload.MessageHash{}, pulseNumber)
	if err != nil {
		return s.error(ctx, errors.Wrap(err, "failed to wrap message meta"))
	}

	{ // count hash and store it for future use
		messageHash := payload.MessageHash{}
		if err := messageHash.Unmarshal(wrapped.ID); err != nil {
			return s.error(ctx, errors.Wrap(err, "failed to unmarshal message hash"))
		}

		key := MessageKey(messageHash)
		if !s.catalog.Create(ctx, key, &s.shared) {
			panic("failed to publish response bucket")
		}

		if err := s.publisher.Publish(bus.TopicOutgoing, msg); err != nil {
			return s.error(ctx, errors.Wrap(err, "failed to publish message"))
		}
	}

	return ctx.Jump(s.stepWaitMoreResponses)
}

func (s *MessageSend) stepWaitMoreResponses(ctx smachine.ExecutionContext) smachine.StateUpdate {
	s.waitStarted = time.Now()
	wakeUpAt := s.waitStarted.Add(s.timeout)

	return ctx.WaitAnyUntil(wakeUpAt).ThenJump(s.stepCheckResult)
}

func (s *MessageSend) stepCheckResult(ctx smachine.ExecutionContext) smachine.StateUpdate {
	var (
		replies = s.shared.replies

		err error
	)
	s.shared.replies = make([]*message.Message, 0, 0)

	if len(replies) == 0 {
		if time.Now().Sub(s.waitStarted) > 200*time.Millisecond {
			ctx.Log().Error("WakeUp before we got any replies", nil)

			wakeUpAt := s.waitStarted.Add(s.timeout)
			return ctx.WaitAnyUntil(wakeUpAt).ThenRepeat()
		}

		// we got no replies == timeout
		replies = nil
		err = errors.New("timeout")
	}

	stateFunc := s.resultCallback(replies, err)
	return ctx.Jump(stateFunc)
}

func (s *MessageSend) Done(ctx smachine.ExecutionContext) *MessageSend {
	// cleanup previous run
	if !s.isStarted {
		panic("not started")
	}
	s.isStarted = false

	if s.key != MessageKey(nil) {
		s.catalog.Finish(ctx, s.key)
		s.key = MessageKey(nil)
	}

	return s
}

func (s *MessageSend) event(ctx smachine.ExecutionContext, replies []*message.Message, err error) smachine.StateUpdate {
	stateFunc := s.resultCallback(replies, err)
	return ctx.Jump(stateFunc)
}

func (s *MessageSend) error(ctx smachine.ExecutionContext, err error) smachine.StateUpdate {
	return s.Done(ctx).event(ctx, nil, err)
}
