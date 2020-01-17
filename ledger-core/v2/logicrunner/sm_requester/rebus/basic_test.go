package rebus

import (
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/assured-ledger/ledger-core/v2/conveyor"
	"github.com/insolar/assured-ledger/ledger-core/v2/conveyor/smachine"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar/jet"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar/payload"
	"github.com/insolar/assured-ledger/ledger-core/v2/instrumentation/inslogger"
	"github.com/insolar/assured-ledger/ledger-core/v2/logicrunner/statemachine"
	"github.com/insolar/assured-ledger/ledger-core/v2/pulse"
	"github.com/insolar/assured-ledger/ledger-core/v2/vanilla/injector"
	"github.com/insolar/assured-ledger/ledger-core/v2/virtual/integration/mimic"
)

type TestingSend struct {
	smachine.StateMachineDeclTemplate

	t      testing.TB
	sender MessageSend
	msg    *message.Message

	err      error
	messages []*message.Message
}

/* --------------- DECLARATION --------------- */

func (s *TestingSend) InjectDependencies(_ smachine.StateMachine, _ smachine.SlotLink, injector *injector.DependencyInjector) {
	s.sender.InjectDependencies(injector)
}

func (s *TestingSend) GetInitStateFor(_ smachine.StateMachine) smachine.InitFunc {
	return s.stepInit
}

/* --------------- INSTANCE --------------- */

func (s *TestingSend) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return s
}

func (s *TestingSend) stepInit(ctx smachine.InitializationContext) smachine.StateUpdate {
	return ctx.Jump(s.stepSend)
}

func (s *TestingSend) stepSend(ctx smachine.ExecutionContext) smachine.StateUpdate {
	s.sender.Message = s.msg

	return s.sender.JumpSelectTarget(ctx, func(messages []*message.Message, err error) smachine.StateFunc {
		s.messages = messages
		s.err = err

		return s.stepFinishWait
	})
}

func (s *TestingSend) stepFinishWait(ctx smachine.ExecutionContext) smachine.StateUpdate {
	assert.NoError(s.t, s.err)
	assert.Len(s.t, s.messages, 1)

	return ctx.Stop()
}

// =====================================================================================================================

//go:generate minimock -i github.com/ThreeDotsLabs/watermill/message.Publisher -o ./ -s _mock.go -g

type Request struct{ msg *message.Message }
type Response struct{ msg *message.Message }

func TestMessage_SendReceive(t *testing.T) {
	var (
		ctx   = inslogger.TestContext(t)
		goCtx = ctx
	)

	slotMachineConfig := smachine.SlotMachineConfig{
		PollingPeriod:     500 * time.Millisecond,
		PollingTruncate:   1 * time.Millisecond,
		SlotPageSize:      1000,
		ScanCountLimit:    100000,
		SlotMachineLogger: statemachine.ConveyorLoggerFactory{},
	}

	pulseConveyorConfig := conveyor.PulseConveyorConfig{
		ConveyorMachineConfig: slotMachineConfig,
		SlotMachineConfig:     slotMachineConfig,
		EventlessSleep:        100 * time.Millisecond,
		MinCachePulseAge:      100,
		MaxPastPulseAge:       1000,
	}

	var (
		messageOut *message.Message
		messageIn  *message.Message
		err        error
	)

	messageOut, err = payload.NewMessage(&payload.GetRequest{})
	require.NoError(t, err)

	messageIn, err = payload.NewMessage(&payload.Replication{})
	require.NoError(t, err)

	defaultHandlers := func(pn pulse.Number, input conveyor.InputEvent) smachine.CreateFunc {
		switch input.(type) {
		case Request:
			return func(ctx smachine.ConstructionContext) smachine.StateMachine {
				ctx.SetContext(goCtx)
				return &TestingSend{
					t: t,
				}
			}
		case Response:
			return func(ctx smachine.ConstructionContext) smachine.StateMachine {
				ctx.SetContext(goCtx)
				return &MessageReceive{
					Message: nil,
				}
			}
		default:
			panic("123")
		}
	}

	jc := jet.NewCoordinatorMock(t)
	p := NewPublisherMock(t)

	pc := conveyor.NewPulseConveyor(ctx, pulseConveyorConfig, defaultHandlers, nil)
	pc.AddDependency(jet.Coordinator(jc))
	pc.AddDependency(message.Publisher(p))

	pulseGenerator := mimic.NewPulseGenerator(10)
	pulseData := pulseGenerator.Generate()
	assert.NoError(t, pc.CommitPulseChange(pulseData.AsRange()))

	assert.NoError(t, pc.AddInput(ctx, pulseData.PulseNumber, Request{msg: messageOut}))
	assert.NoError(t, pc.AddInput(ctx, pulseData.PulseNumber, Response{msg: messageIn}))
}
