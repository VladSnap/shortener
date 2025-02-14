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

// CreateShortLink mocks base method.
func (m *MockShortLinkRepo) CreateShortLink(arg0 *data.ShortLinkData) (*data.ShortLinkData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateShortLink", arg0)
	ret0, _ := ret[0].(*data.ShortLinkData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateShortLink indicates an expected call of CreateShortLink.
func (mr *MockShortLinkRepoMockRecorder) CreateShortLink(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateShortLink", reflect.TypeOf((*MockShortLinkRepo)(nil).CreateShortLink), arg0)
}

// GetURL mocks base method.
func (m *MockShortLinkRepo) GetURL(arg0 string) (*data.ShortLinkData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURL", arg0)
	ret0, _ := ret[0].(*data.ShortLinkData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetURL indicates an expected call of GetURL.
func (mr *MockShortLinkRepoMockRecorder) GetURL(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURL", reflect.TypeOf((*MockShortLinkRepo)(nil).GetURL), arg0)
}
