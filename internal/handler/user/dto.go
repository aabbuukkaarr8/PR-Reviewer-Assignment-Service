package user

import (
	usersrv "github.com/aabbuukkaarr8/PRService/internal/service/user"
)

type PullRequestShort struct {
	PullRequestID   string `json:"pull_request_id" binding:"required"`
	PullRequestName string `json:"pull_request_name" binding:"required"`
	AuthorID        string `json:"author_id" binding:"required"`
	Status          string `json:"status" binding:"required,oneof=OPEN MERGED"`
}

type SetIsActiveRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	IsActive *bool  `json:"is_active" binding:"required"`
}

type User struct {
	UserID   string `json:"user_id" binding:"required"`
	Username string `json:"username" binding:"required"`
	TeamName string `json:"team_name" binding:"required"`
	IsActive bool   `json:"is_active" binding:"required"`
}

type SetIsActiveResponse struct {
	User User `json:"user"`
}

type GetReviewResponse struct {
	UserID       string             `json:"user_id"`
	PullRequests []PullRequestShort `json:"pull_requests"`
}

func (p *PullRequestShort) FillFromService(s usersrv.PullRequestShort) {
	p.PullRequestID = s.PullRequestID
	p.PullRequestName = s.PullRequestName
	p.AuthorID = s.AuthorID
	p.Status = s.Status
}

func (u *User) FillFromService(s usersrv.User) {
	u.UserID = s.UserID
	u.Username = s.Username
	u.TeamName = s.TeamName
	u.IsActive = s.IsActive
}
