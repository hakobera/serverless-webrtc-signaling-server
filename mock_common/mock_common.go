// Code generated by MockGen. DO NOT EDIT.
// Source: common/interface.go

// Package mock_common is a generated GoMock package.
package mock_common

import (
	gomock "github.com/golang/mock/gomock"
	common "github.com/hakobera/serverless-webrtc-signaling-server/common"
	reflect "reflect"
)

// MockApiGatewayManagementAPI is a mock of ApiGatewayManagementAPI interface
type MockApiGatewayManagementAPI struct {
	ctrl     *gomock.Controller
	recorder *MockApiGatewayManagementAPIMockRecorder
}

// MockApiGatewayManagementAPIMockRecorder is the mock recorder for MockApiGatewayManagementAPI
type MockApiGatewayManagementAPIMockRecorder struct {
	mock *MockApiGatewayManagementAPI
}

// NewMockApiGatewayManagementAPI creates a new mock instance
func NewMockApiGatewayManagementAPI(ctrl *gomock.Controller) *MockApiGatewayManagementAPI {
	mock := &MockApiGatewayManagementAPI{ctrl: ctrl}
	mock.recorder = &MockApiGatewayManagementAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockApiGatewayManagementAPI) EXPECT() *MockApiGatewayManagementAPIMockRecorder {
	return m.recorder
}

// PostToConnection mocks base method
func (m *MockApiGatewayManagementAPI) PostToConnection(connectionID, body string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostToConnection", connectionID, body)
	ret0, _ := ret[0].(error)
	return ret0
}

// PostToConnection indicates an expected call of PostToConnection
func (mr *MockApiGatewayManagementAPIMockRecorder) PostToConnection(connectionID, body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostToConnection", reflect.TypeOf((*MockApiGatewayManagementAPI)(nil).PostToConnection), connectionID, body)
}

// MockDB is a mock of DB interface
type MockDB struct {
	ctrl     *gomock.Controller
	recorder *MockDBMockRecorder
}

// MockDBMockRecorder is the mock recorder for MockDB
type MockDBMockRecorder struct {
	mock *MockDB
}

// NewMockDB creates a new mock instance
func NewMockDB(ctrl *gomock.Controller) *MockDB {
	mock := &MockDB{ctrl: ctrl}
	mock.recorder = &MockDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDB) EXPECT() *MockDBMockRecorder {
	return m.recorder
}

// Table mocks base method
func (m *MockDB) Table(name string) common.Table {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Table", name)
	ret0, _ := ret[0].(common.Table)
	return ret0
}

// Table indicates an expected call of Table
func (mr *MockDBMockRecorder) Table(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Table", reflect.TypeOf((*MockDB)(nil).Table), name)
}

// TxPut mocks base method
func (m *MockDB) TxPut(items ...common.TableItem) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range items {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "TxPut", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// TxPut indicates an expected call of TxPut
func (mr *MockDBMockRecorder) TxPut(items ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TxPut", reflect.TypeOf((*MockDB)(nil).TxPut), items...)
}

// RoomsTable mocks base method
func (m *MockDB) RoomsTable() common.Table {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RoomsTable")
	ret0, _ := ret[0].(common.Table)
	return ret0
}

// RoomsTable indicates an expected call of RoomsTable
func (mr *MockDBMockRecorder) RoomsTable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RoomsTable", reflect.TypeOf((*MockDB)(nil).RoomsTable))
}

// ConnectionsTable mocks base method
func (m *MockDB) ConnectionsTable() common.Table {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConnectionsTable")
	ret0, _ := ret[0].(common.Table)
	return ret0
}

// ConnectionsTable indicates an expected call of ConnectionsTable
func (mr *MockDBMockRecorder) ConnectionsTable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConnectionsTable", reflect.TypeOf((*MockDB)(nil).ConnectionsTable))
}

// MockTable is a mock of Table interface
type MockTable struct {
	ctrl     *gomock.Controller
	recorder *MockTableMockRecorder
}

// MockTableMockRecorder is the mock recorder for MockTable
type MockTableMockRecorder struct {
	mock *MockTable
}

// NewMockTable creates a new mock instance
func NewMockTable(ctrl *gomock.Controller) *MockTable {
	mock := &MockTable{ctrl: ctrl}
	mock.recorder = &MockTableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTable) EXPECT() *MockTableMockRecorder {
	return m.recorder
}

// FindOne mocks base method
func (m *MockTable) FindOne(column string, key, out interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOne", column, key, out)
	ret0, _ := ret[0].(error)
	return ret0
}

// FindOne indicates an expected call of FindOne
func (mr *MockTableMockRecorder) FindOne(column, key, out interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOne", reflect.TypeOf((*MockTable)(nil).FindOne), column, key, out)
}

// Put mocks base method
func (m *MockTable) Put(row interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", row)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put
func (mr *MockTableMockRecorder) Put(row interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockTable)(nil).Put), row)
}

// Delete mocks base method
func (m *MockTable) Delete(column string, key interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", column, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockTableMockRecorder) Delete(column, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTable)(nil).Delete), column, key)
}
