// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	ping "payslip-generation-system/internal/entity/ping"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockPingServiceProvider is a mock of PingServiceProvider interface.
type MockPingServiceProvider struct {
	ctrl     *gomock.Controller
	recorder *MockPingServiceProviderMockRecorder
}

// MockPingServiceProviderMockRecorder is the mock recorder for MockPingServiceProvider.
type MockPingServiceProviderMockRecorder struct {
	mock *MockPingServiceProvider
}

// NewMockPingServiceProvider creates a new mock instance.
func NewMockPingServiceProvider(ctrl *gomock.Controller) *MockPingServiceProvider {
	mock := &MockPingServiceProvider{ctrl: ctrl}
	mock.recorder = &MockPingServiceProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPingServiceProvider) EXPECT() *MockPingServiceProviderMockRecorder {
	return m.recorder
}

// Ping mocks base method.
func (m *MockPingServiceProvider) Ping(ctx context.Context) (ping.PingPong, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", ctx)
	ret0, _ := ret[0].(ping.PingPong)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Ping indicates an expected call of Ping.
func (mr *MockPingServiceProviderMockRecorder) Ping(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockPingServiceProvider)(nil).Ping), ctx)
}
