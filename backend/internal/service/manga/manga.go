package manga

import (
	"encoding/json"
	"errors"
	"log/slog"
	"strings"

	"github.com/fojnk/Task-Test-devBack/configs"
	"github.com/fojnk/Task-Test-devBack/internal/models"
	"github.com/fojnk/Task-Test-devBack/internal/service/manga/readmanga"
	"github.com/fojnk/Task-Test-devBack/pkg/tools"
	"github.com/goware/urlx"
)

type MangaService struct {
}

func NewMangaService() *MangaService {
	return &MangaService{}
}

func (m *MangaService) GetMangaList() ([]models.Manga, error) {
	return readmanga.GetMangaList()
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

func (m *MangaService) DownloadManga(downloadOpts models.DownloadOpts, mangaUrl string) string {
	url, _ := urlx.Parse(mangaUrl)
	host, _, _ := urlx.SplitHostPort(url)

	downloadOpts.MangaURL = strings.Split(url.String(), "?")[0]
	downloadOpts.SavePath = strings.Trim(url.Path, "/")

	if tools.CheckSource(configs.Cfg.CurrentURLs.MangaLib, host) {
		go readmanga.DownloadManga(downloadOpts)
	}

	resp := make(map[string]interface{})
	resp["status"] = "OK"

	respData, _ := json.Marshal(resp)

	return string(respData)
}
