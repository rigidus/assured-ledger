//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package logoutput

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/insolar/assured-ledger/ledger-core/v2/log/logcommon"
	"github.com/insolar/assured-ledger/ledger-core/v2/log/outputsyslog"
)

type LogOutput uint8

const (
	StdErrOutput LogOutput = iota
	SysLogOutput
)

func (l LogOutput) String() string {
	switch l {
	case StdErrOutput:
		return "stderr"
	case SysLogOutput:
		return "syslog"
	}
	return string(l)
}

func OpenLogBareOutput(output LogOutput, param string) (logcommon.BareOutput, error) {
	switch output {
	case StdErrOutput:
		w := os.Stderr
		return logcommon.BareOutput{
			Writer:         w,
			FlushFn:        w.Sync,
			ProtectedClose: true,
		}, nil
	case SysLogOutput:
		executableName := filepath.Base(os.Args[0])
		w, err := outputsyslog.ConnectSyslogByParam(param, executableName)
		if err != nil {
			return logcommon.BareOutput{}, err
		}
		return logcommon.BareOutput{
			Writer:         w,
			FlushFn:        w.Flush,
			ProtectedClose: false,
		}, nil
	default:
		return logcommon.BareOutput{}, errors.New("unknown output " + output.String())
	}
}