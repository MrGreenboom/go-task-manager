package service

import (
	"context"
	"testing"

	"github.com/MrGreenboom/go-task-manager/internal/model"
)

type fakeTaskRepo struct {
	created bool
}

func (f *fakeTaskRepo) Create(ctx context.Context, t *model.Task) (int64, error) {
	f.created = true
	return 1, nil
}
func (f *fakeTaskRepo) GetByID(ctx context.Context, userID, id int64) (*model.Task, error) {
	return nil, nil
}
func (f *fakeTaskRepo) List(ctx context.Context, userID int64) ([]model.Task, error) {
	return nil, nil
}
func (f *fakeTaskRepo) Update(ctx context.Context, userID int64, t *model.Task) error {
	return nil
}
func (f *fakeTaskRepo) Delete(ctx context.Context, userID, id int64) error {
	return nil
}

func TestTaskService_Create_OK(t *testing.T) {
	repo := &fakeTaskRepo{}
	svc := NewTaskService(repo)

	task := &model.Task{
		Title:  "test",
		UserID: 1,
	}

	_, err := svc.Create(context.Background(), task)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !repo.created {
		t.Fatal("expected Create to be called")
	}
}

func TestTaskService_Create_NoTitle(t *testing.T) {
	repo := &fakeTaskRepo{}
	svc := NewTaskService(repo)

	task := &model.Task{UserID: 1}
	_, err := svc.Create(context.Background(), task)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestTaskService_Create_NoUserID(t *testing.T) {
	repo := &fakeTaskRepo{}
	svc := NewTaskService(repo)

	task := &model.Task{Title: "x"}
	_, err := svc.Create(context.Background(), task)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
