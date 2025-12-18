package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"user-api/internal/models"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, name string, dob time.Time) (*models.User, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
	Update(ctx context.Context, id int64, name string, dob time.Time) (*models.User, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]*models.User, error)
	Count(ctx context.Context) (int64, error)
}

type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// Create inserts a new user into the database
func (r *userRepository) Create(ctx context.Context, name string, dob time.Time) (*models.User, error) {
	query := `
		INSERT INTO users (name, dob)
		VALUES ($1, $2)
		RETURNING id, name, dob, created_at, updated_at
	`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, name, dob).Scan(
		&user.ID,
		&user.Name,
		&user.DOB,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	query := `
		SELECT id, name, dob, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.DOB,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

// Update modifies an existing user
func (r *userRepository) Update(ctx context.Context, id int64, name string, dob time.Time) (*models.User, error) {
	query := `
		UPDATE users
		SET name = $2, dob = $3
		WHERE id = $1
		RETURNING id, name, dob, created_at, updated_at
	`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, id, name, dob).Scan(
		&user.ID,
		&user.Name,
		&user.DOB,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

// Delete removes a user from the database
func (r *userRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// List retrieves users with pagination
func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*models.User, error) {
	query := `
		SELECT id, name, dob, created_at, updated_at
		FROM users
		ORDER BY id
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.DOB,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Count returns the total number of users
func (r *userRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM users`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
