// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/assured-ledger/blob/master/LICENSE.md.

package proc_test

import (
	"context"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"github.com/insolar/assured-ledger/ledger-core/v2/insolar"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar/bus"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar/flow"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar/gen"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar/jet"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar/payload"
	insolarPulse "github.com/insolar/assured-ledger/ledger-core/v2/insolar/pulse"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar/record"
	"github.com/insolar/assured-ledger/ledger-core/v2/instrumentation/inslogger"
	"github.com/insolar/assured-ledger/ledger-core/v2/ledger/drop"
	"github.com/insolar/assured-ledger/ledger-core/v2/ledger/light/executor"
	"github.com/insolar/assured-ledger/ledger-core/v2/ledger/light/proc"
	"github.com/insolar/assured-ledger/ledger-core/v2/ledger/object"
	"github.com/insolar/assured-ledger/ledger-core/v2/pulse"
)

func TestHotObjects_Proceed(t *testing.T) {
	ctx := flow.TestContextWithPulse(inslogger.TestContext(t), pulse.MinTimePulse+10)
	mc := minimock.NewController(t)

	var (
		drops       *drop.ModifierMock
		indexes     *object.MemoryIndexModifierMock
		jetStorage  *jet.StorageMock
		jetFetcher  *executor.JetFetcherMock
		jetReleaser *executor.JetReleaserMock
		coordinator *jet.CoordinatorMock
		calculator  *insolarPulse.CalculatorMock
		sender      *bus.SenderMock
		registry    *executor.MetricsRegistryMock
	)

	setup := func(mc minimock.MockController) {
		drops = drop.NewModifierMock(mc)
		indexes = object.NewMemoryIndexModifierMock(mc)
		jetStorage = jet.NewStorageMock(mc)
		jetFetcher = executor.NewJetFetcherMock(mc)
		jetReleaser = executor.NewJetReleaserMock(mc)
		coordinator = jet.NewCoordinatorMock(mc)
		calculator = insolarPulse.NewCalculatorMock(mc)
		sender = bus.NewSenderMock(mc)
		registry = executor.NewMetricsRegistryMock(mc)
	}

	t.Run("basic ok", func(t *testing.T) {
		setup(mc)
		defer mc.Finish()

		expectedPulse := insolar.Pulse{
			PulseNumber: pulse.MinTimePulse + 10,
		}
		expectedJetID := gen.JetID()
		expectedObjJetID := expectedJetID
		meta := payload.Meta{}
		expectedDrop := drop.Drop{
			Pulse: expectedPulse.PulseNumber,
			JetID: expectedJetID,
			Split: false,
		}
		idxs := []record.Index{
			{
				ObjID: gen.ID(),
				// this is hack, PendingRecords in record.Index should be always empty
				PendingRecords: []insolar.ID{},
			},
		}

		drops.SetMock.Inspect(func(ctx context.Context, drop drop.Drop) {
			assert.Equal(t, expectedDrop, drop, "didn't set drop")
		}).Return(nil)

		jetStorage.UpdateMock.Inspect(func(ctx context.Context, pulse insolar.PulseNumber, actual bool, ids ...insolar.JetID) {
			assert.Equal(t, expectedPulse.PulseNumber, pulse, "wrong pulse received")
			assert.Equal(t, expectedJetID, ids[0], "wrong jetID received")
		}).Return(nil)

		calculator.BackwardsMock.Return(insolar.Pulse{}, insolarPulse.ErrNotFound)

		jetStorage.ForIDMock.Return(expectedObjJetID, false)

		indexes.SetMock.Inspect(func(ctx context.Context, pn insolar.PulseNumber, index record.Index) {
			assert.Equal(t, idxs[0], index)
		}).Return()

		jetFetcher.ReleaseMock.Inspect(func(ctx context.Context, jetID insolar.JetID, pulse insolar.PulseNumber) {
			assert.Equal(t, expectedJetID, jetID)
		}).Return()
		jetReleaser.UnlockMock.Inspect(func(ctx context.Context, pulse insolar.PulseNumber, jetID insolar.JetID) {
			assert.Equal(t, expectedJetID, jetID)
		}).Return(nil)

		expectedToHeavyMsg, _ := payload.NewMessage(&payload.GotHotConfirmation{
			JetID: expectedJetID,
			Pulse: expectedPulse.PulseNumber,
			Split: expectedDrop.Split,
		})

		sender.SendRoleMock.Inspect(func(ctx context.Context, msg *message.Message, role insolar.DynamicRole, object insolar.Reference) {
			assert.Equal(t, expectedToHeavyMsg.Payload, msg.Payload)
		}).Return(make(chan *message.Message), func() {})

		// start test
		p := proc.NewHotObjects(meta, expectedPulse.PulseNumber, expectedJetID, expectedDrop, idxs, 10)
		p.Dep(drops, indexes, jetStorage, jetFetcher, jetReleaser, coordinator, calculator, sender, registry)

		err := p.Proceed(ctx)
		assert.NoError(t, err)
	})

	t.Run("ok, with pendings", func(t *testing.T) {
		setup(mc)
		defer mc.Finish()

		currentPulse := insolar.Pulse{
			PulseNumber: pulse.MinTimePulse + 100,
		}
		abandonedRequestPulse := insolar.Pulse{
			PulseNumber: pulse.MinTimePulse,
		}
		thresholdAbandonedRequestPulse := insolar.Pulse{
			PulseNumber: pulse.MinTimePulse + 80,
		}

		expectedJetID := gen.JetID()
		expectedObjJetID := expectedJetID
		expectedObjectID := gen.ID()
		meta := payload.Meta{}
		expectedDrop := drop.Drop{
			Pulse: currentPulse.PulseNumber,
			JetID: expectedJetID,
			Split: false,
		}
		idxs := []record.Index{
			{
				ObjID: expectedObjectID,
				Lifeline: record.Lifeline{
					EarliestOpenRequest: &abandonedRequestPulse.PulseNumber,
				},
				// this is hack, PendingRecords in record.Index should be always empty
				PendingRecords: []insolar.ID{},
			},
		}

		drops.SetMock.Inspect(func(ctx context.Context, drop drop.Drop) {
			assert.Equal(t, expectedDrop, drop, "didn't set drop")
		}).Return(nil)

		jetStorage.UpdateMock.Inspect(func(ctx context.Context, pulse insolar.PulseNumber, actual bool, ids ...insolar.JetID) {
			assert.Equal(t, currentPulse.PulseNumber, pulse, "wrong pulse received")
			assert.Equal(t, expectedJetID, ids[0], "wrong jetID received")
		}).Return(nil)

		calculator.BackwardsMock.Return(thresholdAbandonedRequestPulse, nil)
		jetStorage.ForIDMock.Return(expectedObjJetID, false)

		indexes.SetMock.Inspect(func(ctx context.Context, pn insolar.PulseNumber, index record.Index) {
			assert.Equal(t, idxs[0], index)
		}).Return()

		jetFetcher.ReleaseMock.Inspect(func(ctx context.Context, jetID insolar.JetID, pulse insolar.PulseNumber) {
			assert.Equal(t, expectedJetID, jetID)
		}).Return()
		jetReleaser.UnlockMock.Inspect(func(ctx context.Context, pulse insolar.PulseNumber, jetID insolar.JetID) {
			assert.Equal(t, expectedJetID, jetID)
		}).Return(nil)

		expectedToHeavyMsg, _ := payload.NewMessage(&payload.GotHotConfirmation{
			JetID: expectedJetID,
			Pulse: currentPulse.PulseNumber,
			Split: expectedDrop.Split,
		})

		expectedToVirtualMsg, _ := payload.NewMessage(&payload.AbandonedRequestsNotification{
			ObjectID: expectedObjectID,
		})

		sender.SendRoleMock.Inspect(func(ctx context.Context, msg *message.Message, role insolar.DynamicRole, object insolar.Reference) {
			if role == insolar.DynamicRoleHeavyExecutor {
				assert.Equal(t, expectedToHeavyMsg.Payload, msg.Payload)
				return
			} else if role == insolar.DynamicRoleVirtualExecutor {
				assert.Equal(t, expectedToVirtualMsg.Payload, msg.Payload)
				return
			}
			assert.True(t, false, "didn't receive at least 2 messages")
		}).Return(make(chan *message.Message), func() {})

		registry.SetOldestAbandonedRequestAgeMock.Return()

		// start test
		p := proc.NewHotObjects(meta, currentPulse.PulseNumber, expectedJetID, expectedDrop, idxs, 10)
		p.Dep(drops, indexes, jetStorage, jetFetcher, jetReleaser, coordinator, calculator, sender, registry)

		err := p.Proceed(ctx)
		assert.NoError(t, err)
	})

	t.Run("ok, with pending limitation to 0", func(t *testing.T) {
		setup(mc)
		defer mc.Finish()

		currentPulse := insolar.Pulse{
			PulseNumber: pulse.MinTimePulse + 100,
		}
		abandonedRequestPulse := insolar.Pulse{
			PulseNumber: pulse.MinTimePulse,
		}
		thresholdAbandonedRequestPulse := insolar.Pulse{
			PulseNumber: pulse.MinTimePulse + 80,
		}

		expectedJetID := gen.JetID()
		expectedObjJetID := expectedJetID
		expectedObjectID := gen.ID()
		meta := payload.Meta{}
		expectedDrop := drop.Drop{
			Pulse: currentPulse.PulseNumber,
			JetID: expectedJetID,
			Split: false,
		}
		idxs := []record.Index{
			{
				ObjID: expectedObjectID,
				Lifeline: record.Lifeline{
					EarliestOpenRequest: &abandonedRequestPulse.PulseNumber,
				},
				// this is hack, PendingRecords in record.Index should be always empty
				PendingRecords: []insolar.ID{},
			},
		}

		drops.SetMock.Inspect(func(ctx context.Context, drop drop.Drop) {
			assert.Equal(t, expectedDrop, drop, "didn't set drop")
		}).Return(nil)

		jetStorage.UpdateMock.Inspect(func(ctx context.Context, pulse insolar.PulseNumber, actual bool, ids ...insolar.JetID) {
			assert.Equal(t, currentPulse.PulseNumber, pulse, "wrong pulse received")
			assert.Equal(t, expectedJetID, ids[0], "wrong jetID received")
		}).Return(nil)

		calculator.BackwardsMock.Return(thresholdAbandonedRequestPulse, nil)
		jetStorage.ForIDMock.Return(expectedObjJetID, false)

		indexes.SetMock.Inspect(func(ctx context.Context, pn insolar.PulseNumber, index record.Index) {
			assert.Equal(t, idxs[0], index)
		}).Return()

		jetFetcher.ReleaseMock.Inspect(func(ctx context.Context, jetID insolar.JetID, pulse insolar.PulseNumber) {
			assert.Equal(t, expectedJetID, jetID)
		}).Return()
		jetReleaser.UnlockMock.Inspect(func(ctx context.Context, pulse insolar.PulseNumber, jetID insolar.JetID) {
			assert.Equal(t, expectedJetID, jetID)
		}).Return(nil)

		expectedToHeavyMsg, _ := payload.NewMessage(&payload.GotHotConfirmation{
			JetID: expectedJetID,
			Pulse: currentPulse.PulseNumber,
			Split: expectedDrop.Split,
		})

		sender.SendRoleMock.Inspect(func(ctx context.Context, msg *message.Message, role insolar.DynamicRole, object insolar.Reference) {
			assert.Equal(t, expectedToHeavyMsg.Payload, msg.Payload)
			assert.Equal(t, insolar.DynamicRoleHeavyExecutor, role)
		}).Return(make(chan *message.Message), func() {})

		// start test
		p := proc.NewHotObjects(meta, currentPulse.PulseNumber, expectedJetID, expectedDrop, idxs, 0)
		p.Dep(drops, indexes, jetStorage, jetFetcher, jetReleaser, coordinator, calculator, sender, registry)

		err := p.Proceed(ctx)
		assert.NoError(t, err)
	})
}
