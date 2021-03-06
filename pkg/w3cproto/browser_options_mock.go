// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/w3c/browser_options.go

// Package w3c is a generated GoMock package.
package w3cproto

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockBrowserOptions is a mock of BrowserOptions interface
type MockBrowserOptions struct {
	ctrl     *gomock.Controller
	recorder *MockBrowserOptionsMockRecorder
}

// MockBrowserOptionsMockRecorder is the mock recorder for MockBrowserOptions
type MockBrowserOptionsMockRecorder struct {
	mock *MockBrowserOptions
}

// NewMockBrowserOptions creates a new mock instance
func NewMockBrowserOptions(ctrl *gomock.Controller) *MockBrowserOptions {
	mock := &MockBrowserOptions{ctrl: ctrl}
	mock.recorder = &MockBrowserOptionsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBrowserOptions) EXPECT() *MockBrowserOptionsMockRecorder {
	return m.recorder
}

// FirstMatch mocks base method
func (m *MockBrowserOptions) FirstMatch() []Capabilities {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FirstMatch")
	ret0, _ := ret[0].([]Capabilities)
	return ret0
}

// FirstMatch indicates an expected call of FirstMatch
func (mr *MockBrowserOptionsMockRecorder) FirstMatch() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FirstMatch", reflect.TypeOf((*MockBrowserOptions)(nil).FirstMatch))
}

// AlwaysMatch mocks base method
func (m *MockBrowserOptions) AlwaysMatch() Capabilities {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AlwaysMatch")
	ret0, _ := ret[0].(Capabilities)
	return ret0
}

// AlwaysMatch indicates an expected call of AlwaysMatch
func (mr *MockBrowserOptionsMockRecorder) AlwaysMatch() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AlwaysMatch", reflect.TypeOf((*MockBrowserOptions)(nil).AlwaysMatch))
}
