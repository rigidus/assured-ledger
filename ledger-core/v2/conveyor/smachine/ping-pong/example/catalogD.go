//
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package example

import (
	"fmt"

	"github.com/insolar/assured-ledger/ledger-core/v2/conveyor/smachine"
	"github.com/insolar/assured-ledger/ledger-core/v2/vanilla/longbits"
)

func CreateCatalogD() CatalogD {
	return &catalogD{}
}

type CatalogD = *catalogD

type GameSharedState struct {
	key         longbits.ByteString
	BallOwner   *StateMachine3
	cnt			int
}

func (p *GameSharedState) GetKey() longbits.ByteString {
	return p.key
}

type GameSharedStateAccessor struct {
	link smachine.SharedDataLink
}

func (v GameSharedStateAccessor) Prepare(fn func(*GameSharedState)) smachine.SharedDataAccessor {
	return v.link.PrepareAccess(func(data interface{}) bool {
		fn(data.(*GameSharedState))
		return false
	})
}

type catalogD struct {
}

func (p *catalogD) Get(ctx smachine.ExecutionContext, key longbits.ByteString) GameSharedStateAccessor {
	if v, ok := p.TryGet(ctx, key); ok {
		return v
	}
	panic(fmt.Sprintf("missing entry: %s", key))
}

func (p *catalogD) TryGet(ctx smachine.ExecutionContext, key longbits.ByteString) (GameSharedStateAccessor, bool) {

	if v := ctx.GetPublishedLink(key); v.IsAssignableTo((*GameSharedState)(nil)) {
		return GameSharedStateAccessor{v}, true
	}
	return GameSharedStateAccessor{}, false
}

func (p *catalogD) GetOrCreate(ctx smachine.ExecutionContext, key longbits.ByteString) GameSharedStateAccessor {
	if v, ok := p.TryGet(ctx, key); ok {
		return v
	}

	ctx.InitChild(func(ctx smachine.ConstructionContext) smachine.StateMachine {
		return &catalogEntryDSM{sharedState: GameSharedState{
			key: key,
			//Mutex: smachine.NewExclusiveWithFlags("", 0), //smachine.QueueAllowsPriority),
			//Mutex: smachine.NewSemaphoreWithFlags(2, "", smachine.QueueAllowsPriority).SyncLink(),
		}}
	})

	return p.Get(ctx, key)
}

type catalogEntryDSM struct {
	smachine.StateMachineDeclTemplate
	sharedState GameSharedState
}



func (sm *catalogEntryDSM) GetInitStateFor(smachine.StateMachine) smachine.InitFunc {
	return sm.Init
}

func (sm *catalogEntryDSM) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return sm
}

func (sm *catalogEntryDSM) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	sdl := ctx.Share(&sm.sharedState, 0)
	if !ctx.Publish(sm.sharedState.key, sdl) {
		return ctx.Stop()
	}
	return ctx.JumpExt(smachine.SlotStep{Transition: sm.State1, Flags: smachine.StepWeak})
}

func (sm *catalogEntryDSM) State1(ctx smachine.ExecutionContext) smachine.StateUpdate {
	return ctx.Sleep().ThenRepeat()
}
