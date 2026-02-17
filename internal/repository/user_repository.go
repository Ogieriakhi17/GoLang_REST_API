package repository

import (
	"context"
	"time"
	"todos_api/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

/*
CreateUser inserts a new user into the database.

This function:
  - Creates a timeout-protected context to prevent hanging queries
  - Inserts the user's email and hashed password
  - Returns the complete created user record including generated fields

Parameters:
  pool - PostgreSQL connection pool
  user - Pointer to User struct containing email and hashed password

Returns:
  *models.User - Newly created user with ID and timestamps populated
  error        - Database error (including duplicate email violations)

Security:
  - Password must already be hashed before calling this function.
  - This function NEVER hashes passwords itself.
  - Prevents storage of plaintext passwords.

Database fields returned:
  - id
  - email
  - password (hashed)
  - created_at
  - updated_at

Common errors:
  - Unique constraint violation (duplicate email)
  - Database connectivity issues
*/
func CreateUser(pool *pgxpool.Pool, user *models.User) (*models.User, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
	INSERT INTO users (email, password)
	VALUES ($1, $2)
	RETURNING id, email, password, created_at, updated_at
	`

	err := pool.QueryRow(ctx, query, user.Email, user.Password).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

/*
GetUserByEmail retrieves a user from the database using their email address.

This function is primarily used for authentication during login.

Parameters:
  pool  - PostgreSQL connection pool
  email - Email address of the user

Returns:
  *models.User - User record if found
  error        - pgx.ErrNoRows if user does not exist

Security:
  - Used during login to retrieve hashed password
  - Password verification occurs outside this function

Authentication flow usage:
  1. Retrieve user by email
  2. Compare hashed password using bcrypt
  3. Generate JWT token if valid

Returned fields:
  - id
  - email
  - password (hashed)
  - created_at
  - updated_at
*/
func GetUserByEmail(pool *pgxpool.Pool, email string) (*models.User, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT id, email, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	var user models.User

	err := pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
/*
GetUserByID retrieves a user by their unique ID.

This function is commonly used for:
  - Token validation
  - Retrieving authenticated user details
  - Authorization checks

Parameters:
  pool - PostgreSQL connection pool
  id   - User ID

Returns:
  *models.User - User record if found
  error        - pgx.ErrNoRows if user does not exist

Security:
  - Used after JWT token verification
  - Ensures authenticated user exists in database

Returned fields:
  - id
  - email
  - password (hashed)
  - created_at
  - updated_at

Common usage flow:
  JWT Token → extract user_id → call GetUserByID → authorize request
*/
func GetUserByID(pool *pgxpool.Pool, id int) (*models.User, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT id, email, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var user models.User

	err := pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
