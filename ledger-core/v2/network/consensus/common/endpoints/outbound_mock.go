package endpoints

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

import (
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"

	"github.com/insolar/assured-ledger/ledger-core/v2/insolar"
	"github.com/insolar/assured-ledger/ledger-core/v2/vanilla/longbits"
)

// OutboundMock implements Outbound
type OutboundMock struct {
	t minimock.Tester

	funcAsByteString          func() (b1 longbits.ByteString)
	inspectFuncAsByteString   func()
	afterAsByteStringCounter  uint64
	beforeAsByteStringCounter uint64
	AsByteStringMock          mOutboundMockAsByteString

	funcCanAccept          func(connection Inbound) (b1 bool)
	inspectFuncCanAccept   func(connection Inbound)
	afterCanAcceptCounter  uint64
	beforeCanAcceptCounter uint64
	CanAcceptMock          mOutboundMockCanAccept

	funcGetEndpointType          func() (n1 NodeEndpointType)
	inspectFuncGetEndpointType   func()
	afterGetEndpointTypeCounter  uint64
	beforeGetEndpointTypeCounter uint64
	GetEndpointTypeMock          mOutboundMockGetEndpointType

	funcGetIPAddress          func() (i1 IPAddress)
	inspectFuncGetIPAddress   func()
	afterGetIPAddressCounter  uint64
	beforeGetIPAddressCounter uint64
	GetIPAddressMock          mOutboundMockGetIPAddress

	funcGetNameAddress          func() (n1 Name)
	inspectFuncGetNameAddress   func()
	afterGetNameAddressCounter  uint64
	beforeGetNameAddressCounter uint64
	GetNameAddressMock          mOutboundMockGetNameAddress

	funcGetRelayID          func() (s1 insolar.ShortNodeID)
	inspectFuncGetRelayID   func()
	afterGetRelayIDCounter  uint64
	beforeGetRelayIDCounter uint64
	GetRelayIDMock          mOutboundMockGetRelayID
}

// NewOutboundMock returns a mock for Outbound
func NewOutboundMock(t minimock.Tester) *OutboundMock {
	m := &OutboundMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.AsByteStringMock = mOutboundMockAsByteString{mock: m}

	m.CanAcceptMock = mOutboundMockCanAccept{mock: m}
	m.CanAcceptMock.callArgs = []*OutboundMockCanAcceptParams{}

	m.GetEndpointTypeMock = mOutboundMockGetEndpointType{mock: m}

	m.GetIPAddressMock = mOutboundMockGetIPAddress{mock: m}

	m.GetNameAddressMock = mOutboundMockGetNameAddress{mock: m}

	m.GetRelayIDMock = mOutboundMockGetRelayID{mock: m}

	return m
}

type mOutboundMockAsByteString struct {
	mock               *OutboundMock
	defaultExpectation *OutboundMockAsByteStringExpectation
	expectations       []*OutboundMockAsByteStringExpectation
}

// OutboundMockAsByteStringExpectation specifies expectation struct of the Outbound.AsByteString
type OutboundMockAsByteStringExpectation struct {
	mock *OutboundMock

	results *OutboundMockAsByteStringResults
	Counter uint64
}

// OutboundMockAsByteStringResults contains results of the Outbound.AsByteString
type OutboundMockAsByteStringResults struct {
	b1 longbits.ByteString
}

// Expect sets up expected params for Outbound.AsByteString
func (mmAsByteString *mOutboundMockAsByteString) Expect() *mOutboundMockAsByteString {
	if mmAsByteString.mock.funcAsByteString != nil {
		mmAsByteString.mock.t.Fatalf("OutboundMock.AsByteString mock is already set by Set")
	}

	if mmAsByteString.defaultExpectation == nil {
		mmAsByteString.defaultExpectation = &OutboundMockAsByteStringExpectation{}
	}

	return mmAsByteString
}

// Inspect accepts an inspector function that has same arguments as the Outbound.AsByteString
func (mmAsByteString *mOutboundMockAsByteString) Inspect(f func()) *mOutboundMockAsByteString {
	if mmAsByteString.mock.inspectFuncAsByteString != nil {
		mmAsByteString.mock.t.Fatalf("Inspect function is already set for OutboundMock.AsByteString")
	}

	mmAsByteString.mock.inspectFuncAsByteString = f

	return mmAsByteString
}

// Return sets up results that will be returned by Outbound.AsByteString
func (mmAsByteString *mOutboundMockAsByteString) Return(b1 longbits.ByteString) *OutboundMock {
	if mmAsByteString.mock.funcAsByteString != nil {
		mmAsByteString.mock.t.Fatalf("OutboundMock.AsByteString mock is already set by Set")
	}

	if mmAsByteString.defaultExpectation == nil {
		mmAsByteString.defaultExpectation = &OutboundMockAsByteStringExpectation{mock: mmAsByteString.mock}
	}
	mmAsByteString.defaultExpectation.results = &OutboundMockAsByteStringResults{b1}
	return mmAsByteString.mock
}

