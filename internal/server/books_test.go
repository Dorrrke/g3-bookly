package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Dorrrke/g3-bookly/internal/domain/models"
	"github.com/Dorrrke/g3-bookly/internal/logger"
	"github.com/Dorrrke/g3-bookly/internal/server/mocks"
	storerrros "github.com/Dorrrke/g3-bookly/internal/storage/errros"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestAllBooks(t *testing.T) {
	logger.Get(false)
	var srv Server
	r := gin.New()
	r.Use(gin.Recovery())
	r.GET("/books", srv.JWTAuthMiddleware(), srv.allBooks)
	httpSrv := httptest.NewServer(r)
	jwt, err := srv.createJWTToken("test")
	assert.NoError(t, err)

	type want struct {
		body       string
		statusCode int
		err        string
	}
	type test struct {
		name     string
		books    []models.Book
		request  string
		jwt      string
		mockFlag bool
		err      error
		want     want
	}
	tests := []test{
		{
			name:     "default call",
			mockFlag: true,
			books: []models.Book{
				{
					BID:    "BID1",
					Lable:  "Test Book 1",
					Author: "Test Author 1",
					Desc:   "Terst desc 1",
					Age:    114,
					Count:  1,
				},
				{
					BID:    "BID2",
					Lable:  "Test Book 2",
					Author: "Test Author 2",
					Desc:   "Terst desc 2",
					Age:    115,
					Count:  2,
				},
				{
					BID:    "BID3",
					Lable:  "Test Book 3",
					Author: "Test Author 3",
					Desc:   "Terst desc 3",
					Age:    116,
					Count:  3,
				},
				{
					BID:    "BID4",
					Lable:  "Test Book 4",
					Author: "Test Author 4",
					Desc:   "Terst desc 4",
					Age:    13,
					Count:  4,
				},
			},
			request: "/books",
			jwt:     jwt,
			want: want{
				body:       `[{"bid":"BID1","lable":"Test Book 1","author":"Test Author 1","desc":"Terst desc 1","age":114,"count":1},{"bid":"BID2","lable":"Test Book 2","author":"Test Author 2","desc":"Terst desc 2","age":115,"count":2},{"bid":"BID3","lable":"Test Book 3","author":"Test Author 3","desc":"Terst desc 3","age":116,"count":3},{"bid":"BID4","lable":"Test Book 4","author":"Test Author 4","desc":"Terst desc 4","age":13,"count":4}]`,
				statusCode: http.StatusOK,
			},
		},
		{
			name:     "invalid token call",
			mockFlag: false,
			request:  "/books",
			want: want{
				body:       `invalid token{"error":"User ID not found"}`,
				statusCode: http.StatusUnauthorized,
			},
		},
		{
			name:     "empty books list call",
			mockFlag: true,
			request:  "/books",
			jwt:      jwt,
			err:      storerrros.ErrEmptyBooksList,
			want: want{
				body:       `empty books list`,
				statusCode: http.StatusNotFound,
			},
		},
		{
			name:     "error call",
			mockFlag: true,
			request:  "/books",
			jwt:      jwt,
			err:      errors.New("test err"),
			want: want{
				body:       `{"error":"test err"}`,
				statusCode: http.StatusInternalServerError,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			storMock := mocks.NewStorage(t)
			if tc.mockFlag {
				storMock.On("GetBooks").Return(tc.books, tc.err)
			}
			srv.storage = storMock
			req := resty.New().R()
			req.Method = http.MethodGet
			req.URL = httpSrv.URL + tc.request
			req.SetHeader("Authorization", tc.jwt)
			resp, err := req.Send()
			assert.NoError(t, err)
			assert.Equal(t, tc.want.statusCode, resp.StatusCode())
			assert.Equal(t, tc.want.body, string(resp.Body()))
		})
	}
}
