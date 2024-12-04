package repository

import (
	"github.com/fojnk/Task-Test-devBack/internal/models"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user models.User) (string, error)
	GetUser(guid string) (models.User, error)
	GetUserTokens(guid string) ([]models.Token, error)
	SaveRefreshToken(guid, token_hash string) (int, error)
	RemoveToken(id int) error
}

type Respository struct {
	Authorization
}

func NewRepository(db *sqlx.DB) *Respository {
	return &Respository{
		Authorization: NewAuthPostgres(db),
	}
}