//Set uses given function f to mock the Outbound.AsByteString method
func (mmAsByteString *mOutboundMockAsByteString) Set(f func() (b1 longbits.ByteString)) *OutboundMock {
	if mmAsByteString.defaultExpectation != nil {
		mmAsByteString.mock.t.Fatalf("Default expectation is already set for the Outbound.AsByteString method")
	}

	if len(mmAsByteString.expectations) > 0 {
		mmAsByteString.mock.t.Fatalf("Some expectations are already set for the Outbound.AsByteString method")
	}

	mmAsByteString.mock.funcAsByteString = f
	return mmAsByteString.mock
}

// AsByteString implements Outbound
func (mmAsByteString *OutboundMock) AsByteString() (b1 longbits.ByteString) {
	mm_atomic.AddUint64(&mmAsByteString.beforeAsByteStringCounter, 1)
	defer mm_atomic.AddUint64(&mmAsByteString.afterAsByteStringCounter, 1)

	if mmAsByteString.inspectFuncAsByteString != nil {
		mmAsByteString.inspectFuncAsByteString()
	}

	if mmAsByteString.AsByteStringMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmAsByteString.AsByteStringMock.defaultExpectation.Counter, 1)

		mm_results := mmAsByteString.AsByteStringMock.defaultExpectation.results
		if mm_results == nil {
			mmAsByteString.t.Fatal("No results are set for the OutboundMock.AsByteString")
		}
		return (*mm_results).b1
	}
	if mmAsByteString.funcAsByteString != nil {
		return mmAsByteString.funcAsByteString()
	}
	mmAsByteString.t.Fatalf("Unexpected call to OutboundMock.AsByteString.")
	return
}

// AsByteStringAfterCounter returns a count of finished OutboundMock.AsByteString invocations
func (mmAsByteString *OutboundMock) AsByteStringAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmAsByteString.afterAsByteStringCounter)
}

// AsByteStringBeforeCounter returns a count of OutboundMock.AsByteString invocations
func (mmAsByteString *OutboundMock) AsByteStringBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmAsByteString.beforeAsByteStringCounter)
}

// MinimockAsByteStringDone returns true if the count of the AsByteString invocations corresponds
// the number of defined expectations
func (m *OutboundMock) MinimockAsByteStringDone() bool {
	for _, e := range m.AsByteStringMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.AsByteStringMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterAsByteStringCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcAsByteString != nil && mm_atomic.LoadUint64(&m.afterAsByteStringCounter) < 1 {
		return false
	}
	return true
}

// MinimockAsByteStringInspect logs each unmet expectation
func (m *OutboundMock) MinimockAsByteStringInspect() {
	for _, e := range m.AsByteStringMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Error("Expected call to OutboundMock.AsByteString")
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.AsByteStringMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterAsByteStringCounter) < 1 {
		m.t.Error("Expected call to OutboundMock.AsByteString")
	}
	// if func was set then invocations count should be greater than zero
	if m.funcAsByteString != nil && mm_atomic.LoadUint64(&m.afterAsByteStringCounter) < 1 {
		m.t.Error("Expected call to OutboundMock.AsByteString")
	}
}

type mOutboundMockCanAccept struct {
	mock               *OutboundMock
	defaultExpectation *OutboundMockCanAcceptExpectation
	expectations       []*OutboundMockCanAcceptExpectation

	callArgs []*OutboundMockCanAcceptParams
	mutex    sync.RWMutex
}

// OutboundMockCanAcceptExpectation specifies expectation struct of the Outbound.CanAccept
type OutboundMockCanAcceptExpectation struct {
	mock    *OutboundMock
	params  *OutboundMockCanAcceptParams
	results *OutboundMockCanAcceptResults
	Counter uint64
}

// OutboundMockCanAcceptParams contains parameters of the Outbound.CanAccept
type OutboundMockCanAcceptParams struct {
	connection Inbound
}

// OutboundMockCanAcceptResults contains results of the Outbound.CanAccept
type OutboundMockCanAcceptResults struct {
	b1 bool
}

// Expect sets up expected params for Outbound.CanAccept
func (mmCanAccept *mOutboundMockCanAccept) Expect(connection Inbound) *mOutboundMockCanAccept {
	if mmCanAccept.mock.funcCanAccept != nil {
		mmCanAccept.mock.t.Fatalf("OutboundMock.CanAccept mock is already set by Set")
	}

	if mmCanAccept.defaultExpectation == nil {
		mmCanAccept.defaultExpectation = &OutboundMockCanAcceptExpectation{}
	}

	mmCanAccept.defaultExpectation.params = &OutboundMockCanAcceptParams{connection}
	for _, e := range mmCanAccept.expectations {
		if minimock.Equal(e.params, mmCanAccept.defaultExpectation.params) {
			mmCanAccept.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmCanAccept.defaultExpectation.params)
		}
	}

	return mmCanAccept
}

