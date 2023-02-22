// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/devldavydov/promytheus/internal/server/storage (interfaces: Storage)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	metric "github.com/devldavydov/promytheus/internal/common/metric"
	storage "github.com/devldavydov/promytheus/internal/server/storage"
	gomock "github.com/golang/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockStorage) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockStorageMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStorage)(nil).Close))
}

// GetAllMetrics mocks base method.
func (m *MockStorage) GetAllMetrics() ([]storage.StorageItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllMetrics")
	ret0, _ := ret[0].([]storage.StorageItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllMetrics indicates an expected call of GetAllMetrics.
func (mr *MockStorageMockRecorder) GetAllMetrics() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllMetrics", reflect.TypeOf((*MockStorage)(nil).GetAllMetrics))
}

// GetCounterMetric mocks base method.
func (m *MockStorage) GetCounterMetric(arg0 string) (metric.Counter, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCounterMetric", arg0)
	ret0, _ := ret[0].(metric.Counter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCounterMetric indicates an expected call of GetCounterMetric.
func (mr *MockStorageMockRecorder) GetCounterMetric(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCounterMetric", reflect.TypeOf((*MockStorage)(nil).GetCounterMetric), arg0)
}

// GetGaugeMetric mocks base method.
func (m *MockStorage) GetGaugeMetric(arg0 string) (metric.Gauge, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGaugeMetric", arg0)
	ret0, _ := ret[0].(metric.Gauge)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGaugeMetric indicates an expected call of GetGaugeMetric.
func (mr *MockStorageMockRecorder) GetGaugeMetric(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGaugeMetric", reflect.TypeOf((*MockStorage)(nil).GetGaugeMetric), arg0)
}

// Ping mocks base method.
func (m *MockStorage) Ping() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(bool)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockStorageMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockStorage)(nil).Ping))
}

// SetCounterMetric mocks base method.
func (m *MockStorage) SetCounterMetric(arg0 string, arg1 metric.Counter) (metric.Counter, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetCounterMetric", arg0, arg1)
	ret0, _ := ret[0].(metric.Counter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetCounterMetric indicates an expected call of SetCounterMetric.
func (mr *MockStorageMockRecorder) SetCounterMetric(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCounterMetric", reflect.TypeOf((*MockStorage)(nil).SetCounterMetric), arg0, arg1)
}

// SetGaugeMetric mocks base method.
func (m *MockStorage) SetGaugeMetric(arg0 string, arg1 metric.Gauge) (metric.Gauge, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetGaugeMetric", arg0, arg1)
	ret0, _ := ret[0].(metric.Gauge)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SetGaugeMetric indicates an expected call of SetGaugeMetric.
func (mr *MockStorageMockRecorder) SetGaugeMetric(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetGaugeMetric", reflect.TypeOf((*MockStorage)(nil).SetGaugeMetric), arg0, arg1)
}

// SetMetrics mocks base method.
func (m *MockStorage) SetMetrics(arg0 []storage.StorageItem) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetMetrics", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetMetrics indicates an expected call of SetMetrics.
func (mr *MockStorageMockRecorder) SetMetrics(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMetrics", reflect.TypeOf((*MockStorage)(nil).SetMetrics), arg0)
}
