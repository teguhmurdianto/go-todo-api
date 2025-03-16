package services

import (
	"errors"
	"fmt"

	"github.com/teguh/go-todo-api/internal/app/models"
	"github.com/teguh/go-todo-api/internal/app/repositories"
)

// TodoService handles business logic for todos
type TodoService struct {
	repo *repositories.TodoRepository
}

// NewTodoService creates a new TodoService
func NewTodoService() *TodoService {
	return &TodoService{
		repo: repositories.NewTodoRepository(),
	}
}

// CreateTodo creates a new todo
func (s *TodoService) CreateTodo(create models.TodoCreate) (*models.Todo, error) {
	// Validate input
	if create.Title == "" {
		return nil, errors.New("title is required")
	}

	// Create the todo model
	todo, err := models.NewTodo(create)
	if err != nil {
		return nil, fmt.Errorf("failed to create todo model: %w", err)
	}

	// Save to database
	if err := s.repo.Create(todo); err != nil {
		return nil, fmt.Errorf("failed to save todo: %w", err)
	}

	return todo, nil
}

// GetTodoByID retrieves a todo by its ID
func (s *TodoService) GetTodoByID(id string) (*models.Todo, error) {
	todo, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}
	if todo == nil {
		return nil, errors.New("todo not found")
	}
	return todo, nil
}

// GetAllTodos retrieves all todos with optional filtering
func (s *TodoService) GetAllTodos(completed *bool) ([]*models.Todo, error) {
	todos, err := s.repo.GetAll(completed)
	if err != nil {
		return nil, fmt.Errorf("failed to get todos: %w", err)
	}
	return todos, nil
}

// UpdateTodo updates a todo
func (s *TodoService) UpdateTodo(id string, update models.TodoUpdate) (*models.Todo, error) {
	// Validate that the todo exists
	exists, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to check if todo exists: %w", err)
	}
	if exists == nil {
		return nil, errors.New("todo not found")
	}

	// Update the todo
	updated, err := s.repo.Update(id, &update)
	if err != nil {
		return nil, fmt.Errorf("failed to update todo: %w", err)
	}
	if updated == nil {
		return nil, errors.New("todo not found")
	}

	return updated, nil
}

// DeleteTodo deletes a todo
func (s *TodoService) DeleteTodo(id string) error {
	// Validate that the todo exists
	exists, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to check if todo exists: %w", err)
	}
	if exists == nil {
		return errors.New("todo not found")
	}

	// Delete the todo
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	return nil
}
