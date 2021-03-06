package network

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar"
	mm_network "github.com/insolar/assured-ledger/ledger-core/v2/network"
)

// GatewayerMock implements network.Gatewayer
type GatewayerMock struct {
	t minimock.Tester

	funcGateway          func() (g1 mm_network.Gateway)
	inspectFuncGateway   func()
	afterGatewayCounter  uint64
	beforeGatewayCounter uint64
	GatewayMock          mGatewayerMockGateway

	funcSwitchState          func(ctx context.Context, state insolar.NetworkState, pulse insolar.Pulse)
	inspectFuncSwitchState   func(ctx context.Context, state insolar.NetworkState, pulse insolar.Pulse)
	afterSwitchStateCounter  uint64
	beforeSwitchStateCounter uint64
	SwitchStateMock          mGatewayerMockSwitchState
}

// NewGatewayerMock returns a mock for network.Gatewayer
func NewGatewayerMock(t minimock.Tester) *GatewayerMock {
	m := &GatewayerMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GatewayMock = mGatewayerMockGateway{mock: m}

	m.SwitchStateMock = mGatewayerMockSwitchState{mock: m}
	m.SwitchStateMock.callArgs = []*GatewayerMockSwitchStateParams{}

	return m
}

type mGatewayerMockGateway struct {
	mock               *GatewayerMock
	defaultExpectation *GatewayerMockGatewayExpectation
	expectations       []*GatewayerMockGatewayExpectation
}

// GatewayerMockGatewayExpectation specifies expectation struct of the Gatewayer.Gateway
type GatewayerMockGatewayExpectation struct {
	mock *GatewayerMock

	results *GatewayerMockGatewayResults
	Counter uint64
}

// GatewayerMockGatewayResults contains results of the Gatewayer.Gateway
type GatewayerMockGatewayResults struct {
	g1 mm_network.Gateway
}

// Expect sets up expected params for Gatewayer.Gateway
func (mmGateway *mGatewayerMockGateway) Expect() *mGatewayerMockGateway {
	if mmGateway.mock.funcGateway != nil {
		mmGateway.mock.t.Fatalf("GatewayerMock.Gateway mock is already set by Set")
	}

	if mmGateway.defaultExpectation == nil {
		mmGateway.defaultExpectation = &GatewayerMockGatewayExpectation{}
	}

	return mmGateway
}

// Inspect accepts an inspector function that has same arguments as the Gatewayer.Gateway
func (mmGateway *mGatewayerMockGateway) Inspect(f func()) *mGatewayerMockGateway {
	if mmGateway.mock.inspectFuncGateway != nil {
		mmGateway.mock.t.Fatalf("Inspect function is already set for GatewayerMock.Gateway")
	}

	mmGateway.mock.inspectFuncGateway = f

	return mmGateway
}

// Return sets up results that will be returned by Gatewayer.Gateway
func (mmGateway *mGatewayerMockGateway) Return(g1 mm_network.Gateway) *GatewayerMock {
	if mmGateway.mock.funcGateway != nil {
		mmGateway.mock.t.Fatalf("GatewayerMock.Gateway mock is already set by Set")
	}

	if mmGateway.defaultExpectation == nil {
		mmGateway.defaultExpectation = &GatewayerMockGatewayExpectation{mock: mmGateway.mock}
	}
	mmGateway.defaultExpectation.results = &GatewayerMockGatewayResults{g1}
	return mmGateway.mock
}

//Set uses given function f to mock the Gatewayer.Gateway method
func (mmGateway *mGatewayerMockGateway) Set(f func() (g1 mm_network.Gateway)) *GatewayerMock {
	if mmGateway.defaultExpectation != nil {
		mmGateway.mock.t.Fatalf("Default expectation is already set for the Gatewayer.Gateway method")
	}

	if len(mmGateway.expectations) > 0 {
		mmGateway.mock.t.Fatalf("Some expectations are already set for the Gatewayer.Gateway method")
	}

	mmGateway.mock.funcGateway = f
	return mmGateway.mock
}

// Gateway implements network.Gatewayer
func (mmGateway *GatewayerMock) Gateway() (g1 mm_network.Gateway) {
	mm_atomic.AddUint64(&mmGateway.beforeGatewayCounter, 1)
	defer mm_atomic.AddUint64(&mmGateway.afterGatewayCounter, 1)

	if mmGateway.inspectFuncGateway != nil {
		mmGateway.inspectFuncGateway()
	}

	if mmGateway.GatewayMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGateway.GatewayMock.defaultExpectation.Counter, 1)

		mm_results := mmGateway.GatewayMock.defaultExpectation.results
		if mm_results == nil {
			mmGateway.t.Fatal("No results are set for the GatewayerMock.Gateway")
		}
		return (*mm_results).g1
	}
	if mmGateway.funcGateway != nil {
		return mmGateway.funcGateway()
	}
	mmGateway.t.Fatalf("Unexpected call to GatewayerMock.Gateway.")
	return
}

