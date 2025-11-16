package pullrequest

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	prrepo "github.com/aabbuukkaarr8/PRService/internal/repository/pullrequest"
	"github.com/aabbuukkaarr8/PRService/internal/repository/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) PRExists(ctx context.Context, pullRequestID string) (bool, error) {
	args := m.Called(ctx, pullRequestID)
	return args.Bool(0), args.Error(1)
}

func (m *mockRepo) GetUser(ctx context.Context, userID string) (user.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return user.User{}, args.Error(1)
	}
	return args.Get(0).(user.User), args.Error(1)
}

func (m *mockRepo) GetActiveTeamMembers(ctx context.Context, teamName string, excludeUserID string) ([]user.User, error) {
	args := m.Called(ctx, teamName, excludeUserID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]user.User), args.Error(1)
}

func (m *mockRepo) CreatePullRequest(ctx context.Context, request *prrepo.CreatePullRequest) (prrepo.PullRequest, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return prrepo.PullRequest{}, args.Error(1)
	}
	return args.Get(0).(prrepo.PullRequest), args.Error(1)
}

func (m *mockRepo) GetPullRequest(ctx context.Context, pullRequestID string) (prrepo.PullRequest, error) {
	args := m.Called(ctx, pullRequestID)
	if args.Get(0) == nil {
		return prrepo.PullRequest{}, args.Error(1)
	}
	return args.Get(0).(prrepo.PullRequest), args.Error(1)
}

func (m *mockRepo) MergePullRequest(ctx context.Context, pullRequestID string) (prrepo.PullRequest, error) {
	args := m.Called(ctx, pullRequestID)
	if args.Get(0) == nil {
		return prrepo.PullRequest{}, args.Error(1)
	}
	return args.Get(0).(prrepo.PullRequest), args.Error(1)
}

func (m *mockRepo) UpdatePullRequestReviewers(ctx context.Context, pullRequestID string, assignedReviewers []string) (prrepo.PullRequest, error) {
	args := m.Called(ctx, pullRequestID, assignedReviewers)
	if args.Get(0) == nil {
		return prrepo.PullRequest{}, args.Error(1)
	}
	return args.Get(0).(prrepo.PullRequest), args.Error(1)
}

func (m *mockRepo) GetReviewerStats(ctx context.Context) ([]prrepo.ReviewerStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]prrepo.ReviewerStats), args.Error(1)
}

func (m *mockRepo) GetPRStats(ctx context.Context) (prrepo.PRStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return prrepo.PRStats{}, args.Error(1)
	}
	return args.Get(0).(prrepo.PRStats), args.Error(1)
}

func (m *mockRepo) GetOpenPRsByReviewers(ctx context.Context, reviewerIDs []string) ([]prrepo.OpenPRWithReviewer, error) {
	args := m.Called(ctx, reviewerIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]prrepo.OpenPRWithReviewer), args.Error(1)
}

