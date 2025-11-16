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

func TestService_SetIsActive(t *testing.T) {
	tests := []struct {
		name           string
		userID          string
		isActive        bool
		setupMock      func(*mockRepoForUser)
		expectedError  error
		validateResult func(*testing.T, User)
	}{
		{
			name:    "successful set to active",
			userID:   "user-001",
			isActive: true,
			setupMock: func(m *mockRepoForUser) {
				m.On("GetUser", mock.Anything, "user-001").Return(user.User{
					UserID:   "user-001",
					Username: "alice",
					TeamName: "backend",
					IsActive: false,
				}, nil)
				m.On("UpdateUserIsActive", mock.Anything, "user-001", true).Return(nil)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, u User) {
				assert.Equal(t, "user-001", u.UserID)
				assert.Equal(t, "alice", u.Username)
				assert.Equal(t, "backend", u.TeamName)
				assert.True(t, u.IsActive)
			},
		},
		{
			name:    "successful set to inactive",
			userID:   "user-002",
			isActive: false,
			setupMock: func(m *mockRepoForUser) {
				m.On("GetUser", mock.Anything, "user-002").Return(user.User{
					UserID:   "user-002",
					Username: "bob",
					TeamName: "frontend",
					IsActive: true,
				}, nil)
				m.On("UpdateUserIsActive", mock.Anything, "user-002", false).Return(nil)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, u User) {
				assert.Equal(t, "user-002", u.UserID)
				assert.False(t, u.IsActive)
			},
		},
		{
			name:    "user not found",
			userID:   "user-999",
			isActive: true,
			setupMock: func(m *mockRepoForUser) {
				m.On("GetUser", mock.Anything, "user-999").Return(user.User{}, sql.ErrNoRows)
			},
			expectedError: ErrUserNotFound,
			validateResult: nil,
		},
		{
			name:    "error getting user",
			userID:   "user-001",
			isActive: true,
			setupMock: func(m *mockRepoForUser) {
				m.On("GetUser", mock.Anything, "user-001").Return(user.User{}, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
			validateResult: nil,
		},
		{
			name:    "error updating user",
			userID:   "user-001",
			isActive: true,
			setupMock: func(m *mockRepoForUser) {
				m.On("GetUser", mock.Anything, "user-001").Return(user.User{
					UserID:   "user-001",
					Username: "alice",
					TeamName: "backend",
					IsActive: false,
				}, nil)
				m.On("UpdateUserIsActive", mock.Anything, "user-001", true).Return(errors.New("database error"))
			},
			expectedError: errors.New("database error"),
			validateResult: nil,
		},
		{
			name:    "idempotent update - already active",
			userID:   "user-003",
			isActive: true,
			setupMock: func(m *mockRepoForUser) {
				m.On("GetUser", mock.Anything, "user-003").Return(user.User{
					UserID:   "user-003",
					Username: "charlie",
					TeamName: "devops",
					IsActive: true,
				}, nil)
				m.On("UpdateUserIsActive", mock.Anything, "user-003", true).Return(nil)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, u User) {
				assert.Equal(t, "user-003", u.UserID)
				assert.True(t, u.IsActive)
			},
		},
		{
			name:    "idempotent update - already inactive",
			userID:   "user-004",
			isActive: false,
			setupMock: func(m *mockRepoForUser) {
				m.On("GetUser", mock.Anything, "user-004").Return(user.User{
					UserID:   "user-004",
					Username: "david",
					TeamName: "qa",
					IsActive: false,
				}, nil)
				m.On("UpdateUserIsActive", mock.Anything, "user-004", false).Return(nil)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, u User) {
				assert.Equal(t, "user-004", u.UserID)
				assert.False(t, u.IsActive)
			},
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
			result, err := service.SetIsActive(ctx, tt.userID, tt.isActive)

			if tt.expectedError != nil {
				assert.Error(t, err)
				if errors.Is(tt.expectedError, ErrUserNotFound) {
					assert.ErrorIs(t, err, ErrUserNotFound)
				} else {
					assert.Contains(t, err.Error(), tt.expectedError.Error())
				}
				assert.Equal(t, User{}, result)
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

