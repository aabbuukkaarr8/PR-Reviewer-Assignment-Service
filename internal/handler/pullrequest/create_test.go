package pullrequest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aabbuukkaarr8/PRService/internal/api/models"
	prsrv "github.com/aabbuukkaarr8/PRService/internal/service/pullrequest"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockService struct {
	mock.Mock
}

func (m *mockService) CreatePullRequest(ctx context.Context, request prsrv.CreatePullRequest) (prsrv.PullRequest, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return prsrv.PullRequest{}, args.Error(1)
	}
	return args.Get(0).(prsrv.PullRequest), args.Error(1)
}

func (m *mockService) MergePullRequest(ctx context.Context, pullRequestID string) (prsrv.PullRequest, error) {
	args := m.Called(ctx, pullRequestID)
	if args.Get(0) == nil {
		return prsrv.PullRequest{}, args.Error(1)
	}
	return args.Get(0).(prsrv.PullRequest), args.Error(1)
}

func (m *mockService) ReassignReviewer(ctx context.Context, pullRequestID, oldUserID string) (prsrv.PullRequest, string, error) {
	args := m.Called(ctx, pullRequestID, oldUserID)
	if args.Get(0) == nil {
		return prsrv.PullRequest{}, "", args.Error(2)
	}
	return args.Get(0).(prsrv.PullRequest), args.String(1), args.Error(2)
}

func TestHandler_CreatePullRequest(t *testing.T) {
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
			name: "successful creation",
			requestBody: CreateRequest{
				PullRequestID:   "pr-001",
				PullRequestName: "Test PR",
				AuthorID:        "user-001",
			},
			setupMock: func(m *mockService) {
				now := time.Now()
				m.On("CreatePullRequest", mock.Anything, prsrv.CreatePullRequest{
					PullRequestId:   "pr-001",
					PullRequestName: "Test PR",
					AuthorId:        "user-001",
				}).Return(prsrv.PullRequest{
					PullRequestID:     "pr-001",
					PullRequestName: "Test PR",
					AuthorID:       "user-001",
					Status:         "OPEN",
					AssignedReviewers: []string{"user-002", "user-003"},
					CreatedAt:      &now,
					MergedAt:       nil,
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedError:  "",
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response CreatePullRequestResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "pr-001", response.PR.PullRequestID)
				assert.Equal(t, "Test PR", response.PR.PullRequestName)
				assert.Equal(t, "user-001", response.PR.AuthorID)
				assert.Equal(t, "OPEN", response.PR.Status)
				assert.Len(t, response.PR.AssignedReviewers, 2)
			},
		},
		{
			name: "invalid JSON - missing required field",
			requestBody: map[string]interface{}{
				"pull_request_id": "pr-002",
				"author_id":       "user-001",
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
			name: "PR already exists",
			requestBody: CreateRequest{
				PullRequestID:   "pr-003",
				PullRequestName: "Test PR 3",
				AuthorID:        "user-001",
			},
			setupMock: func(m *mockService) {
				m.On("CreatePullRequest", mock.Anything, prsrv.CreatePullRequest{
					PullRequestId:   "pr-003",
					PullRequestName: "Test PR 3",
					AuthorId:        "user-001",
				}).Return(prsrv.PullRequest{}, prsrv.ErrPRExists)
			},
			expectedStatus: http.StatusConflict,
			expectedError:  string(models.PREXISTS),
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, w.Body.String(), string(models.PREXISTS))
				assert.Contains(t, w.Body.String(), "PR id already exists")
			},
		},
		{
			name: "author not found",
			requestBody: CreateRequest{
				PullRequestID:   "pr-004",
				PullRequestName: "Test PR 4",
				AuthorID:        "user-999",
			},
			setupMock: func(m *mockService) {
				m.On("CreatePullRequest", mock.Anything, prsrv.CreatePullRequest{
					PullRequestId:   "pr-004",
					PullRequestName: "Test PR 4",
					AuthorId:        "user-999",
				}).Return(prsrv.PullRequest{}, prsrv.ErrNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  string(models.NOTFOUND),
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, w.Body.String(), string(models.NOTFOUND))
				assert.Contains(t, w.Body.String(), "author or team not found")
			},
		},
		{
			name: "internal server error",
			requestBody: CreateRequest{
				PullRequestID:   "pr-005",
				PullRequestName: "Test PR 5",
				AuthorID:        "user-001",
			},
			setupMock: func(m *mockService) {
				m.On("CreatePullRequest", mock.Anything, prsrv.CreatePullRequest{
					PullRequestId:   "pr-005",
					PullRequestName: "Test PR 5",
					AuthorId:        "user-001",
				}).Return(prsrv.PullRequest{}, errors.New("database connection failed"))
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
			name: "empty pull_request_id",
			requestBody: CreateRequest{
				PullRequestID:   "",
				PullRequestName: "Test PR",
				AuthorID:        "user-001",
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
			name: "empty author_id",
			requestBody: CreateRequest{
				PullRequestID:   "pr-006",
				PullRequestName: "Test PR",
				AuthorID:        "",
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

			handler := &Handler{
				service: mockSvc,
			}

			router := gin.New()
			router.POST("/pullRequest/create", handler.CreatePullRequest)

			var bodyBytes []byte
			var err error

			if str, ok := tt.requestBody.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req, err := http.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewBuffer(bodyBytes))
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

