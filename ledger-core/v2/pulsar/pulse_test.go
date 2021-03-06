// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/assured-ledger/blob/master/LICENSE.md.

package pulsar

import (
	"testing"

	"github.com/insolar/assured-ledger/ledger-core/v2/configuration"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar"
	"github.com/insolar/assured-ledger/ledger-core/v2/pulsar/pulsartestutils"
	"github.com/stretchr/testify/require"
)

func TestNewPulse(t *testing.T) {
	generator := &pulsartestutils.MockEntropyGenerator{}
	previousPulse := insolar.PulseNumber(876)
	expectedPulse := previousPulse + insolar.PulseNumber(configuration.NewPulsar().NumberDelta)

	result := NewPulse(configuration.NewPulsar().NumberDelta, previousPulse, generator)

	require.Equal(t, result.Entropy[:], pulsartestutils.MockEntropy[:])
	require.Equal(t, result.PulseNumber, insolar.PulseNumber(expectedPulse))
}
