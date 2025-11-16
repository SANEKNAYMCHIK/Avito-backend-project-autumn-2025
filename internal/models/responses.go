package models

import "time"

type TeamMember struct {
	IsActive bool   `json:"is_active"`
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}

type TeamResponse struct {
	Members  []TeamMember `json:"members"`
	TeamName string       `json:"team_name"`
}

type UserResponse struct {
	IsActive bool   `json:"is_active"`
	TeamName string `json:"team_name"`
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}

type PullRequestResponse struct {
	AssignedReviewers []string   `json:"assigned_reviewers"`
	AuthorId          string     `json:"author_id"`
	CreatedAt         *time.Time `json:"createdAt"`
	MergedAt          *time.Time `json:"mergedAt"`
	PullRequestId     string     `json:"pull_request_id"`
	PullRequestName   string     `json:"pull_request_name"`
	Status            string     `json:"status"`
}

type PullRequestShort struct {
	AuthorId        string `json:"author_id"`
	PullRequestId   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	Status          string `json:"status"`
}

type ReassignResponse struct {
	PR         PullRequestResponse `json:"pr"`
	ReplacedBy string              `json:"replaced_by"`
}

type UserPRsResponse struct {
	UserID       string             `json:"user_id"`
	PullRequests []PullRequestShort `json:"pull_requests"`
}
