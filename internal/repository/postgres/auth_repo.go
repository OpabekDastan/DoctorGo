package postgres

import (
	"log"
	"strings"

	"github.com/jmoiron/sqlx"

	"doctor_go/internal/model"
)

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	query := `
		SELECT id, email, password_hash, role, first_name, last_name, is_active, created_at
		FROM users
		WHERE LOWER(email) = LOWER($1)
	`

	email = strings.TrimSpace(email)
	log.Printf("AUTH REPO: lookup email=%s", email)

	if err := r.db.Get(&user, query, email); err != nil {
		log.Printf("AUTH REPO: lookup error = %v", err)
		return nil, err
	}

	log.Printf("AUTH REPO: found user_id=%d email=%s", user.ID, user.Email)
	return &user, nil
}

func (r *AuthRepository) GetUserByID(userID int64) (*model.User, error) {
	var user model.User
	query := `
		SELECT id, email, password_hash, role, first_name, last_name, is_active, created_at
		FROM users
		WHERE id = $1
	`

	if err := r.db.Get(&user, query, userID); err != nil {
		log.Printf("AUTH REPO: get by id error = %v", err)
		return nil, err
	}

	return &user, nil
}
