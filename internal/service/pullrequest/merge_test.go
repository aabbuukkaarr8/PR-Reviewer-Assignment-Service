package pullrequest

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	prrepo "github.com/aabbuukkaarr8/PRService/internal/repository/pullrequest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_MergePullRequest(t *testing.T) {
	tests := []struct {
		name           string
		pullRequestID  string
		setupMock      func(*mockRepo)
		expectedError  error
		validateResult func(*testing.T, PullRequest)
	}{
		{
			name:          "successful merge of open PR",
			pullRequestID: "pr-001",
			setupMock: func(m *mockRepo) {
				createdAt := time.Now()
				m.On("GetPullRequest", mock.Anything, "pr-001").Return(prrepo.PullRequest{
					PullRequestID:     "pr-001",
					PullRequestName:   "Test PR",
					AuthorID:          "user-001",
					Status:            "OPEN",
					AssignedReviewers: []string{"user-002", "user-003"},
					CreatedAt:         &createdAt,
					MergedAt:          nil,
				}, nil)
				mergedAt := time.Now()
				m.On("MergePullRequest", mock.Anything, "pr-001").Return(prrepo.PullRequest{
					PullRequestID:     "pr-001",
					PullRequestName:   "Test PR",
					AuthorID:          "user-001",
					Status:            "MERGED",
					AssignedReviewers: []string{"user-002", "user-003"},
					CreatedAt:         &createdAt,
					MergedAt:          &mergedAt,
				}, nil)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, pr PullRequest) {
				assert.Equal(t, "pr-001", pr.PullRequestID)
				assert.Equal(t, "MERGED", pr.Status)
				assert.NotNil(t, pr.MergedAt)
				assert.Len(t, pr.AssignedReviewers, 2)
			},
		},
		{
			name:          "idempotent merge - PR already merged",
			pullRequestID: "pr-002",
			setupMock: func(m *mockRepo) {
				createdAt := time.Now()
				mergedAt := time.Now()
				m.On("GetPullRequest", mock.Anything, "pr-002").Return(prrepo.PullRequest{
					PullRequestID:     "pr-002",
					PullRequestName:   "Test PR 2",
					AuthorID:          "user-001",
					Status:            "MERGED",
					AssignedReviewers: []string{"user-002"},
					CreatedAt:         &createdAt,
					MergedAt:          &mergedAt,
				}, nil)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, pr PullRequest) {
				assert.Equal(t, "pr-002", pr.PullRequestID)
				assert.Equal(t, "MERGED", pr.Status)
				assert.NotNil(t, pr.MergedAt)
			},
		},
		{
			name:          "PR not found",
			pullRequestID: "pr-999",
			setupMock: func(m *mockRepo) {
				m.On("GetPullRequest", mock.Anything, "pr-999").Return(prrepo.PullRequest{}, sql.ErrNoRows)
			},
			expectedError: ErrNotFound,
			validateResult: nil,
		},
		{
			name:          "error getting PR",
			pullRequestID: "pr-003",
			setupMock: func(m *mockRepo) {
				m.On("GetPullRequest", mock.Anything, "pr-003").Return(prrepo.PullRequest{}, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
			validateResult: nil,
		},
		{
			name:          "error merging PR",
			pullRequestID: "pr-004",
			setupMock: func(m *mockRepo) {
				createdAt := time.Now()
				m.On("GetPullRequest", mock.Anything, "pr-004").Return(prrepo.PullRequest{
					PullRequestID:     "pr-004",
					PullRequestName:   "Test PR 4",
					AuthorID:          "user-001",
					Status:            "OPEN",
					AssignedReviewers: []string{"user-002"},
					CreatedAt:         &createdAt,
					MergedAt:          nil,
				}, nil)
				m.On("MergePullRequest", mock.Anything, "pr-004").Return(prrepo.PullRequest{}, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
			validateResult: nil,
		},
		{
			name:          "successful merge with empty reviewers",
			pullRequestID: "pr-005",
			setupMock: func(m *mockRepo) {
				createdAt := time.Now()
				m.On("GetPullRequest", mock.Anything, "pr-005").Return(prrepo.PullRequest{
					PullRequestID:     "pr-005",
					PullRequestName:   "Test PR 5",
					AuthorID:          "user-001",
					Status:            "OPEN",
					AssignedReviewers: []string{},
					CreatedAt:         &createdAt,
					MergedAt:          nil,
				}, nil)
				mergedAt := time.Now()
				m.On("MergePullRequest", mock.Anything, "pr-005").Return(prrepo.PullRequest{
					PullRequestID:     "pr-005",
					PullRequestName:   "Test PR 5",
					AuthorID:          "user-001",
					Status:            "MERGED",
					AssignedReviewers: []string{},
					CreatedAt:         &createdAt,
					MergedAt:          &mergedAt,
				}, nil)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, pr PullRequest) {
				assert.Equal(t, "pr-005", pr.PullRequestID)
				assert.Equal(t, "MERGED", pr.Status)
				assert.Empty(t, pr.AssignedReviewers)
				assert.NotNil(t, pr.MergedAt)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockRepo)
			tt.setupMock(mockRepo)

			service := &Service{
				repo: mockRepo,
			}

			ctx := context.Background()
			result, err := service.MergePullRequest(ctx, tt.pullRequestID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				if errors.Is(tt.expectedError, ErrNotFound) {
					assert.ErrorIs(t, err, ErrNotFound)
				} else {
					assert.Contains(t, err.Error(), tt.expectedError.Error())
				}
				assert.Equal(t, PullRequest{}, result)
			} else {
				assert.NoError(t, err)
				if tt.validateResult != nil {
					tt.validateResult(t, result)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

