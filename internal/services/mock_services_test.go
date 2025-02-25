// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/VladSnap/shortener/internal/services (interfaces: ShortLinkRepo)

// Package services is a generated GoMock package.
package services

import (
	context "context"
	reflect "reflect"

	data "github.com/VladSnap/shortener/internal/data"
	gomock "github.com/golang/mock/gomock"
)

// MockShortLinkRepo is a mock of ShortLinkRepo interface.
type MockShortLinkRepo struct {
	ctrl     *gomock.Controller
	recorder *MockShortLinkRepoMockRecorder
}

// MockShortLinkRepoMockRecorder is the mock recorder for MockShortLinkRepo.
type MockShortLinkRepoMockRecorder struct {
	mock *MockShortLinkRepo
}

// NewMockShortLinkRepo creates a new mock instance.
func NewMockShortLinkRepo(ctrl *gomock.Controller) *MockShortLinkRepo {
	mock := &MockShortLinkRepo{ctrl: ctrl}
	mock.recorder = &MockShortLinkRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockShortLinkRepo) EXPECT() *MockShortLinkRepoMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockShortLinkRepo) Add(arg0 context.Context, arg1 *data.ShortLinkData) (*data.ShortLinkData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", arg0, arg1)
	ret0, _ := ret[0].(*data.ShortLinkData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add.
func (mr *MockShortLinkRepoMockRecorder) Add(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockShortLinkRepo)(nil).Add), arg0, arg1)
}

// AddBatch mocks base method.
func (m *MockShortLinkRepo) AddBatch(arg0 context.Context, arg1 []*data.ShortLinkData) ([]*data.ShortLinkData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddBatch", arg0, arg1)
	ret0, _ := ret[0].([]*data.ShortLinkData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddBatch indicates an expected call of AddBatch.
func (mr *MockShortLinkRepoMockRecorder) AddBatch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddBatch", reflect.TypeOf((*MockShortLinkRepo)(nil).AddBatch), arg0, arg1)
}

// DeleteBatch mocks base method.
func (m *MockShortLinkRepo) DeleteBatch(arg0 context.Context, arg1 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBatch", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBatch indicates an expected call of DeleteBatch.
func (mr *MockShortLinkRepoMockRecorder) DeleteBatch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBatch", reflect.TypeOf((*MockShortLinkRepo)(nil).DeleteBatch), arg0, arg1)
}

// Get mocks base method.
func (m *MockShortLinkRepo) Get(arg0 context.Context, arg1 string) (*data.ShortLinkData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*data.ShortLinkData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockShortLinkRepoMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockShortLinkRepo)(nil).Get), arg0, arg1)
}

// GetAllByUserID mocks base method.
func (m *MockShortLinkRepo) GetAllByUserID(arg0 context.Context, arg1 string) ([]*data.ShortLinkData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllByUserID", arg0, arg1)
	ret0, _ := ret[0].([]*data.ShortLinkData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllByUserID indicates an expected call of GetAllByUserID.
func (mr *MockShortLinkRepoMockRecorder) GetAllByUserID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllByUserID", reflect.TypeOf((*MockShortLinkRepo)(nil).GetAllByUserID), arg0, arg1)
}
