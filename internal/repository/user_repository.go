package repository

import (
	"LiveChat/internal/model"
	"context"
	"database/sql"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type repository struct {
	db DBTX
}

func NewRepository(db DBTX) model.Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	var lastInsertID int64
	query := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id"

	err := r.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password).Scan(&lastInsertID)
	if err != nil {
		return &model.User{}, err
	}

	user.ID = lastInsertID
	return user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	u := model.User{}

	query := "SELECT * FROM users WHERE email=$1"

	err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Username, &u.Email, &u.Password)
	if err != nil {
		return &model.User{}, err
	}

	return &u, nil

}
