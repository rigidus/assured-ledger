// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/assured-ledger/blob/master/LICENSE.md.

package gateway

import (
	"context"

	"github.com/insolar/assured-ledger/ledger-core/v2/insolar"
	"github.com/insolar/assured-ledger/ledger-core/v2/network"
)

func newWaitConsensus(b *Base) *WaitConsensus {
	return &WaitConsensus{b, make(chan insolar.Pulse, 1)}
}

type WaitConsensus struct {
	*Base

	consensusFinished chan insolar.Pulse
}

func (g *WaitConsensus) Run(ctx context.Context, pulse insolar.Pulse) {
	select {
	case <-g.bootstrapTimer.C:
		g.FailState(ctx, bootstrapTimeoutMessage)
	case newPulse := <-g.consensusFinished:
		g.Gatewayer.SwitchState(ctx, insolar.WaitMajority, newPulse)
	}
}

func (g *WaitConsensus) GetState() insolar.NetworkState {
	return insolar.WaitConsensus
}

func (g *WaitConsensus) OnConsensusFinished(ctx context.Context, report network.Report) {
	g.consensusFinished <- EnsureGetPulse(ctx, g.PulseAccessor, report.PulseNumber)
	close(g.consensusFinished)
}
