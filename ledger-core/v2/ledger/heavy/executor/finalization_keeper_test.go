// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/assured-ledger/blob/master/LICENSE.md.

package executor

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/insolar/assured-ledger/ledger-core/v2/insolar"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar/pulse"
	insolarPulse "github.com/insolar/assured-ledger/ledger-core/v2/insolar/pulse"
)

func TestFinalizationKeeper_WeAreTooYoung(t *testing.T) {
	testPulse := insolar.GenesisPulse.PulseNumber
	jkMock := NewJetKeeperMock(t)
	jkMock.TopSyncPulseMock.Expect().Return(testPulse + 1)

	calcMock := insolarPulse.NewCalculatorMock(t)
	calcMock.BackwardsMock.Set(func(p context.Context, p1 insolar.PulseNumber, p2 int) (r insolar.Pulse, r1 error) {
		require.Equal(t, testPulse, p1)

		return insolar.Pulse{}, pulse.ErrNotFound
	})

	fk := NewFinalizationKeeperDefault(jkMock, calcMock, 100)
	err := fk.OnPulse(context.Background(), testPulse)
	require.NoError(t, err)
}

func TestFinalizationKeeper_CalculatorReturnError(t *testing.T) {
	testPulse := insolar.GenesisPulse.PulseNumber
	jkMock := NewJetKeeperMock(t)
	jkMock.TopSyncPulseMock.Expect().Return(testPulse + 1)

	testError := errors.New("Test_CalculatorReturnError")

	calcMock := insolarPulse.NewCalculatorMock(t)
	calcMock.BackwardsMock.Set(func(p context.Context, p1 insolar.PulseNumber, p2 int) (r insolar.Pulse, r1 error) {
		require.Equal(t, testPulse, p1)

		return insolar.Pulse{}, testError
	})

	fk := NewFinalizationKeeperDefault(jkMock, calcMock, 100)
	err := fk.OnPulse(context.Background(), testPulse)
	require.Contains(t, err.Error(), testError.Error())
}

func TestFinalizationKeeper_OldCurrentPulse(t *testing.T) {
	testPulse := insolar.GenesisPulse.PulseNumber
	jkMock := NewJetKeeperMock(t)
	jkMock.TopSyncPulseMock.Expect().Return(testPulse + 1)

	limit := 100

	calcMock := insolarPulse.NewCalculatorMock(t)
	calcMock.BackwardsMock.Return(insolar.Pulse{PulseNumber: testPulse + insolar.PulseNumber(limit)}, nil)

	fk := NewFinalizationKeeperDefault(jkMock, calcMock, limit)
	err := fk.OnPulse(context.Background(), testPulse)
	require.EqualError(t, err, "Current pulse ( 65537 ) is less than last confirmed ( 65538 )")
}

func TestFinalizationKeeper_LimitExceeded(t *testing.T) {

	// Remove require panic, when INS-3121 is fixed
	testBody := func() {
		testPulse := insolar.GenesisPulse.PulseNumber
		limit := 10
		jkMock := NewJetKeeperMock(t)
		jkMock.TopSyncPulseMock.Expect().Return(testPulse)

		calcMock := insolarPulse.NewCalculatorMock(t)
		calcMock.BackwardsMock.Set(func(p context.Context, p1 insolar.PulseNumber, p2 int) (r insolar.Pulse, r1 error) {
			return insolar.Pulse{PulseNumber: p1 - insolar.PulseNumber(p2)}, nil
		})

		fk := NewFinalizationKeeperDefault(jkMock, calcMock, limit)
		err := fk.OnPulse(context.Background(), testPulse+insolar.PulseNumber(limit*10))
		require.Contains(t, err.Error(), "last finalized pulse falls behind too much")
	}
	require.Panics(t, testBody)

}

func TestFinalizationKeeper_HappyPath(t *testing.T) {
	testPulse := insolar.GenesisPulse.PulseNumber
	limit := 10
	jkMock := NewJetKeeperMock(t)
	jkMock.TopSyncPulseMock.Expect().Return(testPulse)

	calcMock := insolarPulse.NewCalculatorMock(t)
	calcMock.BackwardsMock.Return(insolar.Pulse{PulseNumber: testPulse - 1}, nil)

	fk := NewFinalizationKeeperDefault(jkMock, calcMock, limit)
	err := fk.OnPulse(context.Background(), testPulse+insolar.PulseNumber(limit))
	require.NoError(t, err)
}
