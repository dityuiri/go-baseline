// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/dityuiri/go-baseline/repository (interfaces: IPlaceholderProducer)

// Package repository_mock is a generated GoMock package.
package repository_mock

import (
	context "context"
	reflect "reflect"

	model "github.com/dityuiri/go-baseline/model"
	gomock "github.com/golang/mock/gomock"
)

// MockIPlaceholderProducer is a mock of IPlaceholderProducer interface.
type MockIPlaceholderProducer struct {
	ctrl     *gomock.Controller
	recorder *MockIPlaceholderProducerMockRecorder
}

// MockIPlaceholderProducerMockRecorder is the mock recorder for MockIPlaceholderProducer.
type MockIPlaceholderProducerMockRecorder struct {
	mock *MockIPlaceholderProducer
}

// NewMockIPlaceholderProducer creates a new mock instance.
func NewMockIPlaceholderProducer(ctrl *gomock.Controller) *MockIPlaceholderProducer {
	mock := &MockIPlaceholderProducer{ctrl: ctrl}
	mock.recorder = &MockIPlaceholderProducerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIPlaceholderProducer) EXPECT() *MockIPlaceholderProducerMockRecorder {
	return m.recorder
}

// ProducePlaceholderRecord mocks base method.
func (m *MockIPlaceholderProducer) ProducePlaceholderRecord(arg0 context.Context, arg1 model.PlaceholderMessage) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProducePlaceholderRecord", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProducePlaceholderRecord indicates an expected call of ProducePlaceholderRecord.
func (mr *MockIPlaceholderProducerMockRecorder) ProducePlaceholderRecord(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProducePlaceholderRecord", reflect.TypeOf((*MockIPlaceholderProducer)(nil).ProducePlaceholderRecord), arg0, arg1)
}
