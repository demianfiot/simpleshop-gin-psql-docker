package handler

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"prac/pkg/service"
	mock_service "prac/pkg/service/mocks"
	"prac/todo"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestHandler_SignUp(t *testing.T) {
	type mockBahavior func(s *mock_service.MockAuthorization, user todo.User)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            todo.User
		mockBehavior         mockBahavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"name":"Test","email":"test","password":"qwerty","role":"customer"}`,
			inputUser: todo.User{
				Name:         "Test",
				Email:        "test",
				PasswordHash: "qwerty",
				Role:         "customer",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user todo.User) {
				s.EXPECT().
					CreateUser(gomock.Any(), user).
					Return(1, nil)
			},

			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:                 "Empty Fields",
			inputBody:            `{"password":"qwerty","role":"customer"}`,
			mockBehavior:         func(s *mock_service.MockAuthorization, user todo.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service Failure",
			inputBody: `{"name":"Test","email":"test","password":"qwerty","role":"customer"}`,
			inputUser: todo.User{
				Name:         "Test",
				Email:        "test",
				PasswordHash: "qwerty",
				Role:         "customer",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user todo.User) {
				s.EXPECT().
					CreateUser(gomock.Any(), user).
					Return(1, errors.New("service failure"))
			},

			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"service failure"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			services := service.Service{Authorization: auth}
			handler := NewHandler(services)

			r := gin.New()
			r.POST("/auth/signup", handler.SignUp)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/signup", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})

	}
}
