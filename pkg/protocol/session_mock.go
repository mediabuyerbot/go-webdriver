// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/protocol/session.go

// Package protocol is a generated GoMock package.
package protocol

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockSession is a mock of Session interface
type MockSession struct {
	ctrl     *gomock.Controller
	recorder *MockSessionMockRecorder
}

// MockSessionMockRecorder is the mock recorder for MockSession
type MockSessionMockRecorder struct {
	mock *MockSession
}

// NewMockSession creates a new mock instance
func NewMockSession(ctrl *gomock.Controller) *MockSession {
	mock := &MockSession{ctrl: ctrl}
	mock.recorder = &MockSessionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSession) EXPECT() *MockSessionMockRecorder {
	return m.recorder
}

// ID mocks base method
func (m *MockSession) ID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ID indicates an expected call of ID
func (mr *MockSessionMockRecorder) ID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ID", reflect.TypeOf((*MockSession)(nil).ID))
}

// Capabilities mocks base method
func (m *MockSession) Capabilities() Capabilities {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Capabilities")
	ret0, _ := ret[0].(Capabilities)
	return ret0
}

// Capabilities indicates an expected call of Capabilities
func (mr *MockSessionMockRecorder) Capabilities() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Capabilities", reflect.TypeOf((*MockSession)(nil).Capabilities))
}

// Status mocks base method
func (m *MockSession) Status(arg0 context.Context) (Status, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Status", arg0)
	ret0, _ := ret[0].(Status)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Status indicates an expected call of Status
func (mr *MockSessionMockRecorder) Status(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Status", reflect.TypeOf((*MockSession)(nil).Status), arg0)
}

// Delete mocks base method
func (m *MockSession) Delete(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockSessionMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockSession)(nil).Delete), arg0)
}

// MockOptions is a mock of Options interface
type MockOptions struct {
	ctrl     *gomock.Controller
	recorder *MockOptionsMockRecorder
}

// MockOptionsMockRecorder is the mock recorder for MockOptions
type MockOptionsMockRecorder struct {
	mock *MockOptions
}

// NewMockOptions creates a new mock instance
func NewMockOptions(ctrl *gomock.Controller) *MockOptions {
	mock := &MockOptions{ctrl: ctrl}
	mock.recorder = &MockOptionsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockOptions) EXPECT() *MockOptionsMockRecorder {
	return m.recorder
}

// Proxy mocks base method
func (m *MockOptions) Proxy() *Proxy {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Proxy")
	ret0, _ := ret[0].(*Proxy)
	return ret0
}

// Proxy indicates an expected call of Proxy
func (mr *MockOptionsMockRecorder) Proxy() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Proxy", reflect.TypeOf((*MockOptions)(nil).Proxy))
}

// FirstMatch mocks base method
func (m *MockOptions) FirstMatch() []O {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FirstMatch")
	ret0, _ := ret[0].([]O)
	return ret0
}

// FirstMatch indicates an expected call of FirstMatch
func (mr *MockOptionsMockRecorder) FirstMatch() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FirstMatch", reflect.TypeOf((*MockOptions)(nil).FirstMatch))
}

// AlwaysMatch mocks base method
func (m *MockOptions) AlwaysMatch() O {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AlwaysMatch")
	ret0, _ := ret[0].(O)
	return ret0
}

// AlwaysMatch indicates an expected call of AlwaysMatch
func (mr *MockOptionsMockRecorder) AlwaysMatch() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AlwaysMatch", reflect.TypeOf((*MockOptions)(nil).AlwaysMatch))
}
