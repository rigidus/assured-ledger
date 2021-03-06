package descriptor

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
	"github.com/insolar/assured-ledger/ledger-core/v2/insolar"
)

// CacheMock implements Cache
type CacheMock struct {
	t minimock.Tester

	funcByPrototypeRef          func(ctx context.Context, protoRef insolar.Reference) (p1 PrototypeDescriptor, c2 CodeDescriptor, err error)
	inspectFuncByPrototypeRef   func(ctx context.Context, protoRef insolar.Reference)
	afterByPrototypeRefCounter  uint64
	beforeByPrototypeRefCounter uint64
	ByPrototypeRefMock          mCacheMockByPrototypeRef

	funcRegisterCallback          func(cb CacheCallbackType)
	inspectFuncRegisterCallback   func(cb CacheCallbackType)
	afterRegisterCallbackCounter  uint64
	beforeRegisterCallbackCounter uint64
	RegisterCallbackMock          mCacheMockRegisterCallback
}

// NewCacheMock returns a mock for Cache
func NewCacheMock(t minimock.Tester) *CacheMock {
	m := &CacheMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ByPrototypeRefMock = mCacheMockByPrototypeRef{mock: m}
	m.ByPrototypeRefMock.callArgs = []*CacheMockByPrototypeRefParams{}

	m.RegisterCallbackMock = mCacheMockRegisterCallback{mock: m}
	m.RegisterCallbackMock.callArgs = []*CacheMockRegisterCallbackParams{}

	return m
}

type mCacheMockByPrototypeRef struct {
	mock               *CacheMock
	defaultExpectation *CacheMockByPrototypeRefExpectation
	expectations       []*CacheMockByPrototypeRefExpectation

	callArgs []*CacheMockByPrototypeRefParams
	mutex    sync.RWMutex
}

// CacheMockByPrototypeRefExpectation specifies expectation struct of the Cache.ByPrototypeRef
type CacheMockByPrototypeRefExpectation struct {
	mock    *CacheMock
	params  *CacheMockByPrototypeRefParams
	results *CacheMockByPrototypeRefResults
	Counter uint64
}

// CacheMockByPrototypeRefParams contains parameters of the Cache.ByPrototypeRef
type CacheMockByPrototypeRefParams struct {
	ctx      context.Context
	protoRef insolar.Reference
}

// CacheMockByPrototypeRefResults contains results of the Cache.ByPrototypeRef
type CacheMockByPrototypeRefResults struct {
	p1  PrototypeDescriptor
	c2  CodeDescriptor
	err error
}

// Expect sets up expected params for Cache.ByPrototypeRef
func (mmByPrototypeRef *mCacheMockByPrototypeRef) Expect(ctx context.Context, protoRef insolar.Reference) *mCacheMockByPrototypeRef {
	if mmByPrototypeRef.mock.funcByPrototypeRef != nil {
		mmByPrototypeRef.mock.t.Fatalf("CacheMock.ByPrototypeRef mock is already set by Set")
	}

	if mmByPrototypeRef.defaultExpectation == nil {
		mmByPrototypeRef.defaultExpectation = &CacheMockByPrototypeRefExpectation{}
	}

	mmByPrototypeRef.defaultExpectation.params = &CacheMockByPrototypeRefParams{ctx, protoRef}
	for _, e := range mmByPrototypeRef.expectations {
		if minimock.Equal(e.params, mmByPrototypeRef.defaultExpectation.params) {
			mmByPrototypeRef.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmByPrototypeRef.defaultExpectation.params)
		}
	}

	return mmByPrototypeRef
}

// Inspect accepts an inspector function that has same arguments as the Cache.ByPrototypeRef
func (mmByPrototypeRef *mCacheMockByPrototypeRef) Inspect(f func(ctx context.Context, protoRef insolar.Reference)) *mCacheMockByPrototypeRef {
	if mmByPrototypeRef.mock.inspectFuncByPrototypeRef != nil {
		mmByPrototypeRef.mock.t.Fatalf("Inspect function is already set for CacheMock.ByPrototypeRef")
	}

	mmByPrototypeRef.mock.inspectFuncByPrototypeRef = f

	return mmByPrototypeRef
}

// Return sets up results that will be returned by Cache.ByPrototypeRef
func (mmByPrototypeRef *mCacheMockByPrototypeRef) Return(p1 PrototypeDescriptor, c2 CodeDescriptor, err error) *CacheMock {
	if mmByPrototypeRef.mock.funcByPrototypeRef != nil {
		mmByPrototypeRef.mock.t.Fatalf("CacheMock.ByPrototypeRef mock is already set by Set")
	}

	if mmByPrototypeRef.defaultExpectation == nil {
		mmByPrototypeRef.defaultExpectation = &CacheMockByPrototypeRefExpectation{mock: mmByPrototypeRef.mock}
	}
	mmByPrototypeRef.defaultExpectation.results = &CacheMockByPrototypeRefResults{p1, c2, err}
	return mmByPrototypeRef.mock
}

