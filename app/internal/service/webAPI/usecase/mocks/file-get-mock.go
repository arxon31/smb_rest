// Code generated by MockGen. DO NOT EDIT.
// Source: file-get.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockFileDownloader is a mock of FileDownloader interface.
type MockFileDownloader struct {
	ctrl     *gomock.Controller
	recorder *MockFileDownloaderMockRecorder
}

// MockFileDownloaderMockRecorder is the mock recorder for MockFileDownloader.
type MockFileDownloaderMockRecorder struct {
	mock *MockFileDownloader
}

// NewMockFileDownloader creates a new mock instance.
func NewMockFileDownloader(ctrl *gomock.Controller) *MockFileDownloader {
	mock := &MockFileDownloader{ctrl: ctrl}
	mock.recorder = &MockFileDownloaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileDownloader) EXPECT() *MockFileDownloaderMockRecorder {
	return m.recorder
}

// GetFile mocks base method.
func (m *MockFileDownloader) GetFile(ctx context.Context, path string) (string, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadFile", ctx, path)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetFile indicates an expected call of GetFile.
func (mr *MockFileDownloaderMockRecorder) GetFile(ctx, path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadFile", reflect.TypeOf((*MockFileDownloader)(nil).GetFile), ctx, path)
}
