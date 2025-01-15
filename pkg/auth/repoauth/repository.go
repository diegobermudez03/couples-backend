package repoauth

import (
	"time"

	"github.com/diegobermudez03/couples-backend/pkg/auth/domainauth"
	"github.com/google/uuid"
)

type AuthPostgresRepo struct {
}

func NewAuthPostgresRepo() domainauth.AuthRepository{
	return &AuthPostgresRepo{}
}

func (r *AuthPostgresRepo) CreateUserAuth(id uuid.UUID, email string, hash string) error{
	return nil
}

func (r *AuthPostgresRepo) CreateSession(id uuid.UUID, userId uuid.UUID, token string, device string, os string, expiresAt time.Time) error{
	return nil
}