// Inspect accepts an inspector function that has same arguments as the Outbound.CanAccept
func (mmCanAccept *mOutboundMockCanAccept) Inspect(f func(connection Inbound)) *mOutboundMockCanAccept {
	if mmCanAccept.mock.inspectFuncCanAccept != nil {
		mmCanAccept.mock.t.Fatalf("Inspect function is already set for OutboundMock.CanAccept")
	}

	mmCanAccept.mock.inspectFuncCanAccept = f

	return mmCanAccept
}

// Return sets up results that will be returned by Outbound.CanAccept
func (mmCanAccept *mOutboundMockCanAccept) Return(b1 bool) *OutboundMock {
	if mmCanAccept.mock.funcCanAccept != nil {
		mmCanAccept.mock.t.Fatalf("OutboundMock.CanAccept mock is already set by Set")
	}

	if mmCanAccept.defaultExpectation == nil {
		mmCanAccept.defaultExpectation = &OutboundMockCanAcceptExpectation{mock: mmCanAccept.mock}
	}
	mmCanAccept.defaultExpectation.results = &OutboundMockCanAcceptResults{b1}
	return mmCanAccept.mock
}

//Set uses given function f to mock the Outbound.CanAccept method
func (mmCanAccept *mOutboundMockCanAccept) Set(f func(connection Inbound) (b1 bool)) *OutboundMock {
	if mmCanAccept.defaultExpectation != nil {
		mmCanAccept.mock.t.Fatalf("Default expectation is already set for the Outbound.CanAccept method")
	}

	if len(mmCanAccept.expectations) > 0 {
		mmCanAccept.mock.t.Fatalf("Some expectations are already set for the Outbound.CanAccept method")
	}

	mmCanAccept.mock.funcCanAccept = f
	return mmCanAccept.mock
}

// When sets expectation for the Outbound.CanAccept which will trigger the result defined by the following
// Then helper
func (mmCanAccept *mOutboundMockCanAccept) When(connection Inbound) *OutboundMockCanAcceptExpectation {
	if mmCanAccept.mock.funcCanAccept != nil {
		mmCanAccept.mock.t.Fatalf("OutboundMock.CanAccept mock is already set by Set")
	}

	expectation := &OutboundMockCanAcceptExpectation{
		mock:   mmCanAccept.mock,
		params: &OutboundMockCanAcceptParams{connection},
	}
	mmCanAccept.expectations = append(mmCanAccept.expectations, expectation)
	return expectation
}

// Then sets up Outbound.CanAccept return parameters for the expectation previously defined by the When method
func (e *OutboundMockCanAcceptExpectation) Then(b1 bool) *OutboundMock {
	e.results = &OutboundMockCanAcceptResults{b1}
	return e.mock
}

// CanAccept implements Outbound
func (mmCanAccept *OutboundMock) CanAccept(connection Inbound) (b1 bool) {
	mm_atomic.AddUint64(&mmCanAccept.beforeCanAcceptCounter, 1)
	defer mm_atomic.AddUint64(&mmCanAccept.afterCanAcceptCounter, 1)

	if mmCanAccept.inspectFuncCanAccept != nil {
		mmCanAccept.inspectFuncCanAccept(connection)
	}

	mm_params := &OutboundMockCanAcceptParams{connection}

	// Record call args
	mmCanAccept.CanAcceptMock.mutex.Lock()
	mmCanAccept.CanAcceptMock.callArgs = append(mmCanAccept.CanAcceptMock.callArgs, mm_params)
	mmCanAccept.CanAcceptMock.mutex.Unlock()

	for _, e := range mmCanAccept.CanAcceptMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.b1
		}
	}

	if mmCanAccept.CanAcceptMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmCanAccept.CanAcceptMock.defaultExpectation.Counter, 1)
		mm_want := mmCanAccept.CanAcceptMock.defaultExpectation.params
		mm_got := OutboundMockCanAcceptParams{connection}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmCanAccept.t.Errorf("OutboundMock.CanAccept got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmCanAccept.CanAcceptMock.defaultExpectation.results
		if mm_results == nil {
			mmCanAccept.t.Fatal("No results are set for the OutboundMock.CanAccept")
		}
		return (*mm_results).b1
	}
	if mmCanAccept.funcCanAccept != nil {
		return mmCanAccept.funcCanAccept(connection)
	}
	mmCanAccept.t.Fatalf("Unexpected call to OutboundMock.CanAccept. %v", connection)
	return
}

// CanAcceptAfterCounter returns a count of finished OutboundMock.CanAccept invocations
func (mmCanAccept *OutboundMock) CanAcceptAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmCanAccept.afterCanAcceptCounter)
}

// CanAcceptBeforeCounter returns a count of OutboundMock.CanAccept invocations
func (mmCanAccept *OutboundMock) CanAcceptBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmCanAccept.beforeCanAcceptCounter)
}

// Calls returns a list of arguments used in each call to OutboundMock.CanAccept.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmCanAccept *mOutboundMockCanAccept) Calls() []*OutboundMockCanAcceptParams {
	mmCanAccept.mutex.RLock()

	argCopy := make([]*OutboundMockCanAcceptParams, len(mmCanAccept.callArgs))
	copy(argCopy, mmCanAccept.callArgs)

	mmCanAccept.mutex.RUnlock()

	return argCopy
}

