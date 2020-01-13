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

package sm_request

import (
	"github.com/insolar/assured-ledger/ledger-core/v2/conveyor/smachine"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar/payload"
	"github.com/insolar/assured-ledger/ledger-core/v2/vanilla/injector"
)

type StateMachineExecutorResults struct {
	// input arguments
	Meta    *payload.Meta
	Payload *payload.ExecutorResults
}

var declExecutorResults smachine.StateMachineDeclaration = &declarationExecutorResults{}

type declarationExecutorResults struct {
	smachine.StateMachineDeclTemplate
}

/* -------- Declaration ------------- */

func (declarationExecutorResults) GetInitStateFor(sm smachine.StateMachine) smachine.InitFunc {
	s := sm.(*StateMachineExecutorResults)
	return s.Init
}

func (declarationExecutorResults) InjectDependencies(sm smachine.StateMachine, _ smachine.SlotLink, injector *injector.DependencyInjector) {
	_ = sm.(*StateMachineExecutorResults)
}

/* -------- Instance ------------- */

func (s *StateMachineExecutorResults) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return declExecutorResults
}

func (s *StateMachineExecutorResults) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	return ctx.Stop()
}
