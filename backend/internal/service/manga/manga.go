package manga

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/fojnk/Task-Test-devBack/internal/models"
	"github.com/fojnk/Task-Test-devBack/internal/service/manga/mangalib"
)

type MangaService struct {
}

func NewMangaService() *MangaService {
	return &MangaService{}
}

func (m *MangaService) GetMangaList() (bytes.Buffer, error) {
	return mangalib.GetMangaList()
}

func (m *MangaService) GetChaptersList(mangaUrl string) (string, error) {
	var err error
	var isMtr bool
	var userHash string
	var rawChaptersList mangalib.ChaptersRawData
	var transList []models.RMTranslators

	rawChaptersList, err = mangalib.GetChaptersListFromApi(mangaUrl)
	if err != nil {
		slog.Error(
			"Ошибка при получении списка глав",
			slog.String("Message", err.Error()),
		)

		return "", errors.New("при получении списка глав произошла ошибка. Подробности в лог-файле")
	}

	resp := make(map[string]interface{})

	if len(rawChaptersList.List) > 0 {
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

// func (m *MangaService) DownloadManga(downloadOpts models.DownloadOpts, mangaUrl string) string {
// 	url, _ := urlx.Parse(mangaUrl)
// 	host, _, _ := urlx.SplitHostPort(url)

// 	downloadOpts.MangaURL = strings.Split(url.String(), "?")[0]
// 	downloadOpts.SavePath = strings.Trim(url.Path, "/")

// 	if tools.CheckSource(configs.Cfg.CurrentURLs.MangaLib, host) {
// 		go mangalib.DownloadManga(downloadOpts)
// 	}

// 	resp := make(map[string]interface{})
// 	resp["status"] = "OK"

// 	respData, _ := json.Marshal(resp)

// 	return string(respData)
// }