// MinimockCanAcceptDone returns true if the count of the CanAccept invocations corresponds
// the number of defined expectations
func (m *OutboundMock) MinimockCanAcceptDone() bool {
	for _, e := range m.CanAcceptMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.CanAcceptMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterCanAcceptCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcCanAccept != nil && mm_atomic.LoadUint64(&m.afterCanAcceptCounter) < 1 {
		return false
	}
	return true
}

// MinimockCanAcceptInspect logs each unmet expectation
func (m *OutboundMock) MinimockCanAcceptInspect() {
	for _, e := range m.CanAcceptMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to OutboundMock.CanAccept with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.CanAcceptMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterCanAcceptCounter) < 1 {
		if m.CanAcceptMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to OutboundMock.CanAccept")
		} else {
			m.t.Errorf("Expected call to OutboundMock.CanAccept with params: %#v", *m.CanAcceptMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcCanAccept != nil && mm_atomic.LoadUint64(&m.afterCanAcceptCounter) < 1 {
		m.t.Error("Expected call to OutboundMock.CanAccept")
	}
}

type mOutboundMockGetEndpointType struct {
	mock               *OutboundMock
	defaultExpectation *OutboundMockGetEndpointTypeExpectation
	expectations       []*OutboundMockGetEndpointTypeExpectation
}

// OutboundMockGetEndpointTypeExpectation specifies expectation struct of the Outbound.GetEndpointType
type OutboundMockGetEndpointTypeExpectation struct {
	mock *OutboundMock

	results *OutboundMockGetEndpointTypeResults
	Counter uint64
}

// OutboundMockGetEndpointTypeResults contains results of the Outbound.GetEndpointType
type OutboundMockGetEndpointTypeResults struct {
	n1 NodeEndpointType
}

// Expect sets up expected params for Outbound.GetEndpointType
func (mmGetEndpointType *mOutboundMockGetEndpointType) Expect() *mOutboundMockGetEndpointType {
	if mmGetEndpointType.mock.funcGetEndpointType != nil {
		mmGetEndpointType.mock.t.Fatalf("OutboundMock.GetEndpointType mock is already set by Set")
	}

	if mmGetEndpointType.defaultExpectation == nil {
		mmGetEndpointType.defaultExpectation = &OutboundMockGetEndpointTypeExpectation{}
	}

	return mmGetEndpointType
}

// Inspect accepts an inspector function that has same arguments as the Outbound.GetEndpointType
func (mmGetEndpointType *mOutboundMockGetEndpointType) Inspect(f func()) *mOutboundMockGetEndpointType {
	if mmGetEndpointType.mock.inspectFuncGetEndpointType != nil {
		mmGetEndpointType.mock.t.Fatalf("Inspect function is already set for OutboundMock.GetEndpointType")
	}

	mmGetEndpointType.mock.inspectFuncGetEndpointType = f

	return mmGetEndpointType
}

// Return sets up results that will be returned by Outbound.GetEndpointType
func (mmGetEndpointType *mOutboundMockGetEndpointType) Return(n1 NodeEndpointType) *OutboundMock {
	if mmGetEndpointType.mock.funcGetEndpointType != nil {
		mmGetEndpointType.mock.t.Fatalf("OutboundMock.GetEndpointType mock is already set by Set")
	}

	if mmGetEndpointType.defaultExpectation == nil {
		mmGetEndpointType.defaultExpectation = &OutboundMockGetEndpointTypeExpectation{mock: mmGetEndpointType.mock}
	}
	mmGetEndpointType.defaultExpectation.results = &OutboundMockGetEndpointTypeResults{n1}
	return mmGetEndpointType.mock
}

//Set uses given function f to mock the Outbound.GetEndpointType method
func (mmGetEndpointType *mOutboundMockGetEndpointType) Set(f func() (n1 NodeEndpointType)) *OutboundMock {
	if mmGetEndpointType.defaultExpectation != nil {
		mmGetEndpointType.mock.t.Fatalf("Default expectation is already set for the Outbound.GetEndpointType method")
	}

	if len(mmGetEndpointType.expectations) > 0 {
		mmGetEndpointType.mock.t.Fatalf("Some expectations are already set for the Outbound.GetEndpointType method")
	}

	mmGetEndpointType.mock.funcGetEndpointType = f
	return mmGetEndpointType.mock
}

// GetEndpointType implements Outbound
func (mmGetEndpointType *OutboundMock) GetEndpointType() (n1 NodeEndpointType) {
	mm_atomic.AddUint64(&mmGetEndpointType.beforeGetEndpointTypeCounter, 1)
	defer mm_atomic.AddUint64(&mmGetEndpointType.afterGetEndpointTypeCounter, 1)

	if mmGetEndpointType.inspectFuncGetEndpointType != nil {
		mmGetEndpointType.inspectFuncGetEndpointType()
	}

	if mmGetEndpointType.GetEndpointTypeMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGetEndpointType.GetEndpointTypeMock.defaultExpectation.Counter, 1)

		mm_results := mmGetEndpointType.GetEndpointTypeMock.defaultExpectation.results
		if mm_results == nil {
			mmGetEndpointType.t.Fatal("No results are set for the OutboundMock.GetEndpointType")
		}
		return (*mm_results).n1
	}
	if mmGetEndpointType.funcGetEndpointType != nil {
		return mmGetEndpointType.funcGetEndpointType()
	}
	mmGetEndpointType.t.Fatalf("Unexpected call to OutboundMock.GetEndpointType.")
	return
}

