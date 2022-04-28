// Code generated by MockGen. DO NOT EDIT.
// Source: library/lock/lock.go

// Package redis is a generated GoMock package.
package redis

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockLocker is a mock of Locker interface.
type MockLocker struct {
	ctrl     *gomock.Controller
	recorder *MockLockerMockRecorder
}

// MockLockerMockRecorder is the mock recorder for MockLocker.
type MockLockerMockRecorder struct {
	mock *MockLocker
}

// NewMockLocker creates a new mock instance.
func NewMockLocker(ctrl *gomock.Controller) *MockLocker {
	mock := &MockLocker{ctrl: ctrl}
	mock.recorder = &MockLockerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLocker) EXPECT() *MockLockerMockRecorder {
	return m.recorder
}

// Lock mocks base method.
func (m *MockLocker) Lock(ctx context.Context, key string, random interface{}, duration time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Lock", ctx, key, random, duration)
	ret0, _ := ret[0].(error)
	return ret0
}

// Lock indicates an expected call of Lock.
func (mr *MockLockerMockRecorder) Lock(ctx, key, random, duration interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Lock", reflect.TypeOf((*MockLocker)(nil).Lock), ctx, key, random, duration)
}

// Unlock mocks base method.
func (m *MockLocker) Unlock(ctx context.Context, key string, random interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unlock", ctx, key, random)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unlock indicates an expected call of Unlock.
func (mr *MockLockerMockRecorder) Unlock(ctx, key, random interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unlock", reflect.TypeOf((*MockLocker)(nil).Unlock), ctx, key, random)
}
