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
	MergedAt          *time.Time `json:"mergedAt"`
	PullRequestId     string     `json:"pull_request_id"`
	PullRequestName   string     `json:"pull_request_name"`
	Status            string     `json:"status"`
}

type PullRequestShort struct {
	AssignedReviewers []string `json:"assigned_reviewers"`
	AuthorId          string   `json:"author_id"`
	PullRequestId     string   `json:"pull_request_id"`
	PullRequestName   string   `json:"pull_request_name"`
	Status            string   `json:"status"`
}

type ReassignResponse struct {
	PR         PullRequestResponse `json:"pr"`
	ReplacedBy string              `json:"replaced_by"`
}

type PullRequestReview struct {
	AuthorId        string `json:"author_id"`
	PullRequestId   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	Status          string `json:"status"`
}

type UserPRsResponse struct {
	UserID       string              `json:"user_id"`
	PullRequests []PullRequestReview `json:"pull_requests"`
}

type CreateTeamRequest struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

type CreatePRRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type MergePRRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

type ReassignRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_user_id"`
}

type SetActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}
