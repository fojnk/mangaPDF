package repository

import (
	"fmt"

	"github.com/fojnk/Task-Test-devBack/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (a *AuthPostgres) CreateUser(user models.User) (string, error) {
	var id string
	query := fmt.Sprintf("INSERT INTO %s (guid, username, email, password_hash) values ($1, $2, $3, $4) RETURNING guid", usersTable)
	row := a.db.QueryRow(query, user.Guid, user.Username, user.Email, user.Password)

	if err := row.Scan(&id); err != nil {
		return "", err
	}

	return id, nil
}

func (a *AuthPostgres) GetUser(guid string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE guid=$1", usersTable)

	err := a.db.Get(&user, query, guid)
	return user, err
}

func (a *AuthPostgres) GetUserTokens(guid string) ([]models.Token, error) {
	var tokens []models.Token
	query := fmt.Sprintf("SELECT id, user_id, token_hash FROM %s WHERE user_id=$1", userTokens)

	err := a.db.Select(&tokens, query, guid)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (a *AuthPostgres) SaveRefreshToken(guid, token_hash string) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (user_id, token_hash) values ($1, $2) RETURNING id", userTokens)
	row := a.db.QueryRow(query, guid, token_hash)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	logrus.Info(query)

	return id, nil
}

func (a *AuthPostgres) RemoveToken(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", userTokens)
	_, err := a.db.Exec(query, id)
	return err
}
