package mock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	"reflect"
	b "test_task/Server/database_interface"
)

// MockService is a mock of Service interface
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// CreateBook  mocks base method
func (m *MockService) CreateBook(arg0 context.Context, arg1 b.Book) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBook", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockServiceMockRecorder) CreateBook(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBook", reflect.TypeOf((*MockService)(nil).CreateBook), arg0, arg1)
}

// UpdateBook  mocks base method
func (m *MockService) UpdateBook(arg0 context.Context, arg1 int, arg2 b.Book) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBook", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockServiceMockRecorder) UpdateBook(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBook", reflect.TypeOf((*MockService)(nil).UpdateBook), arg0, arg1, arg2)
}

// GetBook  mocks base method
func (m *MockService) GetBook(arg0 context.Context, arg1 int) (b.BookWithGenre, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBook", arg0, arg1)
	ret0, _ := ret[0].(b.BookWithGenre)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockServiceMockRecorder) GetBook(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBook", reflect.TypeOf((*MockService)(nil).GetBook), arg0, arg1)
}

// GetAllBooks  mocks base method
func (m *MockService) GetAllBooks(arg0 context.Context, arg1 b.AllBooksRequest) ([]b.Book, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllBooks", arg0, arg1)
	ret0, _ := ret[0].([]b.Book)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockServiceMockRecorder) GetAllBooks(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllBooks", reflect.TypeOf((*MockService)(nil).GetAllBooks), arg0, arg1)
}

// DeleteBook  mocks base method
func (m *MockService) DeleteBook(arg0 context.Context, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBook", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockServiceMockRecorder) DeleteBook(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBook", reflect.TypeOf((*MockService)(nil).DeleteBook), arg0, arg1)
}
