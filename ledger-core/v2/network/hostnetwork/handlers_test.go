// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/assured-ledger/blob/master/LICENSE.md.

package hostnetwork

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/fortytw2/leaktest"

	"github.com/insolar/assured-ledger/ledger-core/v2/instrumentation/inslogger"
	"github.com/insolar/assured-ledger/ledger-core/v2/network/hostnetwork/packet"
)

func TestNewStreamHandler(t *testing.T) {
	defer leaktest.Check(t)()

	requestHandler := func(ctx context.Context, p *packet.ReceivedPacket) {
		inslogger.FromContext(ctx).Info("requestHandler")
	}

	h := NewStreamHandler(requestHandler, nil)

	con1, _ := net.Pipe()

	done := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		h.HandleStream(ctx, "127.0.0.1:8080", con1)
		done <- struct{}{}
	}()

	cancel()
	// con2.Close()

	select {
	case <-done:
		return
	case <-time.After(time.Second * 5):
		t.Fail()
	}
}