// GetEndpointTypeAfterCounter returns a count of finished OutboundMock.GetEndpointType invocations
func (mmGetEndpointType *OutboundMock) GetEndpointTypeAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetEndpointType.afterGetEndpointTypeCounter)
}

// GetEndpointTypeBeforeCounter returns a count of OutboundMock.GetEndpointType invocations
func (mmGetEndpointType *OutboundMock) GetEndpointTypeBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetEndpointType.beforeGetEndpointTypeCounter)
}

// MinimockGetEndpointTypeDone returns true if the count of the GetEndpointType invocations corresponds
// the number of defined expectations
func (m *OutboundMock) MinimockGetEndpointTypeDone() bool {
	for _, e := range m.GetEndpointTypeMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetEndpointTypeMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetEndpointTypeCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetEndpointType != nil && mm_atomic.LoadUint64(&m.afterGetEndpointTypeCounter) < 1 {
		return false
	}
	return true
}

// MinimockGetEndpointTypeInspect logs each unmet expectation
func (m *OutboundMock) MinimockGetEndpointTypeInspect() {
	for _, e := range m.GetEndpointTypeMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Error("Expected call to OutboundMock.GetEndpointType")
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetEndpointTypeMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetEndpointTypeCounter) < 1 {
		m.t.Error("Expected call to OutboundMock.GetEndpointType")
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetEndpointType != nil && mm_atomic.LoadUint64(&m.afterGetEndpointTypeCounter) < 1 {
		m.t.Error("Expected call to OutboundMock.GetEndpointType")
	}
}

type mOutboundMockGetIPAddress struct {
	mock               *OutboundMock
	defaultExpectation *OutboundMockGetIPAddressExpectation
	expectations       []*OutboundMockGetIPAddressExpectation
}

// OutboundMockGetIPAddressExpectation specifies expectation struct of the Outbound.GetIPAddress
type OutboundMockGetIPAddressExpectation struct {
	mock *OutboundMock

	results *OutboundMockGetIPAddressResults
	Counter uint64
}

// OutboundMockGetIPAddressResults contains results of the Outbound.GetIPAddress
type OutboundMockGetIPAddressResults struct {
	i1 IPAddress
}

// Expect sets up expected params for Outbound.GetIPAddress
func (mmGetIPAddress *mOutboundMockGetIPAddress) Expect() *mOutboundMockGetIPAddress {
	if mmGetIPAddress.mock.funcGetIPAddress != nil {
		mmGetIPAddress.mock.t.Fatalf("OutboundMock.GetIPAddress mock is already set by Set")
	}

	if mmGetIPAddress.defaultExpectation == nil {
		mmGetIPAddress.defaultExpectation = &OutboundMockGetIPAddressExpectation{}
	}

	return mmGetIPAddress
}

// Inspect accepts an inspector function that has same arguments as the Outbound.GetIPAddress
func (mmGetIPAddress *mOutboundMockGetIPAddress) Inspect(f func()) *mOutboundMockGetIPAddress {
	if mmGetIPAddress.mock.inspectFuncGetIPAddress != nil {
		mmGetIPAddress.mock.t.Fatalf("Inspect function is already set for OutboundMock.GetIPAddress")
	}

	mmGetIPAddress.mock.inspectFuncGetIPAddress = f

	return mmGetIPAddress
}

// Return sets up results that will be returned by Outbound.GetIPAddress
func (mmGetIPAddress *mOutboundMockGetIPAddress) Return(i1 IPAddress) *OutboundMock {
	if mmGetIPAddress.mock.funcGetIPAddress != nil {
		mmGetIPAddress.mock.t.Fatalf("OutboundMock.GetIPAddress mock is already set by Set")
	}

	if mmGetIPAddress.defaultExpectation == nil {
		mmGetIPAddress.defaultExpectation = &OutboundMockGetIPAddressExpectation{mock: mmGetIPAddress.mock}
	}
	mmGetIPAddress.defaultExpectation.results = &OutboundMockGetIPAddressResults{i1}
	return mmGetIPAddress.mock
}

