package mocks

import (
	reflect "reflect"

	models "github.com/VladSnap/shortener/internal/data/models"
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

// CreateShortLink mocks base method.
func (m *MockShortLinkRepo) CreateShortLink(arg0 *models.ShortLinkData) (*models.ShortLinkData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateShortLink", arg0)
	ret0, _ := ret[0].(*models.ShortLinkData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateShortLink indicates an expected call of CreateShortLink.
func (mr *MockShortLinkRepoMockRecorder) CreateShortLink(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateShortLink",
		reflect.TypeOf((*MockShortLinkRepo)(nil).CreateShortLink), arg0)
}

// GetURL mocks base method.
func (m *MockShortLinkRepo) GetURL(arg0 string) (*models.ShortLinkData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURL", arg0)
	ret0, _ := ret[0].(*models.ShortLinkData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetURL indicates an expected call of GetURL.
func (mr *MockShortLinkRepoMockRecorder) GetURL(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURL",
		reflect.TypeOf((*MockShortLinkRepo)(nil).GetURL), arg0)
}
