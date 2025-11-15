package users

// User структура пользователя для repository слоя
type User struct {
	UserID   string
	Username string
	TeamName string
	IsActive bool
}

// PullRequestShort короткая версия PR для списков
type PullRequestShort struct {
	PullRequestID   string
	PullRequestName string
	AuthorID        string
	Status          string // OPEN or MERGED
}