// GatewayAfterCounter returns a count of finished GatewayerMock.Gateway invocations
func (mmGateway *GatewayerMock) GatewayAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGateway.afterGatewayCounter)
}

// GatewayBeforeCounter returns a count of GatewayerMock.Gateway invocations
func (mmGateway *GatewayerMock) GatewayBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGateway.beforeGatewayCounter)
}

// MinimockGatewayDone returns true if the count of the Gateway invocations corresponds
// the number of defined expectations
func (m *GatewayerMock) MinimockGatewayDone() bool {
	for _, e := range m.GatewayMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GatewayMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGatewayCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGateway != nil && mm_atomic.LoadUint64(&m.afterGatewayCounter) < 1 {
		return false
	}
	return true
}

// MinimockGatewayInspect logs each unmet expectation
func (m *GatewayerMock) MinimockGatewayInspect() {
	for _, e := range m.GatewayMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Error("Expected call to GatewayerMock.Gateway")
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GatewayMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGatewayCounter) < 1 {
		m.t.Error("Expected call to GatewayerMock.Gateway")
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGateway != nil && mm_atomic.LoadUint64(&m.afterGatewayCounter) < 1 {
		m.t.Error("Expected call to GatewayerMock.Gateway")
	}
}

type mGatewayerMockSwitchState struct {
	mock               *GatewayerMock
	defaultExpectation *GatewayerMockSwitchStateExpectation
	expectations       []*GatewayerMockSwitchStateExpectation

	callArgs []*GatewayerMockSwitchStateParams
	mutex    sync.RWMutex
}

// GatewayerMockSwitchStateExpectation specifies expectation struct of the Gatewayer.SwitchState
type GatewayerMockSwitchStateExpectation struct {
	mock   *GatewayerMock
	params *GatewayerMockSwitchStateParams

	Counter uint64
}

// GatewayerMockSwitchStateParams contains parameters of the Gatewayer.SwitchState
type GatewayerMockSwitchStateParams struct {
	ctx   context.Context
	state insolar.NetworkState
	pulse insolar.Pulse
}

// Expect sets up expected params for Gatewayer.SwitchState
func (mmSwitchState *mGatewayerMockSwitchState) Expect(ctx context.Context, state insolar.NetworkState, pulse insolar.Pulse) *mGatewayerMockSwitchState {
	if mmSwitchState.mock.funcSwitchState != nil {
		mmSwitchState.mock.t.Fatalf("GatewayerMock.SwitchState mock is already set by Set")
	}

	if mmSwitchState.defaultExpectation == nil {
		mmSwitchState.defaultExpectation = &GatewayerMockSwitchStateExpectation{}
	}

	mmSwitchState.defaultExpectation.params = &GatewayerMockSwitchStateParams{ctx, state, pulse}
	for _, e := range mmSwitchState.expectations {
		if minimock.Equal(e.params, mmSwitchState.defaultExpectation.params) {
			mmSwitchState.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmSwitchState.defaultExpectation.params)
		}
	}

	return mmSwitchState
}

// Inspect accepts an inspector function that has same arguments as the Gatewayer.SwitchState
func (mmSwitchState *mGatewayerMockSwitchState) Inspect(f func(ctx context.Context, state insolar.NetworkState, pulse insolar.Pulse)) *mGatewayerMockSwitchState {
	if mmSwitchState.mock.inspectFuncSwitchState != nil {
		mmSwitchState.mock.t.Fatalf("Inspect function is already set for GatewayerMock.SwitchState")
	}

	mmSwitchState.mock.inspectFuncSwitchState = f

	return mmSwitchState
}

// Return sets up results that will be returned by Gatewayer.SwitchState
func (mmSwitchState *mGatewayerMockSwitchState) Return() *GatewayerMock {
	if mmSwitchState.mock.funcSwitchState != nil {
		mmSwitchState.mock.t.Fatalf("GatewayerMock.SwitchState mock is already set by Set")
	}

	if mmSwitchState.defaultExpectation == nil {
		mmSwitchState.defaultExpectation = &GatewayerMockSwitchStateExpectation{mock: mmSwitchState.mock}
	}

	return mmSwitchState.mock
}

