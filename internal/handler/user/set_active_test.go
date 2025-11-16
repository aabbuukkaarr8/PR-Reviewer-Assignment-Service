package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aabbuukkaarr8/PRService/internal/api/models"
	usersrv "github.com/aabbuukkaarr8/PRService/internal/service/user"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/sirupsen/logrus"
)

func TestHandler_SetIsActive(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*mockService)
		expectedStatus int
		expectedError  string
		validateBody   func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "successful set to active",
			requestBody: SetIsActiveRequest{
				UserID:   "user-001",
				IsActive: boolPtr(true),
			},
			setupMock: func(m *mockService) {
				m.On("SetIsActive", mock.Anything, "user-001", true).Return(usersrv.User{
					UserID:   "user-001",
					Username: "alice",
					TeamName: "backend",
					IsActive: true,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response SetIsActiveResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "user-001", response.User.UserID)
				assert.Equal(t, "alice", response.User.Username)
				assert.True(t, response.User.IsActive)
			},
		},
		{
			name: "successful set to inactive",
			requestBody: SetIsActiveRequest{
				UserID:   "user-002",
				IsActive: boolPtr(false),
			},
			setupMock: func(m *mockService) {
				m.On("SetIsActive", mock.Anything, "user-002", false).Return(usersrv.User{
					UserID:   "user-002",
					Username: "bob",
					TeamName: "frontend",
					IsActive: false,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response SetIsActiveResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "user-002", response.User.UserID)
				assert.False(t, response.User.IsActive)
			},
		},
		{
			name: "invalid JSON - missing user_id",
			requestBody: map[string]interface{}{
				"is_active": true,
			},
			setupMock:      func(m *mockService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_REQUEST",
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, w.Body.String(), "INVALID_REQUEST")
			},
		},
		{
			name: "invalid JSON - missing is_active",
			requestBody: map[string]interface{}{
				"user_id": "user-001",
			},
			setupMock:      func(m *mockService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_REQUEST",
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, w.Body.String(), "INVALID_REQUEST")
			},
		},
		{
			name: "invalid JSON - empty body",
			requestBody: map[string]interface{}{},
			setupMock:      func(m *mockService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_REQUEST",
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, w.Body.String(), "INVALID_REQUEST")
			},
		},
		{
			name: "user not found",
			requestBody: SetIsActiveRequest{
				UserID:   "user-999",
				IsActive: boolPtr(true),
			},
			setupMock: func(m *mockService) {
				m.On("SetIsActive", mock.Anything, "user-999", true).Return(usersrv.User{}, usersrv.ErrUserNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  string(models.NOTFOUND),
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, w.Body.String(), string(models.NOTFOUND))
				assert.Contains(t, w.Body.String(), "user not found")
			},
		},
		{
			name: "internal server error",
			requestBody: SetIsActiveRequest{
				UserID:   "user-001",
				IsActive: boolPtr(true),
			},
			setupMock: func(m *mockService) {
				m.On("SetIsActive", mock.Anything, "user-001", true).Return(usersrv.User{}, errors.New("database connection failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "INTERNAL_ERROR",
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, w.Body.String(), "INTERNAL_ERROR")
				assert.Contains(t, w.Body.String(), "database connection failed")
			},
		},
		{
			name: "invalid JSON format",
			requestBody: "invalid json string",
			setupMock:      func(m *mockService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_REQUEST",
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, w.Body.String(), "INVALID_REQUEST")
			},
		},
		{
			name: "empty user_id",
			requestBody: SetIsActiveRequest{
				UserID:   "",
				IsActive: boolPtr(true),
			},
			setupMock:      func(m *mockService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_REQUEST",
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, w.Body.String(), "INVALID_REQUEST")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(mockService)
			tt.setupMock(mockSvc)

			logger := logrus.New()
			logger.SetLevel(logrus.ErrorLevel)
			handler := &Handler{
				service: mockSvc,
				logger:  logger,
			}

			router := gin.New()
			router.POST("/users/setIsActive", handler.SetIsActive)

			var bodyBytes []byte
			var err error

			if str, ok := tt.requestBody.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req, err := http.NewRequest(http.MethodPost, "/users/setIsActive", bytes.NewBuffer(bodyBytes))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.validateBody != nil {
				tt.validateBody(t, w)
			}

			if tt.expectedError != "" {
				assert.Contains(t, w.Body.String(), tt.expectedError)
			}

			mockSvc.AssertExpectations(t)
		})
	}
}

func boolPtr(b bool) *bool {
	return &b
}

