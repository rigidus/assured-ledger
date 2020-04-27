// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/assured-ledger/blob/master/LICENSE.md.

package small

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/insolar/assured-ledger/ledger-core/v2/conveyor"
	"github.com/insolar/assured-ledger/ledger-core/v2/conveyor/smachine"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar/gen"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar/jet"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar/payload"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar/pulse"
	"github.com/insolar/assured-ledger/ledger-core/v2/instrumentation/inslogger"
	"github.com/insolar/assured-ledger/ledger-core/v2/network/messagesender"
	messagesenderadapter "github.com/insolar/assured-ledger/ledger-core/v2/network/messagesender/adapter"
	pulse2 "github.com/insolar/assured-ledger/ledger-core/v2/pulse"
	"github.com/insolar/assured-ledger/ledger-core/v2/reference"
	"github.com/insolar/assured-ledger/ledger-core/v2/runner"
	"github.com/insolar/assured-ledger/ledger-core/v2/runner/adapter"
	"github.com/insolar/assured-ledger/ledger-core/v2/runner/executor"
	"github.com/insolar/assured-ledger/ledger-core/v2/testutils"
	"github.com/insolar/assured-ledger/ledger-core/v2/vanilla/longbits"
	"github.com/insolar/assured-ledger/ledger-core/v2/virtual/descriptor"
	"github.com/insolar/assured-ledger/ledger-core/v2/virtual/execute"
	"github.com/insolar/assured-ledger/ledger-core/v2/virtual/integration/utils"
)

func Test_Constructor_Increment_Pending_Counters(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()

	ctx := inslogger.TestContext(t)

	machineConfig := smachine.SlotMachineConfig{
		PollingPeriod:   500 * time.Millisecond,
		PollingTruncate: 1 * time.Millisecond,
		SlotPageSize:    1000,
		ScanCountLimit:  100000,
	}

	pd := pulse2.NewFirstPulsarData(10, longbits.Bits256{})
	pn := pd.PulseNumber
	smExecute := execute.SMExecute{}

	pulseConveyor := conveyor.NewPulseConveyor(ctx, conveyor.PulseConveyorConfig{
		ConveyorMachineConfig: machineConfig,
		SlotMachineConfig:     machineConfig,
		EventlessSleep:        100 * time.Millisecond,
		MinCachePulseAge:      100,
		MaxPastPulseAge:       1000,
	}, func(inputPn pulse2.Number, input conveyor.InputEvent) smachine.CreateFunc {
		require.Equal(t, pn, inputPn)
		return func(ctx smachine.ConstructionContext) smachine.StateMachine {
			smExecute.Payload = input.(*payload.VCallRequest)
			return &smExecute
		}
	}, nil)

	executorMock := testutils.NewMachineLogicExecutorMock(mc)
	executorMock.CallConstructorMock.Set(
		func(ctx context.Context, callContext *insolar.LogicCallContext, code insolar.Reference,
			name string, args insolar.Arguments) (objectState []byte, result insolar.Arguments, err error) {
			sMachine := smExecute
			fmt.Printf("%v", sMachine)
			return nil, []byte("345"), nil
		})
	runnerService := runner.NewService()
	require.NoError(t, runnerService.Init())
	manager := executor.NewManager()
	err := manager.RegisterExecutor(insolar.MachineTypeBuiltin, executorMock)
	require.NoError(t, err)
	runnerService.Manager = manager
	cacheMock := descriptor.NewCacheMock(t)
	runnerService.Cache = cacheMock
	cacheMock.ByPrototypeRefMock.Return(
		descriptor.NewPrototypeDescriptor(gen.Reference(), gen.ID(), gen.Reference()),
		descriptor.NewCodeDescriptor(nil, insolar.MachineTypeBuiltin, gen.Reference()),
		nil,
	)
	runnerAdapter := adapter.CreateRunnerServiceAdapter(ctx, runnerService)
	pulseConveyor.AddDependency(runnerAdapter)

	publisherMock := &utils.PublisherMock{}
	jetCoordinatorMock := jet.NewCoordinatorMock(t).
		MeMock.Return(gen.Reference()).
		QueryRoleMock.Return([]insolar.Reference{gen.Reference()}, nil)
	pulses := pulse.NewStorageMem()
	messageSender := messagesender.NewDefaultService(publisherMock, jetCoordinatorMock, pulses)
	messageSenderAdapter := messagesenderadapter.CreateMessageSendService(ctx, messageSender)
	pulseConveyor.AddDependency(messageSenderAdapter)

	emerChan := make(chan struct{})
	pulseConveyor.StartWorker(emerChan, func() {})
	defer func() {
		close(emerChan)
	}()

	require.NoError(t, pulseConveyor.CommitPulseChange(pd.AsRange()))

	prototype := gen.Reference()
	pl := payload.VCallRequest{
		Polymorph:           uint32(payload.TypeVCallRequest),
		CallType:            payload.CTConstructor,
		CallFlags:           0,
		CallAsOf:            0,
		Caller:              insolar.Reference{},
		Callee:              gen.Reference(),
		CallSiteDeclaration: prototype,
		CallSiteMethod:      "test",
		CallSequence:        0,
		CallReason:          insolar.Reference{},
		RootTX:              insolar.Reference{},
		CallTX:              insolar.Reference{},
		CallRequestFlags:    0,
		KnownCalleeIncoming: insolar.Reference{},
		EntryHeadHash:       nil,
		CallOutgoing:        reference.Local{},
		Arguments:           nil,
	}

	testIsDone := make(chan struct{}, 0)

	err = pulseConveyor.AddInput(ctx, pd.PulseNumber, &pl)
	require.NoError(t, err)
	<-testIsDone
}
