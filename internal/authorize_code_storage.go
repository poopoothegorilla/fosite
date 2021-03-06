// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/ory-am/fosite/handler/core (interfaces: AuthorizeCodeStorage)

package internal

import (
	gomock "github.com/golang/mock/gomock"
	fosite "github.com/ory-am/fosite"
	context "golang.org/x/net/context"
)

// Mock of AuthorizeCodeStorage interface
type MockAuthorizeCodeStorage struct {
	ctrl     *gomock.Controller
	recorder *_MockAuthorizeCodeStorageRecorder
}

// Recorder for MockAuthorizeCodeStorage (not exported)
type _MockAuthorizeCodeStorageRecorder struct {
	mock *MockAuthorizeCodeStorage
}

func NewMockAuthorizeCodeStorage(ctrl *gomock.Controller) *MockAuthorizeCodeStorage {
	mock := &MockAuthorizeCodeStorage{ctrl: ctrl}
	mock.recorder = &_MockAuthorizeCodeStorageRecorder{mock}
	return mock
}

func (_m *MockAuthorizeCodeStorage) EXPECT() *_MockAuthorizeCodeStorageRecorder {
	return _m.recorder
}

func (_m *MockAuthorizeCodeStorage) CreateAuthorizeCodeSession(_param0 context.Context, _param1 string, _param2 fosite.Requester) error {
	ret := _m.ctrl.Call(_m, "CreateAuthorizeCodeSession", _param0, _param1, _param2)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockAuthorizeCodeStorageRecorder) CreateAuthorizeCodeSession(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateAuthorizeCodeSession", arg0, arg1, arg2)
}

func (_m *MockAuthorizeCodeStorage) DeleteAuthorizeCodeSession(_param0 context.Context, _param1 string) error {
	ret := _m.ctrl.Call(_m, "DeleteAuthorizeCodeSession", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockAuthorizeCodeStorageRecorder) DeleteAuthorizeCodeSession(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteAuthorizeCodeSession", arg0, arg1)
}

func (_m *MockAuthorizeCodeStorage) GetAuthorizeCodeSession(_param0 context.Context, _param1 string, _param2 interface{}) (fosite.Requester, error) {
	ret := _m.ctrl.Call(_m, "GetAuthorizeCodeSession", _param0, _param1, _param2)
	ret0, _ := ret[0].(fosite.Requester)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockAuthorizeCodeStorageRecorder) GetAuthorizeCodeSession(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetAuthorizeCodeSession", arg0, arg1, arg2)
}
