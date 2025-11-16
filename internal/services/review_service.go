package services

import (
	"fmt"
	"math/rand"
	"slices"
	"time"

	"github.com/SANEKNAYMCHIK/Avito-backend-project-autumn-2025/internal/errors"
	"github.com/SANEKNAYMCHIK/Avito-backend-project-autumn-2025/internal/models"
	"github.com/SANEKNAYMCHIK/Avito-backend-project-autumn-2025/internal/repositories"
	"github.com/lib/pq"
)

type reviewService struct {
	repo *repositories.Repository
}

func NewReviewService(repo *repositories.Repository) ReviewService {
	return &reviewService{repo: repo}
}

func (s *reviewService) CreateTeam(teamName string, members []models.TeamMember) (*models.TeamResponse, error) {
	exists, err := s.repo.Team.TeamExists(teamName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.NewTeamExists(teamName)
	}
	team := &models.Team{
		ID:   generateID(),
		Name: teamName,
	}
	users := make([]models.User, len(members))
	for i, member := range members {
		users[i] = models.User{
			ID:       member.UserId,
			Username: member.Username,
			IsActive: member.IsActive,
		}
	}
	err = s.repo.Team.CreateTeam(team, users)
	if err != nil {
		return nil, err
	}
	return &models.TeamResponse{
		TeamName: teamName,
		Members:  members,
	}, nil
}

func (s *reviewService) GetTeam(teamName string) (*models.TeamResponse, error) {
	team, err := s.repo.Team.GetTeamByName(teamName)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, errors.NewNotFound()
	}
	users, err := s.repo.Team.GetTeamUsers(team.ID)
	if err != nil {
		return nil, err
	}
	members := make([]models.TeamMember, len(users))
	for i, user := range users {
		members[i] = models.TeamMember{
			UserId:   user.ID,
			Username: user.Username,
			IsActive: user.IsActive,
		}
	}
	return &models.TeamResponse{
		TeamName: team.Name,
		Members:  members,
	}, nil
}

func (s *reviewService) SetUserActive(userID string, isActive bool) (*models.UserResponse, error) {
	user, err := s.repo.User.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.NewNotFound()
	}
	err = s.repo.User.UpdateUser(userID, isActive)
	if err != nil {
		return nil, err
	}
	user.IsActive = isActive
	team, err := s.repo.User.GetUserTeam(userID)
	if err != nil {
		return nil, err
	}
	return &models.UserResponse{
		UserId:   user.ID,
		Username: user.Username,
		TeamName: team.Name,
		IsActive: user.IsActive,
	}, nil
}

func (s *reviewService) CreatePR(prID, title, authorID string) (*models.PullRequestResponse, error) {
	author, err := s.repo.User.GetUserByID(authorID)
	if err != nil {
		return nil, err
	}
	if author == nil {
		return nil, errors.NewNotFound()
	}
	reviewers, err := s.autoAssignReviewers(authorID)
	if err != nil {
		return nil, err
	}
	pr := &models.PullRequest{
		ID:        prID,
		Title:     title,
		AuthorID:  authorID,
		Status:    "OPEN",
		Reviewers: pq.StringArray(reviewers),
		CreatedAt: time.Now(),
	}
	err = s.repo.PR.CreatePR(pr)
	if err != nil {
		return nil, err
	}
	return &models.PullRequestResponse{
		PullRequestId:     pr.ID,
		PullRequestName:   pr.Title,
		AuthorId:          pr.AuthorID,
		Status:            pr.Status,
		AssignedReviewers: pr.Reviewers,
		CreatedAt:         &pr.CreatedAt,
	}, nil
}

