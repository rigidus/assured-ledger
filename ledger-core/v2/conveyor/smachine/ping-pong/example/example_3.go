///
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
///

package example

import (
	"fmt"
	"time"

	"github.com/insolar/assured-ledger/ledger-core/v2/vanilla/injector"
	"github.com/insolar/assured-ledger/ledger-core/v2/vanilla/longbits"

	"github.com/insolar/assured-ledger/ledger-core/v2/conveyor/smachine"
)

type StateMachine3 struct {
	serviceA *ServiceAdapterA
	catalogC CatalogC
	catalogD CatalogD

	mutex   smachine.SyncLink
	testKey longbits.ByteString
	waitKey longbits.ByteString
	result  string
	count   int
}

var IterationCount uint64
var Limiter = smachine.NewFixedSemaphore(1000, "global")

/* -------- Declaration ------------- */

var declarationStateMachine3 smachine.StateMachineDeclaration = &StateMachine3Declaration{}

type StateMachine3Declaration struct {
	smachine.StateMachineDeclTemplate
}

func (StateMachine3Declaration) InjectDependencies(sm smachine.StateMachine, _ smachine.SlotLink, injector *injector.DependencyInjector) {
	s := sm.(*StateMachine3)
	injector.MustInject(&s.serviceA)
	injector.MustInject(&s.catalogC)
}

func (StateMachine3Declaration) GetInitStateFor(sm smachine.StateMachine) smachine.InitFunc {
	s := sm.(*StateMachine3)
	return s.Init
}

/* -------- Instance ------------- */

func (s *StateMachine3) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return declarationStateMachine3
}

func (s *StateMachine3) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	s.testKey = longbits.WrapStr("make-pair")
	fmt.Printf("init: %v | %v\n", ctx.SlotLink(), time.Now())
	return ctx.Jump(s.Start)
}

func (s *StateMachine3) Start(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if v, ok := s.catalogC.TryGet(ctx, s.testKey) ; ok {
		myCustomSharedStateAccessor := v
		mySharedDataAccessor := myCustomSharedStateAccessor.Prepare(func(state *CustomSharedState) {
			state.SecondPlayer = s
			key := fmt.Sprintf("%s.%s", state.FirstPlayer, state.SecondPlayer)
			s.waitKey = longbits.ByteString(key)
			fmt.Printf("Start:%p (Shared), Second = %p, Pair=%p \n", s, s, state.FirstPlayer)
		})
		mySharedAccessReport := mySharedDataAccessor.TryUse(ctx)
		return smachine.RepeatOrJumpElse(ctx, mySharedAccessReport, s.WaitForGame, s.Wrong)
	} else {
		myCustomSharedStateAccessor := s.catalogC.GetOrCreate(ctx, s.testKey)
		mySharedDataAccessor := myCustomSharedStateAccessor.Prepare(func(state *CustomSharedState) {
			state.FirstPlayer = s
			fmt.Printf("Start:%p (no Shared), set First = %p \n", s, s)
		})
		mySharedAccessReport := mySharedDataAccessor.TryUse(ctx)
		return smachine.RepeatOrJumpElse(ctx, mySharedAccessReport, s.WaitForSecond, s.Wrong)
	}
}

func (s *StateMachine3) WaitForSecond(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if v, ok := s.catalogC.TryGet(ctx, s.testKey) ; ok {
		myCustomSharedStateAccessor := v
		switcher := false
		mySharedDataAccessor := myCustomSharedStateAccessor.Prepare(func(state *CustomSharedState) {
			if nil != state.SecondPlayer {
				switcher = true
				fmt.Printf("Wait for Second:%p, Second = %p \n", s, state.SecondPlayer)
			} else {
				fmt.Printf("Wait for Second:%p, Waiting = %p \n", s, state.SecondPlayer)
			}
		})
		mySharedAccessReport := mySharedDataAccessor.TryUse(ctx)
		if switcher {
			// TODO: making new SharedObject for Game and BallOwner initialization
			myGameSharedStateAccessor := s.catalogD.GetOrCreate(ctx, s.waitKey)
			mySharedDataAccessor := myGameSharedStateAccessor.Prepare(func(state *GameSharedState) {
				state.BallOwner = s
				fmt.Printf("Start:%p (no Shared), set First = %p \n", s, s)
			})
			mySharedAccessReport := mySharedDataAccessor.TryUse(ctx)



			// Go to Game
			return smachine.RepeatOrJumpElse(ctx, mySharedAccessReport, s.Game, s.Wrong)
		} else {
			// No second player, repeat
			return smachine.RepeatOrJumpElse(ctx, mySharedAccessReport, s.WaitForSecond, s.Wrong)
		}
	} else {
		fmt.Printf("WaitForSecond: Something wrong \n")
		return ctx.Jump(s.Wrong)
	}
}

func (s *StateMachine3) WaitForGame(ctx smachine.ExecutionContext) smachine.StateUpdate {
	fmt.Printf("WaitForGame: %p \n", s)
	if v, ok := s.catalogC.TryGet(ctx, s.waitKey) ; ok {
		fmt.Printf("%v \n", v)
		return ctx.Jump(s.Game)
	} else {
		// repeat
		return ctx.Sleep().ThenRepeat()
		// return ctx.Jump(s.WaitForGame)
	}
}

func (s *StateMachine3) Game(ctx smachine.ExecutionContext) smachine.StateUpdate {
	fmt.Printf("Game %p \n", s)
	return ctx.Jump(s.GameOver)
}

func (s *StateMachine3) GameOver(ctx smachine.ExecutionContext) smachine.StateUpdate {
	fmt.Printf("GameOver %p \n", s)
	return ctx.Jump(s.GameOver)
	//return ctx.Stop()
}

func (s *StateMachine3) Wrong(ctx smachine.ExecutionContext) smachine.StateUpdate {
	return ctx.WaitAnyUntil(time.Now().Add(time.Second)).ThenJump(s.Wrong)
	//return ctx.Stop()
}
