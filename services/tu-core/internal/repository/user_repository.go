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
	SchoolName string // <-- Tambahan untuk nama sekolah
}

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (*User, error)
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
		       COALESCE(u.sekolah_id::text, ''), COALESCE(s.nama, '')
		FROM user_login u
		LEFT JOIN tabel_sekolah s ON u.sekolah_id = s.sekolah_id
		WHERE u.username = $1 AND u.is_active = true
	`
	var u User
	err := r.db.QueryRow(ctx, query, username).
		Scan(&u.UserID, &u.Password, &u.Name, &u.Role, &u.SchoolID, &u.SchoolName)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
