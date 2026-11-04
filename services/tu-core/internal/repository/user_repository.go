package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	UserID     string
	Password   string
	Name       string
	Role       string
	SchoolID   string
	SchoolName string
	IsActive   bool
}

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*User, error)
	ListUsersBySchool(ctx context.Context, schoolID string) ([]User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUserRole(ctx context.Context, userID, role string) error
	ToggleUserStatus(ctx context.Context, userID string, isActive bool) error
	UpdatePassword(ctx context.Context, userID, passwordHash string) error // <-- BARU
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*User, error) {
	query := `
		SELECT u.user_id, u.password_hash, u.username, u.role, 
		       COALESCE(u.sekolah_id::text, ''), COALESCE(s.nama, ''), u.is_active
		FROM user_login u
		LEFT JOIN tabel_sekolah s ON u.sekolah_id = s.sekolah_id
		WHERE u.username = $1 AND u.is_active = true
	`
	var u User
	err := r.db.QueryRow(ctx, query, username).
		Scan(&u.UserID, &u.Password, &u.Name, &u.Role, &u.SchoolID, &u.SchoolName, &u.IsActive)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) ListUsersBySchool(ctx context.Context, schoolID string) ([]User, error) {
	query := `
		SELECT user_id, username, role, COALESCE(sekolah_id::text, ''), is_active
		FROM user_login
		WHERE sekolah_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.db.Query(ctx, query, schoolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.UserID, &u.Name, &u.Role, &u.SchoolID, &u.IsActive); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO user_login (user_id, username, password_hash, role, is_active, sekolah_id)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, query, user.UserID, user.Name, user.Password, user.Role, user.IsActive, user.SchoolID)
	return err
}

func (r *userRepository) UpdateUserRole(ctx context.Context, userID, role string) error {
	query := `UPDATE user_login SET role = $1, updated_at = NOW() WHERE user_id = $2`
	_, err := r.db.Exec(ctx, query, role, userID)
	return err
}

func (r *userRepository) ToggleUserStatus(ctx context.Context, userID string, isActive bool) error {
	query := `UPDATE user_login SET is_active = $1, updated_at = NOW() WHERE user_id = $2`
	_, err := r.db.Exec(ctx, query, isActive, userID)
	return err
}

func (r *userRepository) UpdatePassword(ctx context.Context, userID, passwordHash string) error {
	query := `UPDATE user_login SET password_hash = $1, updated_at = NOW() WHERE user_id = $2`
	_, err := r.db.Exec(ctx, query, passwordHash, userID)
	return err
}
