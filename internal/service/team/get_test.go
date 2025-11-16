package team

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/aabbuukkaarr8/PRService/internal/repository/team"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_GetTeam(t *testing.T) {
	tests := []struct {
		name           string
		teamName       string
		setupMock      func(*mockRepo)
		expectedError  error
		validateResult func(*testing.T, Team)
	}{
		{
			name:     "successful get with multiple users",
			teamName: "backend",
			setupMock: func(m *mockRepo) {
				m.On("GetTeam", mock.Anything, "backend").Return("backend", []team.User{
					{UserID: "user-001", Username: "alice", IsActive: true},
					{UserID: "user-002", Username: "bob", IsActive: true},
					{UserID: "user-003", Username: "charlie", IsActive: false},
				}, nil)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, team Team) {
				assert.Equal(t, "backend", team.TeamName)
				assert.Len(t, team.Members, 3)
				assert.Equal(t, "user-001", team.Members[0].UserID)
				assert.Equal(t, "alice", team.Members[0].Username)
				assert.True(t, team.Members[0].IsActive)
				assert.Equal(t, "user-003", team.Members[2].UserID)
				assert.False(t, team.Members[2].IsActive)
			},
		},
		{
			name:     "successful get with single user",
			teamName:  "frontend",
			setupMock: func(m *mockRepo) {
				m.On("GetTeam", mock.Anything, "frontend").Return("frontend", []team.User{
					{UserID: "user-004", Username: "david", IsActive: true},
				}, nil)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, team Team) {
				assert.Equal(t, "frontend", team.TeamName)
				assert.Len(t, team.Members, 1)
				assert.Equal(t, "user-004", team.Members[0].UserID)
				assert.Equal(t, "david", team.Members[0].Username)
			},
		},
		{
			name:     "successful get with no users",
			teamName: "empty-team",
			setupMock: func(m *mockRepo) {
				m.On("GetTeam", mock.Anything, "empty-team").Return("empty-team", []team.User{}, nil)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, team Team) {
				assert.Equal(t, "empty-team", team.TeamName)
				assert.Empty(t, team.Members)
			},
		},
		{
			name:     "team not found",
			teamName: "nonexistent",
			setupMock: func(m *mockRepo) {
				m.On("GetTeam", mock.Anything, "nonexistent").Return("", nil, sql.ErrNoRows)
			},
			expectedError: ErrTeamNotFound,
			validateResult: nil,
		},
		{
			name:     "error getting team",
			teamName: "backend",
			setupMock: func(m *mockRepo) {
				m.On("GetTeam", mock.Anything, "backend").Return("", nil, errors.New("database error"))
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
			result, err := service.GetTeam(ctx, tt.teamName)

			if tt.expectedError != nil {
				assert.Error(t, err)
				if errors.Is(tt.expectedError, ErrTeamNotFound) {
					assert.ErrorIs(t, err, ErrTeamNotFound)
				} else {
					assert.Contains(t, err.Error(), tt.expectedError.Error())
				}
				assert.Equal(t, Team{}, result)
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