//Set uses given function f to mock the Cache.ByPrototypeRef method
func (mmByPrototypeRef *mCacheMockByPrototypeRef) Set(f func(ctx context.Context, protoRef insolar.Reference) (p1 PrototypeDescriptor, c2 CodeDescriptor, err error)) *CacheMock {
	if mmByPrototypeRef.defaultExpectation != nil {
		mmByPrototypeRef.mock.t.Fatalf("Default expectation is already set for the Cache.ByPrototypeRef method")
	}

	if len(mmByPrototypeRef.expectations) > 0 {
		mmByPrototypeRef.mock.t.Fatalf("Some expectations are already set for the Cache.ByPrototypeRef method")
	}

	mmByPrototypeRef.mock.funcByPrototypeRef = f
	return mmByPrototypeRef.mock
}

// When sets expectation for the Cache.ByPrototypeRef which will trigger the result defined by the following
// Then helper
func (mmByPrototypeRef *mCacheMockByPrototypeRef) When(ctx context.Context, protoRef insolar.Reference) *CacheMockByPrototypeRefExpectation {
	if mmByPrototypeRef.mock.funcByPrototypeRef != nil {
		mmByPrototypeRef.mock.t.Fatalf("CacheMock.ByPrototypeRef mock is already set by Set")
	}

	expectation := &CacheMockByPrototypeRefExpectation{
		mock:   mmByPrototypeRef.mock,
		params: &CacheMockByPrototypeRefParams{ctx, protoRef},
	}
	mmByPrototypeRef.expectations = append(mmByPrototypeRef.expectations, expectation)
	return expectation
}

// Then sets up Cache.ByPrototypeRef return parameters for the expectation previously defined by the When method
func (e *CacheMockByPrototypeRefExpectation) Then(p1 PrototypeDescriptor, c2 CodeDescriptor, err error) *CacheMock {
	e.results = &CacheMockByPrototypeRefResults{p1, c2, err}
	return e.mock
}

// ByPrototypeRef implements Cache
func (mmByPrototypeRef *CacheMock) ByPrototypeRef(ctx context.Context, protoRef insolar.Reference) (p1 PrototypeDescriptor, c2 CodeDescriptor, err error) {
	mm_atomic.AddUint64(&mmByPrototypeRef.beforeByPrototypeRefCounter, 1)
	defer mm_atomic.AddUint64(&mmByPrototypeRef.afterByPrototypeRefCounter, 1)

	if mmByPrototypeRef.inspectFuncByPrototypeRef != nil {
		mmByPrototypeRef.inspectFuncByPrototypeRef(ctx, protoRef)
	}

	mm_params := &CacheMockByPrototypeRefParams{ctx, protoRef}

	// Record call args
	mmByPrototypeRef.ByPrototypeRefMock.mutex.Lock()
	mmByPrototypeRef.ByPrototypeRefMock.callArgs = append(mmByPrototypeRef.ByPrototypeRefMock.callArgs, mm_params)
	mmByPrototypeRef.ByPrototypeRefMock.mutex.Unlock()

	for _, e := range mmByPrototypeRef.ByPrototypeRefMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.p1, e.results.c2, e.results.err
		}
	}

	if mmByPrototypeRef.ByPrototypeRefMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmByPrototypeRef.ByPrototypeRefMock.defaultExpectation.Counter, 1)
		mm_want := mmByPrototypeRef.ByPrototypeRefMock.defaultExpectation.params
		mm_got := CacheMockByPrototypeRefParams{ctx, protoRef}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmByPrototypeRef.t.Errorf("CacheMock.ByPrototypeRef got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmByPrototypeRef.ByPrototypeRefMock.defaultExpectation.results
		if mm_results == nil {
			mmByPrototypeRef.t.Fatal("No results are set for the CacheMock.ByPrototypeRef")
		}
		return (*mm_results).p1, (*mm_results).c2, (*mm_results).err
	}
	if mmByPrototypeRef.funcByPrototypeRef != nil {
		return mmByPrototypeRef.funcByPrototypeRef(ctx, protoRef)
	}
	mmByPrototypeRef.t.Fatalf("Unexpected call to CacheMock.ByPrototypeRef. %v %v", ctx, protoRef)
	return
}

// ByPrototypeRefAfterCounter returns a count of finished CacheMock.ByPrototypeRef invocations
func (mmByPrototypeRef *CacheMock) ByPrototypeRefAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmByPrototypeRef.afterByPrototypeRefCounter)
}

