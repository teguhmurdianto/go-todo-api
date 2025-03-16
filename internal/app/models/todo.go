package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Todo represents a todo item
type Todo struct {
	ID          string       `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Completed   bool         `json:"completed"`
	Priority    int          `json:"priority"`
	DueDate     sql.NullTime `json:"-"`
	DueDateStr  string       `json:"due_date,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// TodoCreate represents the data needed to create a new todo
type TodoCreate struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	Priority    int    `json:"priority"`
	DueDate     string `json:"due_date,omitempty"`
}

// TodoUpdate represents the data needed to update a todo
type TodoUpdate struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Completed   *bool   `json:"completed,omitempty"`
	Priority    *int    `json:"priority,omitempty"`
	DueDate     *string `json:"due_date,omitempty"`
}

// NewTodo creates a new Todo with default values
func NewTodo(create TodoCreate) (*Todo, error) {
	todo := &Todo{
		ID:          uuid.New().String(),
		Title:       create.Title,
		Description: create.Description,
		Completed:   false,
		Priority:    create.Priority,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Parse due date if provided
	if create.DueDate != "" {
		dueDate, err := time.Parse(time.RFC3339, create.DueDate)
		if err != nil {
			return nil, err
		}
		todo.DueDate = sql.NullTime{
			Time:  dueDate,
			Valid: true,
		}
		todo.DueDateStr = create.DueDate
	}

	return todo, nil
}

// FormatDates formats the dates for JSON response
func (t *Todo) FormatDates() {
	if t.DueDate.Valid {
		t.DueDateStr = t.DueDate.Time.Format(time.RFC3339)
	}
}
