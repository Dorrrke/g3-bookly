package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Dorrrke/g3-bookly/internal/domain/models"
	"github.com/Dorrrke/g3-bookly/internal/logger"
	"github.com/Dorrrke/g3-bookly/internal/server/mocks"
	storerrros "github.com/Dorrrke/g3-bookly/internal/storage/errros"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	testUUID := "test-uuid-134-qwer43"
	testToken, err := createJWTToken(testUUID)
	assert.NoError(t, err)
	logger.Get(false)
	vaid := validator.New()
	var srv Server
	srv.valid = vaid
	r := gin.New()
	r.Use(gin.Recovery())
	r.POST("/register", srv.register)
	httpSrv := httptest.NewServer(r)
	type want struct {
		body       string
		statusCode int
		header     string
	}
	type test struct {
		name    string
		user    models.User
		uuid    string
		request string
		method  string
		mock    bool
		err     error
		want    want
	}

	tests := []test{
		{
			name: "successful call",
			user: models.User{
				Email: "testuser123@yandex.ru",
				Pass:  "qwerty12345678",
				Age:   22,
			},
			uuid:    testUUID,
			request: "/register",
			method:  http.MethodPost,
			mock:    true,
			want: want{
				body:       testUUID,
				statusCode: http.StatusCreated,
				header:     testToken,
			},
		},
		{
			name: "invalid e-mail",
			user: models.User{
				Email: "testuser123@",
				Pass:  "qwerty12345678",
				Age:   22,
			},
			uuid:    testUUID,
			request: "/register",
			method:  http.MethodPost,
			mock:    false,
			want: want{
				body:       `{"error":"Key: 'User.Email' Error:Field validation for 'Email' failed on the 'email' tag"}`,
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "invalid password",
			user: models.User{
				Email: "testuser123@yandex.ru",
				Pass:  "qwerty",
				Age:   22,
			},
			uuid:    testUUID,
			request: "/register",
			method:  http.MethodPost,
			mock:    false,
			want: want{
				body:       `{"error":"Key: 'User.Pass' Error:Field validation for 'Pass' failed on the 'min' tag"}`,
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "invalid age",
			user: models.User{
				Email: "testuser123@yandex.ru",
				Pass:  "qwerty12345687",
				Age:   10,
			},
			uuid:    testUUID,
			request: "/register",
			method:  http.MethodPost,
			mock:    false,
			want: want{
				body:       `{"error":"Key: 'User.Age' Error:Field validation for 'Age' failed on the 'gte' tag"}`,
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "user exists",
			user: models.User{
				Email: "testuser123@yandex.ru",
				Pass:  "qwerty12345687",
				Age:   33,
			},
			uuid:    testUUID,
			request: "/register",
			method:  http.MethodPost,
			err:     storerrros.ErrUserExists,
			mock:    true,
			want: want{
				body:       "user alredy exists",
				statusCode: http.StatusConflict,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			storMock := mocks.NewStorage(t)
			if tc.mock {
				storMock.On("SaveUser", tc.user).Return(tc.uuid, tc.err)
				srv.storage = storMock
			}
			req := resty.New().R()
			req.Method = tc.method
			req.URL = httpSrv.URL + tc.request
			body, err := json.Marshal(tc.user)
			assert.NoError(t, err)
			req.Body = body
			resp, err := req.Send()
			assert.NoError(t, err)
			header := resp.Header().Get("Authorization")
			respBody := string(resp.Body())
			assert.Equal(t, tc.want.statusCode, resp.StatusCode())
			assert.Equal(t, tc.want.header, header)
			assert.Equal(t, tc.want.body, respBody)
		})
	}
}

func TestLogin(t *testing.T) {
	testUUID := "test-uuid-134-qwer43"
	testToken, err := createJWTToken(testUUID)
	assert.NoError(t, err)
	logger.Get(false)
	vaid := validator.New()
	var srv Server
	srv.valid = vaid
	r := gin.New()
	r.Use(gin.Recovery())
	r.POST("/login", srv.login)
	httpSrv := httptest.NewServer(r)
	type want struct {
		body       string
		statusCode int
		header     string
	}
	type test struct {
		name    string
		user    models.User
		uuid    string
		request string
		method  string
		mock    bool
		err     error
		want    want
	}

	tests := []test{
		{
			name: "successful call",
			user: models.User{
				Email: "testuser123@yandex.ru",
				Pass:  "qwerty12345678",
				Age:   22,
			},
			uuid:    testUUID,
			request: "/login",
			method:  http.MethodPost,
			mock:    true,
			want: want{
				body:       "user " + testUUID + " are logined",
				statusCode: http.StatusOK,
				header:     testToken,
			},
		},
		{
			name: "invalid e-mail",
			user: models.User{
				Email: "testuser123@",
				Pass:  "qwerty12345678",
				Age:   22,
			},
			uuid:    testUUID,
			request: "/login",
			method:  http.MethodPost,
			mock:    false,
			want: want{
				body:       `{"error":"Key: 'User.Email' Error:Field validation for 'Email' failed on the 'email' tag"}`,
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "invalid password",
			user: models.User{
				Email: "testuser123@yandex.ru",
				Pass:  "qwerty",
				Age:   22,
			},
			uuid:    testUUID,
			request: "/login",
			method:  http.MethodPost,
			mock:    false,
			want: want{
				body:       `{"error":"Key: 'User.Pass' Error:Field validation for 'Pass' failed on the 'min' tag"}`,
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "invalid age",
			user: models.User{
				Email: "testuser123@yandex.ru",
				Pass:  "qwerty12345687",
				Age:   1,
			},
			uuid:    testUUID,
			request: "/login",
			method:  http.MethodPost,
			mock:    false,
			want: want{
				body:       `{"error":"Key: 'User.Age' Error:Field validation for 'Age' failed on the 'gte' tag"}`,
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "user no exists",
			user: models.User{
				Email: "testuser123@yandex.ru",
				Pass:  "qwerty12345687",
				Age:   33,
			},
			uuid:    testUUID,
			request: "/login",
			method:  http.MethodPost,
			err:     storerrros.ErrUserNoExist,
			mock:    true,
			want: want{
				body:       "invalid login or password: %!w(*errors.errorString=&{user does not exists})",
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "invalid password",
			user: models.User{
				Email: "testuser123@yandex.ru",
				Pass:  "qwerty12345687",
				Age:   33,
			},
			uuid:    testUUID,
			request: "/login",
			method:  http.MethodPost,
			err:     storerrros.ErrInvalidPassword,
			mock:    true,
			want: want{
				body:       storerrros.ErrInvalidPassword.Error(),
				statusCode: http.StatusUnauthorized,
			},
		},
		{
			name: "internal error",
			user: models.User{
				Email: "testuser123@yandex.ru",
				Pass:  "qwerty12345687",
				Age:   33,
			},
			uuid:    testUUID,
			request: "/login",
			method:  http.MethodPost,
			err:     errors.New("internal error"),
			mock:    true,
			want: want{
				body:       `{"error":"internal error"}`,
				statusCode: http.StatusInternalServerError,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			storMock := mocks.NewStorage(t)
			if tc.mock {
				storMock.On("ValidUser", tc.user).Return(tc.uuid, tc.err)
				srv.storage = storMock
			}
			req := resty.New().R()
			req.Method = tc.method
			req.URL = httpSrv.URL + tc.request
			body, err := json.Marshal(tc.user)
			assert.NoError(t, err)
			req.Body = body
			resp, err := req.Send()
			assert.NoError(t, err)
			header := resp.Header().Get("Authorization")
			respBody := string(resp.Body())
			assert.Equal(t, tc.want.statusCode, resp.StatusCode())
			assert.Equal(t, tc.want.header, header)
			assert.Equal(t, tc.want.body, respBody)
		})
	}
}

func BenchmarkLogin(b *testing.B) {
	testUUID := "test-uuid-134-qwer43"
	logger.Get(false)
	vaid := validator.New()
	var srv Server
	srv.valid = vaid
	r := gin.New()
	r.Use(gin.Recovery())
	r.POST("/login", srv.login)
	httpSrv := httptest.NewServer(r)
	user := models.User{
		Email: "testuser123@yandex.ru",
		Pass:  "qwerty12345687",
		Age:   33,
	}

	storMock := mocks.NewStorage(b)
	storMock.On("ValidUser", user).Return(testUUID, nil)
	srv.storage = storMock
	req := resty.New().R()
	req.Method = http.MethodPost
	req.URL = httpSrv.URL + "/login"
	body, err := json.Marshal(user)
	assert.NoError(b, err)
	req.Body = body

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req.Send()
	}
}
