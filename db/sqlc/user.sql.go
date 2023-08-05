// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: user.sql

package db

import (
	"context"
	"database/sql"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  email, 
  hashed_pass,
  username, 
  full_name
) VALUES (
  $1, $2, $3, $4
)
RETURNING username, email, hashed_pass, full_name, password_changed_at, created_at, is_email_activated
`

type CreateUserParams struct {
	Email      string `json:"email"`
	HashedPass string `json:"hashed_pass"`
	Username   string `json:"username"`
	FullName   string `json:"full_name"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (Users, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Email,
		arg.HashedPass,
		arg.Username,
		arg.FullName,
	)
	var i Users
	err := row.Scan(
		&i.Username,
		&i.Email,
		&i.HashedPass,
		&i.FullName,
		&i.PasswordChangedAt,
		&i.CreatedAt,
		&i.IsEmailActivated,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT username, email, hashed_pass, full_name, password_changed_at, created_at, is_email_activated FROM users
WHERE username = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, username string) (Users, error) {
	row := q.db.QueryRowContext(ctx, getUser, username)
	var i Users
	err := row.Scan(
		&i.Username,
		&i.Email,
		&i.HashedPass,
		&i.FullName,
		&i.PasswordChangedAt,
		&i.CreatedAt,
		&i.IsEmailActivated,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET hashed_pass = COALESCE($1, hashed_pass),
    password_changed_at = COALESCE($2, password_changed_at),
    email = COALESCE($3, email),
    full_name = COALESCE($4, full_name),
    is_email_activated = COALESCE($5, is_email_activated)
WHERE username = $6
RETURNING username, email, hashed_pass, full_name, password_changed_at, created_at, is_email_activated
`

type UpdateUserParams struct {
	HashedPass        sql.NullString `json:"hashed_pass"`
	PasswordChangedAt sql.NullTime   `json:"password_changed_at"`
	Email             sql.NullString `json:"email"`
	FullName          sql.NullString `json:"full_name"`
	IsEmailActivated  sql.NullBool   `json:"is_email_activated"`
	Username          string         `json:"username"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (Users, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.HashedPass,
		arg.PasswordChangedAt,
		arg.Email,
		arg.FullName,
		arg.IsEmailActivated,
		arg.Username,
	)
	var i Users
	err := row.Scan(
		&i.Username,
		&i.Email,
		&i.HashedPass,
		&i.FullName,
		&i.PasswordChangedAt,
		&i.CreatedAt,
		&i.IsEmailActivated,
	)
	return i, err
}
