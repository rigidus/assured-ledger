// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/assured-ledger/blob/master/LICENSE.md.

// +build slowtest

package api_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fortytw2/leaktest"
	"github.com/stretchr/testify/require"

	"github.com/insolar/assured-ledger/ledger-core/v2/application/api"
	"github.com/insolar/assured-ledger/ledger-core/v2/configuration"
	"github.com/insolar/assured-ledger/ledger-core/v2/instrumentation/inslogger"
)

func waitForStatus(t *testing.T, nc *api.NetworkChecker, expected bool) {
	ctx := context.Background()
	var available bool
	for i := 0; i < 10; i++ {
		available = nc.IsAvailable(ctx)
		if available == expected {
			return
		}
		time.Sleep(time.Second)
	}
	require.Fail(t, "Status not passed, expected: ", expected)
}

func TestAvailabilityChecker_UpdateStatus(t *testing.T) {
	cfg := configuration.NewLog()
	cfg.Level = "Debug"
	ctx, _ := inslogger.InitNodeLogger(context.Background(), cfg, "", "")

	defer leaktest.Check(t)()

	counter := 0

	keeper := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch counter {
			case 0:
				w.WriteHeader(http.StatusOK)
				_, err := fmt.Fprintln(w, "{\"available\": false}")
				require.NoError(t, err)
			case 1:
				w.WriteHeader(http.StatusBadRequest)
			case 2:
				w.WriteHeader(http.StatusOK)
				_, err := fmt.Fprintln(w, "{\"test\": }")
				require.NoError(t, err)
			default:
				w.WriteHeader(http.StatusOK)
				_, err := fmt.Fprintln(w, "{\"available\": true}")
				require.NoError(t, err)
			}
			counter += 1
		}))

	var checkPeriod uint = 1
	config := configuration.AvailabilityChecker{
		Enabled:        true,
		KeeperURL:      keeper.URL,
		RequestTimeout: 2,
		CheckPeriod:    checkPeriod,
	}

	nc := api.NewNetworkChecker(config)
	require.False(t, nc.IsAvailable(ctx))

	defer nc.Stop()
	err := nc.Start(ctx)
	require.NoError(t, err)

	// counter = 0
	require.False(t, nc.IsAvailable(ctx))
	time.Sleep(time.Duration(checkPeriod))

	// counter = 1
	require.False(t, nc.IsAvailable(ctx))
	time.Sleep(time.Duration(checkPeriod))

	// counter = 2, bad response body
	require.False(t, nc.IsAvailable(ctx))

	// counter default
	waitForStatus(t, nc, true)

	keeper.Close()
	waitForStatus(t, nc, false)
}
