/*
 * Copyright 2020 Insolar Network Ltd.
 * All rights reserved.
 * This material is licensed under the Insolar License version 1.0,
 * available at https://github.com/insolar/assured-ledger/blob/master/LICENSE.md.
 */

package smachine

import (
	"github.com/insolar/assured-ledger/ledger-core/v2/vanilla/throw"
)

type subroutineMarker struct {
	slotId SlotID
	_      [0]func() // force incomparable by value
}

func (s *Slot) newSubroutineMarker() *subroutineMarker {
	return &subroutineMarker{slotId: s.GetSlotID()}
}

func (s *Slot) getSubroutineMarker() *subroutineMarker {
	if s.subroutineStack != nil {
		return s.subroutineStack.childMarker
	}
	return nil
}

func (s *Slot) ensureNoSubroutine() {
	if s.hasSubroutine() {
		panic(throw.IllegalState())
	}
}

func (m *SlotMachine) handleUnsafeSubroutineUpdate(su StateUpdate, producedBy *subroutineMarker) StateUpdate {
	panic(throw.NotImplemented())
	// TODO handleUnsafeSubroutineUpdate
}

func (s *Slot) hasSubroutine() bool {
	return s.subroutineStack != nil
}

func (s *Slot) checkSubroutineMarker(marker *subroutineMarker) (isCurrent, isValid bool) {
	switch {
	case marker == nil:
		return s.subroutineStack == nil, true
	case s.subroutineStack == nil:
		//
	case s.GetSlotID() != marker.slotId:
		//
	case s.subroutineStack.childMarker == marker:
		return true, true
	default:
		for sbs := s.subroutineStack.subroutineStack; sbs != nil; sbs = sbs.subroutineStack {
			if sbs.childMarker == marker {
				return false, true
			}
		}
	}
	return false, false
}

func (s *Slot) prepareSubroutineStart(ssm SubroutineStateMachine, exitFn SubroutineExitFunc) SlotStep {
	if exitFn == nil {
		panic(throw.IllegalValue())
	}
	if ssm == nil {
		panic(throw.IllegalValue())
	}
	decl := ssm.GetStateMachineDeclaration()
	if decl == nil {
		panic(throw.IllegalValue())
	}
	initFn := ssm.GetSubroutineInitState()
	if initFn == nil {
		panic(throw.IllegalValue())
	}

	return SlotStep{Transition: func(ctx ExecutionContext) StateUpdate {
		ec := ctx.(*executionContext)
		slot := ec.s
		m := slot.machine

		prev := slot.slotDeclarationData
		stackedSdd := &slotStackedDeclarationData{prev, prev.defMigrate,
			exitFn, slot.newSubroutineMarker(),
			prev.subroutineStack != nil && prev.subroutineStack.hasMigrates,
		}
		if stackedSdd.stackMigrate != nil {
			stackedSdd.hasMigrates = true
		}

		inheritable := slot.inheritable
		slot.slotDeclarationData.subroutineStack = stackedSdd
		slot.slotDeclarationData.declaration = decl
		slot.slotDeclarationData.shadowMigrate = nil
		slot.slotDeclarationData.defTerminate = nil

		// get injects sorted out
		var localInjects []interface{}
		slot.inheritable, localInjects = m.prepareInjects(nil, slot.NewLink(), ssm, InheritResolvedDependencies,
			true, nil, inheritable)

		// Step Logger
		m.prepareStepLogger(slot, ssm, ctx.Log().GetTracerId())

		// shadow migrate for injected dependencies
		slot.shadowMigrate = buildShadowMigrator(localInjects, slot.declaration.GetShadowMigrateFor(ssm))

		return initFn.defaultInit(ctx)
	}}
}

var defaultSubroutineStartDecl = StepDeclaration{stepDeclExt: stepDeclExt{Name: "<subroutineStart>"}}
var defaultSubroutineExitDecl = StepDeclaration{stepDeclExt: stepDeclExt{Name: "<subroutineExit>"}}

func (s *Slot) prepareSubroutineExit(err error, _ FixedSlotWorker) {
	prev := s.subroutineStack
	if prev == nil {
		panic(throw.IllegalState())
	}

	td := TerminationData{
		Slot:   s.NewStepLink(),
		Parent: s.parent,
		Result: s.defResult,
		Error:  err,
		worker: nil, // not applicable for subroutine exit
	}

	returnFn := prev.returnFn
	s.slotDeclarationData = prev.slotDeclarationData
	s.step = SlotStep{Transition: func(ctx ExecutionContext) StateUpdate {
		return returnFn(ctx, td)
	}}
	s.stepDecl = &defaultSubroutineExitDecl
}
