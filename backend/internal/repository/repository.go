package repository

import (
	"github.com/fojnk/Task-Test-devBack/internal/models"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(id int) (models.User, error)
	GetUserSessions(id int) ([]models.Session, error)
	CreateSession(user_id int, session models.Session) (int, error)
	GetUserByUsername(username, password string) (models.User, error)
	ClearSession(id int) error
}

type Respository struct {
	Authorization
}

func NewRepository(db *sqlx.DB) *Respository {
	return &Respository{
		Authorization: NewAuthPostgres(db),
	}
}
