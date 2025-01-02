package service

import (
	"github.com/fojnk/Task-Test-devBack/internal/models"
	"github.com/fojnk/Task-Test-devBack/internal/repository"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GenerateTokens(user_id int) (string, string, int64, int64, error)
	Refresh(accessToken, refreshToken string) (string, string, int64, int64, error)
	GetUserById(id int) (models.User, error)
	GetUserByUsername(username, password string) (models.User, error)
	ParseToken(string) (int, string, error)
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Respository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
