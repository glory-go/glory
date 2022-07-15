// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package mock_config is a generated GoMock package.
package mock_config

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockConfigCenter is a mock of ConfigCenter interface.
type MockConfigCenter struct {
	ctrl     *gomock.Controller
	recorder *MockConfigCenterMockRecorder
}

// MockConfigCenterMockRecorder is the mock recorder for MockConfigCenter.
type MockConfigCenterMockRecorder struct {
	mock *MockConfigCenter
}

// NewMockConfigCenter creates a new mock instance.
func NewMockConfigCenter(ctrl *gomock.Controller) *MockConfigCenter {
	mock := &MockConfigCenter{ctrl: ctrl}
	mock.recorder = &MockConfigCenterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConfigCenter) EXPECT() *MockConfigCenterMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockConfigCenter) Get(params ...any) (string, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range params {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockConfigCenterMockRecorder) Get(params ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockConfigCenter)(nil).Get), params...)
}

// Init mocks base method.
func (m *MockConfigCenter) Init(config map[string]any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Init", config)
	ret0, _ := ret[0].(error)
	return ret0
}

// Init indicates an expected call of Init.
func (mr *MockConfigCenterMockRecorder) Init(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockConfigCenter)(nil).Init), config)
}

// Name mocks base method.
func (m *MockConfigCenter) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockConfigCenterMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockConfigCenter)(nil).Name))
}

// MockComponent is a mock of Component interface.
type MockComponent struct {
	ctrl     *gomock.Controller
	recorder *MockComponentMockRecorder
}

// MockComponentMockRecorder is the mock recorder for MockComponent.
type MockComponentMockRecorder struct {
	mock *MockComponent
}

// NewMockComponent creates a new mock instance.
func NewMockComponent(ctrl *gomock.Controller) *MockComponent {
	mock := &MockComponent{ctrl: ctrl}
	mock.recorder = &MockComponentMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockComponent) EXPECT() *MockComponentMockRecorder {
	return m.recorder
}

// Init mocks base method.
func (m *MockComponent) Init(config map[string]any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Init", config)
	ret0, _ := ret[0].(error)
	return ret0
}

// Init indicates an expected call of Init.
func (mr *MockComponentMockRecorder) Init(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockComponent)(nil).Init), config)
}

// Name mocks base method.
func (m *MockComponent) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockComponentMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockComponent)(nil).Name))
}