func (m *mockRepo) BulkDeactivateTeamUsers(ctx context.Context, teamName string) ([]string, error) {
	args := m.Called(ctx, teamName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockRepo) BulkUpdatePullRequestReviewers(ctx context.Context, updates []prrepo.PRReviewerUpdate) error {
	args := m.Called(ctx, updates)
	return args.Error(0)
}

func TestService_CreatePullRequest(t *testing.T) {
	tests := []struct {
		name           string
		request        CreatePullRequest
		setupMock      func(*mockRepo)
		expectedError  error
		validateResult func(*testing.T, PullRequest)
	}{
		{
			name: "successful creation with 2 reviewers",
			request: CreatePullRequest{
				PullRequestId:   "pr-001",
				PullRequestName:  "Test PR",
				AuthorId:         "user-001",
			},
			setupMock: func(m *mockRepo) {
				m.On("PRExists", mock.Anything, "pr-001").Return(false, nil)
				m.On("GetUser", mock.Anything, "user-001").Return(user.User{
					UserID:   "user-001",
					Username: "alice",
					TeamName: "backend",
					IsActive: true,
				}, nil)
				m.On("GetActiveTeamMembers", mock.Anything, "backend", "user-001").Return([]user.User{
					{UserID: "user-002", Username: "bob", TeamName: "backend", IsActive: true},
					{UserID: "user-003", Username: "charlie", TeamName: "backend", IsActive: true},
					{UserID: "user-004", Username: "david", TeamName: "backend", IsActive: true},
				}, nil)
				m.On("CreatePullRequest", mock.Anything, mock.AnythingOfType("*pullrequest.CreatePullRequest")).Return(
					prrepo.PullRequest{
						PullRequestID:     "pr-001",
						PullRequestName:   "Test PR",
						AuthorID:          "user-001",
						Status:            "OPEN",
						AssignedReviewers: []string{"user-002", "user-003"},
					}, nil)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, pr PullRequest) {
				assert.Equal(t, "pr-001", pr.PullRequestID)
				assert.Equal(t, "Test PR", pr.PullRequestName)
				assert.Equal(t, "user-001", pr.AuthorID)
				assert.Equal(t, "OPEN", pr.Status)
				assert.Len(t, pr.AssignedReviewers, 2)
				assert.Contains(t, pr.AssignedReviewers, "user-002")
				assert.Contains(t, pr.AssignedReviewers, "user-003")
			},
		},
		{
			name: "successful creation with 1 reviewer (only 1 available)",
			request: CreatePullRequest{
				PullRequestId:   "pr-002",
				PullRequestName: "Test PR 2",
				AuthorId:         "user-001",
			},
			setupMock: func(m *mockRepo) {
				m.On("PRExists", mock.Anything, "pr-002").Return(false, nil)
				m.On("GetUser", mock.Anything, "user-001").Return(user.User{
					UserID:   "user-001",
					Username: "alice",
					TeamName: "backend",
					IsActive: true,
				}, nil)
				m.On("GetActiveTeamMembers", mock.Anything, "backend", "user-001").Return([]user.User{
					{UserID: "user-002", Username: "bob", TeamName: "backend", IsActive: true},
				}, nil)
				m.On("CreatePullRequest", mock.Anything, mock.AnythingOfType("*pullrequest.CreatePullRequest")).Return(
					prrepo.PullRequest{
						PullRequestID:     "pr-002",
						PullRequestName:   "Test PR 2",
						AuthorID:          "user-001",
						Status:            "OPEN",
						AssignedReviewers: []string{"user-002"},
					}, nil)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, pr PullRequest) {
				assert.Equal(t, "pr-002", pr.PullRequestID)
				assert.Len(t, pr.AssignedReviewers, 1)
				assert.Equal(t, "user-002", pr.AssignedReviewers[0])
			},
		},
		{
			name: "successful creation with no reviewers (no team members)",
			request: CreatePullRequest{
				PullRequestId:   "pr-003",
				PullRequestName: "Test PR 3",
				AuthorId:         "user-001",
			},
			setupMock: func(m *mockRepo) {
				m.On("PRExists", mock.Anything, "pr-003").Return(false, nil)
				m.On("GetUser", mock.Anything, "user-001").Return(user.User{
					UserID:   "user-001",
					Username: "alice",
					TeamName: "backend",
					IsActive: true,
				}, nil)
				m.On("GetActiveTeamMembers", mock.Anything, "backend", "user-001").Return([]user.User{}, nil)
				m.On("CreatePullRequest", mock.Anything, mock.AnythingOfType("*pullrequest.CreatePullRequest")).Return(
					prrepo.PullRequest{
						PullRequestID:     "pr-003",
						PullRequestName:   "Test PR 3",
						AuthorID:          "user-001",
						Status:            "OPEN",
						AssignedReviewers: []string{},
					}, nil)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, pr PullRequest) {
				assert.Equal(t, "pr-003", pr.PullRequestID)
				assert.Empty(t, pr.AssignedReviewers)
			},
		},
		{
			name: "PR already exists",
			request: CreatePullRequest{
				PullRequestId:   "pr-004",
				PullRequestName: "Test PR 4",
				AuthorId:         "user-001",
			},
			setupMock: func(m *mockRepo) {
				m.On("PRExists", mock.Anything, "pr-004").Return(true, nil)
			},
			expectedError: ErrPRExists,
			validateResult: nil,
		},
		{
			name: "author not found",
			request: CreatePullRequest{
				PullRequestId:   "pr-005",
				PullRequestName: "Test PR 5",
				AuthorId:         "user-999",
			},
			setupMock: func(m *mockRepo) {
				m.On("PRExists", mock.Anything, "pr-005").Return(false, nil)
				m.On("GetUser", mock.Anything, "user-999").Return(user.User{}, sql.ErrNoRows)
			},
			expectedError: ErrNotFound,
			validateResult: nil,
		},
		{
			name: "error checking PR existence",
			request: CreatePullRequest{
				PullRequestId:   "pr-006",
				PullRequestName: "Test PR 6",
				AuthorId:         "user-001",
			},
			setupMock: func(m *mockRepo) {
				m.On("PRExists", mock.Anything, "pr-006").Return(false, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
			validateResult: nil,
		},
		{
			name: "error getting user",
			request: CreatePullRequest{
				PullRequestId:   "pr-007",
				PullRequestName: "Test PR 7",
				AuthorId:         "user-001",
			},
			setupMock: func(m *mockRepo) {
				m.On("PRExists", mock.Anything, "pr-007").Return(false, nil)
				m.On("GetUser", mock.Anything, "user-001").Return(user.User{}, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
			validateResult: nil,
		},
		{
			name: "error getting team members",
			request: CreatePullRequest{
				PullRequestId:   "pr-008",
				PullRequestName: "Test PR 8",
				AuthorId:         "user-001",
			},
			setupMock: func(m *mockRepo) {
				m.On("PRExists", mock.Anything, "pr-008").Return(false, nil)
				m.On("GetUser", mock.Anything, "user-001").Return(user.User{
					UserID:   "user-001",
					Username: "alice",
					TeamName: "backend",
					IsActive: true,
				}, nil)
				m.On("GetActiveTeamMembers", mock.Anything, "backend", "user-001").Return(nil, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
			validateResult: nil,
		},
		{
			name: "error creating PR in database",
			request: CreatePullRequest{
				PullRequestId:   "pr-009",
				PullRequestName: "Test PR 9",
				AuthorId:         "user-001",
			},
			setupMock: func(m *mockRepo) {
				m.On("PRExists", mock.Anything, "pr-009").Return(false, nil)
				m.On("GetUser", mock.Anything, "user-001").Return(user.User{
					UserID:   "user-001",
					Username: "alice",
					TeamName: "backend",
					IsActive: true,
				}, nil)
				m.On("GetActiveTeamMembers", mock.Anything, "backend", "user-001").Return([]user.User{
					{UserID: "user-002", Username: "bob", TeamName: "backend", IsActive: true},
				}, nil)
				m.On("CreatePullRequest", mock.Anything, mock.AnythingOfType("*pullrequest.CreatePullRequest")).Return(nil, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
			validateResult: nil,
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
			result, err := service.CreatePullRequest(ctx, tt.request)

			if tt.expectedError != nil {
				assert.Error(t, err)
				if errors.Is(tt.expectedError, ErrPRExists) {
					assert.ErrorIs(t, err, ErrPRExists)
				} else if errors.Is(tt.expectedError, ErrNotFound) {
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

