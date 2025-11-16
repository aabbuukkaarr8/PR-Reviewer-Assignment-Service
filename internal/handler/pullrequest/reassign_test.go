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

func TestHandler_ReassignReviewer(t *testing.T) {
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
			name: "successful reassignment",
			requestBody: ReassignReviewerRequest{
				PullRequestID: "pr-001",
				OldUserID:     "user-002",
			},
			setupMock: func(m *mockService) {
				now := time.Now()
				m.On("ReassignReviewer", mock.Anything, "pr-001", "user-002").Return(prsrv.PullRequest{
					PullRequestID:     "pr-001",
					PullRequestName:   "Test PR",
					AuthorID:          "user-001",
					Status:            "OPEN",
					AssignedReviewers: []string{"user-004", "user-003"},
					CreatedAt:         &now,
					MergedAt:          nil,
				}, "user-004", nil)
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response ReassignReviewerResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "pr-001", response.PR.PullRequestID)
				assert.Equal(t, "user-004", response.ReplacedBy)
				assert.Len(t, response.PR.AssignedReviewers, 2)
			},
		},
		{
			name: "invalid JSON - missing required field",
			requestBody: map[string]interface{}{
				"pull_request_id": "pr-002",
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
			name: "PR or user not found",
			requestBody: ReassignReviewerRequest{
				PullRequestID: "pr-999",
				OldUserID:     "user-999",
			},
			setupMock: func(m *mockService) {
				m.On("ReassignReviewer", mock.Anything, "pr-999", "user-999").Return(prsrv.PullRequest{}, "", prsrv.ErrNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  string(models.NOTFOUND),
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, w.Body.String(), string(models.NOTFOUND))
				assert.Contains(t, w.Body.String(), "PR or user not found")
			},
		},
		{
			name: "PR is merged",
			requestBody: ReassignReviewerRequest{
				PullRequestID: "pr-003",
				OldUserID:     "user-002",
			},
			setupMock: func(m *mockService) {
				m.On("ReassignReviewer", mock.Anything, "pr-003", "user-002").Return(prsrv.PullRequest{}, "", prsrv.ErrPRMerged)
			},
			expectedStatus: http.StatusConflict,
			expectedError:  string(models.PRMERGED),
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, w.Body.String(), string(models.PRMERGED))
				assert.Contains(t, w.Body.String(), "cannot reassign on merged PR")
			},
		},
		{
			name: "reviewer not assigned",
			requestBody: ReassignReviewerRequest{
				PullRequestID: "pr-004",
				OldUserID:     "user-999",
			},
			setupMock: func(m *mockService) {
				m.On("ReassignReviewer", mock.Anything, "pr-004", "user-999").Return(prsrv.PullRequest{}, "", prsrv.ErrNotAssigned)
			},
			expectedStatus: http.StatusConflict,
			expectedError:  string(models.NOTASSIGNED),
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, w.Body.String(), string(models.NOTASSIGNED))
				assert.Contains(t, w.Body.String(), "reviewer is not assigned to this PR")
			},
		},
		{
			name: "no replacement candidate",
			requestBody: ReassignReviewerRequest{
				PullRequestID: "pr-005",
				OldUserID:     "user-002",
			},
			setupMock: func(m *mockService) {
				m.On("ReassignReviewer", mock.Anything, "pr-005", "user-002").Return(prsrv.PullRequest{}, "", prsrv.ErrNoCandidate)
			},
			expectedStatus: http.StatusConflict,
			expectedError:  string(models.NOCANDIDATE),
			validateBody: func(t *testing.T, w *httptest.ResponseRecorder) {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, w.Body.String(), string(models.NOCANDIDATE))
				assert.Contains(t, w.Body.String(), "no active replacement candidate in team")
			},
		},
		{
			name: "internal server error",
			requestBody: ReassignReviewerRequest{
				PullRequestID: "pr-006",
				OldUserID:     "user-002",
			},
			setupMock: func(m *mockService) {
				m.On("ReassignReviewer", mock.Anything, "pr-006", "user-002").Return(prsrv.PullRequest{}, "", errors.New("database connection failed"))
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
			requestBody: ReassignReviewerRequest{
				PullRequestID: "",
				OldUserID:     "user-002",
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
			name: "empty old_user_id",
			requestBody: ReassignReviewerRequest{
				PullRequestID: "pr-007",
				OldUserID:     "",
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
			router.POST("/pullRequest/reassign", handler.ReassignReviewer)

			var bodyBytes []byte
			var err error

			if str, ok := tt.requestBody.(string); ok {
				bodyBytes = []byte(str)
			} else {
				bodyBytes, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req, err := http.NewRequest(http.MethodPost, "/pullRequest/reassign", bytes.NewBuffer(bodyBytes))
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

