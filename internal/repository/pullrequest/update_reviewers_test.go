package pullrequest

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/aabbuukkaarr8/PRService/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestRepository_UpdatePullRequestReviewers(t *testing.T) {
	tests := []struct {
		name                string
		pullRequestID       string
		assignedReviewers   []string
		setupMock           func(mock sqlmock.Sqlmock)
		expectedResult      PullRequest
		expectedError       error
	}{
		{
			name:              "successful update with 2 reviewers",
			pullRequestID:     "pr-001",
			assignedReviewers: []string{"user-002", "user-003"},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE pullrequests SET assigned_reviewers`).
					WithArgs(sqlmock.AnyArg(), "pr-001").
					WillReturnResult(sqlmock.NewResult(1, 1))

				createdAt := time.Now()
				rows := sqlmock.NewRows([]string{"pull_request_id", "pull_request_name", "author_id", "status", "assigned_reviewers", "created_at", "merged_at"}).
					AddRow("pr-001", "Test PR", "user-001", "OPEN", pq.Array([]string{"user-002", "user-003"}), createdAt, nil)
				mock.ExpectQuery(`SELECT pull_request_id, pull_request_name, author_id, status, assigned_reviewers, created_at, merged_at`).
					WithArgs("pr-001").
					WillReturnRows(rows)
			},
			expectedResult: PullRequest{
				PullRequestID:     "pr-001",
				PullRequestName:   "Test PR",
				AuthorID:          "user-001",
				Status:            "OPEN",
				AssignedReviewers: []string{"user-002", "user-003"},
			},
			expectedError: nil,
		},
		{
			name:              "successful update with 1 reviewer",
			pullRequestID:     "pr-002",
			assignedReviewers: []string{"user-002"},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE pullrequests SET assigned_reviewers`).
					WithArgs(sqlmock.AnyArg(), "pr-002").
					WillReturnResult(sqlmock.NewResult(1, 1))

				createdAt := time.Now()
				rows := sqlmock.NewRows([]string{"pull_request_id", "pull_request_name", "author_id", "status", "assigned_reviewers", "created_at", "merged_at"}).
					AddRow("pr-002", "Test PR 2", "user-001", "OPEN", pq.Array([]string{"user-002"}), createdAt, nil)
				mock.ExpectQuery(`SELECT pull_request_id, pull_request_name, author_id, status, assigned_reviewers, created_at, merged_at`).
					WithArgs("pr-002").
					WillReturnRows(rows)
			},
			expectedResult: PullRequest{
				PullRequestID:     "pr-002",
				PullRequestName:   "Test PR 2",
				AuthorID:          "user-001",
				Status:            "OPEN",
				AssignedReviewers: []string{"user-002"},
			},
			expectedError: nil,
		},
		{
			name:              "successful update with empty reviewers",
			pullRequestID:     "pr-003",
			assignedReviewers: []string{},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE pullrequests SET assigned_reviewers`).
					WithArgs(sqlmock.AnyArg(), "pr-003").
					WillReturnResult(sqlmock.NewResult(1, 1))

				createdAt := time.Now()
				rows := sqlmock.NewRows([]string{"pull_request_id", "pull_request_name", "author_id", "status", "assigned_reviewers", "created_at", "merged_at"}).
					AddRow("pr-003", "Test PR 3", "user-001", "OPEN", pq.Array([]string{}), createdAt, nil)
				mock.ExpectQuery(`SELECT pull_request_id, pull_request_name, author_id, status, assigned_reviewers, created_at, merged_at`).
					WithArgs("pr-003").
					WillReturnRows(rows)
			},
			expectedResult: PullRequest{
				PullRequestID:     "pr-003",
				PullRequestName:   "Test PR 3",
				AuthorID:          "user-001",
				Status:            "OPEN",
				AssignedReviewers: []string{},
			},
			expectedError: nil,
		},
		{
			name:              "PR not found after update",
			pullRequestID:     "pr-004",
			assignedReviewers: []string{"user-002"},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE pullrequests SET assigned_reviewers`).
					WithArgs(sqlmock.AnyArg(), "pr-004").
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectQuery(`SELECT pull_request_id, pull_request_name, author_id, status, assigned_reviewers, created_at, merged_at`).
					WithArgs("pr-004").
					WillReturnError(sql.ErrNoRows)
			},
			expectedResult: PullRequest{},
			expectedError:  sql.ErrNoRows,
		},
		{
			name:              "database error on update",
			pullRequestID:     "pr-005",
			assignedReviewers: []string{"user-002"},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE pullrequests SET assigned_reviewers`).
					WithArgs(sqlmock.AnyArg(), "pr-005").
					WillReturnError(errors.New("database connection error"))
			},
			expectedResult: PullRequest{},
			expectedError:  errors.New("database connection error"),
		},
		{
			name:              "database error on select",
			pullRequestID:     "pr-006",
			assignedReviewers: []string{"user-002"},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE pullrequests SET assigned_reviewers`).
					WithArgs(sqlmock.AnyArg(), "pr-006").
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectQuery(`SELECT pull_request_id, pull_request_name, author_id, status, assigned_reviewers, created_at, merged_at`).
					WithArgs("pr-006").
					WillReturnError(errors.New("database query error"))
			},
			expectedResult: PullRequest{},
			expectedError:  errors.New("database query error"),
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
			result, err := repo.UpdatePullRequestReviewers(ctx, tt.pullRequestID, tt.assignedReviewers)

			if tt.expectedError != nil {
				assert.Error(t, err)
				if errors.Is(tt.expectedError, sql.ErrNoRows) {
					assert.ErrorIs(t, err, sql.ErrNoRows)
				} else {
					assert.Contains(t, err.Error(), tt.expectedError.Error())
				}
				assert.Equal(t, PullRequest{}, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult.PullRequestID, result.PullRequestID)
				assert.Equal(t, tt.expectedResult.PullRequestName, result.PullRequestName)
				assert.Equal(t, tt.expectedResult.AuthorID, result.AuthorID)
				assert.Equal(t, tt.expectedResult.Status, result.Status)
				assert.Equal(t, tt.expectedResult.AssignedReviewers, result.AssignedReviewers)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

