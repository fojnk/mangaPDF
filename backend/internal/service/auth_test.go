package service

import (
	"testing"
	"time"

	"github.com/fojnk/Task-Test-devBack/internal/repository/mocks"
	"github.com/golang/mock/gomock"
)

func TestTokenHash(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	myAuthService := NewAuthService(mocks.NewMockAuthorization(mockCtrl))

	guid, _ := generateRandSeq()
	pair_key, _ := generateRandSeq()

	token, _ := myAuthService.newJWT(guid, "123.123.123.123", pair_key, 12*time.Hour)

	hash, _ := hashRefreshToken(token)

	if !checkEqHash(hash, token) {
		t.Errorf("eq test failed")
	}

}

func TestCorrectPairOfTokens(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	myAuthService := NewAuthService(mocks.NewMockAuthorization(mockCtrl))

	guid, _ := generateRandSeq()
	pair_key, _ := generateRandSeq()

	accessToken, _ := myAuthService.newJWT(guid, "123.123.123.123", pair_key, 12*time.Hour)

	refreshToken, _ := myAuthService.newJWT(guid, "123.123.123.123", pair_key, 100*time.Hour)

	_, _, key1, err := myAuthService.parseToken(accessToken)
	if err != nil {
		t.Errorf("parse token faliled")
	}

	_, _, key2, err := myAuthService.parseToken(refreshToken)
	if err != nil {
		t.Errorf("parse token faliled")
	}

	if key1 != key2 {
		t.Errorf("newJWT faliled")
	}

}

func TestIncorrectPairOfTokens(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	myAuthService := NewAuthService(mocks.NewMockAuthorization(mockCtrl))

	guid, _ := generateRandSeq()
	pair_key, _ := generateRandSeq()

	accessToken, _ := myAuthService.newJWT(guid, "123.123.123.123", "12345", 12*time.Hour)

	refreshToken, _ := myAuthService.newJWT(guid, "123.123.123.123", pair_key, 100*time.Hour)

	_, _, key1, err := myAuthService.parseToken(accessToken)
	if err != nil {
		t.Errorf("parse token faliled")
	}

	_, _, key2, err := myAuthService.parseToken(refreshToken)
	if err != nil {
		t.Errorf("parse token faliled")
	}

	if key1 == key2 {
		t.Errorf("newJWT faliled")
	}

}
