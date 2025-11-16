package user

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/aabbuukkaarr8/PRService/internal/repository/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepoForUser struct {
	mock.Mock
}

func (m *mockRepoForUser) GetUser(ctx context.Context, userID string) (user.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return user.User{}, args.Error(1)
	}
	return args.Get(0).(user.User), args.Error(1)
}

func (m *mockRepoForUser) UpdateUserIsActive(ctx context.Context, userID string, isActive bool) error {
	args := m.Called(ctx, userID, isActive)
	return args.Error(0)
}

func (m *mockRepoForUser) GetUserPullRequests(ctx context.Context, userID string) ([]user.PullRequestShort, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]user.PullRequestShort), args.Error(1)
}

func TestService_GetReview(t *testing.T) {
	tests := []struct {
		name           string
		userID          string
		setupMock      func(*mockRepoForUser)
		expectedError  error
		validateResult func(*testing.T, []PullRequestShort)
	}{
		{
			name:  "successful get with multiple PRs",
			userID: "user-001",
			setupMock: func(m *mockRepoForUser) {
				m.On("GetUser", mock.Anything, "user-001").Return(user.User{
					UserID:   "user-001",
					Username: "alice",
					TeamName: "backend",
					IsActive: true,
				}, nil)
				m.On("GetUserPullRequests", mock.Anything, "user-001").Return([]user.PullRequestShort{
					{PullRequestID: "pr-001", PullRequestName: "Test PR 1", AuthorID: "user-002", Status: "OPEN"},
					{PullRequestID: "pr-002", PullRequestName: "Test PR 2", AuthorID: "user-003", Status: "MERGED"},
				}, nil)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, prs []PullRequestShort) {
				assert.Len(t, prs, 2)
				assert.Equal(t, "pr-001", prs[0].PullRequestID)
				assert.Equal(t, "Test PR 1", prs[0].PullRequestName)
				assert.Equal(t, "OPEN", prs[0].Status)
				assert.Equal(t, "pr-002", prs[1].PullRequestID)
				assert.Equal(t, "MERGED", prs[1].Status)
			},
		},
		{
			name:  "successful get with single PR",
			userID: "user-002",
			setupMock: func(m *mockRepoForUser) {
				m.On("GetUser", mock.Anything, "user-002").Return(user.User{
					UserID:   "user-002",
					Username: "bob",
					TeamName: "frontend",
					IsActive: true,
				}, nil)
				m.On("GetUserPullRequests", mock.Anything, "user-002").Return([]user.PullRequestShort{
					{PullRequestID: "pr-003", PullRequestName: "Test PR 3", AuthorID: "user-001", Status: "OPEN"},
				}, nil)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, prs []PullRequestShort) {
				assert.Len(t, prs, 1)
				assert.Equal(t, "pr-003", prs[0].PullRequestID)
			},
		},
		{
			name:  "successful get with no PRs",
			userID: "user-003",
			setupMock: func(m *mockRepoForUser) {
				m.On("GetUser", mock.Anything, "user-003").Return(user.User{
					UserID:   "user-003",
					Username: "charlie",
					TeamName: "devops",
					IsActive: true,
				}, nil)
				m.On("GetUserPullRequests", mock.Anything, "user-003").Return([]user.PullRequestShort{}, nil)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, prs []PullRequestShort) {
				assert.Empty(t, prs)
			},
		},
		{
			name:  "user not found",
			userID: "user-999",
			setupMock: func(m *mockRepoForUser) {
				m.On("GetUser", mock.Anything, "user-999").Return(user.User{}, sql.ErrNoRows)
			},
			expectedError: ErrUserNotFound,
			validateResult: nil,
		},
		{
			name:  "error getting user",
			userID: "user-001",
			setupMock: func(m *mockRepoForUser) {
				m.On("GetUser", mock.Anything, "user-001").Return(user.User{}, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
			validateResult: nil,
		},
		{
			name:  "error getting user PRs",
			userID: "user-001",
			setupMock: func(m *mockRepoForUser) {
				m.On("GetUser", mock.Anything, "user-001").Return(user.User{
					UserID:   "user-001",
					Username: "alice",
					TeamName: "backend",
					IsActive: true,
				}, nil)
				m.On("GetUserPullRequests", mock.Anything, "user-001").Return(nil, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
			validateResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockRepoForUser)
			tt.setupMock(mockRepo)

			service := &Service{
				repo: mockRepo,
			}

			ctx := context.Background()
			result, err := service.GetReview(ctx, tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				if errors.Is(tt.expectedError, ErrUserNotFound) {
					assert.ErrorIs(t, err, ErrUserNotFound)
				} else {
					assert.Contains(t, err.Error(), tt.expectedError.Error())
				}
				assert.Nil(t, result)
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

