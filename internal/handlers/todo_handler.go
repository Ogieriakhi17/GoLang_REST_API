package handlers

import (
	"net/http"
	"strconv"
	"todos_api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateToDoInput struct {
	Title     string `json:"title" binding:"required"`
	Completed bool   `json:"completed"`
}

type UpdateTodoInput struct {
	Title     *string `json: "title"`
	Completed *bool   `json: "completed"`
}

/*
CreateToDoHandler creates a new ToDo for the authenticated user.

This handler:
 1. Extracts the authenticated user's ID from Gin context (set by AuthMiddleware)
 2. Validates and binds the JSON request body
 3. Calls the repository layer to insert the ToDo into the database
 4. Returns the created ToDo with HTTP 201 status

Authentication Required: YES

Possible responses:
  201 Created       - ToDo successfully created
  400 Bad Request   - Invalid JSON or missing required fields
  500 Internal Error - Database or server error
*/
func CreateToDoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input CreateToDoInput
		UserIDInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id does not exist"})
			return
		}

		UserID := UserIDInterface.(string)

		if err := c.ShouldBind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		todo, err := repository.CreateTodo(pool, input.Title, input.Completed, UserID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, todo)
	}
}

/*
GetAllTodosHandler retrieves all ToDos belonging to the authenticated user.

This handler ensures users only see their own ToDos.

Authentication Required: YES

Possible responses:
  200 OK            - Returns list of ToDos
  500 Internal Error - Database or server error
*/
func GetAllTodosHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		UserIDInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id does not exist"})
			return
		}

		UserID := UserIDInterface.(string)

		todos, err := repository.GetAllTodos(pool, UserID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, todos)
	}
}

/*
GetTodoByIDHandler retrieves a specific ToDo by its ID.

Ensures:
  - Valid ID format
  - ToDo belongs to authenticated user

Authentication Required: YES

URL Parameter:
  id (int) - ToDo ID

Possible responses:
  200 OK           - Returns requested ToDo
  400 Bad Request  - Invalid ID format
  404 Not Found    - ToDo does not exist or does not belong to user
  500 Internal Error - Database error
*/
func GetTodoByIDHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		UserIDInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id does not exist"})
			return
		}

		UserID := UserIDInterface.(string)

		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request format"})
			return
		}

		todos, err := repository.GetTodoByID(pool, id, UserID)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "To-Do not found"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, todos)
	}
}

/*

// UpdateTodoHandler updates an existing ToDo.
//
// Supports partial updates:
//   - Title only
//   - Completed only
//   - Both fields
//
// This handler:
//   1. Validates user authentication
//   2. Parses ToDo ID
//   3. Validates request body
//   4. Fetches existing ToDo
//   5. Applies partial updates
//   6. Saves updated ToDo
//
// Authentication Required: YES
//
// Possible responses:
//   200 OK
//   400 Bad Request
//   404 Not Found
//   500 Internal Error
*/
func UpdateTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		UserIDInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id does not exist"})
			return
		}

		UserID := UserIDInterface.(string)

		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
			return
		}

		var input UpdateTodoInput

		if err := c.ShouldBind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if input.Title == nil && input.Completed == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "At least one field is required (title/completed)"})
			return
		}

		existing, err := repository.GetTodoByID(pool, id, UserID)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "ToDo not Found"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		title := existing.Title

		if input.Title != nil {
			title = *input.Title
		}

		completed := existing.Completed
		if input.Completed != nil {
			completed = *input.Completed
		}

		todo, err := repository.UpdateTodo(pool, id, title, completed, UserID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, todo)
	}
}

/*
DeleteTodoHandler deletes a ToDo belonging to the authenticated user.

Ensures users can only delete their own ToDos.

Authentication Required: YES

Possible responses:
  200 OK
  400 Bad Request
  404 Not Found
  500 Internal Error
*/
func DeleteTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		UserIDInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id does not exist"})
			return
		}

		UserID := UserIDInterface.(string)

		idString := c.Param("id")
		id, err := strconv.Atoi(idString)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
			return
		}

		err = repository.DeleteTodo(pool, id, UserID)

		if err != nil {
			if err.Error() == "ToDo with id: "+idString+" not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": "ToDo not Found"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"message": "ToDo successfully deleted"})
	}
}