//Set uses given function f to mock the Outbound.GetIPAddress method
func (mmGetIPAddress *mOutboundMockGetIPAddress) Set(f func() (i1 IPAddress)) *OutboundMock {
	if mmGetIPAddress.defaultExpectation != nil {
		mmGetIPAddress.mock.t.Fatalf("Default expectation is already set for the Outbound.GetIPAddress method")
	}

	if len(mmGetIPAddress.expectations) > 0 {
		mmGetIPAddress.mock.t.Fatalf("Some expectations are already set for the Outbound.GetIPAddress method")
	}

	mmGetIPAddress.mock.funcGetIPAddress = f
	return mmGetIPAddress.mock
}

// GetIPAddress implements Outbound
func (mmGetIPAddress *OutboundMock) GetIPAddress() (i1 IPAddress) {
	mm_atomic.AddUint64(&mmGetIPAddress.beforeGetIPAddressCounter, 1)
	defer mm_atomic.AddUint64(&mmGetIPAddress.afterGetIPAddressCounter, 1)

	if mmGetIPAddress.inspectFuncGetIPAddress != nil {
		mmGetIPAddress.inspectFuncGetIPAddress()
	}

	if mmGetIPAddress.GetIPAddressMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGetIPAddress.GetIPAddressMock.defaultExpectation.Counter, 1)

		mm_results := mmGetIPAddress.GetIPAddressMock.defaultExpectation.results
		if mm_results == nil {
			mmGetIPAddress.t.Fatal("No results are set for the OutboundMock.GetIPAddress")
		}
		return (*mm_results).i1
	}
	if mmGetIPAddress.funcGetIPAddress != nil {
		return mmGetIPAddress.funcGetIPAddress()
	}
	mmGetIPAddress.t.Fatalf("Unexpected call to OutboundMock.GetIPAddress.")
	return
}

// GetIPAddressAfterCounter returns a count of finished OutboundMock.GetIPAddress invocations
func (mmGetIPAddress *OutboundMock) GetIPAddressAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetIPAddress.afterGetIPAddressCounter)
}

// GetIPAddressBeforeCounter returns a count of OutboundMock.GetIPAddress invocations
func (mmGetIPAddress *OutboundMock) GetIPAddressBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetIPAddress.beforeGetIPAddressCounter)
}

// MinimockGetIPAddressDone returns true if the count of the GetIPAddress invocations corresponds
// the number of defined expectations
func (m *OutboundMock) MinimockGetIPAddressDone() bool {
	for _, e := range m.GetIPAddressMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetIPAddressMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetIPAddressCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetIPAddress != nil && mm_atomic.LoadUint64(&m.afterGetIPAddressCounter) < 1 {
		return false
	}
	return true
}

// MinimockGetIPAddressInspect logs each unmet expectation
func (m *OutboundMock) MinimockGetIPAddressInspect() {
	for _, e := range m.GetIPAddressMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Error("Expected call to OutboundMock.GetIPAddress")
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetIPAddressMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetIPAddressCounter) < 1 {
		m.t.Error("Expected call to OutboundMock.GetIPAddress")
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetIPAddress != nil && mm_atomic.LoadUint64(&m.afterGetIPAddressCounter) < 1 {
		m.t.Error("Expected call to OutboundMock.GetIPAddress")
	}
}

type mOutboundMockGetNameAddress struct {
	mock               *OutboundMock
	defaultExpectation *OutboundMockGetNameAddressExpectation
	expectations       []*OutboundMockGetNameAddressExpectation
}

// OutboundMockGetNameAddressExpectation specifies expectation struct of the Outbound.GetNameAddress
type OutboundMockGetNameAddressExpectation struct {
	mock *OutboundMock

	results *OutboundMockGetNameAddressResults
	Counter uint64
}

// OutboundMockGetNameAddressResults contains results of the Outbound.GetNameAddress
type OutboundMockGetNameAddressResults struct {
	n1 Name
}

// Expect sets up expected params for Outbound.GetNameAddress
func (mmGetNameAddress *mOutboundMockGetNameAddress) Expect() *mOutboundMockGetNameAddress {
	if mmGetNameAddress.mock.funcGetNameAddress != nil {
		mmGetNameAddress.mock.t.Fatalf("OutboundMock.GetNameAddress mock is already set by Set")
	}

	if mmGetNameAddress.defaultExpectation == nil {
		mmGetNameAddress.defaultExpectation = &OutboundMockGetNameAddressExpectation{}
	}

	return mmGetNameAddress
}

// Inspect accepts an inspector function that has same arguments as the Outbound.GetNameAddress
func (mmGetNameAddress *mOutboundMockGetNameAddress) Inspect(f func()) *mOutboundMockGetNameAddress {
	if mmGetNameAddress.mock.inspectFuncGetNameAddress != nil {
		mmGetNameAddress.mock.t.Fatalf("Inspect function is already set for OutboundMock.GetNameAddress")
	}

	mmGetNameAddress.mock.inspectFuncGetNameAddress = f

	return mmGetNameAddress
}

