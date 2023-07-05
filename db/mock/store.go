// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/YuanData/aquatrade/db/sqlc (interfaces: Store)

// Package mockdb is a generated GoMock package.
package mockdb

import (
	context "context"
	reflect "reflect"

	db "github.com/YuanData/aquatrade/db/sqlc"
	gomock "github.com/golang/mock/gomock"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// AddTraderBalance mocks base method.
func (m *MockStore) AddTraderBalance(arg0 context.Context, arg1 db.AddTraderBalanceParams) (db.Trader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddTraderBalance", arg0, arg1)
	ret0, _ := ret[0].(db.Trader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddTraderBalance indicates an expected call of AddTraderBalance.
func (mr *MockStoreMockRecorder) AddTraderBalance(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTraderBalance", reflect.TypeOf((*MockStore)(nil).AddTraderBalance), arg0, arg1)
}

// CreatePayment mocks base method.
func (m *MockStore) CreatePayment(arg0 context.Context, arg1 db.CreatePaymentParams) (db.Payment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePayment", arg0, arg1)
	ret0, _ := ret[0].(db.Payment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePayment indicates an expected call of CreatePayment.
func (mr *MockStoreMockRecorder) CreatePayment(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePayment", reflect.TypeOf((*MockStore)(nil).CreatePayment), arg0, arg1)
}

// CreateRecord mocks base method.
func (m *MockStore) CreateRecord(arg0 context.Context, arg1 db.CreateRecordParams) (db.Record, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRecord", arg0, arg1)
	ret0, _ := ret[0].(db.Record)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRecord indicates an expected call of CreateRecord.
func (mr *MockStoreMockRecorder) CreateRecord(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRecord", reflect.TypeOf((*MockStore)(nil).CreateRecord), arg0, arg1)
}

// CreateTrader mocks base method.
func (m *MockStore) CreateTrader(arg0 context.Context, arg1 db.CreateTraderParams) (db.Trader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTrader", arg0, arg1)
	ret0, _ := ret[0].(db.Trader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTrader indicates an expected call of CreateTrader.
func (mr *MockStoreMockRecorder) CreateTrader(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTrader", reflect.TypeOf((*MockStore)(nil).CreateTrader), arg0, arg1)
}

// DeleteTrader mocks base method.
func (m *MockStore) DeleteTrader(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTrader", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTrader indicates an expected call of DeleteTrader.
func (mr *MockStoreMockRecorder) DeleteTrader(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTrader", reflect.TypeOf((*MockStore)(nil).DeleteTrader), arg0, arg1)
}

// GetPayment mocks base method.
func (m *MockStore) GetPayment(arg0 context.Context, arg1 int64) (db.Payment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPayment", arg0, arg1)
	ret0, _ := ret[0].(db.Payment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPayment indicates an expected call of GetPayment.
func (mr *MockStoreMockRecorder) GetPayment(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPayment", reflect.TypeOf((*MockStore)(nil).GetPayment), arg0, arg1)
}

// GetRecord mocks base method.
func (m *MockStore) GetRecord(arg0 context.Context, arg1 int64) (db.Record, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRecord", arg0, arg1)
	ret0, _ := ret[0].(db.Record)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRecord indicates an expected call of GetRecord.
func (mr *MockStoreMockRecorder) GetRecord(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecord", reflect.TypeOf((*MockStore)(nil).GetRecord), arg0, arg1)
}

// GetTrader mocks base method.
func (m *MockStore) GetTrader(arg0 context.Context, arg1 int64) (db.Trader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTrader", arg0, arg1)
	ret0, _ := ret[0].(db.Trader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTrader indicates an expected call of GetTrader.
func (mr *MockStoreMockRecorder) GetTrader(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTrader", reflect.TypeOf((*MockStore)(nil).GetTrader), arg0, arg1)
}

// ListPayments mocks base method.
func (m *MockStore) ListPayments(arg0 context.Context, arg1 db.ListPaymentsParams) ([]db.Payment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPayments", arg0, arg1)
	ret0, _ := ret[0].([]db.Payment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListPayments indicates an expected call of ListPayments.
func (mr *MockStoreMockRecorder) ListPayments(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPayments", reflect.TypeOf((*MockStore)(nil).ListPayments), arg0, arg1)
}

// ListRecords mocks base method.
func (m *MockStore) ListRecords(arg0 context.Context, arg1 db.ListRecordsParams) ([]db.Record, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListRecords", arg0, arg1)
	ret0, _ := ret[0].([]db.Record)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListRecords indicates an expected call of ListRecords.
func (mr *MockStoreMockRecorder) ListRecords(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRecords", reflect.TypeOf((*MockStore)(nil).ListRecords), arg0, arg1)
}

// ListTraders mocks base method.
func (m *MockStore) ListTraders(arg0 context.Context, arg1 db.ListTradersParams) ([]db.Trader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTraders", arg0, arg1)
	ret0, _ := ret[0].([]db.Trader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTraders indicates an expected call of ListTraders.
func (mr *MockStoreMockRecorder) ListTraders(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTraders", reflect.TypeOf((*MockStore)(nil).ListTraders), arg0, arg1)
}

// PaymentTx mocks base method.
func (m *MockStore) PaymentTx(arg0 context.Context, arg1 db.PaymentTxParams) (db.PaymentTxResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PaymentTx", arg0, arg1)
	ret0, _ := ret[0].(db.PaymentTxResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PaymentTx indicates an expected call of PaymentTx.
func (mr *MockStoreMockRecorder) PaymentTx(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PaymentTx", reflect.TypeOf((*MockStore)(nil).PaymentTx), arg0, arg1)
}

// UpdateTrader mocks base method.
func (m *MockStore) UpdateTrader(arg0 context.Context, arg1 db.UpdateTraderParams) (db.Trader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTrader", arg0, arg1)
	ret0, _ := ret[0].(db.Trader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateTrader indicates an expected call of UpdateTrader.
func (mr *MockStoreMockRecorder) UpdateTrader(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTrader", reflect.TypeOf((*MockStore)(nil).UpdateTrader), arg0, arg1)
}