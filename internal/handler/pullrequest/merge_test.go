package pullrequest

import (
	"bytes"
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
	"github.com/sirupsen/logrus"
)

func TestHandler_MergePullRequest(t *testing.T) {
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
			name: "successful merge",
			requestBody: MergeRequest{
				PullRequestID: "pr-001",
			},
			setupMock: func(m *mockService) {
				now := time.Now()
				m.On("MergePullRequest", mock.Anything, "pr-001").Return(prsrv.PullRequest{
					PullRequestID:     "pr-001",
					PullRequestName:   "Test PR",
					AuthorID:          "user-001",
					Status:            "MERGED",
					AssignedReviewers: []string{"user-002", "user-003"},
					CreatedAt:         &now,
					MergedAt:          &now,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response MergePullRequestResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "pr-001", response.PR.PullRequestID)
				assert.Equal(t, "MERGED", response.PR.Status)
				assert.NotNil(t, response.PR.MergedAt)
				assert.Len(t, response.PR.AssignedReviewers, 2)
			},
		},
		{
			name: "invalid JSON - missing required field",
			requestBody: map[string]interface{}{
				"some_field": "value",
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
			name: "PR not found",
			requestBody: MergeRequest{
				PullRequestID: "pr-999",
			},
			setupMock: func(m *mockService) {
				m.On("MergePullRequest", mock.Anything, "pr-999").Return(prsrv.PullRequest{}, prsrv.ErrNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  string(models.NOTFOUND),
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, w.Body.String(), string(models.NOTFOUND))
				assert.Contains(t, w.Body.String(), "PR not found")
			},
		},
		{
			name: "internal server error",
			requestBody: MergeRequest{
				PullRequestID: "pr-002",
			},
			setupMock: func(m *mockService) {
				m.On("MergePullRequest", mock.Anything, "pr-002").Return(prsrv.PullRequest{}, errors.New("database connection failed"))
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
			requestBody: MergeRequest{
				PullRequestID: "",
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
			name: "idempotent merge - PR already merged",
			requestBody: MergeRequest{
				PullRequestID: "pr-003",
			},
			setupMock: func(m *mockService) {
				now := time.Now()
				m.On("MergePullRequest", mock.Anything, "pr-003").Return(prsrv.PullRequest{
					PullRequestID:     "pr-003",
					PullRequestName:   "Test PR 3",
					AuthorID:          "user-001",
					Status:            "MERGED",
					AssignedReviewers: []string{"user-002"},
					CreatedAt:         &now,
					MergedAt:          &now,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response MergePullRequestResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "pr-003", response.PR.PullRequestID)
				assert.Equal(t, "MERGED", response.PR.Status)
				assert.NotNil(t, response.PR.MergedAt)
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
			router.POST("/pullRequest/merge", handler.MergePullRequest)

			var bodyBytes []byte
			var err error

			if str, ok := tt.requestBody.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req, err := http.NewRequest(http.MethodPost, "/pullRequest/merge", bytes.NewBuffer(bodyBytes))
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