// Return sets up results that will be returned by Outbound.GetNameAddress
func (mmGetNameAddress *mOutboundMockGetNameAddress) Return(n1 Name) *OutboundMock {
	if mmGetNameAddress.mock.funcGetNameAddress != nil {
		mmGetNameAddress.mock.t.Fatalf("OutboundMock.GetNameAddress mock is already set by Set")
	}

	if mmGetNameAddress.defaultExpectation == nil {
		mmGetNameAddress.defaultExpectation = &OutboundMockGetNameAddressExpectation{mock: mmGetNameAddress.mock}
	}
	mmGetNameAddress.defaultExpectation.results = &OutboundMockGetNameAddressResults{n1}
	return mmGetNameAddress.mock
}

//Set uses given function f to mock the Outbound.GetNameAddress method
func (mmGetNameAddress *mOutboundMockGetNameAddress) Set(f func() (n1 Name)) *OutboundMock {
	if mmGetNameAddress.defaultExpectation != nil {
		mmGetNameAddress.mock.t.Fatalf("Default expectation is already set for the Outbound.GetNameAddress method")
	}

	if len(mmGetNameAddress.expectations) > 0 {
		mmGetNameAddress.mock.t.Fatalf("Some expectations are already set for the Outbound.GetNameAddress method")
	}

	mmGetNameAddress.mock.funcGetNameAddress = f
	return mmGetNameAddress.mock
}

// GetNameAddress implements Outbound
func (mmGetNameAddress *OutboundMock) GetNameAddress() (n1 Name) {
	mm_atomic.AddUint64(&mmGetNameAddress.beforeGetNameAddressCounter, 1)
	defer mm_atomic.AddUint64(&mmGetNameAddress.afterGetNameAddressCounter, 1)

	if mmGetNameAddress.inspectFuncGetNameAddress != nil {
		mmGetNameAddress.inspectFuncGetNameAddress()
	}

	if mmGetNameAddress.GetNameAddressMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGetNameAddress.GetNameAddressMock.defaultExpectation.Counter, 1)

		mm_results := mmGetNameAddress.GetNameAddressMock.defaultExpectation.results
		if mm_results == nil {
			mmGetNameAddress.t.Fatal("No results are set for the OutboundMock.GetNameAddress")
		}
		return (*mm_results).n1
	}
	if mmGetNameAddress.funcGetNameAddress != nil {
		return mmGetNameAddress.funcGetNameAddress()
	}
	mmGetNameAddress.t.Fatalf("Unexpected call to OutboundMock.GetNameAddress.")
	return
}

// GetNameAddressAfterCounter returns a count of finished OutboundMock.GetNameAddress invocations
func (mmGetNameAddress *OutboundMock) GetNameAddressAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetNameAddress.afterGetNameAddressCounter)
}

// GetNameAddressBeforeCounter returns a count of OutboundMock.GetNameAddress invocations
func (mmGetNameAddress *OutboundMock) GetNameAddressBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetNameAddress.beforeGetNameAddressCounter)
}

// MinimockGetNameAddressDone returns true if the count of the GetNameAddress invocations corresponds
// the number of defined expectations
func (m *OutboundMock) MinimockGetNameAddressDone() bool {
	for _, e := range m.GetNameAddressMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetNameAddressMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetNameAddressCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetNameAddress != nil && mm_atomic.LoadUint64(&m.afterGetNameAddressCounter) < 1 {
		return false
	}
	return true
}

// MinimockGetNameAddressInspect logs each unmet expectation
func (m *OutboundMock) MinimockGetNameAddressInspect() {
	for _, e := range m.GetNameAddressMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Error("Expected call to OutboundMock.GetNameAddress")
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetNameAddressMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetNameAddressCounter) < 1 {
		m.t.Error("Expected call to OutboundMock.GetNameAddress")
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetNameAddress != nil && mm_atomic.LoadUint64(&m.afterGetNameAddressCounter) < 1 {
		m.t.Error("Expected call to OutboundMock.GetNameAddress")
	}
}

type mOutboundMockGetRelayID struct {
	mock               *OutboundMock
	defaultExpectation *OutboundMockGetRelayIDExpectation
	expectations       []*OutboundMockGetRelayIDExpectation
}

// OutboundMockGetRelayIDExpectation specifies expectation struct of the Outbound.GetRelayID
type OutboundMockGetRelayIDExpectation struct {
	mock *OutboundMock

	results *OutboundMockGetRelayIDResults
	Counter uint64
}

// OutboundMockGetRelayIDResults contains results of the Outbound.GetRelayID
type OutboundMockGetRelayIDResults struct {
	s1 insolar.ShortNodeID
}

// Expect sets up expected params for Outbound.GetRelayID
func (mmGetRelayID *mOutboundMockGetRelayID) Expect() *mOutboundMockGetRelayID {
	if mmGetRelayID.mock.funcGetRelayID != nil {
		mmGetRelayID.mock.t.Fatalf("OutboundMock.GetRelayID mock is already set by Set")
	}

	if mmGetRelayID.defaultExpectation == nil {
		mmGetRelayID.defaultExpectation = &OutboundMockGetRelayIDExpectation{}
	}

	return mmGetRelayID
}

