package team

import (
	"context"
	"errors"
	"testing"

	"github.com/aabbuukkaarr8/PRService/internal/repository/team"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) TeamExists(ctx context.Context, teamName string) (bool, error) {
	args := m.Called(ctx, teamName)
	return args.Bool(0), args.Error(1)
}

func (m *mockRepo) CreateTeam(ctx context.Context, teamName string) error {
	args := m.Called(ctx, teamName)
	return args.Error(0)
}

func (m *mockRepo) GetTeam(ctx context.Context, teamName string) (string, []team.User, error) {
	args := m.Called(ctx, teamName)
	if args.Get(0) == nil {
		return "", nil, args.Error(2)
	}
	if args.Get(1) == nil {
		return args.String(0), nil, args.Error(2)
	}
	return args.String(0), args.Get(1).([]team.User), args.Error(2)
}

func (m *mockRepo) UserExists(ctx context.Context, userID string) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}

func (m *mockRepo) CreateUser(ctx context.Context, userID, username, teamName string, isActive bool) error {
	args := m.Called(ctx, userID, username, teamName, isActive)
	return args.Error(0)
}

func (m *mockRepo) UpdateUser(ctx context.Context, userID, username, teamName string, isActive bool) error {
	args := m.Called(ctx, userID, username, teamName, isActive)
	return args.Error(0)
}

func (m *mockRepo) BeginTx(ctx context.Context) (team.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(team.Tx), args.Error(1)
}

type mockTx struct {
	mock.Mock
}

func (m *mockTx) CreateTeam(teamName string) error {
	args := m.Called(teamName)
	return args.Error(0)
}

func (m *mockTx) CreateUser(userID, username, teamName string, isActive bool) error {
	args := m.Called(userID, username, teamName, isActive)
	return args.Error(0)
}

func (m *mockTx) UpdateUser(userID, username, teamName string, isActive bool) error {
	args := m.Called(userID, username, teamName, isActive)
	return args.Error(0)
}

func (m *mockTx) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *mockTx) Rollback() error {
	args := m.Called()
	return args.Error(0)
}

