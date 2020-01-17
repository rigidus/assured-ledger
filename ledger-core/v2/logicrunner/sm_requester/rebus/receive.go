package rebus

import (
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"

	"github.com/insolar/assured-ledger/ledger-core/v2/conveyor/smachine"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar/payload"
	"github.com/insolar/assured-ledger/ledger-core/v2/vanilla/injector"
)

type MessageReceive struct {
	smachine.StateMachineDeclTemplate

	catalog MessageCatalog

	Message *message.Message

	meta payload.Meta
	key  MessageKey
}

/* --------------- DECLARATION --------------- */

func (s *MessageReceive) InjectDependencies(_ smachine.StateMachine, _ smachine.SlotLink, injector *injector.DependencyInjector) {
	injector.MustInject(&s)
}

func (s *MessageReceive) GetInitStateFor(_ smachine.StateMachine) smachine.InitFunc {
	return s.stepInit
}

/* --------------- INSTANCE ------------------ */

func (s *MessageReceive) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return s
}

func (s *MessageReceive) stepInit(ctx smachine.InitializationContext) smachine.StateUpdate {
	// trying to unmarshal payload
	if err := s.meta.Unmarshal(s.Message.Payload); err != nil {
		return ctx.Error(errors.Wrap(err, "failed to unmarshal meta"))
	}

	// not a reply, do nothing
	if s.meta.OriginHash.IsZero() {
		return ctx.Stop()
	}
	s.key = MessageKey(s.meta.OriginHash)

	return ctx.Jump(s.stepFindAwaitingSender)
}

func (s *MessageReceive) stepFindAwaitingSender(ctx smachine.ExecutionContext) smachine.StateUpdate {
	slotLink := ctx.GetPublishedGlobalAlias(s.key)
	if slotLink.IsEmpty() {
		ctx.Log().Warn("nobody expects answer, dropping")
	}

	asyncLog := ctx.LogAsync()

	cb := func(callContext smachine.MachineCallContext) {
		accessor, ok := s.catalog.TryGet(ctx, s.key)
		if !ok {
			asyncLog.Warn("nobody expects this message, already")
			return
		}

		cb := func(bucket *SharedResponseBucket) {
			bucket.replies = append(bucket.replies, s.Message)
		}
		for {
			switch accessor.Prepare(cb).TryUse(ctx).GetDecision() {
			case smachine.Passed:
				asyncLog.Warn("message successfully delivered")
				return
			case smachine.NotPassed:
				time.Sleep(10 * time.Millisecond)
			case smachine.Impossible:
				asyncLog.Warn("nobody expects this message, now")
			}
		}
	}

	if !smachine.ScheduleCallTo(slotLink, cb, false) {
		return ctx.Errorf("failed to schedule call to slotLink, probably inactive now")
	}

	return ctx.Stop()
}
