package team

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aabbuukkaarr8/PRService/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestRepository_CreateTeam(t *testing.T) {
	tests := []struct {
		name          string
		teamName      string
		setupMock     func(mock sqlmock.Sqlmock)
		expectedError error
	}{
		{
			name:     "successful creation",
			teamName: "backend",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO teams \(team_name\) VALUES`).
					WithArgs("backend").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: nil,
		},
		{
			name:     "duplicate team name",
			teamName: "backend",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO teams \(team_name\) VALUES`).
					WithArgs("backend").
					WillReturnError(errors.New("duplicate key value violates unique constraint"))
			},
			expectedError: errors.New("duplicate key value violates unique constraint"),
		},
		{
			name:     "database error",
			teamName: "frontend",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO teams \(team_name\) VALUES`).
					WithArgs("frontend").
					WillReturnError(errors.New("database connection error"))
			},
			expectedError: errors.New("database connection error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock: %v", err)
			}
			defer db.Close()

			tt.setupMock(mock)

			store := store.New()
			store.SetConn(db)

			repo := NewRepository(store)

			ctx := context.Background()
			err = repo.CreateTeam(ctx, tt.teamName)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestRepository_CreateUser(t *testing.T) {
	tests := []struct {
		name          string
		userID        string
		username      string
		teamName      string
		isActive      bool
		setupMock     func(mock sqlmock.Sqlmock)
		expectedError error
	}{
		{
			name:     "successful creation",
			userID:   "user-001",
			username: "alice",
			teamName: "backend",
			isActive: true,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO users \(user_id, username, team_name, is_active\) VALUES`).
					WithArgs("user-001", "alice", "backend", true).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: nil,
		},
		{
			name:     "successful creation with inactive user",
			userID:   "user-002",
			username: "bob",
			teamName: "frontend",
			isActive: false,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO users \(user_id, username, team_name, is_active\) VALUES`).
					WithArgs("user-002", "bob", "frontend", false).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: nil,
		},
		{
			name:     "duplicate user id",
			userID:   "user-001",
			username: "alice",
			teamName: "backend",
			isActive: true,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO users \(user_id, username, team_name, is_active\) VALUES`).
					WithArgs("user-001", "alice", "backend", true).
					WillReturnError(errors.New("duplicate key value violates unique constraint"))
			},
			expectedError: errors.New("duplicate key value violates unique constraint"),
		},
		{
			name:     "database error",
			userID:   "user-003",
			username: "charlie",
			teamName: "devops",
			isActive: true,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO users \(user_id, username, team_name, is_active\) VALUES`).
					WithArgs("user-003", "charlie", "devops", true).
					WillReturnError(errors.New("database connection error"))
			},
			expectedError: errors.New("database connection error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock: %v", err)
			}
			defer db.Close()

			tt.setupMock(mock)

			store := store.New()
			store.SetConn(db)

			repo := NewRepository(store)

			ctx := context.Background()
			err = repo.CreateUser(ctx, tt.userID, tt.username, tt.teamName, tt.isActive)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
