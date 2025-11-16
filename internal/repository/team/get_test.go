package team

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aabbuukkaarr8/PRService/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestRepository_GetTeam(t *testing.T) {
	tests := []struct {
		name           string
		teamName       string
		setupMock      func(mock sqlmock.Sqlmock)
		expectedResult struct {
			teamName string
			users    []User
		}
		expectedError error
	}{
		{
			name:     "successful get with multiple users",
			teamName: "backend",
			setupMock: func(mock sqlmock.Sqlmock) {
				existsRow := sqlmock.NewRows([]string{"exists"}).AddRow(true)
				mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM teams WHERE team_name`).
					WithArgs("backend").
					WillReturnRows(existsRow)

				usersRows := sqlmock.NewRows([]string{"user_id", "username", "is_active"}).
					AddRow("user-001", "alice", true).
					AddRow("user-002", "bob", true).
					AddRow("user-003", "charlie", false)
				mock.ExpectQuery(`SELECT user_id, username, is_active FROM users WHERE team_name`).
					WithArgs("backend").
					WillReturnRows(usersRows)
			},
			expectedResult: struct {
				teamName string
				users    []User
			}{
				teamName: "backend",
				users: []User{
					{UserID: "user-001", Username: "alice", IsActive: true},
					{UserID: "user-002", Username: "bob", IsActive: true},
					{UserID: "user-003", Username: "charlie", IsActive: false},
				},
			},
			expectedError: nil,
		},
		{
			name:     "successful get with single user",
			teamName: "frontend",
			setupMock: func(mock sqlmock.Sqlmock) {
				existsRow := sqlmock.NewRows([]string{"exists"}).AddRow(true)
				mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM teams WHERE team_name`).
					WithArgs("frontend").
					WillReturnRows(existsRow)

				usersRows := sqlmock.NewRows([]string{"user_id", "username", "is_active"}).
					AddRow("user-004", "david", true)
				mock.ExpectQuery(`SELECT user_id, username, is_active FROM users WHERE team_name`).
					WithArgs("frontend").
					WillReturnRows(usersRows)
			},
			expectedResult: struct {
				teamName string
				users    []User
			}{
				teamName: "frontend",
				users: []User{
					{UserID: "user-004", Username: "david", IsActive: true},
				},
			},
			expectedError: nil,
		},
		{
			name:     "successful get with no users",
			teamName: "empty-team",
			setupMock: func(mock sqlmock.Sqlmock) {
				existsRow := sqlmock.NewRows([]string{"exists"}).AddRow(true)
				mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM teams WHERE team_name`).
					WithArgs("empty-team").
					WillReturnRows(existsRow)

				usersRows := sqlmock.NewRows([]string{"user_id", "username", "is_active"})
				mock.ExpectQuery(`SELECT user_id, username, is_active FROM users WHERE team_name`).
					WithArgs("empty-team").
					WillReturnRows(usersRows)
			},
			expectedResult: struct {
				teamName string
				users    []User
			}{
				teamName: "empty-team",
				users:    []User{},
			},
			expectedError: nil,
		},
		{
			name:     "team not found",
			teamName: "nonexistent",
			setupMock: func(mock sqlmock.Sqlmock) {
				existsRow := sqlmock.NewRows([]string{"exists"}).AddRow(false)
				mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM teams WHERE team_name`).
					WithArgs("nonexistent").
					WillReturnRows(existsRow)
			},
			expectedResult: struct {
				teamName string
				users    []User
			}{
				teamName: "",
				users:    nil,
			},
			expectedError: sql.ErrNoRows,
		},
		{
			name:     "database error on exists check",
			teamName: "backend",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM teams WHERE team_name`).
					WithArgs("backend").
					WillReturnError(errors.New("database connection error"))
			},
			expectedResult: struct {
				teamName string
				users    []User
			}{
				teamName: "",
				users:    nil,
			},
			expectedError: errors.New("database connection error"),
		},
		{
			name:     "database error on users query",
			teamName: "backend",
			setupMock: func(mock sqlmock.Sqlmock) {
				existsRow := sqlmock.NewRows([]string{"exists"}).AddRow(true)
				mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM teams WHERE team_name`).
					WithArgs("backend").
					WillReturnRows(existsRow)

				mock.ExpectQuery(`SELECT user_id, username, is_active FROM users WHERE team_name`).
					WithArgs("backend").
					WillReturnError(errors.New("database query error"))
			},
			expectedResult: struct {
				teamName string
				users    []User
			}{
				teamName: "",
				users:    nil,
			},
			expectedError: errors.New("database query error"),
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
			teamName, users, err := repo.GetTeam(ctx, tt.teamName)

			if tt.expectedError != nil {
				assert.Error(t, err)
				if errors.Is(tt.expectedError, sql.ErrNoRows) {
					assert.ErrorIs(t, err, sql.ErrNoRows)
				} else {
					assert.Contains(t, err.Error(), tt.expectedError.Error())
				}
				assert.Equal(t, "", teamName)
				assert.Nil(t, users)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult.teamName, teamName)
				assert.Equal(t, len(tt.expectedResult.users), len(users))
				for i, expectedUser := range tt.expectedResult.users {
					assert.Equal(t, expectedUser.UserID, users[i].UserID)
					assert.Equal(t, expectedUser.Username, users[i].Username)
					assert.Equal(t, expectedUser.IsActive, users[i].IsActive)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

