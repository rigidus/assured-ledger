package sm_requester

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/insolar/assured-ledger/ledger-core/v2/conveyor"
	"github.com/insolar/assured-ledger/ledger-core/v2/conveyor/smachine"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar/api"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar/jet"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar/record"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar/utils"
	"github.com/insolar/assured-ledger/ledger-core/v2/vanilla/injector"
)

// will operate with 2 callbacks:
// Done callback (to notify that we stopped handling messages)
type MessageSenderV2 struct{}

func (MessageSenderV2) SendRole(
	context.Context,
	*message.Message,
	insolar.DynamicRole,
	insolar.Reference,
) func() {
	return nil
}
func (MessageSenderV2) SendTarget(
	context.Context,
	*message.Message,
	insolar.Reference,
) func() {
	return nil
}

type MessageReceiverV2 struct{}

type OutboundRequest struct {
	smachine.StateMachineDeclTemplate

	// sender    s_sender.SenderServiceAdapter
	jetCoordinator jet.Coordinator
	pub            message.Publisher
	pulseSlot      *conveyor.PulseSlot

	Node insolar.Reference

	// Incoming Arguments
	Incoming        *record.IncomingRequest
	ObjectReference insolar.Reference
	Arguments       []byte
	Method          string
}

/* */

func (s *OutboundRequest) InjectDependencies(_ smachine.StateMachine, _ smachine.SlotLink, injector *injector.DependencyInjector) {
}

func (s *OutboundRequest) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return s
}

func (s *OutboundRequest) GetInitStateFor(smachine.StateMachine) smachine.InitFunc {
	return s.stepInit
}

/* */

func (s *OutboundRequest) stepInit(ctx smachine.InitializationContext) smachine.StateUpdate {
	if s.Incoming == nil {
		pulseNumber := s.pulseSlot.PulseData().PulseNumber

		s.Incoming = &record.IncomingRequest{
			Object:       &s.ObjectReference,
			Method:       s.Method,
			Arguments:    s.Arguments,
			APIRequestID: utils.TraceID(ctx.GetContext()),
			APINode:      s.Node,
			Reason:       api.MakeReason(pulseNumber, s.Arguments),
			Immutable:    true,
		}
	}

	return ctx.Jump(s.stepSendRequest)
}

func (s *OutboundRequest) stepSendRequest(ctx smachine.ExecutionContext) smachine.StateUpdate {
}
