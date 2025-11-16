package errors

import "fmt"

type ErrCode string

const (
	CodeTeamExists   ErrCode = "TEAM_EXISTS"
	CodePRExists     ErrCode = "PR_EXISTS"
	CodePRMerged     ErrCode = "PR_MERGED"
	CodeNotAssigned  ErrCode = "NOT_ASSIGNED"
	CodeNoCandidate  ErrCode = "NO_CANDIDATE"
	CodeNotFound     ErrCode = "NOT_FOUND"
	CodeInvalidInput ErrCode = "INVALID_INPUT"
)

type AppError struct {
	Code    ErrCode `json:"code"`
	Message string  `json:"message"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewTeamExists(teamName string) *AppError {
	return &AppError{
		Code:    CodeTeamExists,
		Message: fmt.Sprintf("%s already exists", teamName),
	}
}

func NewPRExists(prID string) *AppError {
	return &AppError{
		Code:    CodePRExists,
		Message: fmt.Sprintf("PR %s already exists", prID),
	}
}

func NewPRMerged() *AppError {
	return &AppError{
		Code:    CodePRMerged,
		Message: "cannot reassign on merged PR",
	}
}

func NewNotAssigned() *AppError {
	return &AppError{
		Code:    CodeNotAssigned,
		Message: "reviewer is not assigned to this PR",
	}
}

func NewNoCandidate() *AppError {
	return &AppError{
		Code:    CodeNoCandidate,
		Message: "no active replacement candidate in team",
	}
}

func NewNotFound() *AppError {
	return &AppError{
		Code:    CodeNotFound,
		Message: "resource not found",
	}
}

func NewInvalidInput(message string) *AppError {
	return &AppError{
		Code:    CodeInvalidInput,
		Message: message,
	}
}

func IsTeamExists(err error) bool {
	return isErrCode(err, CodeTeamExists)
}

func IsPRExists(err error) bool {
	return isErrCode(err, CodePRExists)
}

func IsPRMerged(err error) bool {
	return isErrCode(err, CodePRMerged)
}

func IsNotFound(err error) bool {
	return isErrCode(err, CodeNotFound)
}

func IsNotAssigned(err error) bool {
	return isErrCode(err, CodeNotAssigned)
}

func IsNoCandidate(err error) bool {
	return isErrCode(err, CodeNoCandidate)
}

func isErrCode(err error, code ErrCode) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == code
	}
	return false
}
