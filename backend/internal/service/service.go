package service

import (
	"github.com/fojnk/Task-Test-devBack/internal/models"
	"github.com/fojnk/Task-Test-devBack/internal/repository"
	"github.com/fojnk/Task-Test-devBack/internal/service/auth"
	"github.com/fojnk/Task-Test-devBack/internal/service/manga"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GenerateTokens(user_id int) (string, string, int64, int64, error)
	Refresh(accessToken, refreshToken string) (string, string, int64, int64, error)
	GetUserById(id int) (models.User, error)
	GetUserByUsername(username, password string) (models.User, error)
	ParseToken(string) (int, string, error)
}

type MangaService interface {
	GetChaptersList(mangaUrl string) (string, error)
	DownloadManga(downloadOpts models.DownloadOpts, mangaUrl string) string
	GetMangaList() ([]models.Manga, error)
}

type Service struct {
	Authorization
	MangaService
}

func NewService(repos *repository.Respository) *Service {
	return &Service{
		Authorization: auth.NewAuthService(repos.Authorization),
		MangaService:  manga.NewMangaService(),
	}
}
