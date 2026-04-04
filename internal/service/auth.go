package service

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"doctor_go/internal/model"
	"doctor_go/internal/repository/postgres"
	jwtpkg "doctor_go/pkg/jwt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInactiveUser       = errors.New("user is inactive")
)

type AuthService struct {
	repo      *postgres.AuthRepository
	jwtSecret string
	ttlMin    int
}

func NewAuthService(repo *postgres.AuthRepository, jwtSecret string, ttlMin int) *AuthService {
	return &AuthService{repo: repo, jwtSecret: jwtSecret, ttlMin: ttlMin}
}

type LoginResult struct {
	AccessToken string     `json:"access_token"`
	User        model.User `json:"user"`
}

func (s *AuthService) Login(email, password string) (*LoginResult, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	password = strings.TrimSpace(password)
	log.Printf("AUTH SERVICE: login called email=%s", email)

	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("AUTH SERVICE: user not found for email=%s", email)
			return nil, ErrInvalidCredentials
		}
		log.Printf("AUTH SERVICE: repo error = %v", err)
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		log.Printf("AUTH SERVICE: password mismatch for user_id=%d email=%s", user.ID, user.Email)
		return nil, ErrInvalidCredentials
	}

	log.Printf("AUTH SERVICE: password verified user_id=%d role=%s active=%v", user.ID, user.Role, user.IsActive)

	if !user.IsActive {
		log.Printf("AUTH SERVICE: user inactive")
		return nil, ErrInactiveUser
	}

	token, err := jwtpkg.Generate(s.jwtSecret, user.ID, user.Role, user.Email, s.ttlMin)
	if err != nil {
		log.Printf("AUTH SERVICE: jwt generate error = %v", err)
		return nil, fmt.Errorf("generate jwt: %w", err)
	}

	log.Printf("AUTH SERVICE: token generated for user_id=%d", user.ID)

	return &LoginResult{
		AccessToken: token,
		User:        *user,
	}, nil
}

func (s *AuthService) Me(userID int64) (*model.User, error) {
	return s.repo.GetUserByID(userID)
}
