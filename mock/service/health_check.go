// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/dityuiri/go-baseline/service (interfaces: IHealthCheckService)

// Package service_mock is a generated GoMock package.
package service_mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIHealthCheckService is a mock of IHealthCheckService interface.
type MockIHealthCheckService struct {
	ctrl     *gomock.Controller
	recorder *MockIHealthCheckServiceMockRecorder
}

// MockIHealthCheckServiceMockRecorder is the mock recorder for MockIHealthCheckService.
type MockIHealthCheckServiceMockRecorder struct {
	mock *MockIHealthCheckService
}

// NewMockIHealthCheckService creates a new mock instance.
func NewMockIHealthCheckService(ctrl *gomock.Controller) *MockIHealthCheckService {
	mock := &MockIHealthCheckService{ctrl: ctrl}
	mock.recorder = &MockIHealthCheckServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIHealthCheckService) EXPECT() *MockIHealthCheckServiceMockRecorder {
	return m.recorder
}

// Ping mocks base method.
func (m *MockIHealthCheckService) Ping() map[string]string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(map[string]string)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockIHealthCheckServiceMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockIHealthCheckService)(nil).Ping))
}