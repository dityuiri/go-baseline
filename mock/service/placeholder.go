// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/dityuiri/go-baseline/service (interfaces: IPlaceholderService)

// Package service_mock is a generated GoMock package.
package service_mock

import (
	context "context"
	dto "github.com/dityuiri/go-baseline/model"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIPlaceholderService is a mock of IPlaceholderService interface.
type MockIPlaceholderService struct {
	ctrl     *gomock.Controller
	recorder *MockIPlaceholderServiceMockRecorder
}

// MockIPlaceholderServiceMockRecorder is the mock recorder for MockIPlaceholderService.
type MockIPlaceholderServiceMockRecorder struct {
	mock *MockIPlaceholderService
}

// NewMockIPlaceholderService creates a new mock instance.
func NewMockIPlaceholderService(ctrl *gomock.Controller) *MockIPlaceholderService {
	mock := &MockIPlaceholderService{ctrl: ctrl}
	mock.recorder = &MockIPlaceholderServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIPlaceholderService) EXPECT() *MockIPlaceholderServiceMockRecorder {
	return m.recorder
}

// CreateNewPlaceholder mocks base method.
func (m *MockIPlaceholderService) CreateNewPlaceholder(arg0 context.Context, arg1 dto.PlaceholderCreateRequest) (dto.PlaceholderCreateResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNewPlaceholder", arg0, arg1)
	ret0, _ := ret[0].(dto.PlaceholderCreateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateNewPlaceholder indicates an expected call of CreateNewPlaceholder.
func (mr *MockIPlaceholderServiceMockRecorder) CreateNewPlaceholder(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNewPlaceholder", reflect.TypeOf((*MockIPlaceholderService)(nil).CreateNewPlaceholder), arg0, arg1)
}

// GetPlaceholder mocks base method.
func (m *MockIPlaceholderService) GetPlaceholder(arg0 context.Context, arg1 string) (dto.PlaceholderGetResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPlaceholder", arg0, arg1)
	ret0, _ := ret[0].(dto.PlaceholderGetResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPlaceholder indicates an expected call of GetPlaceholder.
func (mr *MockIPlaceholderServiceMockRecorder) GetPlaceholder(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPlaceholder", reflect.TypeOf((*MockIPlaceholderService)(nil).GetPlaceholder), arg0, arg1)
}