func TestService_CreateTeam(t *testing.T) {
	tests := []struct {
		name           string
		team           Team
		setupMock      func(*mockRepo, *mockTx)
		expectedError  error
		validateResult func(*testing.T, Team)
	}{
		{
			name: "successful creation with new users",
			team: Team{
				TeamName: "backend",
				Members: []TeamMember{
					{UserID: "user-001", Username: "alice", IsActive: true},
					{UserID: "user-002", Username: "bob", IsActive: true},
				},
			},
			setupMock: func(m *mockRepo, tx *mockTx) {
				m.On("TeamExists", mock.Anything, "backend").Return(false, nil)
				m.On("BeginTx", mock.Anything).Return(tx, nil)
				tx.On("CreateTeam", "backend").Return(nil)
				m.On("UserExists", mock.Anything, "user-001").Return(false, nil)
				tx.On("CreateUser", "user-001", "alice", "backend", true).Return(nil)
				m.On("UserExists", mock.Anything, "user-002").Return(false, nil)
				tx.On("CreateUser", "user-002", "bob", "backend", true).Return(nil)
				tx.On("Commit").Return(nil)
				tx.On("Rollback").Return(nil)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, team Team) {
				assert.Equal(t, "backend", team.TeamName)
				assert.Len(t, team.Members, 2)
			},
		},
		{
			name: "successful creation with existing users (update)",
			team: Team{
				TeamName: "frontend",
				Members: []TeamMember{
					{UserID: "user-003", Username: "charlie", IsActive: true},
				},
			},
			setupMock: func(m *mockRepo, tx *mockTx) {
				m.On("TeamExists", mock.Anything, "frontend").Return(false, nil)
				m.On("BeginTx", mock.Anything).Return(tx, nil)
				tx.On("CreateTeam", "frontend").Return(nil)
				m.On("UserExists", mock.Anything, "user-003").Return(true, nil)
				tx.On("UpdateUser", "user-003", "charlie", "frontend", true).Return(nil)
				tx.On("Commit").Return(nil)
				tx.On("Rollback").Return(nil)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, team Team) {
				assert.Equal(t, "frontend", team.TeamName)
				assert.Len(t, team.Members, 1)
			},
		},
		{
			name: "successful creation with mixed new and existing users",
			team: Team{
				TeamName: "devops",
				Members: []TeamMember{
					{UserID: "user-004", Username: "david", IsActive: true},
					{UserID: "user-005", Username: "eve", IsActive: false},
				},
			},
			setupMock: func(m *mockRepo, tx *mockTx) {
				m.On("TeamExists", mock.Anything, "devops").Return(false, nil)
				m.On("BeginTx", mock.Anything).Return(tx, nil)
				tx.On("CreateTeam", "devops").Return(nil)
				m.On("UserExists", mock.Anything, "user-004").Return(false, nil)
				tx.On("CreateUser", "user-004", "david", "devops", true).Return(nil)
				m.On("UserExists", mock.Anything, "user-005").Return(true, nil)
				tx.On("UpdateUser", "user-005", "eve", "devops", false).Return(nil)
				tx.On("Commit").Return(nil)
				tx.On("Rollback").Return(nil)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, team Team) {
				assert.Equal(t, "devops", team.TeamName)
				assert.Len(t, team.Members, 2)
			},
		},
		{
			name: "team already exists",
			team: Team{
				TeamName: "backend",
				Members: []TeamMember{
					{UserID: "user-001", Username: "alice", IsActive: true},
				},
			},
			setupMock: func(m *mockRepo, tx *mockTx) {
				m.On("TeamExists", mock.Anything, "backend").Return(true, nil)
			},
			expectedError: ErrTeamExists,
			validateResult: nil,
		},
		{
			name: "error checking team existence",
			team: Team{
				TeamName: "backend",
				Members: []TeamMember{
					{UserID: "user-001", Username: "alice", IsActive: true},
				},
			},
			setupMock: func(m *mockRepo, tx *mockTx) {
				m.On("TeamExists", mock.Anything, "backend").Return(false, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
			validateResult: nil,
		},
		{
			name: "error beginning transaction",
			team: Team{
				TeamName: "backend",
				Members: []TeamMember{
					{UserID: "user-001", Username: "alice", IsActive: true},
				},
			},
			setupMock: func(m *mockRepo, tx *mockTx) {
				m.On("TeamExists", mock.Anything, "backend").Return(false, nil)
				m.On("BeginTx", mock.Anything).Return((*mockTx)(nil), errors.New("transaction error"))
			},
			expectedError: errors.New("transaction error"),
			validateResult: nil,
		},
		{
			name: "error creating team in transaction",
			team: Team{
				TeamName: "backend",
				Members: []TeamMember{
					{UserID: "user-001", Username: "alice", IsActive: true},
				},
			},
			setupMock: func(m *mockRepo, tx *mockTx) {
				m.On("TeamExists", mock.Anything, "backend").Return(false, nil)
				m.On("BeginTx", mock.Anything).Return(tx, nil)
				tx.On("CreateTeam", "backend").Return(errors.New("database error"))
				tx.On("Rollback").Return(nil)
			},
			expectedError: errors.New("database error"),
			validateResult: nil,
		},
		{
			name: "error checking user existence",
			team: Team{
				TeamName: "backend",
				Members: []TeamMember{
					{UserID: "user-001", Username: "alice", IsActive: true},
				},
			},
			setupMock: func(m *mockRepo, tx *mockTx) {
				m.On("TeamExists", mock.Anything, "backend").Return(false, nil)
				m.On("BeginTx", mock.Anything).Return(tx, nil)
				tx.On("CreateTeam", "backend").Return(nil)
				m.On("UserExists", mock.Anything, "user-001").Return(false, errors.New("database error"))
				tx.On("Rollback").Return(nil)
			},
			expectedError: errors.New("database error"),
			validateResult: nil,
		},
		{
			name: "error creating user in transaction",
			team: Team{
				TeamName: "backend",
				Members: []TeamMember{
					{UserID: "user-001", Username: "alice", IsActive: true},
				},
			},
			setupMock: func(m *mockRepo, tx *mockTx) {
				m.On("TeamExists", mock.Anything, "backend").Return(false, nil)
				m.On("BeginTx", mock.Anything).Return(tx, nil)
				tx.On("CreateTeam", "backend").Return(nil)
				m.On("UserExists", mock.Anything, "user-001").Return(false, nil)
				tx.On("CreateUser", "user-001", "alice", "backend", true).Return(errors.New("database error"))
				tx.On("Rollback").Return(nil)
			},
			expectedError: errors.New("database error"),
			validateResult: nil,
		},
		{
			name: "error updating user in transaction",
			team: Team{
				TeamName: "backend",
				Members: []TeamMember{
					{UserID: "user-001", Username: "alice", IsActive: true},
				},
			},
			setupMock: func(m *mockRepo, tx *mockTx) {
				m.On("TeamExists", mock.Anything, "backend").Return(false, nil)
				m.On("BeginTx", mock.Anything).Return(tx, nil)
				tx.On("CreateTeam", "backend").Return(nil)
				m.On("UserExists", mock.Anything, "user-001").Return(true, nil)
				tx.On("UpdateUser", "user-001", "alice", "backend", true).Return(errors.New("database error"))
				tx.On("Rollback").Return(nil)
			},
			expectedError: errors.New("database error"),
			validateResult: nil,
		},
		{
			name: "error committing transaction",
			team: Team{
				TeamName: "backend",
				Members: []TeamMember{
					{UserID: "user-001", Username: "alice", IsActive: true},
				},
			},
			setupMock: func(m *mockRepo, tx *mockTx) {
				m.On("TeamExists", mock.Anything, "backend").Return(false, nil)
				m.On("BeginTx", mock.Anything).Return(tx, nil)
				tx.On("CreateTeam", "backend").Return(nil)
				m.On("UserExists", mock.Anything, "user-001").Return(false, nil)
				tx.On("CreateUser", "user-001", "alice", "backend", true).Return(nil)
				tx.On("Commit").Return(errors.New("commit error"))
				tx.On("Rollback").Return(nil)
			},
			expectedError: errors.New("commit error"),
			validateResult: nil,
		},
		{
			name: "successful creation with empty members",
			team: Team{
				TeamName: "empty-team",
				Members:  []TeamMember{},
			},
			setupMock: func(m *mockRepo, tx *mockTx) {
				m.On("TeamExists", mock.Anything, "empty-team").Return(false, nil)
				m.On("BeginTx", mock.Anything).Return(tx, nil)
				tx.On("CreateTeam", "empty-team").Return(nil)
				tx.On("Commit").Return(nil)
				tx.On("Rollback").Return(nil)
			},
			expectedError: nil,
			validateResult: func(t *testing.T, team Team) {
				assert.Equal(t, "empty-team", team.TeamName)
				assert.Empty(t, team.Members)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockRepo)
			mockTx := new(mockTx)
			tt.setupMock(mockRepo, mockTx)

			service := &Service{
				repo: mockRepo,
			}

			ctx := context.Background()
			result, err := service.CreateTeam(ctx, tt.team)

			if tt.expectedError != nil {
				assert.Error(t, err)
				if errors.Is(tt.expectedError, ErrTeamExists) {
					assert.ErrorIs(t, err, ErrTeamExists)
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
			mockTx.AssertExpectations(t)
		})
	}
}

