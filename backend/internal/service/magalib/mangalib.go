package mangalib

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/cavaliergopher/grab/v3"
	"github.com/fojnk/Task-Test-devBack/configs"
	"github.com/fojnk/Task-Test-devBack/internal/models"
	"github.com/fojnk/Task-Test-devBack/internal/service/history"
	"github.com/fojnk/Task-Test-devBack/internal/service/pdf"
	"github.com/fojnk/Task-Test-devBack/pkg/tools"
)

type ChaptersRawData struct {
	Chapters struct {
		List []struct {
			ChapterName   string `json:"chapter_name"`
			ChapterNumber string `json:"chapter_number"`
			ChapterVolume int    `json:"chapter_volume"`
		} `json:"list"`
	} `json:"chapters"`
}

type PagesList []struct {
	Page int    `json:"p"`
	URL  string `json:"u"`
}

type Info struct {
	Img struct {
		URL    string `json:"url"`
		Server string `json:"server"`
	} `json:"img"`
	Servers struct {
		Main      string `json:"main"`
		Secondary string `json:"secondary"`
		Compress  string `json:"compress"`
		Fourth    string `json:"fourth"`
	} `json:"servers"`
}

func GetMangaInfo(mangaURL string) (models.MangaInfo, error) {
	var err error
	var mangaInfo models.MangaInfo

	body, err := tools.GetPageCF(mangaURL)
	if err != nil {
		return mangaInfo, err
	}

	chaptersPage, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return mangaInfo, err
	}

	origTitle := chaptersPage.Find(".media-name__alt").Text()

	if origTitle == "" {
		origTitle = chaptersPage.Find(".media-name__main").Text()
	}

	mangaInfo.TitleOrig = origTitle
	mangaInfo.TitleRu = chaptersPage.Find(".media-name__main").Text()

	return mangaInfo, nil
}

func GetChaptersList(mangaURL string) ([]models.ChaptersList, error) {
	var err error
	var chaptersList []models.ChaptersList
	dataRE := regexp.MustCompile(`(?i)window.__DATA__ = {.+};`)

	body, err := tools.GetPageCF(mangaURL)
	if err != nil {
		return chaptersList, err
	}

	buf := new(strings.Builder)
	_, err = io.Copy(buf, body)
	if err != nil {
		return chaptersList, err
	}

	rawData := strings.Trim(dataRE.FindString(buf.String()), "window.__DATA__ = ;")

	chaptersRawData := ChaptersRawData{}
	err = json.Unmarshal([]byte(rawData), &chaptersRawData)
	if err != nil {
		return chaptersList, err
	}

	for _, ch := range chaptersRawData.Chapters.List {
		chapter := models.ChaptersList{
			Title: "(" + strconv.Itoa(ch.ChapterVolume) + "-" + ch.ChapterNumber + ") " + ch.ChapterName,
			Path:  "v" + strconv.Itoa(ch.ChapterVolume) + "/c" + ch.ChapterNumber,
		}

		chaptersList = append(chaptersList, chapter)
	}

	return tools.ReverseList(chaptersList), nil
}

func DownloadManga(downData models.DownloadOpts) error {
	var err error
	var chaptersList []models.ChaptersList
	var saveChapters []string
	savedFilesByVol := make(map[string][]string)

	switch downData.Type {
	case "all":
		chaptersList, err = GetChaptersList(downData.MangaURL)
		if err != nil {
			slog.Error(
				"Ошибка при получении списка глав",
				slog.String("Message", err.Error()),
			)
			return err
		}
		time.Sleep(1 * time.Second)
	case "chapters":
		chaptersRaw := strings.Split(strings.Trim(downData.Chapters, "[] \""), "\",\"")
		for _, ch := range chaptersRaw {
			chapter := models.ChaptersList{
				Path: ch,
			}
			chaptersList = append(chaptersList, chapter)
		}
	}

	chaptersTotal := len(chaptersList)
	chaptersCur := 0

	models.WSChan <- models.WSData{
		Cmd: "initProgress",
		Payload: map[string]interface{}{
			"valNow": 0,
			"valMax": chaptersTotal,
			"width":  0,
		},
	}

	for _, chapter := range chaptersList {
		volume := strings.Split(chapter.Path, "/")[0]

		chSavedFiles, err := DownloadChapter(downData, chapter)
		if err != nil {
			models.WSChan <- models.WSData{
				Cmd: "updateLog",
				Payload: map[string]interface{}{
					"type": "err",
					"text": "-- Ошибка при скачивании главы:" + err.Error(),
				},
			}
		}

		savedFilesByVol[volume] = append(savedFilesByVol[volume], chSavedFiles...)

		chaptersCur++

		saveChapters = append(saveChapters, chapter.Path)

		time.Sleep(time.Duration(configs.Cfg.Mangalib.TimeoutChapter) * time.Millisecond)

		models.WSChan <- models.WSData{
			Cmd: "updateProgress",
			Payload: map[string]interface{}{
				"valNow": chaptersCur,
				"width":  tools.GetPercent(chaptersCur, chaptersTotal),
			},
		}
	}

	chapterPath := path.Join(configs.Cfg.Savepath, downData.SavePath)

	if downData.PDFvol == "1" {
		models.WSChan <- models.WSData{
			Cmd: "updateLog",
			Payload: map[string]interface{}{
				"type": "std",
				"text": "Создаю PDF для томов",
			},
		}

		pdf.CreateVolPDF(chapterPath, savedFilesByVol, downData.Del)
	}

	if downData.PDFall == "1" {
		models.WSChan <- models.WSData{
			Cmd: "updateLog",
			Payload: map[string]interface{}{
				"type": "std",
				"text": "Создаю PDF для манги",
			},
		}

		pdf.CreateMangaPdf(chapterPath, savedFilesByVol, downData.Del)
	}

	mangaID := tools.GetMD5(downData.MangaURL)
	history.SaveHistory(mangaID, saveChapters)

	models.WSChan <- models.WSData{
		Cmd: "downloadComplete",
		Payload: map[string]interface{}{
			"text": "Скачивание завершено!",
		},
	}

	return nil
}

