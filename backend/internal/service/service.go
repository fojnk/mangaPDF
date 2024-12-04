package service

import (
	"github.com/fojnk/Task-Test-devBack/internal/models"
	"github.com/fojnk/Task-Test-devBack/internal/repository"
)

type Authorization interface {
	CreateUser(user models.User) (string, error)
	GenerateTokens(guid, ip string) (string, string, error)
	Refresh(accessToken, refreshToken, ip string) (string, string, error)
	GetUserByGuid(guid string) (models.User, error)
}

type Service struct {
	Authorization
}

func NewService(repos *repository.Respository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
	}
}
