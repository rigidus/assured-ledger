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
	waitKey string //longbits.ByteString
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
	// Looking for make-pair shared object
	if v, ok := s.catalogC.TryGet(ctx, s.testKey) ; ok {
		// Make-pair found
		mySharedAccessReport := v.PrepareWithSecondPlayer(func(state *CustomSharedState) {
			// First Player is already there, I will be Second Player
			// Key for Game Object
			key := fmt.Sprintf("%p.%p", state.FirstPlayer, state.SecondPlayer)
			s.waitKey = key //longbits.ByteString(key)
		},s ).TryUse(ctx)
		return smachine.RepeatOrJumpElse(ctx, mySharedAccessReport, s.WaitForGame, s.Wrong)
	} else {
		// Make-pair not found
		report := s.catalogC.GetOrCreate(ctx, s.testKey).PrepareWithFirstPlayer(func(state *CustomSharedState) {
		}, s).TryUse(ctx)
		return smachine.RepeatOrJumpElse(ctx, report, s.WaitForSecond, s.Wrong)
	}
}

func (s *StateMachine3) WaitForSecond(ctx smachine.ExecutionContext) smachine.StateUpdate {
	// Waiting for Second Player on Make-pair object
	if v, ok := s.catalogC.TryGet(ctx, s.testKey) ; ok {
		switcher := false // Default: No Second player
		mySharedAccessReport := v.Prepare(func(state *CustomSharedState) {
			if "" != state.SecondPlayer {
				switcher = true // Second player is here!
				s.waitKey = fmt.Sprintf("%p.%p", state.FirstPlayer, state.SecondPlayer)
			} else {
			}
		}).TryUse(ctx)
		if switcher {
			// making new SharedObject for Game and BallOwner initialization
			myGameSharedStateAccessor := s.catalogD.GetOrCreate(ctx, longbits.ByteString(s.waitKey))
			mySharedAccessReport := myGameSharedStateAccessor.Prepare(func(state *GameSharedState) {
				state.BallOwner = s // I will be owner
				state.cnt = 0
			}).TryUse(ctx)
			// Go to Game
			return smachine.RepeatOrJumpElse(ctx, mySharedAccessReport, s.Game, s.Wrong)
		} else {
			// No second player, repeat
			return smachine.RepeatOrJumpElse(ctx, mySharedAccessReport, s.WaitForSecond, s.Wrong)
		}
	} else {
		fmt.Printf("Error (WaitForSecond): Something wrong \n")
		return ctx.Jump(s.Wrong)
	}
}

func (s *StateMachine3) WaitForGame(ctx smachine.ExecutionContext) smachine.StateUpdate {
	// Waiting for Game Object
	// fmt.Printf("WaitForGame: %s \n", s.waitKey)
	if _, ok := s.catalogC.TryGet(ctx, longbits.ByteString(s.waitKey)) ; ok {
		return ctx.Jump(s.Game)
	} else {
		// repeat
		// return ctx.Sleep().ThenRepeat() // - Does not works, why? [TODO]
		time.After(2*time.Second) // ...same
		return ctx.Jump(s.WaitForGame)
	}
}

func (s *StateMachine3) Game(ctx smachine.ExecutionContext) smachine.StateUpdate {
	fmt.Printf("Game %p \n", s)
	if v, ok := s.catalogD.TryGet(ctx, longbits.ByteString(s.waitKey)) ; ok {
		endgame := false
		mySharedAccessReport := v.Prepare(func(state *GameSharedState) {
			if s != state.BallOwner {
				state.BallOwner = s
			}
			state.cnt += 1
			if state.cnt > 2 {
				endgame = true
			}
		}).TryUse(ctx)
		if endgame {
			// Go to GameOver
			return smachine.RepeatOrJumpElse(ctx, mySharedAccessReport, s.GameOver, s.Wrong)
		} else {
			// repeat
			return smachine.RepeatOrJumpElse(ctx, mySharedAccessReport, s.Game, s.Wrong)
		}
	} else {
		// Impossible (let it be for dbg)
		return ctx.Jump(s.Wrong)
	}
}

func (s *StateMachine3) GameOver(ctx smachine.ExecutionContext) smachine.StateUpdate {
	return ctx.Stop()
	//return ctx.Jump(s.GameOver)
	//return ctx.Stop()
}

func (s *StateMachine3) Wrong(ctx smachine.ExecutionContext) smachine.StateUpdate {
	fmt.Printf("ErrState: WRONG!!! %p \n", s)
	//return ctx.WaitAnyUntil(time.Now().Add(time.Second)).ThenJump(s.Wrong)
	return ctx.Stop()
}
