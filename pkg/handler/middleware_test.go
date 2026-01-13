package handler

import (
	"errors"
	"net/http/httptest"
	"prac/pkg/service"
	mock_service "prac/pkg/service/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestHandler_userIdentity(t *testing.T) {
	type MockBehavior func(s *mock_service.MockAuthorization, tocken string)

	testTable := []struct {
		name                 string
		header               string
		headerValue          string
		tocken               string
		mockBehavior         MockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			header:      "Authorization",
			headerValue: "Bearer tocken",
			tocken:      "tocken",
			mockBehavior: func(s *mock_service.MockAuthorization, tocken string) {
				s.EXPECT().
					ParseToken(gomock.Any(), tocken). //tocken
					Return(uint(1), "customer", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "1",
		},
		{
			name:                 "No Header",
			header:               "",
			mockBehavior:         func(s *mock_service.MockAuthorization, tocken string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"empty auth header"}`,
		},
		{
			name:                 "Invalid Header",
			header:               "Authorization",
			headerValue:          "Barer tocken",
			mockBehavior:         func(s *mock_service.MockAuthorization, tocken string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"invalid auth header"}`,
		},
		{
			name:        "Invalid Token",
			header:      "Authorization",
			headerValue: "Bearer tocken",
			tocken:      "tocken",
			mockBehavior: func(s *mock_service.MockAuthorization, tocken string) {
				s.EXPECT().
					ParseToken(gomock.Any(), tocken).
					Return(uint(0), "", errors.New("failed to parse token"))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"failed to parse token"}`,
		},
		{
			name:        "Ok",
			header:      "Authorization",
			headerValue: "Bearer tocken",
			tocken:      "tocken",
			mockBehavior: func(s *mock_service.MockAuthorization, tocken string) {
				s.EXPECT().
					ParseToken(gomock.Any(), tocken).
					Return(uint(1), "customer", errors.New("failed to parse tocken"))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"failed to parse tocken"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.tocken)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(*services)

			r := gin.New()
			r.POST("/protected", handler.userIdentity, func(c *gin.Context) {
				userId, _ := c.Get("userId")
				c.String(200, "%v", userId)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/protected", nil)
			req.Header.Set(testCase.header, testCase.headerValue)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}
