package manga

import (
	"encoding/json"
	"errors"
	"log/slog"
	"strings"

	"github.com/fojnk/Task-Test-devBack/internal/models"
	"github.com/fojnk/Task-Test-devBack/internal/service/manga/readmanga"
)

type MangaService struct {
}

func NewMangaService() *MangaService {
	return &MangaService{}
}

func (m *MangaService) GetMangaList(offset int) ([]models.Manga, error) {
	return readmanga.GetMangaList(offset)
}

func (m *MangaService) GetChaptersList(mangaUrl string) (string, error) {
	var err error
	var isMtr bool
	var userHash string
	var rawChaptersList []models.ChaptersList
	chaptersList := make(map[string][]models.ChaptersList)
	var transList []models.RMTranslators

	rawChaptersList, transList, isMtr, userHash, err = readmanga.GetChaptersList(mangaUrl)
	if err != nil {
		slog.Error(
			"Ошибка при получении списка глав",
			slog.String("Message", err.Error()),
		)

		return "", errors.New("при получении списка глав произошла ошибка. Подробности в лог-файле")
	}

	for _, ch := range rawChaptersList {
		parts := strings.Split(ch.Path, "/")
		volNum := strings.TrimLeft(parts[0], "vol")
		chaptersList[volNum] = append(chaptersList[volNum], ch)
	}

	resp := make(map[string]interface{})

	if len(rawChaptersList) > 0 {
		resp["status"] = "success"
		resp["is_mtr"] = isMtr
		resp["user_hash"] = userHash
		resp["payload"] = rawChaptersList
		resp["translators"] = transList
	} else {
		resp["status"] = "error"
		resp["errtext"] = "Глав не найдено. Проверьте правильность ввода адреса манги."
	}

	respData, _ := json.Marshal(resp)

	return string(respData), nil
}

func (m *MangaService) DownloadManga(downloadOpts models.DownloadOpts, mangaUrl string) (string, error) {
	downloadOpts.MangaURL = mangaUrl
	downloadOpts.SavePath = mangaUrl

	return readmanga.DownloadManga(downloadOpts)
}
