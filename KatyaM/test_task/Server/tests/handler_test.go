package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	client2 "test_task/Server/client"
	b "test_task/Server/database_interface"
	"test_task/Server/errors"
	handler2 "test_task/Server/handler"
	mock2 "test_task/Server/mock"
	s "test_task/Server/service"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

type env struct {
	ctrl   *gomock.Controller
	mock   *mock2.MockService
	server *httptest.Server
	client *client2.BookClient
}

func (e *env) stop() {
	e.server.Close()
	e.ctrl.Finish()
}

func newEnv(t *testing.T) (*env, func()) {
	env := &env{}
	env.ctrl = gomock.NewController(t)
	env.mock = mock2.NewMockService(env.ctrl)

	log := zaptest.NewLogger(t)

	mux := http.NewServeMux()

	handler := handler2.NewBookHandler(log, s.BookService{L: log}, env.mock)
	handler.Register(mux)

	env.server = httptest.NewServer(mux)

	env.client = client2.NewBookClient(log, env.server.URL)

	return env, env.stop
}

func TestCreateError(t *testing.T) {
	env, stop := newEnv(t)
	defer stop()
	ctx := context.Background()
	req := b.Book{"b", 3, 5, 5}
	env.mock.EXPECT().CreateBook(gomock.Any(), req).Return(0, errors.ErrDBRequest)
	_, err := env.client.Create(ctx, "b", 3, 5, 5)
	require.Error(t, err)
	require.Contains(t, err.Error(), errors.ErrDBRequest.Error())
}

func TestCreateSameName(t *testing.T) {
	env, stop := newEnv(t)
	defer stop()
	ctx := context.Background()
	req_1 := b.Book{"b", 3, 5, 5}
	req_2 := b.Book{"c", 5, 6, 7}
	id_1 := 1
	id_2 := 2
	env.mock.EXPECT().CreateBook(gomock.Any(), req_1).Times(1).Return(id_1, nil)
	env.mock.EXPECT().CreateBook(gomock.Any(), req_2).Times(1).Return(id_2, nil)
	env.mock.EXPECT().CreateBook(gomock.Any(), req_1).Times(1).Return(0, errors.ErrSameName)
	client_id_1, err := env.client.Create(ctx, "b", 3, 5, 5)
	require.NoError(t, err)
	require.Equal(t, id_1, client_id_1)

	client_id_2, err := env.client.Create(ctx, "c", 5, 6, 7)
	require.NoError(t, err)
	require.Equal(t, id_2, client_id_2)

	_, err = env.client.Create(ctx, "b", 3, 5, 5)
	require.Error(t, err)
	require.Contains(t, err.Error(), errors.ErrSameName.Error())
}

func TestCreateBadRequest(t *testing.T) {
	env, stop := newEnv(t)
	defer stop()
	ctx := context.Background()
	req := b.Book{"", 0, 0, 0}
	env.mock.EXPECT().CreateBook(gomock.Any(), req).Return(0, errors.ErrDBRequest)
	_, err := env.client.Create(ctx, "", 0, 0, 0)
	require.Error(t, err)
	require.Contains(t, err.Error(), errors.ErrDBRequest.Error())
}

func TestUpdateError(t *testing.T) {
	env, stop := newEnv(t)
	defer stop()
	ctx, _ := context.WithCancel(context.Background())
	req := b.Book{"b", 3, 5, 5}
	id := 1
	env.mock.EXPECT().UpdateBook(gomock.Any(), id, req).Return(errors.ErrWhileUpdate)
	err := env.client.Update(ctx, id, "b", 3, 5, 5)
	require.Error(t, err)
	require.Contains(t, err.Error(), errors.ErrWhileUpdate.Error())
}

func TestUpdateNoError(t *testing.T) {
	env, stop := newEnv(t)
	defer stop()
	ctx, _ := context.WithCancel(context.Background())
	req := b.Book{"b", 3, 5, 5}
	id := 1
	env.mock.EXPECT().UpdateBook(gomock.Any(), id, req).Return(nil)
	err := env.client.Update(ctx, id, "b", 3, 5, 5)
	require.NoError(t, err)
}

func TestDeleteError(t *testing.T) {
	env, stop := newEnv(t)
	defer stop()
	ctx := context.Background()
	id := 3
	env.mock.EXPECT().DeleteBook(gomock.Any(), id).Return(errors.ErrDBRequest)
	err := env.client.Delete(ctx, id)
	require.Error(t, err)
	require.Contains(t, err.Error(), errors.ErrDBRequest.Error())
}

func TestDeleteNoError(t *testing.T) {
	env, stop := newEnv(t)
	defer stop()
	ctx := context.Background()
	id := 3
	env.mock.EXPECT().DeleteBook(gomock.Any(), id).Return(nil)
	err := env.client.Delete(ctx, id)
	require.NoError(t, err)
}

func TestGetBookError(t *testing.T) {
	env, stop := newEnv(t)
	defer stop()
	ctx := context.Background()
	id := 3
	env.mock.EXPECT().GetBook(gomock.Any(), id).Return(b.BookWithGenre{}, errors.ErrDBRequest)
	_, err := env.client.GetBook(ctx, id)
	require.Error(t, err)
	require.Contains(t, err.Error(), errors.ErrDBRequest.Error())
}

func TestGetBookNoError(t *testing.T) {
	env, stop := newEnv(t)
	defer stop()
	ctx := context.Background()
	id := 3
	ret := b.BookWithGenre{"aaa", 1, "Classics", 3}
	env.mock.EXPECT().GetBook(gomock.Any(), id).Return(ret, nil)
	book, err := env.client.GetBook(ctx, id)
	require.NoError(t, err)
	require.Equal(t, ret, book)
}

func TestGetAllError(t *testing.T) {
	env, stop := newEnv(t)
	defer stop()
	ctx := context.Background()
	req := b.AllBooksRequest{Genre: 1, Max_price: 10}
	env.mock.EXPECT().GetAllBooks(gomock.Any(), req).Return(nil, errors.ErrDBRequest)
	_, err := env.client.GetAllBooks(ctx, req)
	require.Error(t, err)
	require.Contains(t, err.Error(), errors.ErrDBRequest.Error())
}

func TestGetAllNoError(t *testing.T) {
	env, stop := newEnv(t)
	defer stop()
	ctx := context.Background()
	req := b.AllBooksRequest{Genre: 1, Max_price: 10}
	ans := []b.Book{b.Book{"bbb", 1, 2, 3}}
	env.mock.EXPECT().GetAllBooks(gomock.Any(), req).Return(ans, nil)
	books, err := env.client.GetAllBooks(ctx, req)
	require.NoError(t, err)
	require.Equal(t, ans, books)
}
