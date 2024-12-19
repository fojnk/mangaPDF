package repository

import (
	"fmt"

	"github.com/fojnk/Task-Test-devBack/internal/models"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (a *AuthPostgres) CreateUser(user models.User) (int, error) {
	var id int
	query := fmt.Sprintf(`
		INSERT INTO %s 
		(username, email, password_hash, wallet, role, subscription, end_of_sub) 
		values ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id`, usersTable)
	row := a.db.QueryRow(query, user.Username, user.Email, user.Password,
		user.Wallet, user.Role, user.Subscription, user.EndOfSub)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (a *AuthPostgres) GetUser(id int) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", usersTable)

	err := a.db.Get(&user, query, id)
	return user, err
}

func (a *AuthPostgres) GetUserByUsername(username, password string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE username=$1 AND password_hash=$2", usersTable)

	err := a.db.Get(&user, query, username, password)
	return user, err
}

func (a *AuthPostgres) GetUserSessions(id int) ([]models.Session, error) {
	var sessions []models.Session
	query := fmt.Sprintf(`
		SELECT id, user_id, refresh_token, fingerprint, ip
		FROM %s 
		WHERE user_id=$1`, userSessions)

	err := a.db.Select(&sessions, query, id)
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func (a *AuthPostgres) GetSession(user_id int, refresh_token string) (models.Session, error) {
	var session models.Session
	query := fmt.Sprintf(`
		SELECT id, user_id, refresh_token, fingerprint, ip
		FROM %s 
		WHERE user_id=$1 AND refresh_token=$2`, userSessions)

	err := a.db.Get(&session, query, user_id, refresh_token)
	return session, err
}

func (a *AuthPostgres) CreateSession(user_id int, session models.Session) (int, error) {
	var id int
	query := fmt.Sprintf(`
		INSERT INTO %s 
		(user_id, refresh_token, fingerprint, ip) 
		values ($1, $2, $3, $4) 
		RETURNING id`, userSessions)
	row := a.db.QueryRow(query, user_id, session.RefreshToken, session.Fingerprint, session.Ip)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (a *AuthPostgres) ClearSession(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", userSessions)
	_, err := a.db.Exec(query, id)
	return err
}
