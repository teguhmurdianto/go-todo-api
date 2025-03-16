package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/teguh/go-todo-api/internal/app/models"
	"github.com/teguh/go-todo-api/internal/database"
)

// TodoRepository handles database operations for todos
type TodoRepository struct {
	db *sql.DB
}

// NewTodoRepository creates a new TodoRepository
func NewTodoRepository() *TodoRepository {
	return &TodoRepository{
		db: database.DB,
	}
}

// Create inserts a new todo into the database
func (r *TodoRepository) Create(todo *models.Todo) error {
	query := `
		INSERT INTO todos (id, title, description, completed, priority, due_date, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(
		query,
		todo.ID,
		todo.Title,
		todo.Description,
		todo.Completed,
		todo.Priority,
		todo.DueDate,
		todo.CreatedAt,
		todo.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}

	return nil
}

// GetByID retrieves a todo by its ID
func (r *TodoRepository) GetByID(id string) (*models.Todo, error) {
	query := `
		SELECT id, title, description, completed, priority, due_date, created_at, updated_at
		FROM todos
		WHERE id = ?
	`

	var todo models.Todo
	err := r.db.QueryRow(query, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.Priority,
		&todo.DueDate,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to get todo by ID: %w", err)
	}

	todo.FormatDates()
	return &todo, nil
}

// GetAll retrieves all todos with optional filtering
func (r *TodoRepository) GetAll(completed *bool) ([]*models.Todo, error) {
	var query string
	var args []interface{}

	if completed != nil {
		query = `
			SELECT id, title, description, completed, priority, due_date, created_at, updated_at
			FROM todos
			WHERE completed = ?
			ORDER BY priority DESC, created_at DESC
		`
		args = append(args, *completed)
	} else {
		query = `
			SELECT id, title, description, completed, priority, due_date, created_at, updated_at
			FROM todos
			ORDER BY priority DESC, created_at DESC
		`
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query todos: %w", err)
	}
	defer rows.Close()

	var todos []*models.Todo
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.Completed,
			&todo.Priority,
			&todo.DueDate,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan todo row: %w", err)
		}

		todo.FormatDates()
		todos = append(todos, &todo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating todo rows: %w", err)
	}

	return todos, nil
}

// Update updates a todo in the database
func (r *TodoRepository) Update(id string, update *models.TodoUpdate) (*models.Todo, error) {
	// First get the existing todo
	todo, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}
	if todo == nil {
		return nil, nil // Not found
	}

	// Apply updates if provided
	if update.Title != nil {
		todo.Title = *update.Title
	}
	if update.Description != nil {
		todo.Description = *update.Description
	}
	if update.Completed != nil {
		todo.Completed = *update.Completed
	}
	if update.Priority != nil {
		todo.Priority = *update.Priority
	}
	if update.DueDate != nil {
		if *update.DueDate == "" {
			todo.DueDate = sql.NullTime{Valid: false}
			todo.DueDateStr = ""
		} else {
			dueDate, err := time.Parse(time.RFC3339, *update.DueDate)
			if err != nil {
				return nil, fmt.Errorf("invalid due date format: %w", err)
			}
			todo.DueDate = sql.NullTime{
				Time:  dueDate,
				Valid: true,
			}
			todo.DueDateStr = *update.DueDate
		}
	}

	// Update the updated_at timestamp
	todo.UpdatedAt = time.Now()

	// Perform the update
	query := `
		UPDATE todos
		SET title = ?, description = ?, completed = ?, priority = ?, due_date = ?, updated_at = ?
		WHERE id = ?
	`

	_, err = r.db.Exec(
		query,
		todo.Title,
		todo.Description,
		todo.Completed,
		todo.Priority,
		todo.DueDate,
		todo.UpdatedAt,
		todo.ID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update todo: %w", err)
	}

	return todo, nil
}

// Delete removes a todo from the database
func (r *TodoRepository) Delete(id string) error {
	query := "DELETE FROM todos WHERE id = ?"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil // Not found, but not an error
	}

	return nil
}
