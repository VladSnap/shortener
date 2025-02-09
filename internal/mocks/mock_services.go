package mocks

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
func (m *MockShortLinkRepo) AddBatch(ctx context.Context, links []*data.ShortLinkData) ([]*data.ShortLinkData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddBatch", ctx, links)
	ret0, _ := ret[0].([]*data.ShortLinkData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddBatch indicates an expected call of AddBatch.
func (mr *MockShortLinkRepoMockRecorder) AddBatch(ctx, links interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddBatch",
		reflect.TypeOf((*MockShortLinkRepo)(nil).AddBatch), ctx, links)
}

// CreateShortLink mocks base method.
func (m *MockShortLinkRepo) CreateShortLink(link *data.ShortLinkData) (*data.ShortLinkData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateShortLink", link)
	ret0, _ := ret[0].(*data.ShortLinkData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateShortLink indicates an expected call of CreateShortLink.
func (mr *MockShortLinkRepoMockRecorder) CreateShortLink(link interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateShortLink",
		reflect.TypeOf((*MockShortLinkRepo)(nil).CreateShortLink), link)
}

// GetURL mocks base method.
func (m *MockShortLinkRepo) GetURL(shortID string) (*data.ShortLinkData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURL", shortID)
	ret0, _ := ret[0].(*data.ShortLinkData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetURL indicates an expected call of GetURL.
func (mr *MockShortLinkRepoMockRecorder) GetURL(shortID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURL",
		reflect.TypeOf((*MockShortLinkRepo)(nil).GetURL), shortID)
}