func (s *reviewService) autoAssignReviewers(authorID string) ([]string, error) {
	team, err := s.repo.User.GetUserTeam(authorID)
	if err != nil {
		return nil, err
	}
	activeUsers, err := s.repo.User.GetActiveUsersByTeam(team.ID)
	if err != nil {
		return nil, err
	}
	var candidates []string
	for _, user := range activeUsers {
		if user.ID != authorID {
			candidates = append(candidates, user.ID)
		}
	}
	if len(candidates) == 0 {
		return []string{}, nil
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})
	if len(candidates) > 2 {
		return candidates[:2], nil
	}
	return candidates, nil
}

func (s *reviewService) MergePR(prID string) (*models.PullRequestResponse, error) {
	pr, err := s.repo.PR.GetPRByID(prID)
	if err != nil {
		return nil, err
	}
	if pr.Status == "MERGED" {
		return s.convertPRToResponse(pr), nil
	}
	now := time.Now()
	pr.Status = "MERGED"
	pr.MergedAt = &now
	err = s.repo.PR.UpdatePR(pr)
	if err != nil {
		return nil, err
	}
	return s.convertPRToResponse(pr), nil
}

func (s *reviewService) ReassignReviewer(prID, oldReviewerID string) (*models.ReassignResponse, error) {
	pr, err := s.repo.PR.GetPRByID(prID)
	if err != nil {
		return nil, err
	}
	if pr.Status == "MERGED" {
		return nil, errors.NewPRMerged()
	}
	if !slices.Contains(pr.Reviewers, oldReviewerID) {
		return nil, errors.NewNotAssigned()
	}
	newReviewerID, err := s.findReplacementReviewer(oldReviewerID, pr.Reviewers, pr.AuthorID)
	if err != nil {
		return nil, err
	}
	for i, reviewer := range pr.Reviewers {
		if reviewer == oldReviewerID {
			pr.Reviewers[i] = newReviewerID
			break
		}
	}
	err = s.repo.PR.UpdatePR(pr)
	if err != nil {
		return nil, err
	}
	response := &models.ReassignResponse{
		PR:         *s.convertPRToResponse(pr),
		ReplacedBy: newReviewerID,
	}
	return response, nil
}

func (s *reviewService) findReplacementReviewer(oldReviewerID string, currentReviewers []string, authorID string) (string, error) {
	team, err := s.repo.User.GetUserTeam(oldReviewerID)
	if err != nil {
		return "", err
	}
	activeUsers, err := s.repo.User.GetActiveUsersByTeam(team.ID)
	if err != nil {
		return "", err
	}
	var candidates []string
	for _, user := range activeUsers {
		if user.ID != authorID &&
			user.ID != oldReviewerID &&
			!slices.Contains(currentReviewers, user.ID) {
			candidates = append(candidates, user.ID)
		}
	}
	if len(candidates) == 0 {
		return "", errors.NewNoCandidate()
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return candidates[r.Intn(len(candidates))], nil
}

func (s *reviewService) GetUserReviews(userID string) (*models.UserPRsResponse, error) {
	user, err := s.repo.User.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.NewNotFound()
	}
	prs, err := s.repo.PR.GetPRsByReviewer(userID)
	if err != nil {
		return nil, err
	}
	prShorts := make([]models.PullRequestShort, len(prs))
	for i, pr := range prs {
		prShorts[i] = models.PullRequestShort{
			PullRequestId:   pr.ID,
			PullRequestName: pr.Title,
			AuthorId:        pr.AuthorID,
			Status:          pr.Status,
		}
	}
	return &models.UserPRsResponse{
		UserID:       userID,
		PullRequests: prShorts,
	}, nil
}

func (s *reviewService) convertPRToResponse(pr *models.PullRequest) *models.PullRequestResponse {
	return &models.PullRequestResponse{
		PullRequestId:     pr.ID,
		PullRequestName:   pr.Title,
		AuthorId:          pr.AuthorID,
		Status:            pr.Status,
		AssignedReviewers: pr.Reviewers,
		CreatedAt:         &pr.CreatedAt,
	}
}

func generateID() string {
	return fmt.Sprint(time.Now().UnixNano())
}
