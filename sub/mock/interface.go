// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package mock_sub is a generated GoMock package.
package mock_sub

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	pubsub "github.com/lileio/pubsub"
)

// MockSubProvider is a mock of SubProvider interface.
type MockSubProvider struct {
	ctrl     *gomock.Controller
	recorder *MockSubProviderMockRecorder
}

// MockSubProviderMockRecorder is the mock recorder for MockSubProvider.
type MockSubProviderMockRecorder struct {
	mock *MockSubProvider
}

// NewMockSubProvider creates a new mock instance.
func NewMockSubProvider(ctrl *gomock.Controller) *MockSubProvider {
	mock := &MockSubProvider{ctrl: ctrl}
	mock.recorder = &MockSubProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSubProvider) EXPECT() *MockSubProviderMockRecorder {
	return m.recorder
}

// Init mocks base method.
func (m *MockSubProvider) Init(config map[string]any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Init", config)
	ret0, _ := ret[0].(error)
	return ret0
}

// Init indicates an expected call of Init.
func (mr *MockSubProviderMockRecorder) Init(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockSubProvider)(nil).Init), config)
}

// Name mocks base method.
func (m *MockSubProvider) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockSubProviderMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockSubProvider)(nil).Name))
}

// Run mocks base method.
func (m *MockSubProvider) Run() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Run")
	ret0, _ := ret[0].(error)
	return ret0
}

// Run indicates an expected call of Run.
func (mr *MockSubProviderMockRecorder) Run() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockSubProvider)(nil).Run))
}

// Subscribe mocks base method.
func (m *MockSubProvider) Subscribe(ctx context.Context, instance, topic string, h pubsub.MsgHandler) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Subscribe", ctx, instance, topic, h)
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockSubProviderMockRecorder) Subscribe(ctx, instance, topic, h interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockSubProvider)(nil).Subscribe), ctx, instance, topic, h)
}
