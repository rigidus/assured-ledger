// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/assured-ledger/blob/master/LICENSE.md.

package gateway

import (
	"context"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/insolar/assured-ledger/ledger-core/v2/certificate"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar/gen"
	"github.com/insolar/assured-ledger/ledger-core/v2/network"
	"github.com/insolar/assured-ledger/ledger-core/v2/network/node"
	"github.com/insolar/assured-ledger/ledger-core/v2/pulse"
	mock "github.com/insolar/assured-ledger/ledger-core/v2/testutils/network"
	"github.com/stretchr/testify/assert"
)

func TestWaitMajority_MajorityNotHappenedInETA(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()
	defer mc.Wait(time.Minute)

	nodeKeeper := mock.NewNodeKeeperMock(mc)
	nodeKeeper.GetAccessorMock.Set(func(p1 insolar.PulseNumber) (a1 network.Accessor) {
		accessor := mock.NewAccessorMock(mc)
		accessor.GetWorkingNodesMock.Set(func() (na1 []insolar.NetworkNode) {
			return []insolar.NetworkNode{}
		})
		return accessor
	})

	cert := &certificate.Certificate{MajorityRule: 4}

	b := createBase(mc)
	b.CertificateManager = certificate.NewCertificateManager(cert)
	b.NodeKeeper = nodeKeeper

	waitMajority := newWaitMajority(b)
	assert.Equal(t, insolar.WaitMajority, waitMajority.GetState())
	gatewayer := mock.NewGatewayerMock(mc)
	gatewayer.GatewayMock.Set(func() network.Gateway {
		return waitMajority
	})
	waitMajority.Gatewayer = gatewayer
	waitMajority.bootstrapETA = time.Millisecond
	waitMajority.bootstrapTimer = time.NewTimer(waitMajority.bootstrapETA)

	waitMajority.Run(context.Background(), *insolar.EphemeralPulse)
}

func TestWaitMajority_MajorityHappenedInETA(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()
	defer mc.Wait(time.Minute)

	gatewayer := mock.NewGatewayerMock(mc)
	gatewayer.SwitchStateMock.Set(func(ctx context.Context, state insolar.NetworkState, pulse insolar.Pulse) {
		assert.Equal(t, insolar.WaitMinRoles, state)
	})

	ref := gen.Reference()
	nodeKeeper := mock.NewNodeKeeperMock(mc)
	accessor1 := mock.NewAccessorMock(mc)
	accessor1.GetWorkingNodesMock.Set(func() (na1 []insolar.NetworkNode) {
		return []insolar.NetworkNode{}
	})
	accessor2 := mock.NewAccessorMock(mc)
	accessor2.GetWorkingNodesMock.Set(func() (na1 []insolar.NetworkNode) {
		n := node.NewNode(ref, insolar.StaticRoleHeavyMaterial, nil, "127.0.0.1:123", "")
		return []insolar.NetworkNode{n}
	})
	nodeKeeper.GetAccessorMock.Set(func(p insolar.PulseNumber) (a1 network.Accessor) {
		if p == pulse.MinTimePulse {
			return accessor1
		}
		return accessor2
	})

	discoveryNode := certificate.BootstrapNode{NodeRef: ref.String()}
	cert := &certificate.Certificate{MajorityRule: 1, BootstrapNodes: []certificate.BootstrapNode{discoveryNode}}
	pulseAccessor := mock.NewPulseAccessorMock(mc)
	pulseAccessor.GetPulseMock.Set(func(ctx context.Context, p1 insolar.PulseNumber) (p2 insolar.Pulse, err error) {
		p := *insolar.GenesisPulse
		p.PulseNumber += 10
		return p, nil
	})
	waitMajority := newWaitMajority(&Base{
		CertificateManager: certificate.NewCertificateManager(cert),
		NodeKeeper:         nodeKeeper,
		PulseAccessor:      pulseAccessor,
	})
	waitMajority.Gatewayer = gatewayer
	waitMajority.bootstrapETA = time.Second * 2
	waitMajority.bootstrapTimer = time.NewTimer(waitMajority.bootstrapETA)

	go waitMajority.Run(context.Background(), *insolar.EphemeralPulse)
	time.Sleep(100 * time.Millisecond)

	waitMajority.OnConsensusFinished(context.Background(), network.Report{PulseNumber: pulse.MinTimePulse + 10})
}
