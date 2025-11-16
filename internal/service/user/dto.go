package user

import "github.com/aabbuukkaarr8/PRService/internal/repository/user"

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

// RepoUser структура пользователя из repository (для передачи данных между repository и service)

func (m *User) FillFromDB(dbu *user.User) {
	m.UserID = dbu.UserID
	m.Username = dbu.Username
	m.TeamName = dbu.TeamName
	m.IsActive = dbu.IsActive
}
