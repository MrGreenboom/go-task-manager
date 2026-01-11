package service

import (
	"context"
	"errors"
	"strings"

	"github.com/MrGreenboom/go-task-manager/internal/model"
)

type TaskRepo interface {
	Create(ctx context.Context, t *model.Task) (int64, error)
	GetByID(ctx context.Context, userID, id int64) (*model.Task, error)
	List(ctx context.Context, userID int64) ([]model.Task, error)
	Update(ctx context.Context, userID int64, t *model.Task) error
	Delete(ctx context.Context, userID, id int64) error
}

type TaskService struct {
	repo TaskRepo
}

func NewTaskService(repo TaskRepo) *TaskService{
	return &TaskService{repo: repo}
}

func (s *TaskService) Create(ctx context.Context, t *model.Task) (int64, error) {
	t.Title = strings.TrimSpace(t.Title)
	t.Description = strings.TrimSpace(t.Description)

	if t.Title == "" {
		return 0, errors.New("title is required")
	}
	if t.Status == "" {
		t.Status = "new"
	}
	if t.UserID <= 0 {
		return 0, errors.New("user_id is required")
	}

	return s.repo.Create(ctx, t)
}

func (s *TaskService) GetByID(ctx context.Context, userID, id int64) (*model.Task, error) {
	if userID <= 0 || id <= 0 {
		return nil, errors.New("invalid ids")
	}
	return s.repo.GetByID(ctx, userID, id)
}

func (s *TaskService) List(ctx context.Context, userID int64) ([]model.Task, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user id")
	}
	return s.repo.List(ctx, userID)
}

func (s *TaskService) Update(ctx context.Context, userID int64, t *model.Task) error {
	t.Title = strings.TrimSpace(t.Title)
	t.Description = strings.TrimSpace(t.Description)

	if userID <= 0 || t.ID <= 0 {
		return errors.New("invalid ids")
	}
	if t.Title == "" {
		return errors.New("title is required")
	}
	if t.Status == "" {
		t.Status = "new"
	}

	return s.repo.Update(ctx, userID, t)
}

func (s *TaskService) Delete(ctx context.Context, userID, id int64) error {
	if userID <= 0 || id <= 0 {
		return errors.New("invalid ids")
	}
	return s.repo.Delete(ctx, userID, id)
}
