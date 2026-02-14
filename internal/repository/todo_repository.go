package repository

import (
	"context"
	"fmt"
	"time"
	"todos_api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

/*
CreateTodo inserts a new ToDo into the database for a specific user.

This function:
  - Creates a context with a 5-second timeout to prevent long-running queries
  - Inserts the ToDo record into the todos table
  - Returns the newly created ToDo including auto-generated fields

Parameters:
  pool      - PostgreSQL connection pool
  title     - Title of the ToDo
  completed - Initial completion status
  userID    - ID of the user who owns the ToDo

Returns:
  *models.ToDo - The created ToDo object
  error        - Database or execution error

Security:
  The userID ensures the ToDo is associated with the correct authenticated user.

Database fields returned:
  - id
  - title
  - completed
  - created_at
  - updated_at
  - user_id
*/
func CreateTodo(pool *pgxpool.Pool, title string, completed bool, userID string) (*models.ToDo, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		INSERT INTO todos (title, completed, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, title, completed, created_at, updated_at, user_id
	`
	var todo models.ToDo

	var err error = pool.QueryRow(ctx, query, title, completed, userID).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.UserID,
	)

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

/*
GetAllTodos retrieves all ToDos belonging to a specific user.

This function:
  - Uses a timeout-protected context
  - Queries all ToDos filtered by user_id
  - Orders results by creation time (newest first)

Parameters:
  pool   - PostgreSQL connection pool
  userID - ID of the authenticated user

Returns:
  []models.ToDo - Slice of ToDos belonging to the user
  error         - Database error

Security:
  Ensures users can only retrieve their own ToDos via WHERE user_id clause.
*/
func GetAllTodos(pool *pgxpool.Pool, userID string) ([]models.ToDo, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
	SELECT id, title, completed, created_at, updated_at, user_id
	FROM todos
	WHERE user_id = $1
	ORDER BY created_at DESC
	`
	rows, err := pool.Query(ctx, query, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var todos []models.ToDo = []models.ToDo{}

	for rows.Next() {
		var todo models.ToDo

		err = rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
			&todo.UserID,
		)

		if err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

/*
GetTodoByID retrieves a specific ToDo by its ID and owner.

This function ensures:
  - ToDo exists
  - ToDo belongs to the authenticated user

Parameters:
  pool   - PostgreSQL connection pool
  id     - ToDo ID
  userID - Owner user ID

Returns:
  *models.ToDo - Retrieved ToDo
  error        - Not found or database error

Security:
  Uses BOTH id AND user_id to prevent unauthorized access to other users' ToDos.
*/
func GetTodoByID(pool *pgxpool.Pool, id int, userID string) (*models.ToDo, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
	SELECT id, title, completed, created_at, updated_at, user_id
	FROM todos
	WHERE id = $1 AND user_id = $2
	`
	var todo models.ToDo

	var err error = pool.QueryRow(ctx, query, id, userID).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.UserID,
	)

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

/*
UpdateTodo modifies an existing ToDo.

This function:
  - Updates title and completion status
  - Updates the updated_at timestamp automatically
  - Ensures only the owner can update the ToDo

Parameters:
  pool      - PostgreSQL connection pool
  id        - ToDo ID
  title     - Updated title
  completed - Updated completion status
  userID    - Owner user ID

Returns:
  *models.ToDo - Updated ToDo object
  error        - Database or authorization error

Security:
  Prevents unauthorized updates by validating user ownership.
*/
func UpdateTodo(pool *pgxpool.Pool, id int, title string, completed bool, userID string) (*models.ToDo, error) {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
	UPDATE todos
	SET title = $1, completed = $2, updated_at = CURRENT_TIMESTAMP
	WHERE id = $3 AND user_id = $4
	RETURNING id, title, completed, created_at, updated_at, user_id
	`
	var todo models.ToDo

	var err error = pool.QueryRow(ctx, query, title, completed, id, userID).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.UserID,
	)

	if err != nil {
		return nil, err
	}

	return &todo, nil
}

/*
DeleteTodo removes a ToDo from the database.

This function:
  - Ensures only the owner can delete the ToDo
  - Uses Exec since no row is returned
  - Checks RowsAffected to confirm deletion occurred

Parameters:
  pool   - PostgreSQL connection pool
  id     - ToDo ID
  userID - Owner user ID

Returns:
  error - nil if successful, error otherwise

Security:
  Prevents users from deleting ToDos they do not own.
*/
func DeleteTodo(pool *pgxpool.Pool, id int, userID string) error {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
	DELETE FROM todos
	WHERE id = $1 AND user_id = $2
	`
	var commandTag, err = pool.Exec(ctx, query, id, userID)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("ToDo with id: %v not found", id)
	}

	return nil
}