// ByPrototypeRefBeforeCounter returns a count of CacheMock.ByPrototypeRef invocations
func (mmByPrototypeRef *CacheMock) ByPrototypeRefBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmByPrototypeRef.beforeByPrototypeRefCounter)
}

// Calls returns a list of arguments used in each call to CacheMock.ByPrototypeRef.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmByPrototypeRef *mCacheMockByPrototypeRef) Calls() []*CacheMockByPrototypeRefParams {
	mmByPrototypeRef.mutex.RLock()

	argCopy := make([]*CacheMockByPrototypeRefParams, len(mmByPrototypeRef.callArgs))
	copy(argCopy, mmByPrototypeRef.callArgs)

	mmByPrototypeRef.mutex.RUnlock()

	return argCopy
}

// MinimockByPrototypeRefDone returns true if the count of the ByPrototypeRef invocations corresponds
// the number of defined expectations
func (m *CacheMock) MinimockByPrototypeRefDone() bool {
	for _, e := range m.ByPrototypeRefMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ByPrototypeRefMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterByPrototypeRefCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcByPrototypeRef != nil && mm_atomic.LoadUint64(&m.afterByPrototypeRefCounter) < 1 {
		return false
	}
	return true
}

// MinimockByPrototypeRefInspect logs each unmet expectation
func (m *CacheMock) MinimockByPrototypeRefInspect() {
	for _, e := range m.ByPrototypeRefMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to CacheMock.ByPrototypeRef with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.ByPrototypeRefMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterByPrototypeRefCounter) < 1 {
		if m.ByPrototypeRefMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to CacheMock.ByPrototypeRef")
		} else {
			m.t.Errorf("Expected call to CacheMock.ByPrototypeRef with params: %#v", *m.ByPrototypeRefMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcByPrototypeRef != nil && mm_atomic.LoadUint64(&m.afterByPrototypeRefCounter) < 1 {
		m.t.Error("Expected call to CacheMock.ByPrototypeRef")
	}
}

type mCacheMockRegisterCallback struct {
	mock               *CacheMock
	defaultExpectation *CacheMockRegisterCallbackExpectation
	expectations       []*CacheMockRegisterCallbackExpectation

	callArgs []*CacheMockRegisterCallbackParams
	mutex    sync.RWMutex
}

// CacheMockRegisterCallbackExpectation specifies expectation struct of the Cache.RegisterCallback
type CacheMockRegisterCallbackExpectation struct {
	mock   *CacheMock
	params *CacheMockRegisterCallbackParams

	Counter uint64
}

// CacheMockRegisterCallbackParams contains parameters of the Cache.RegisterCallback
type CacheMockRegisterCallbackParams struct {
	cb CacheCallbackType
}

// Expect sets up expected params for Cache.RegisterCallback
func (mmRegisterCallback *mCacheMockRegisterCallback) Expect(cb CacheCallbackType) *mCacheMockRegisterCallback {
	if mmRegisterCallback.mock.funcRegisterCallback != nil {
		mmRegisterCallback.mock.t.Fatalf("CacheMock.RegisterCallback mock is already set by Set")
	}

	if mmRegisterCallback.defaultExpectation == nil {
		mmRegisterCallback.defaultExpectation = &CacheMockRegisterCallbackExpectation{}
	}

	mmRegisterCallback.defaultExpectation.params = &CacheMockRegisterCallbackParams{cb}
	for _, e := range mmRegisterCallback.expectations {
		if minimock.Equal(e.params, mmRegisterCallback.defaultExpectation.params) {
			mmRegisterCallback.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmRegisterCallback.defaultExpectation.params)
		}
	}

	return mmRegisterCallback
}

// Inspect accepts an inspector function that has same arguments as the Cache.RegisterCallback
func (mmRegisterCallback *mCacheMockRegisterCallback) Inspect(f func(cb CacheCallbackType)) *mCacheMockRegisterCallback {
	if mmRegisterCallback.mock.inspectFuncRegisterCallback != nil {
		mmRegisterCallback.mock.t.Fatalf("Inspect function is already set for CacheMock.RegisterCallback")
	}

	mmRegisterCallback.mock.inspectFuncRegisterCallback = f

	return mmRegisterCallback
}

// Return sets up results that will be returned by Cache.RegisterCallback
func (mmRegisterCallback *mCacheMockRegisterCallback) Return() *CacheMock {
	if mmRegisterCallback.mock.funcRegisterCallback != nil {
		mmRegisterCallback.mock.t.Fatalf("CacheMock.RegisterCallback mock is already set by Set")
	}

	if mmRegisterCallback.defaultExpectation == nil {
		mmRegisterCallback.defaultExpectation = &CacheMockRegisterCallbackExpectation{mock: mmRegisterCallback.mock}
	}

	return mmRegisterCallback.mock
}