//Set uses given function f to mock the Gatewayer.SwitchState method
func (mmSwitchState *mGatewayerMockSwitchState) Set(f func(ctx context.Context, state insolar.NetworkState, pulse insolar.Pulse)) *GatewayerMock {
	if mmSwitchState.defaultExpectation != nil {
		mmSwitchState.mock.t.Fatalf("Default expectation is already set for the Gatewayer.SwitchState method")
	}

	if len(mmSwitchState.expectations) > 0 {
		mmSwitchState.mock.t.Fatalf("Some expectations are already set for the Gatewayer.SwitchState method")
	}

	mmSwitchState.mock.funcSwitchState = f
	return mmSwitchState.mock
}

// SwitchState implements network.Gatewayer
func (mmSwitchState *GatewayerMock) SwitchState(ctx context.Context, state insolar.NetworkState, pulse insolar.Pulse) {
	mm_atomic.AddUint64(&mmSwitchState.beforeSwitchStateCounter, 1)
	defer mm_atomic.AddUint64(&mmSwitchState.afterSwitchStateCounter, 1)

	if mmSwitchState.inspectFuncSwitchState != nil {
		mmSwitchState.inspectFuncSwitchState(ctx, state, pulse)
	}

	mm_params := &GatewayerMockSwitchStateParams{ctx, state, pulse}

	// Record call args
	mmSwitchState.SwitchStateMock.mutex.Lock()
	mmSwitchState.SwitchStateMock.callArgs = append(mmSwitchState.SwitchStateMock.callArgs, mm_params)
	mmSwitchState.SwitchStateMock.mutex.Unlock()

	for _, e := range mmSwitchState.SwitchStateMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return
		}
	}

	if mmSwitchState.SwitchStateMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmSwitchState.SwitchStateMock.defaultExpectation.Counter, 1)
		mm_want := mmSwitchState.SwitchStateMock.defaultExpectation.params
		mm_got := GatewayerMockSwitchStateParams{ctx, state, pulse}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmSwitchState.t.Errorf("GatewayerMock.SwitchState got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		return

	}
	if mmSwitchState.funcSwitchState != nil {
		mmSwitchState.funcSwitchState(ctx, state, pulse)
		return
	}
	mmSwitchState.t.Fatalf("Unexpected call to GatewayerMock.SwitchState. %v %v %v", ctx, state, pulse)

}

// SwitchStateAfterCounter returns a count of finished GatewayerMock.SwitchState invocations
func (mmSwitchState *GatewayerMock) SwitchStateAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmSwitchState.afterSwitchStateCounter)
}

// SwitchStateBeforeCounter returns a count of GatewayerMock.SwitchState invocations
func (mmSwitchState *GatewayerMock) SwitchStateBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmSwitchState.beforeSwitchStateCounter)
}

// Calls returns a list of arguments used in each call to GatewayerMock.SwitchState.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmSwitchState *mGatewayerMockSwitchState) Calls() []*GatewayerMockSwitchStateParams {
	mmSwitchState.mutex.RLock()

	argCopy := make([]*GatewayerMockSwitchStateParams, len(mmSwitchState.callArgs))
	copy(argCopy, mmSwitchState.callArgs)

	mmSwitchState.mutex.RUnlock()

	return argCopy
}

// MinimockSwitchStateDone returns true if the count of the SwitchState invocations corresponds
// the number of defined expectations
func (m *GatewayerMock) MinimockSwitchStateDone() bool {
	for _, e := range m.SwitchStateMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.SwitchStateMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterSwitchStateCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcSwitchState != nil && mm_atomic.LoadUint64(&m.afterSwitchStateCounter) < 1 {
		return false
	}
	return true
}

// MinimockSwitchStateInspect logs each unmet expectation
func (m *GatewayerMock) MinimockSwitchStateInspect() {
	for _, e := range m.SwitchStateMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to GatewayerMock.SwitchState with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.SwitchStateMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterSwitchStateCounter) < 1 {
		if m.SwitchStateMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to GatewayerMock.SwitchState")
		} else {
			m.t.Errorf("Expected call to GatewayerMock.SwitchState with params: %#v", *m.SwitchStateMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcSwitchState != nil && mm_atomic.LoadUint64(&m.afterSwitchStateCounter) < 1 {
		m.t.Error("Expected call to GatewayerMock.SwitchState")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *GatewayerMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockGatewayInspect()

		m.MinimockSwitchStateInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *GatewayerMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *GatewayerMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockGatewayDone() &&
		m.MinimockSwitchStateDone()
}
