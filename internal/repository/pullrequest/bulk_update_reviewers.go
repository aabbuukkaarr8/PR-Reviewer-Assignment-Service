package pullrequest

import (
	"context"

	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type PRReviewerUpdate struct {
	PullRequestID     string
	AssignedReviewers []string
}

func (r *Repository) BulkUpdatePullRequestReviewers(ctx context.Context, updates []PRReviewerUpdate) error {
	if len(updates) == 0 {
		return nil
	}

	tx, err := r.store.GetConn().BeginTx(ctx, nil)
	if err != nil {
		logrus.WithError(err).WithField("updates_count", len(updates)).Error("Database error: failed to begin transaction for bulk update")
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx,
		`UPDATE pullrequests SET assigned_reviewers = $1 WHERE pull_request_id = $2`)
	if err != nil {
		logrus.WithError(err).Error("Database error: failed to prepare statement for bulk update")
		return err
	}
	defer stmt.Close()

	for _, update := range updates {
		_, err := stmt.ExecContext(ctx, pq.Array(update.AssignedReviewers), update.PullRequestID)
		if err != nil {
			logrus.WithError(err).WithField("pull_request_id", update.PullRequestID).Error("Database error: failed to update PR reviewers")
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		logrus.WithError(err).WithField("updates_count", len(updates)).Error("Database error: failed to commit bulk update transaction")
		return err
	}

	return nil
}