//Set uses given function f to mock the Cache.RegisterCallback method
func (mmRegisterCallback *mCacheMockRegisterCallback) Set(f func(cb CacheCallbackType)) *CacheMock {
	if mmRegisterCallback.defaultExpectation != nil {
		mmRegisterCallback.mock.t.Fatalf("Default expectation is already set for the Cache.RegisterCallback method")
	}

	if len(mmRegisterCallback.expectations) > 0 {
		mmRegisterCallback.mock.t.Fatalf("Some expectations are already set for the Cache.RegisterCallback method")
	}

	mmRegisterCallback.mock.funcRegisterCallback = f
	return mmRegisterCallback.mock
}

// RegisterCallback implements Cache
func (mmRegisterCallback *CacheMock) RegisterCallback(cb CacheCallbackType) {
	mm_atomic.AddUint64(&mmRegisterCallback.beforeRegisterCallbackCounter, 1)
	defer mm_atomic.AddUint64(&mmRegisterCallback.afterRegisterCallbackCounter, 1)

	if mmRegisterCallback.inspectFuncRegisterCallback != nil {
		mmRegisterCallback.inspectFuncRegisterCallback(cb)
	}

	mm_params := &CacheMockRegisterCallbackParams{cb}

	// Record call args
	mmRegisterCallback.RegisterCallbackMock.mutex.Lock()
	mmRegisterCallback.RegisterCallbackMock.callArgs = append(mmRegisterCallback.RegisterCallbackMock.callArgs, mm_params)
	mmRegisterCallback.RegisterCallbackMock.mutex.Unlock()

	for _, e := range mmRegisterCallback.RegisterCallbackMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return
		}
	}

	if mmRegisterCallback.RegisterCallbackMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmRegisterCallback.RegisterCallbackMock.defaultExpectation.Counter, 1)
		mm_want := mmRegisterCallback.RegisterCallbackMock.defaultExpectation.params
		mm_got := CacheMockRegisterCallbackParams{cb}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmRegisterCallback.t.Errorf("CacheMock.RegisterCallback got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		return

	}
	if mmRegisterCallback.funcRegisterCallback != nil {
		mmRegisterCallback.funcRegisterCallback(cb)
		return
	}
	mmRegisterCallback.t.Fatalf("Unexpected call to CacheMock.RegisterCallback. %v", cb)

}

// RegisterCallbackAfterCounter returns a count of finished CacheMock.RegisterCallback invocations
func (mmRegisterCallback *CacheMock) RegisterCallbackAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmRegisterCallback.afterRegisterCallbackCounter)
}

// RegisterCallbackBeforeCounter returns a count of CacheMock.RegisterCallback invocations
func (mmRegisterCallback *CacheMock) RegisterCallbackBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmRegisterCallback.beforeRegisterCallbackCounter)
}

// Calls returns a list of arguments used in each call to CacheMock.RegisterCallback.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmRegisterCallback *mCacheMockRegisterCallback) Calls() []*CacheMockRegisterCallbackParams {
	mmRegisterCallback.mutex.RLock()

	argCopy := make([]*CacheMockRegisterCallbackParams, len(mmRegisterCallback.callArgs))
	copy(argCopy, mmRegisterCallback.callArgs)

	mmRegisterCallback.mutex.RUnlock()

	return argCopy
}

// MinimockRegisterCallbackDone returns true if the count of the RegisterCallback invocations corresponds
// the number of defined expectations
func (m *CacheMock) MinimockRegisterCallbackDone() bool {
	for _, e := range m.RegisterCallbackMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.RegisterCallbackMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterRegisterCallbackCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRegisterCallback != nil && mm_atomic.LoadUint64(&m.afterRegisterCallbackCounter) < 1 {
		return false
	}
	return true
}

// MinimockRegisterCallbackInspect logs each unmet expectation
func (m *CacheMock) MinimockRegisterCallbackInspect() {
	for _, e := range m.RegisterCallbackMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to CacheMock.RegisterCallback with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.RegisterCallbackMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterRegisterCallbackCounter) < 1 {
		if m.RegisterCallbackMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to CacheMock.RegisterCallback")
		} else {
			m.t.Errorf("Expected call to CacheMock.RegisterCallback with params: %#v", *m.RegisterCallbackMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRegisterCallback != nil && mm_atomic.LoadUint64(&m.afterRegisterCallbackCounter) < 1 {
		m.t.Error("Expected call to CacheMock.RegisterCallback")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *CacheMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockByPrototypeRefInspect()

		m.MinimockRegisterCallbackInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *CacheMock) MinimockWait(timeout mm_time.Duration) {
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

func (m *CacheMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockByPrototypeRefDone() &&
		m.MinimockRegisterCallbackDone()
}
