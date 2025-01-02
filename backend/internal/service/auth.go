package service

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"

	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fojnk/Task-Test-devBack/internal/models"
	"github.com/fojnk/Task-Test-devBack/internal/repository"
	"golang.org/x/exp/rand"
)

const (
	key  = "jfaopajsfojadsf"
	salt = "fkdsajl3214u98ujkj"
)

type AuthService struct {
	repo repository.Authorization
}

type tokenClaims struct {
	jwt.StandardClaims
	Id  int
	Key string
}

func NewAuthService(repos repository.Authorization) *AuthService {
	return &AuthService{repo: repos}
}

func (a *AuthService) CreateUser(user models.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return a.repo.CreateUser(user)
}

func (a *AuthService) GenerateTokens(user_id int) (string, string, int64, int64, error) {
	user, err := a.repo.GetUser(user_id)
	if err != nil {
		return "", "", 0, 0, err
	}

	pair_key, _ := generateRandSeq()

	accessToken, expA, err := a.newJWT(user_id, pair_key, 12*time.Hour)
	if err != nil {
		return "", "", 0, 0, err
	}

	refreshToken, expR, err := a.newJWT(user_id, pair_key, 1000*time.Hour)
	if err != nil {
		return "", "", 0, 0, err
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(refreshToken))

	newSession := models.Session{
		RefreshToken: encoded,
		Fingerprint:  "",
		Ip:           "",
	}

	if _, err := a.repo.CreateSession(user.Id, newSession); err != nil {
		return "", "", 0, 0, err
	}

	return accessToken, encoded, expA, expR, err
}

func (a *AuthService) newJWT(user_id int, pair_key string, expTime time.Duration) (string, int64, error) {
	exp := time.Now().Add(expTime).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: exp,
			IssuedAt:  time.Now().Unix(),
		},
		user_id,
		pair_key,
	})

	retr, err := token.SignedString([]byte(key))

	return retr, exp, err
}

func (a *AuthService) GetUserById(id int) (models.User, error) {
	return a.repo.GetUser(id)
}

func (a *AuthService) GetUserByUsername(username, password string) (models.User, error) {
	password_hash := generatePasswordHash(password)
	return a.repo.GetUserByUsername(username, password_hash)
}

func (a *AuthService) ParseToken(acessToken string) (int, string, error) {
	token, err := jwt.ParseWithClaims(acessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(key), nil
	})

	if err != nil {
		return 0, "", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, "", errors.New("bad claims format")
	}
	return claims.Id, claims.Key, nil
}

func (s *AuthService) Refresh(accessToken, refreshToken string) (string, string, int64, int64, error) {
	id, pair_key1, err := s.ParseToken(accessToken)
	if err != nil {
		return "", "", 0, 0, err
	}

	session, err := s.repo.GetSession(id, refreshToken)
	if err != nil {
		return "", "", 0, 0, errors.New("unknown refresh token")
	}

	decoded, _ := base64.StdEncoding.DecodeString(refreshToken)

	_, pair_key2, err := s.ParseToken(string(decoded))
	if err != nil {
		return "", "", 0, 0, err
	}

	if pair_key1 != pair_key2 {
		return "", "", 0, 0, errors.New("incorrect token pair")
	}

	// if ip != lastIp {
	// 	if err := email.SendHtmlEmail([]string{user.Email}, "Авторизация с другого IP",
	// 		"Если это вы авторизовались недавно - просто проигнорируйте это сообщение."); err != nil {
	// 		fmt.Println(err.Error())
	// 	}

	// 	logrus.Printf("send warning email to user")
	// }

	s.repo.ClearSession(session.Id)

	return s.GenerateTokens(id)
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func generateRandSeq() (string, error) {
	b := make([]byte, 32)

	randVal := rand.NewSource(uint64(time.Now().Unix()))
	r := rand.New(randVal)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
