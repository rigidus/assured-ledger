// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/assured-ledger/blob/master/LICENSE.md.

package gateway

import (
	"context"
	"testing"
	"time"

	"github.com/insolar/assured-ledger/ledger-core/v2/insolar/gen"
	"github.com/insolar/assured-ledger/ledger-core/v2/network/node"
	"github.com/stretchr/testify/require"

	"github.com/gojuno/minimock/v3"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar"
	"github.com/insolar/assured-ledger/ledger-core/v2/network"
	"github.com/insolar/assured-ledger/ledger-core/v2/pulse"
	mock "github.com/insolar/assured-ledger/ledger-core/v2/testutils/network"
	"github.com/stretchr/testify/assert"
)

func createBase(mc *minimock.Controller) *Base {
	b := &Base{}

	op := mock.NewOriginProviderMock(mc)
	op.GetOriginMock.Set(func() insolar.NetworkNode {
		return node.NewNode(gen.Reference(), insolar.StaticRoleVirtual, nil, "127.0.0.1:123", "")
	})

	aborter := network.NewAborterMock(mc)
	aborter.AbortMock.Set(func(ctx context.Context, reason string) {
		require.Contains(mc, reason, bootstrapTimeoutMessage)
	})

	b.Aborter = aborter
	b.OriginProvider = op
	return b
}

func TestWaitPulsar_PulseNotArrivedInETA(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()
	defer mc.Wait(time.Minute)

	waitPulsar := newWaitPulsar(createBase(mc))
	assert.Equal(t, insolar.WaitPulsar, waitPulsar.GetState())
	gatewayer := mock.NewGatewayerMock(mc)
	waitPulsar.Gatewayer = gatewayer
	gatewayer.GatewayMock.Set(func() network.Gateway {
		return waitPulsar
	})

	waitPulsar.bootstrapETA = time.Millisecond
	waitPulsar.bootstrapTimer = time.NewTimer(waitPulsar.bootstrapETA)

	waitPulsar.Run(context.Background(), *insolar.EphemeralPulse)
}

func TestWaitPulsar_PulseArrivedInETA(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()
	defer mc.Wait(time.Minute)

	gatewayer := mock.NewGatewayerMock(mc)
	gatewayer.SwitchStateMock.Set(func(ctx context.Context, state insolar.NetworkState, pulse insolar.Pulse) {
		assert.Equal(t, insolar.CompleteNetworkState, state)
	})

	pulseAccessor := mock.NewPulseAccessorMock(mc)
	pulseAccessor.GetPulseMock.Set(func(ctx context.Context, p1 insolar.PulseNumber) (p2 insolar.Pulse, err error) {
		p := *insolar.GenesisPulse
		p.PulseNumber += 10
		return p, nil
	})

	waitPulsar := newWaitPulsar(&Base{
		PulseAccessor: pulseAccessor,
	})
	waitPulsar.Gatewayer = gatewayer
	waitPulsar.bootstrapETA = time.Second * 2
	waitPulsar.bootstrapTimer = time.NewTimer(waitPulsar.bootstrapETA)

	go waitPulsar.Run(context.Background(), *insolar.EphemeralPulse)
	time.Sleep(100 * time.Millisecond)

	waitPulsar.OnConsensusFinished(context.Background(), network.Report{PulseNumber: pulse.MinTimePulse + 10})
}
