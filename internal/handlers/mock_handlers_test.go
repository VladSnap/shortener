// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/VladSnap/shortener/internal/handlers (interfaces: ShorterService)

// Package handlers is a generated GoMock package.
package handlers

import (
	context "context"
	reflect "reflect"

	services "github.com/VladSnap/shortener/internal/services"
	gomock "github.com/golang/mock/gomock"
)

// MockShorterService is a mock of ShorterService interface.
type MockShorterService struct {
	ctrl     *gomock.Controller
	recorder *MockShorterServiceMockRecorder
}

// MockShorterServiceMockRecorder is the mock recorder for MockShorterService.
type MockShorterServiceMockRecorder struct {
	mock *MockShorterService
}

// NewMockShorterService creates a new mock instance.
func NewMockShorterService(ctrl *gomock.Controller) *MockShorterService {
	mock := &MockShorterService{ctrl: ctrl}
	mock.recorder = &MockShorterServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockShorterService) EXPECT() *MockShorterServiceMockRecorder {
	return m.recorder
}

// CreateShortLink mocks base method.
func (m *MockShorterService) CreateShortLink(arg0 context.Context, arg1, arg2 string) (*services.ShortedLink, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateShortLink", arg0, arg1, arg2)
	ret0, _ := ret[0].(*services.ShortedLink)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateShortLink indicates an expected call of CreateShortLink.
func (mr *MockShorterServiceMockRecorder) CreateShortLink(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateShortLink", reflect.TypeOf((*MockShorterService)(nil).CreateShortLink), arg0, arg1, arg2)
}

// CreateShortLinkBatch mocks base method.
func (m *MockShorterService) CreateShortLinkBatch(arg0 context.Context, arg1 []*services.OriginalLink, arg2 string) ([]*services.ShortedLink, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateShortLinkBatch", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*services.ShortedLink)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateShortLinkBatch indicates an expected call of CreateShortLinkBatch.
func (mr *MockShorterServiceMockRecorder) CreateShortLinkBatch(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateShortLinkBatch", reflect.TypeOf((*MockShorterService)(nil).CreateShortLinkBatch), arg0, arg1, arg2)
}

// DeleteBatch mocks base method.
func (m *MockShorterService) DeleteBatch(arg0 context.Context, arg1 []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBatch", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBatch indicates an expected call of DeleteBatch.
func (mr *MockShorterServiceMockRecorder) DeleteBatch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBatch", reflect.TypeOf((*MockShorterService)(nil).DeleteBatch), arg0, arg1)
}

// GetAllByUserID mocks base method.
func (m *MockShorterService) GetAllByUserID(arg0 context.Context, arg1 string) ([]*services.ShortedLink, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllByUserID", arg0, arg1)
	ret0, _ := ret[0].([]*services.ShortedLink)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllByUserID indicates an expected call of GetAllByUserID.
func (mr *MockShorterServiceMockRecorder) GetAllByUserID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllByUserID", reflect.TypeOf((*MockShorterService)(nil).GetAllByUserID), arg0, arg1)
}

// GetURL mocks base method.
func (m *MockShorterService) GetURL(arg0 context.Context, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURL", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetURL indicates an expected call of GetURL.
func (mr *MockShorterServiceMockRecorder) GetURL(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURL", reflect.TypeOf((*MockShorterService)(nil).GetURL), arg0, arg1)
}
