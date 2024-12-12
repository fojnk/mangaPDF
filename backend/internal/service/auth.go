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
	"github.com/fojnk/Task-Test-devBack/pkg/email"
	"github.com/sirupsen/logrus"
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
	Ip  string
	Key string
}

func NewAuthService(repos repository.Authorization) *AuthService {
	return &AuthService{repo: repos}
}

func (a *AuthService) CreateUser(user models.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return a.repo.CreateUser(user)
}

func (a *AuthService) GenerateTokens(user_id int, ip string) (string, string, error) {
	user, err := a.repo.GetUser(user_id)
	if err != nil {
		return "", "", err
	}

	pair_key, _ := generateRandSeq()

	accessToken, err := a.newJWT(user_id, ip, pair_key, 12*time.Hour)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := a.newJWT(user_id, ip, pair_key, 1000*time.Hour)
	if err != nil {
		return "", "", err
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(refreshToken))

	newSession := models.Session{
		RefreshToken: encoded,
		Fingerprint:  "",
		Ip:           ip,
	}

	if _, err := a.repo.CreateSession(user.Id, newSession); err != nil {
		return "", "", err
	}

	return accessToken, encoded, err
}

func (a *AuthService) newJWT(user_id int, ip, pair_key string, expTime time.Duration) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expTime).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user_id,
		ip,
		pair_key,
	})

	return token.SignedString([]byte(key))
}

func (a *AuthService) GetUserById(id int) (models.User, error) {
	return a.repo.GetUser(id)
}

func (a *AuthService) GetUserByUsername(username, password string) (models.User, error) {
	return a.repo.GetUserByUsername(username, password)
}

func (a *AuthService) parseToken(acessToken string) (int, string, string, error) {
	token, err := jwt.ParseWithClaims(acessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(key), nil
	})

	if err != nil {
		return 0, "", "", err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, "", "", errors.New("bad claims format")
	}
	return claims.Id, claims.Ip, claims.Key, nil
}

func (s *AuthService) Refresh(accessToken, refreshToken, ip string) (string, string, error) {
	id, lastIp, pair_key1, err := s.parseToken(accessToken)
	if err != nil {
		return "", "", err
	}

	user, err := s.repo.GetUser(id)
	if err != nil {
		return "", "", err
	}

	if err != nil {
		return "", "", err
	}

	sessions, err := s.repo.GetUserSessions(id)
	if err != nil {
		return "", "", err
	}

	exists := false
	var tokenId int
	var decoded []byte

	for _, session := range sessions {
		if session.RefreshToken == refreshToken {
			decoded, _ = base64.StdEncoding.DecodeString(refreshToken)
			exists = true
			tokenId = session.Id
			break
		}
	}

	if !exists {
		return "", "", errors.New("unknown refresh token")
	}

	_, _, pair_key2, err := s.parseToken(string(decoded))
	if err != nil {
		return "", "", err
	}

	if pair_key1 != pair_key2 {
		return "", "", errors.New("incorrect token pair")
	}

	if ip != lastIp {
		if err := email.SendHtmlEmail([]string{user.Email}, "Авторизация с другого IP",
			"Если это вы авторизовались недавно - просто проигнорируйте это сообщение."); err != nil {
			fmt.Println(err.Error())
		}

		logrus.Printf("send warning email to user")
	}

	s.repo.ClearSession(tokenId)

	return s.GenerateTokens(id, ip)
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
