package user

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aabbuukkaarr8/PRService/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestRepository_GetUser(t *testing.T) {
	tests := []struct {
		name           string
		userID          string
		setupMock      func(mock sqlmock.Sqlmock)
		expectedResult User
		expectedError  error
	}{
		{
			name:  "successful get",
			userID: "user-001",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"user_id", "username", "team_name", "is_active"}).
					AddRow("user-001", "alice", "backend", true)
				mock.ExpectQuery(`SELECT user_id, username, team_name, is_active FROM users WHERE user_id`).
					WithArgs("user-001").
					WillReturnRows(rows)
			},
			expectedResult: User{
				UserID:   "user-001",
				Username: "alice",
				TeamName: "backend",
				IsActive: true,
			},
			expectedError: nil,
		},
		{
			name:  "user not found",
			userID: "user-999",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT user_id, username, team_name, is_active FROM users WHERE user_id`).
					WithArgs("user-999").
					WillReturnError(sql.ErrNoRows)
			},
			expectedResult: User{},
			expectedError:  sql.ErrNoRows,
		},
		{
			name:  "database error",
			userID: "user-001",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT user_id, username, team_name, is_active FROM users WHERE user_id`).
					WithArgs("user-001").
					WillReturnError(errors.New("database connection error"))
			},
			expectedResult: User{},
			expectedError:  errors.New("database connection error"),
		},
		{
			name:  "successful get with inactive user",
			userID: "user-002",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"user_id", "username", "team_name", "is_active"}).
					AddRow("user-002", "bob", "frontend", false)
				mock.ExpectQuery(`SELECT user_id, username, team_name, is_active FROM users WHERE user_id`).
					WithArgs("user-002").
					WillReturnRows(rows)
			},
			expectedResult: User{
				UserID:   "user-002",
				Username: "bob",
				TeamName: "frontend",
				IsActive: false,
			},
			expectedError: nil,
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
			result, err := repo.GetUser(ctx, tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				if errors.Is(tt.expectedError, sql.ErrNoRows) {
					assert.ErrorIs(t, err, sql.ErrNoRows)
				} else {
					assert.Contains(t, err.Error(), tt.expectedError.Error())
				}
				assert.Equal(t, User{}, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult.UserID, result.UserID)
				assert.Equal(t, tt.expectedResult.Username, result.Username)
				assert.Equal(t, tt.expectedResult.TeamName, result.TeamName)
				assert.Equal(t, tt.expectedResult.IsActive, result.IsActive)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestRepository_GetUserPullRequests(t *testing.T) {
	tests := []struct {
		name           string
		userID          string
		setupMock      func(mock sqlmock.Sqlmock)
		expectedResult []PullRequestShort
		expectedError  error
	}{
		{
			name:  "successful get with multiple PRs",
			userID: "user-001",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"pull_request_id", "pull_request_name", "author_id", "status"}).
					AddRow("pr-001", "Test PR 1", "user-002", "OPEN").
					AddRow("pr-002", "Test PR 2", "user-003", "MERGED")
				mock.ExpectQuery(`SELECT pull_request_id, pull_request_name, author_id, status`).
					WithArgs("user-001").
					WillReturnRows(rows)
			},
			expectedResult: []PullRequestShort{
				{PullRequestID: "pr-001", PullRequestName: "Test PR 1", AuthorID: "user-002", Status: "OPEN"},
				{PullRequestID: "pr-002", PullRequestName: "Test PR 2", AuthorID: "user-003", Status: "MERGED"},
			},
			expectedError: nil,
		},
		{
			name:  "successful get with single PR",
			userID: "user-002",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"pull_request_id", "pull_request_name", "author_id", "status"}).
					AddRow("pr-003", "Test PR 3", "user-001", "OPEN")
				mock.ExpectQuery(`SELECT pull_request_id, pull_request_name, author_id, status`).
					WithArgs("user-002").
					WillReturnRows(rows)
			},
			expectedResult: []PullRequestShort{
				{PullRequestID: "pr-003", PullRequestName: "Test PR 3", AuthorID: "user-001", Status: "OPEN"},
			},
			expectedError: nil,
		},
		{
			name:  "successful get with no PRs",
			userID: "user-003",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"pull_request_id", "pull_request_name", "author_id", "status"})
				mock.ExpectQuery(`SELECT pull_request_id, pull_request_name, author_id, status`).
					WithArgs("user-003").
					WillReturnRows(rows)
			},
			expectedResult: []PullRequestShort{},
			expectedError:  nil,
		},
		{
			name:  "database error",
			userID: "user-001",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT pull_request_id, pull_request_name, author_id, status`).
					WithArgs("user-001").
					WillReturnError(errors.New("database connection error"))
			},
			expectedResult: nil,
			expectedError:  errors.New("database connection error"),
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
			result, err := repo.GetUserPullRequests(ctx, tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedResult), len(result))
				for i, expectedPR := range tt.expectedResult {
					assert.Equal(t, expectedPR.PullRequestID, result[i].PullRequestID)
					assert.Equal(t, expectedPR.PullRequestName, result[i].PullRequestName)
					assert.Equal(t, expectedPR.AuthorID, result[i].AuthorID)
					assert.Equal(t, expectedPR.Status, result[i].Status)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestRepository_UpdateUserIsActive(t *testing.T) {
	tests := []struct {
		name          string
		userID         string
		isActive       bool
		setupMock     func(mock sqlmock.Sqlmock)
		expectedError error
	}{
		{
			name:    "successful update to active",
			userID:   "user-001",
			isActive: true,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE users SET is_active`).
					WithArgs(true, "user-001").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: nil,
		},
		{
			name:    "successful update to inactive",
			userID:   "user-002",
			isActive: false,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE users SET is_active`).
					WithArgs(false, "user-002").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: nil,
		},
		{
			name:    "user not found",
			userID:   "user-999",
			isActive: true,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE users SET is_active`).
					WithArgs(true, "user-999").
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedError: nil,
		},
		{
			name:    "database error",
			userID:   "user-001",
			isActive: true,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE users SET is_active`).
					WithArgs(true, "user-001").
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
			err = repo.UpdateUserIsActive(ctx, tt.userID, tt.isActive)

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

