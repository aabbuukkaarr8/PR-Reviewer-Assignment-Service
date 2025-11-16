package user

type User struct {
	UserID   string
	Username string
	TeamName string
	IsActive bool
}

type PullRequestShort struct {
	PullRequestID   string
	PullRequestName string
	AuthorID        string
	Status          string
}
