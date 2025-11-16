package user

import (
	"context"
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

type mockService struct {
	mock.Mock
}

func (m *mockService) SetIsActive(ctx context.Context, userID string, isActive bool) (usersrv.User, error) {
	args := m.Called(ctx, userID, isActive)
	if args.Get(0) == nil {
		return usersrv.User{}, args.Error(1)
	}
	return args.Get(0).(usersrv.User), args.Error(1)
}

func (m *mockService) GetReview(ctx context.Context, userID string) ([]usersrv.PullRequestShort, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]usersrv.PullRequestShort), args.Error(1)
}

func TestHandler_GetReview(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		queryParams     string
		setupMock      func(*mockService)
		expectedStatus int
		expectedError  string
		validateBody   func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:        "successful get with multiple PRs",
			queryParams: "user_id=user-001",
			setupMock: func(m *mockService) {
				m.On("GetReview", mock.Anything, "user-001").Return([]usersrv.PullRequestShort{
					{PullRequestID: "pr-001", PullRequestName: "Test PR 1", AuthorID: "user-002", Status: "OPEN"},
					{PullRequestID: "pr-002", PullRequestName: "Test PR 2", AuthorID: "user-003", Status: "MERGED"},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response GetReviewResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "user-001", response.UserID)
				assert.Len(t, response.PullRequests, 2)
				assert.Equal(t, "pr-001", response.PullRequests[0].PullRequestID)
				assert.Equal(t, "OPEN", response.PullRequests[0].Status)
			},
		},
		{
			name:        "successful get with single PR",
			queryParams: "user_id=user-002",
			setupMock: func(m *mockService) {
				m.On("GetReview", mock.Anything, "user-002").Return([]usersrv.PullRequestShort{
					{PullRequestID: "pr-003", PullRequestName: "Test PR 3", AuthorID: "user-001", Status: "OPEN"},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response GetReviewResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "user-002", response.UserID)
				assert.Len(t, response.PullRequests, 1)
			},
		},
		{
			name:        "successful get with no PRs",
			queryParams: "user_id=user-003",
			setupMock: func(m *mockService) {
				m.On("GetReview", mock.Anything, "user-003").Return([]usersrv.PullRequestShort{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response GetReviewResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "user-003", response.UserID)
				assert.Empty(t, response.PullRequests)
			},
		},
		{
			name:        "missing user_id parameter",
			queryParams: "",
			setupMock:   func(m *mockService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_REQUEST",
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, w.Body.String(), "INVALID_REQUEST")
				assert.Contains(t, w.Body.String(), "user_id query parameter is required")
			},
		},
		{
			name:        "empty user_id parameter",
			queryParams: "user_id=",
			setupMock:   func(m *mockService) {},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "INVALID_REQUEST",
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, w.Body.String(), "INVALID_REQUEST")
				assert.Contains(t, w.Body.String(), "user_id query parameter is required")
			},
		},
		{
			name:        "user not found",
			queryParams: "user_id=user-999",
			setupMock: func(m *mockService) {
				m.On("GetReview", mock.Anything, "user-999").Return(nil, usersrv.ErrUserNotFound)
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
			name:        "internal server error",
			queryParams: "user_id=user-001",
			setupMock: func(m *mockService) {
				m.On("GetReview", mock.Anything, "user-001").Return(nil, errors.New("database connection failed"))
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
			router.GET("/users/getReview", handler.GetReview)

			url := "/users/getReview"
			if tt.queryParams != "" {
				url += "?" + tt.queryParams
			}

			req, err := http.NewRequest(http.MethodGet, url, nil)
			assert.NoError(t, err)

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

