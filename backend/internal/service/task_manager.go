package service

import (
	"sync"

	"github.com/fojnk/Task-Test-devBack/internal/models"
	"github.com/fojnk/Task-Test-devBack/internal/service/manga"
	"github.com/google/uuid"
)

type TaskStatus struct {
	Status string // "pending", "processing", "ready", "error"
	File   string // путь к сгенерированному файлу
	Error  string
}

type Manager struct {
	tasks        sync.Map
	mangaService *manga.MangaService
}

func NewTastManager(mangaService *manga.MangaService) *Manager {
	return &Manager{
		mangaService: mangaService,
	}
}

func (m *Manager) CreateTask(input models.DownloadOpts, mangaName string) string {
	taskID := uuid.New().String()
	m.tasks.Store(taskID, TaskStatus{Status: "processing"})

	go func() {
		out, err := m.mangaService.DownloadManga(input, mangaName)
		if err != nil {
			m.tasks.Store(taskID, TaskStatus{Status: "error", Error: err.Error()})
			return
		}

		m.tasks.Store(taskID, TaskStatus{Status: "ready", File: out})
	}()

	return taskID
}

func (m *Manager) GetStatus(taskID string) (TaskStatus, bool) {
	val, ok := m.tasks.Load(taskID)
	if !ok {
		return TaskStatus{}, false
	}
	return val.(TaskStatus), true
}
