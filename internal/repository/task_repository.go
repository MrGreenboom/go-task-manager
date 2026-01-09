package repository

import (
	"context"

	"github.com/MrGreenboom/go-task-manager/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, t *model.Task) (int64, error) {
	query := `
		INSERT INTO tasks (user_id, title, description, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`
	var id int64
	err := r.db.QueryRow(ctx, query, t.UserID, t.Title, t.Description, t.Status).Scan(&id)
	return id, err
}

func (r *TaskRepository) GetByID(ctx context.Context, id int64) (*model.Task, error) {
	query := `
		SELECT id, user_id, title, description, status, created_at, updated_at
		FROM tasks
		WHERE id = $1;
	`
	row := r.db.QueryRow(ctx, query, id)

	var t model.Task
	if err := row.Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TaskRepository) List(ctx context.Context) ([]model.Task, error) {
	query := `
		SELECT id, user_id, title, description, status, created_at, updated_at
		FROM tasks
		ORDER BY id DESC;
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]model.Task, 0)
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

func (r *TaskRepository) Update(ctx context.Context, t *model.Task) error {
	query := `
		UPDATE tasks
		SET title=$1, description=$2, status=$3, updated_at=now()
		WHERE id=$4;
	`
	_, err := r.db.Exec(ctx, query, t.Title, t.Description, t.Status, t.ID)
	return err
}

func (r *TaskRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM tasks WHERE id=$1;`, id)
	return err
}