// Inspect accepts an inspector function that has same arguments as the Outbound.GetRelayID
func (mmGetRelayID *mOutboundMockGetRelayID) Inspect(f func()) *mOutboundMockGetRelayID {
	if mmGetRelayID.mock.inspectFuncGetRelayID != nil {
		mmGetRelayID.mock.t.Fatalf("Inspect function is already set for OutboundMock.GetRelayID")
	}

	mmGetRelayID.mock.inspectFuncGetRelayID = f

	return mmGetRelayID
}

// Return sets up results that will be returned by Outbound.GetRelayID
func (mmGetRelayID *mOutboundMockGetRelayID) Return(s1 insolar.ShortNodeID) *OutboundMock {
	if mmGetRelayID.mock.funcGetRelayID != nil {
		mmGetRelayID.mock.t.Fatalf("OutboundMock.GetRelayID mock is already set by Set")
	}

	if mmGetRelayID.defaultExpectation == nil {
		mmGetRelayID.defaultExpectation = &OutboundMockGetRelayIDExpectation{mock: mmGetRelayID.mock}
	}
	mmGetRelayID.defaultExpectation.results = &OutboundMockGetRelayIDResults{s1}
	return mmGetRelayID.mock
}

//Set uses given function f to mock the Outbound.GetRelayID method
func (mmGetRelayID *mOutboundMockGetRelayID) Set(f func() (s1 insolar.ShortNodeID)) *OutboundMock {
	if mmGetRelayID.defaultExpectation != nil {
		mmGetRelayID.mock.t.Fatalf("Default expectation is already set for the Outbound.GetRelayID method")
	}

	if len(mmGetRelayID.expectations) > 0 {
		mmGetRelayID.mock.t.Fatalf("Some expectations are already set for the Outbound.GetRelayID method")
	}

	mmGetRelayID.mock.funcGetRelayID = f
	return mmGetRelayID.mock
}

// GetRelayID implements Outbound
func (mmGetRelayID *OutboundMock) GetRelayID() (s1 insolar.ShortNodeID) {
	mm_atomic.AddUint64(&mmGetRelayID.beforeGetRelayIDCounter, 1)
	defer mm_atomic.AddUint64(&mmGetRelayID.afterGetRelayIDCounter, 1)

	if mmGetRelayID.inspectFuncGetRelayID != nil {
		mmGetRelayID.inspectFuncGetRelayID()
	}

	if mmGetRelayID.GetRelayIDMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGetRelayID.GetRelayIDMock.defaultExpectation.Counter, 1)

		mm_results := mmGetRelayID.GetRelayIDMock.defaultExpectation.results
		if mm_results == nil {
			mmGetRelayID.t.Fatal("No results are set for the OutboundMock.GetRelayID")
		}
		return (*mm_results).s1
	}
	if mmGetRelayID.funcGetRelayID != nil {
		return mmGetRelayID.funcGetRelayID()
	}
	mmGetRelayID.t.Fatalf("Unexpected call to OutboundMock.GetRelayID.")
	return
}

// GetRelayIDAfterCounter returns a count of finished OutboundMock.GetRelayID invocations
func (mmGetRelayID *OutboundMock) GetRelayIDAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetRelayID.afterGetRelayIDCounter)
}

// GetRelayIDBeforeCounter returns a count of OutboundMock.GetRelayID invocations
func (mmGetRelayID *OutboundMock) GetRelayIDBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetRelayID.beforeGetRelayIDCounter)
}

// MinimockGetRelayIDDone returns true if the count of the GetRelayID invocations corresponds
// the number of defined expectations
func (m *OutboundMock) MinimockGetRelayIDDone() bool {
	for _, e := range m.GetRelayIDMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetRelayIDMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetRelayIDCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetRelayID != nil && mm_atomic.LoadUint64(&m.afterGetRelayIDCounter) < 1 {
		return false
	}
	return true
}

// MinimockGetRelayIDInspect logs each unmet expectation
func (m *OutboundMock) MinimockGetRelayIDInspect() {
	for _, e := range m.GetRelayIDMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Error("Expected call to OutboundMock.GetRelayID")
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetRelayIDMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetRelayIDCounter) < 1 {
		m.t.Error("Expected call to OutboundMock.GetRelayID")
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetRelayID != nil && mm_atomic.LoadUint64(&m.afterGetRelayIDCounter) < 1 {
		m.t.Error("Expected call to OutboundMock.GetRelayID")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *OutboundMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockAsByteStringInspect()

		m.MinimockCanAcceptInspect()

		m.MinimockGetEndpointTypeInspect()

		m.MinimockGetIPAddressInspect()

		m.MinimockGetNameAddressInspect()

		m.MinimockGetRelayIDInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *OutboundMock) MinimockWait(timeout mm_time.Duration) {
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

func (m *OutboundMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockAsByteStringDone() &&
		m.MinimockCanAcceptDone() &&
		m.MinimockGetEndpointTypeDone() &&
		m.MinimockGetIPAddressDone() &&
		m.MinimockGetNameAddressDone() &&
		m.MinimockGetRelayIDDone()
}
