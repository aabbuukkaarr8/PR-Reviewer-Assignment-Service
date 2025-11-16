package pullrequest

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aabbuukkaarr8/PRService/internal/api/models"
	"github.com/aabbuukkaarr8/PRService/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestRepository_CreatePullRequest(t *testing.T) {
	tests := []struct {
		name           string
		request        *CreatePullRequest
		setupMock      func(mock sqlmock.Sqlmock)
		expectedResult PullRequest
		expectedError  error
	}{
		{
			name: "successful creation",
			request: &CreatePullRequest{
				PullRequestId:     "pr-001",
				PullRequestName:   "Test PR",
				AuthorId:          "user-001",
				Status:            models.PullRequestStatusOPEN,
				AssignedReviewers: []string{"user-002", "user-003"},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO pullrequests`).
					WithArgs("pr-001", "Test PR", "user-001", "OPEN", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
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
			name: "successful creation with empty reviewers",
			request: &CreatePullRequest{
				PullRequestId:     "pr-002",
				PullRequestName:   "Test PR 2",
				AuthorId:          "user-001",
				Status:            models.PullRequestStatusOPEN,
				AssignedReviewers: []string{},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO pullrequests`).
					WithArgs("pr-002", "Test PR 2", "user-001", "OPEN", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedResult: PullRequest{
				PullRequestID:     "pr-002",
				PullRequestName:   "Test PR 2",
				AuthorID:          "user-001",
				Status:            "OPEN",
				AssignedReviewers: []string{},
			},
			expectedError: nil,
		},
		{
			name: "database error",
			request: &CreatePullRequest{
				PullRequestId:     "pr-003",
				PullRequestName:   "Test PR 3",
				AuthorId:          "user-001",
				Status:            models.PullRequestStatusOPEN,
				AssignedReviewers: []string{"user-002"},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO pullrequests`).
					WithArgs("pr-003", "Test PR 3", "user-001", "OPEN", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("database connection error"))
			},
			expectedResult: PullRequest{},
			expectedError:  errors.New("database connection error"),
		},
		{
			name: "duplicate key error",
			request: &CreatePullRequest{
				PullRequestId:     "pr-004",
				PullRequestName:   "Test PR 4",
				AuthorId:          "user-001",
				Status:            models.PullRequestStatusOPEN,
				AssignedReviewers: []string{"user-002"},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO pullrequests`).
					WithArgs("pr-004", "Test PR 4", "user-001", "OPEN", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("duplicate key value violates unique constraint"))
			},
			expectedResult: PullRequest{},
			expectedError:  errors.New("duplicate key value violates unique constraint"),
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
			result, err := repo.CreatePullRequest(ctx, tt.request)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Equal(t, PullRequest{}, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult.PullRequestID, result.PullRequestID)
				assert.Equal(t, tt.expectedResult.PullRequestName, result.PullRequestName)
				assert.Equal(t, tt.expectedResult.AuthorID, result.AuthorID)
				assert.Equal(t, tt.expectedResult.Status, result.Status)
				assert.Equal(t, tt.expectedResult.AssignedReviewers, result.AssignedReviewers)
				assert.NotNil(t, result.CreatedAt)
				assert.Nil(t, result.MergedAt)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

