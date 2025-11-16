package pullrequest

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	prrepo "github.com/aabbuukkaarr8/PRService/internal/repository/pullrequest"
	"github.com/aabbuukkaarr8/PRService/internal/repository/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_ReassignReviewer(t *testing.T) {
	tests := []struct {
		name           string
		pullRequestID  string
		oldUserID      string
		setupMock      func(*mockRepo)
		expectedError  error
		validateResult func(*testing.T, PullRequest, string)
	}{
		{
			name:          "successful reassignment",
			pullRequestID: "pr-001",
			oldUserID:     "user-002",
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
				m.On("GetUser", mock.Anything, "user-002").Return(user.User{
					UserID:   "user-002",
					Username: "bob",
					TeamName: "backend",
					IsActive: true,
				}, nil)
				m.On("GetActiveTeamMembers", mock.Anything, "backend", "user-002").Return([]user.User{
					{UserID: "user-004", Username: "david", TeamName: "backend", IsActive: true},
					{UserID: "user-005", Username: "eve", TeamName: "backend", IsActive: true},
				}, nil)
				m.On("UpdatePullRequestReviewers", mock.Anything, "pr-001", mock.AnythingOfType("[]string")).Return(prrepo.PullRequest{
					PullRequestID:     "pr-001",
					PullRequestName:   "Test PR",
					AuthorID:          "user-001",
					Status:            "OPEN",
					AssignedReviewers: []string{"user-004", "user-003"},
					CreatedAt:         &createdAt,
					MergedAt:          nil,
				}, nil)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, pr PullRequest, replacedBy string) {
				assert.Equal(t, "pr-001", pr.PullRequestID)
				assert.Equal(t, "OPEN", pr.Status)
				assert.Len(t, pr.AssignedReviewers, 2)
				assert.Contains(t, pr.AssignedReviewers, "user-003")
				assert.NotEmpty(t, replacedBy)
				assert.Contains(t, []string{"user-004", "user-005"}, replacedBy)
				assert.Contains(t, pr.AssignedReviewers, replacedBy)
			},
		},
		{
			name:          "PR not found",
			pullRequestID: "pr-999",
			oldUserID:     "user-002",
			setupMock: func(m *mockRepo) {
				m.On("GetPullRequest", mock.Anything, "pr-999").Return(prrepo.PullRequest{}, sql.ErrNoRows)
			},
			expectedError: ErrNotFound,
			validateResult: nil,
		},
		{
			name:          "PR is merged",
			pullRequestID: "pr-002",
			oldUserID:     "user-002",
			setupMock: func(m *mockRepo) {
				createdAt := time.Now()
				mergedAt := time.Now()
				m.On("GetPullRequest", mock.Anything, "pr-002").Return(prrepo.PullRequest{
					PullRequestID:     "pr-002",
					PullRequestName:   "Test PR 2",
					AuthorID:          "user-001",
					Status:            "MERGED",
					AssignedReviewers: []string{"user-002", "user-003"},
					CreatedAt:         &createdAt,
					MergedAt:          &mergedAt,
				}, nil)
			},
			expectedError: ErrPRMerged,
			validateResult: nil,
		},
		{
			name:          "reviewer not assigned",
			pullRequestID: "pr-003",
			oldUserID:     "user-999",
			setupMock: func(m *mockRepo) {
				createdAt := time.Now()
				m.On("GetPullRequest", mock.Anything, "pr-003").Return(prrepo.PullRequest{
					PullRequestID:     "pr-003",
					PullRequestName:   "Test PR 3",
					AuthorID:          "user-001",
					Status:            "OPEN",
					AssignedReviewers: []string{"user-002", "user-003"},
					CreatedAt:         &createdAt,
					MergedAt:          nil,
				}, nil)
			},
			expectedError: ErrNotAssigned,
			validateResult: nil,
		},
		{
			name:          "old reviewer not found",
			pullRequestID: "pr-004",
			oldUserID:     "user-002",
			setupMock: func(m *mockRepo) {
				createdAt := time.Now()
				m.On("GetPullRequest", mock.Anything, "pr-004").Return(prrepo.PullRequest{
					PullRequestID:     "pr-004",
					PullRequestName:   "Test PR 4",
					AuthorID:          "user-001",
					Status:            "OPEN",
					AssignedReviewers: []string{"user-002", "user-003"},
					CreatedAt:         &createdAt,
					MergedAt:          nil,
				}, nil)
				m.On("GetUser", mock.Anything, "user-002").Return(user.User{}, sql.ErrNoRows)
			},
			expectedError: ErrNotFound,
			validateResult: nil,
		},
		{
			name:          "no replacement candidates",
			pullRequestID: "pr-005",
			oldUserID:     "user-002",
			setupMock: func(m *mockRepo) {
				createdAt := time.Now()
				m.On("GetPullRequest", mock.Anything, "pr-005").Return(prrepo.PullRequest{
					PullRequestID:     "pr-005",
					PullRequestName:   "Test PR 5",
					AuthorID:          "user-001",
					Status:            "OPEN",
					AssignedReviewers: []string{"user-002", "user-003"},
					CreatedAt:         &createdAt,
					MergedAt:          nil,
				}, nil)
				m.On("GetUser", mock.Anything, "user-002").Return(user.User{
					UserID:   "user-002",
					Username: "bob",
					TeamName: "backend",
					IsActive: true,
				}, nil)
				m.On("GetActiveTeamMembers", mock.Anything, "backend", "user-002").Return([]user.User{}, nil)
			},
			expectedError: ErrNoCandidate,
			validateResult: nil,
		},
		{
			name:          "no candidates after excluding author and other reviewers",
			pullRequestID: "pr-006",
			oldUserID:     "user-002",
			setupMock: func(m *mockRepo) {
				createdAt := time.Now()
				m.On("GetPullRequest", mock.Anything, "pr-006").Return(prrepo.PullRequest{
					PullRequestID:     "pr-006",
					PullRequestName:   "Test PR 6",
					AuthorID:          "user-001",
					Status:            "OPEN",
					AssignedReviewers: []string{"user-002", "user-003"},
					CreatedAt:         &createdAt,
					MergedAt:          nil,
				}, nil)
				m.On("GetUser", mock.Anything, "user-002").Return(user.User{
					UserID:   "user-002",
					Username: "bob",
					TeamName: "backend",
					IsActive: true,
				}, nil)
				m.On("GetActiveTeamMembers", mock.Anything, "backend", "user-002").Return([]user.User{
					{UserID: "user-001", Username: "alice", TeamName: "backend", IsActive: true},
					{UserID: "user-003", Username: "charlie", TeamName: "backend", IsActive: true},
				}, nil)
			},
			expectedError: ErrNoCandidate,
			validateResult: nil,
		},
		{
			name:          "error getting PR",
			pullRequestID: "pr-007",
			oldUserID:     "user-002",
			setupMock: func(m *mockRepo) {
				m.On("GetPullRequest", mock.Anything, "pr-007").Return(prrepo.PullRequest{}, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
			validateResult: nil,
		},
		{
			name:          "error getting old reviewer",
			pullRequestID: "pr-008",
			oldUserID:     "user-002",
			setupMock: func(m *mockRepo) {
				createdAt := time.Now()
				m.On("GetPullRequest", mock.Anything, "pr-008").Return(prrepo.PullRequest{
					PullRequestID:     "pr-008",
					PullRequestName:   "Test PR 8",
					AuthorID:          "user-001",
					Status:            "OPEN",
					AssignedReviewers: []string{"user-002", "user-003"},
					CreatedAt:         &createdAt,
					MergedAt:          nil,
				}, nil)
				m.On("GetUser", mock.Anything, "user-002").Return(user.User{}, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
			validateResult: nil,
		},
		{
			name:          "error getting team members",
			pullRequestID: "pr-009",
			oldUserID:     "user-002",
			setupMock: func(m *mockRepo) {
				createdAt := time.Now()
				m.On("GetPullRequest", mock.Anything, "pr-009").Return(prrepo.PullRequest{
					PullRequestID:     "pr-009",
					PullRequestName:   "Test PR 9",
					AuthorID:          "user-001",
					Status:            "OPEN",
					AssignedReviewers: []string{"user-002", "user-003"},
					CreatedAt:         &createdAt,
					MergedAt:          nil,
				}, nil)
				m.On("GetUser", mock.Anything, "user-002").Return(user.User{
					UserID:   "user-002",
					Username: "bob",
					TeamName: "backend",
					IsActive: true,
				}, nil)
				m.On("GetActiveTeamMembers", mock.Anything, "backend", "user-002").Return(nil, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
			validateResult: nil,
		},
		{
			name:          "error updating reviewers",
			pullRequestID: "pr-010",
			oldUserID:     "user-002",
			setupMock: func(m *mockRepo) {
				createdAt := time.Now()
				m.On("GetPullRequest", mock.Anything, "pr-010").Return(prrepo.PullRequest{
					PullRequestID:     "pr-010",
					PullRequestName:   "Test PR 10",
					AuthorID:          "user-001",
					Status:            "OPEN",
					AssignedReviewers: []string{"user-002", "user-003"},
					CreatedAt:         &createdAt,
					MergedAt:          nil,
				}, nil)
				m.On("GetUser", mock.Anything, "user-002").Return(user.User{
					UserID:   "user-002",
					Username: "bob",
					TeamName: "backend",
					IsActive: true,
				}, nil)
				m.On("GetActiveTeamMembers", mock.Anything, "backend", "user-002").Return([]user.User{
					{UserID: "user-004", Username: "david", TeamName: "backend", IsActive: true},
				}, nil)
				m.On("UpdatePullRequestReviewers", mock.Anything, "pr-010", mock.Anything).Return(prrepo.PullRequest{}, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
			validateResult: nil,
		},
		{
			name:          "successful reassignment with single reviewer",
			pullRequestID: "pr-011",
			oldUserID:     "user-002",
			setupMock: func(m *mockRepo) {
				createdAt := time.Now()
				m.On("GetPullRequest", mock.Anything, "pr-011").Return(prrepo.PullRequest{
					PullRequestID:     "pr-011",
					PullRequestName:   "Test PR 11",
					AuthorID:          "user-001",
					Status:            "OPEN",
					AssignedReviewers: []string{"user-002"},
					CreatedAt:         &createdAt,
					MergedAt:          nil,
				}, nil)
				m.On("GetUser", mock.Anything, "user-002").Return(user.User{
					UserID:   "user-002",
					Username: "bob",
					TeamName: "backend",
					IsActive: true,
				}, nil)
				m.On("GetActiveTeamMembers", mock.Anything, "backend", "user-002").Return([]user.User{
					{UserID: "user-004", Username: "david", TeamName: "backend", IsActive: true},
				}, nil)
				m.On("UpdatePullRequestReviewers", mock.Anything, "pr-011", []string{"user-004"}).Return(prrepo.PullRequest{
					PullRequestID:     "pr-011",
					PullRequestName:   "Test PR 11",
					AuthorID:          "user-001",
					Status:            "OPEN",
					AssignedReviewers: []string{"user-004"},
					CreatedAt:         &createdAt,
					MergedAt:          nil,
				}, nil)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, pr PullRequest, replacedBy string) {
				assert.Equal(t, "pr-011", pr.PullRequestID)
				assert.Len(t, pr.AssignedReviewers, 1)
				assert.Equal(t, "user-004", pr.AssignedReviewers[0])
				assert.Equal(t, "user-004", replacedBy)
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
			result, replacedBy, err := service.ReassignReviewer(ctx, tt.pullRequestID, tt.oldUserID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				if errors.Is(tt.expectedError, ErrNotFound) {
					assert.ErrorIs(t, err, ErrNotFound)
				} else if errors.Is(tt.expectedError, ErrPRMerged) {
					assert.ErrorIs(t, err, ErrPRMerged)
				} else if errors.Is(tt.expectedError, ErrNotAssigned) {
					assert.ErrorIs(t, err, ErrNotAssigned)
				} else if errors.Is(tt.expectedError, ErrNoCandidate) {
					assert.ErrorIs(t, err, ErrNoCandidate)
				} else {
					assert.Contains(t, err.Error(), tt.expectedError.Error())
				}
				assert.Equal(t, PullRequest{}, result)
				assert.Empty(t, replacedBy)
			} else {
				assert.NoError(t, err)
				if tt.validateResult != nil {
					tt.validateResult(t, result, replacedBy)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

