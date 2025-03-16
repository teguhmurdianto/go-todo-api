package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/teguh/go-todo-api/internal/app/models"
	"github.com/teguh/go-todo-api/internal/app/services"
	"github.com/teguh/go-todo-api/pkg/utils"
)

// TodoHandler handles HTTP requests for todos
type TodoHandler struct {
	service *services.TodoService
}

// NewTodoHandler creates a new TodoHandler
func NewTodoHandler() *TodoHandler {
	return &TodoHandler{
		service: services.NewTodoService(),
	}
}

// RegisterRoutes registers the routes for todos
func (h *TodoHandler) RegisterRoutes(router fiber.Router) {
	todos := router.Group("/todos")

	todos.Post("/", h.CreateTodo)
	todos.Get("/", h.GetAllTodos)
	todos.Get("/:id", h.GetTodoByID)
	todos.Patch("/:id", h.UpdateTodo)
	todos.Delete("/:id", h.DeleteTodo)
}

// CreateTodo handles the creation of a new todo
// @Summary Create a new todo
// @Description Create a new todo item
// @Tags todos
// @Accept json
// @Produce json
// @Param todo body models.TodoCreate true "Todo to create"
// @Success 201 {object} models.Todo
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /todos [post]
func (h *TodoHandler) CreateTodo(c *fiber.Ctx) error {
	var input models.TodoCreate
	if err := c.BodyParser(&input); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	todo, err := h.service.CreateTodo(input)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(todo)
}

// GetAllTodos handles retrieving all todos
// @Summary Get all todos
// @Description Get all todo items, optionally filtered by completion status
// @Tags todos
// @Produce json
// @Param completed query boolean false "Filter by completion status"
// @Success 200 {array} models.Todo
// @Failure 500 {object} utils.ErrorResponse
// @Router /todos [get]
func (h *TodoHandler) GetAllTodos(c *fiber.Ctx) error {
	var completed *bool
	if c.Query("completed") != "" {
		completedVal := c.QueryBool("completed")
		completed = &completedVal
	}

	todos, err := h.service.GetAllTodos(completed)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(todos)
}

// GetTodoByID handles retrieving a todo by ID
// @Summary Get a todo by ID
// @Description Get a todo item by its ID
// @Tags todos
// @Produce json
// @Param id path string true "Todo ID"
// @Success 200 {object} models.Todo
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /todos/{id} [get]
func (h *TodoHandler) GetTodoByID(c *fiber.Ctx) error {
	id := c.Params("id")
	todo, err := h.service.GetTodoByID(id)
	if err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "todo not found" {
			status = fiber.StatusNotFound
		}
		return utils.SendError(c, status, err.Error())
	}

	return c.JSON(todo)
}

// UpdateTodo handles updating a todo
// @Summary Update a todo
// @Description Update a todo item by its ID
// @Tags todos
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Param todo body models.TodoUpdate true "Todo update data"
// @Success 200 {object} models.Todo
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /todos/{id} [patch]
func (h *TodoHandler) UpdateTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	var input models.TodoUpdate
	if err := c.BodyParser(&input); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body")
	}

	todo, err := h.service.UpdateTodo(id, input)
	if err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "todo not found" {
			status = fiber.StatusNotFound
		}
		return utils.SendError(c, status, err.Error())
	}

	return c.JSON(todo)
}

// DeleteTodo handles deleting a todo
// @Summary Delete a todo
// @Description Delete a todo item by its ID
// @Tags todos
// @Produce json
// @Param id path string true "Todo ID"
// @Success 204 "No Content"
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /todos/{id} [delete]
func (h *TodoHandler) DeleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.service.DeleteTodo(id)
	if err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "todo not found" {
			status = fiber.StatusNotFound
		}
		return utils.SendError(c, status, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