func DownloadChapter(downData models.DownloadOpts, curChapter models.ChaptersList) ([]string, error) {
	var err error

	models.WSChan <- models.WSData{
		Cmd: "updateLog",
		Payload: map[string]interface{}{
			"type": "std",
			"text": "Скачиваю главу: " + curChapter.Path,
		},
	}

	infoRE := regexp.MustCompile(`(?i)window.__info = {.+};`)
	pagesRE := regexp.MustCompile(`(?i)window.__pg = \[{.+}\];`)

	chapterURL := strings.TrimRight(downData.MangaURL, "/") + "/" + curChapter.Path

	body, err := tools.GetPageCF(chapterURL)
	if err != nil {
		return nil, err
	}

	buf := new(strings.Builder)
	_, err = io.Copy(buf, body)
	if err != nil {
		return nil, err
	}

	rawInfo := strings.Trim(infoRE.FindString(buf.String()), "window.__info = ;")
	rawPages := strings.Trim(pagesRE.FindString(buf.String()), "window.__pg = ;")

	info := Info{}
	pages := PagesList{}

	err = json.Unmarshal([]byte(rawInfo), &info)
	if err != nil {
		slog.Error(
			"Ошибка при распаковке данных (Info)",
			slog.String("Message", err.Error()),
		)
		return nil, err
	}

	err = json.Unmarshal([]byte(rawPages), &pages)
	if err != nil {
		slog.Error(
			"Ошибка при распаковке данных (Pages)",
			slog.String("Message", err.Error()),
		)
		return nil, err
	}

	chapterPath := path.Join(configs.Cfg.Savepath, downData.SavePath, curChapter.Path)

	if _, err := os.Stat(chapterPath); os.IsNotExist(err) {
		os.MkdirAll(chapterPath, 0755)
	}

	serversList := make(map[string]string)
	serversList["compress"] = info.Servers.Compress
	serversList["main"] = info.Servers.Main
	serversList["fourth"] = info.Servers.Fourth
	serversList["secondary"] = info.Servers.Secondary

	var savedFiles []string

	for _, page := range pages {
		isFail := false

		for _, s := range serversList {
			if s == "" {
				continue
			}

			imgURL := s + info.Img.URL + page.URL

			client := grab.NewClient()
			client.UserAgent = configs.Cfg.UserAgent
			req, err := grab.NewRequest(chapterPath, imgURL)
			req.HTTPRequest.Header.Set("Referer", chapterURL)
			if err != nil {
				slog.Error(
					"Ошибка при создании запроса страницы",
					slog.String("Message", err.Error()),
				)
				isFail = true
				continue
			}

			resp := client.Do(req)
			if resp.Err() != nil {
				slog.Error(
					"Ошибка при скачивании страницы",
					slog.String("Message", resp.Err().Error()),
				)
				isFail = true
				continue
			}

			savedFiles = append(savedFiles, resp.Filename)
			isFail = false
			break
		}

		if isFail {
			models.WSChan <- models.WSData{
				Cmd: "updateLog",
				Payload: map[string]interface{}{
					"type": "err",
					"text": "-- Ошибка при скачивании страницы:" + page.URL,
				},
			}
		}

		time.Sleep(time.Duration(configs.Cfg.Mangalib.TimeoutImage) * time.Millisecond)
	}

	if downData.CBZ == "1" {
		models.WSChan <- models.WSData{
			Cmd: "updateLog",
			Payload: map[string]interface{}{
				"type": "std",
				"text": "- Создаю CBZ для главы",
			},
		}

		tools.CreateCBZ(chapterPath)
	}

	if downData.PDFch == "1" {
		models.WSChan <- models.WSData{
			Cmd: "updateLog",
			Payload: map[string]interface{}{
				"type": "std",
				"text": "- Создаю PDF для главы",
			},
		}

		pdf.CreatePDF(chapterPath, savedFiles)
	}

	if downData.PDFvol != "1" && downData.PDFall != "1" && downData.Del == "1" {
		err := os.RemoveAll(chapterPath)
		if err != nil {
			slog.Error(
				"Ошибка при удалении файлов",
				slog.String("Message", err.Error()),
			)
		}
	}

	return savedFiles, nil
}
