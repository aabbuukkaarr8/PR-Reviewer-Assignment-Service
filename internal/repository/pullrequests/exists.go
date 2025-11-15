package pullrequests

import "context"

// PRExists проверяет, существует ли PR
func (r *Repository) PRExists(ctx context.Context, pullRequestID string) (bool, error) {
	var exists bool
	err := r.store.GetConn().QueryRowContext(ctx,
		"SELECT EXISTS(SELECT 1 FROM pullrequests WHERE pull_request_id = $1)",
		pullRequestID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
