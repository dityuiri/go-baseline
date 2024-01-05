// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/dityuiri/go-baseline/adapter/db (interfaces: ITransaction)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	db "github.com/dityuiri/go-baseline/adapter/db"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockITransaction is a mock of ITransaction interface.
type MockITransaction struct {
	ctrl     *gomock.Controller
	recorder *MockITransactionMockRecorder
}

// MockITransactionMockRecorder is the mock recorder for MockITransaction.
type MockITransactionMockRecorder struct {
	mock *MockITransaction
}

// NewMockITransaction creates a new mock instance.
func NewMockITransaction(ctrl *gomock.Controller) *MockITransaction {
	mock := &MockITransaction{ctrl: ctrl}
	mock.recorder = &MockITransactionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockITransaction) EXPECT() *MockITransactionMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockITransaction) Add(arg0 interface{}, arg1 ...string) (*uuid.UUID, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Add", varargs...)
	ret0, _ := ret[0].(*uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add.
func (mr *MockITransactionMockRecorder) Add(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockITransaction)(nil).Add), varargs...)
}

// Commit mocks base method.
func (m *MockITransaction) Commit() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Commit")
	ret0, _ := ret[0].(error)
	return ret0
}

// Commit indicates an expected call of Commit.
func (mr *MockITransactionMockRecorder) Commit() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*MockITransaction)(nil).Commit))
}

// Execute mocks base method.
func (m *MockITransaction) Execute(arg0 string, arg1 ...interface{}) (db.IResult, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Execute", varargs...)
	ret0, _ := ret[0].(db.IResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute.
func (mr *MockITransactionMockRecorder) Execute(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockITransaction)(nil).Execute), varargs...)
}

// ExecuteContext mocks base method.
func (m *MockITransaction) ExecuteContext(arg0 context.Context, arg1 string, arg2 ...interface{}) (db.IResult, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ExecuteContext", varargs...)
	ret0, _ := ret[0].(db.IResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExecuteContext indicates an expected call of ExecuteContext.
func (mr *MockITransactionMockRecorder) ExecuteContext(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecuteContext", reflect.TypeOf((*MockITransaction)(nil).ExecuteContext), varargs...)
}

// Get mocks base method.
func (m *MockITransaction) Get(arg0 interface{}, arg1 ...string) (interface{}, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockITransactionMockRecorder) Get(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockITransaction)(nil).Get), varargs...)
}

// GetRows mocks base method.
func (m *MockITransaction) GetRows(arg0 interface{}, arg1 ...string) (interface{}, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetRows", varargs...)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRows indicates an expected call of GetRows.
func (mr *MockITransactionMockRecorder) GetRows(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRows", reflect.TypeOf((*MockITransaction)(nil).GetRows), varargs...)
}

// GetRowsStatement mocks base method.
func (m *MockITransaction) GetRowsStatement(arg0 interface{}, arg1 string, arg2 ...string) (interface{}, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetRowsStatement", varargs...)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRowsStatement indicates an expected call of GetRowsStatement.
func (mr *MockITransactionMockRecorder) GetRowsStatement(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRowsStatement", reflect.TypeOf((*MockITransaction)(nil).GetRowsStatement), varargs...)
}

// GetStatement mocks base method.
func (m *MockITransaction) GetStatement(arg0 interface{}, arg1 string, arg2 ...string) (interface{}, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetStatement", varargs...)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStatement indicates an expected call of GetStatement.
func (mr *MockITransactionMockRecorder) GetStatement(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStatement", reflect.TypeOf((*MockITransaction)(nil).GetStatement), varargs...)
}

// New mocks base method.
func (m *MockITransaction) New(arg0 interface{}) (*uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "New", arg0)
	ret0, _ := ret[0].(*uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// New indicates an expected call of New.
func (mr *MockITransactionMockRecorder) New(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "New", reflect.TypeOf((*MockITransaction)(nil).New), arg0)
}

// Query mocks base method.
func (m *MockITransaction) Query(arg0 string, arg1 ...interface{}) (db.IRows, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Query", varargs...)
	ret0, _ := ret[0].(db.IRows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query.
func (mr *MockITransactionMockRecorder) Query(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockITransaction)(nil).Query), varargs...)
}

// QueryContext mocks base method.
func (m *MockITransaction) QueryContext(arg0 context.Context, arg1 string, arg2 ...interface{}) (db.IRows, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryContext", varargs...)
	ret0, _ := ret[0].(db.IRows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryContext indicates an expected call of QueryContext.
func (mr *MockITransactionMockRecorder) QueryContext(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryContext", reflect.TypeOf((*MockITransaction)(nil).QueryContext), varargs...)
}

// QueryRow mocks base method.
func (m *MockITransaction) QueryRow(arg0 string, arg1 ...interface{}) db.IRow {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryRow", varargs...)
	ret0, _ := ret[0].(db.IRow)
	return ret0
}

// QueryRow indicates an expected call of QueryRow.
func (mr *MockITransactionMockRecorder) QueryRow(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryRow", reflect.TypeOf((*MockITransaction)(nil).QueryRow), varargs...)
}

// QueryRowContext mocks base method.
func (m *MockITransaction) QueryRowContext(arg0 context.Context, arg1 string, arg2 ...interface{}) db.IRow {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryRowContext", varargs...)
	ret0, _ := ret[0].(db.IRow)
	return ret0
}

// QueryRowContext indicates an expected call of QueryRowContext.
func (mr *MockITransactionMockRecorder) QueryRowContext(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryRowContext", reflect.TypeOf((*MockITransaction)(nil).QueryRowContext), varargs...)
}

// Remove mocks base method.
func (m *MockITransaction) Remove(arg0 interface{}, arg1 ...string) (*uuid.UUID, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Remove", varargs...)
	ret0, _ := ret[0].(*uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Remove indicates an expected call of Remove.
func (mr *MockITransactionMockRecorder) Remove(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockITransaction)(nil).Remove), varargs...)
}

// Rollback mocks base method.
func (m *MockITransaction) Rollback() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Rollback")
	ret0, _ := ret[0].(error)
	return ret0
}

// Rollback indicates an expected call of Rollback.
func (mr *MockITransactionMockRecorder) Rollback() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rollback", reflect.TypeOf((*MockITransaction)(nil).Rollback))
}

// Update mocks base method.
func (m *MockITransaction) Update(arg0 interface{}, arg1 ...string) (*uuid.UUID, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Update", varargs...)
	ret0, _ := ret[0].(*uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockITransactionMockRecorder) Update(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockITransaction)(nil).Update), varargs...)
}
