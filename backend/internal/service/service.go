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
	DownloadManga(downloadOpts models.DownloadOpts, mangaUrl string) (string, error)
	GetMangaList(offset int) ([]models.Manga, error)
}

type TaskService interface {
	CreateTask(input models.DownloadOpts, mangaName string) string
	GetStatus(taskID string) (TaskStatus, bool)
}

type Service struct {
	Authorization
	MangaService
	TaskService
}

func NewService(repos *repository.Respository) *Service {
	mangaService := manga.NewMangaService()
	return &Service{
		Authorization: auth.NewAuthService(repos.Authorization),
		MangaService:  mangaService,
		TaskService:   NewTastManager(mangaService),
	}
}
