package readmanga

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"log/slog"
	"math/rand"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/cavaliergopher/grab/v3"
	"github.com/fojnk/Task-Test-devBack/configs"
	"github.com/fojnk/Task-Test-devBack/internal/models"
	"github.com/fojnk/Task-Test-devBack/internal/service/pdf"
	"github.com/fojnk/Task-Test-devBack/pkg/tools"
	"github.com/goware/urlx"
)

type ServersList []struct {
	Path string `json:"path"`
	Res  bool   `json:"res"`
}

const baseServerUrl = "https://t.readmanga.io"

func GetMangaList() ([]models.Manga, error) {
	var err error

	pageBody, err := tools.GetPageCF("https://t.readmanga.io/list")
	if err != nil {
		return []models.Manga{}, err
	}

	doc, err := goquery.NewDocumentFromReader(pageBody)
	if err != nil {
		log.Fatal(err)
	}

	var mangas []models.Manga

	doc.Find(".tile").Each(func(i int, s *goquery.Selection) {
		title := s.Find("h3 a").AttrOr("title", "")
		rating := s.Find(".compact-rate").AttrOr("title", "")

		id := ""
		s.Find("a").EachWithBreak(func(_ int, a *goquery.Selection) bool {
			href, exists := a.Attr("href")
			if exists && strings.HasPrefix(href, "/") {
				id = strings.TrimPrefix(href, "/")
				return false
			}
			return true
		})

		img := ""
		imgSel := s.Find("img")
		img, _ = imgSel.Attr("data-original")
		if img == "" {
			img, _ = imgSel.Attr("src")
		}

		desc := s.Find(".manga-description").Text()
		desc = strings.TrimSpace(desc)

		if title != "" && id != "" {
			manga := models.Manga{
				Title:       title,
				Id:          id,
				Image:       img,
				Rating:      rating,
				Description: desc,
			}
			mangas = append(mangas, manga)
		}
	})

	return mangas, nil
}

func GetMangaInfo(mangaName string) (models.MangaInfo, error) {
	var err error
	var mangaInfo models.MangaInfo

	pageBody, err := tools.GetPageCF(baseServerUrl + mangaName)
	if err != nil {
		return mangaInfo, err
	}

	chaptersPage, err := goquery.NewDocumentFromReader(pageBody)
	if err != nil {
		return mangaInfo, err
	}

	origTitle := chaptersPage.Find(".original-name").Text()

	if origTitle == "" && chaptersPage.Find(".eng-name").Text() != "" {
		origTitle = chaptersPage.Find(".eng-name").Text()
	}

	if origTitle == "" {
		origTitle = chaptersPage.Find(".name").Text()
	}

	mangaInfo.TitleOrig = origTitle
	mangaInfo.TitleRu = chaptersPage.Find(".name").Text()

	return mangaInfo, nil
}

func GetChaptersList(mangaName string) ([]models.ChaptersList, []models.RMTranslators, bool, string, error) {
	var err error
	var chaptersList []models.ChaptersList
	var transList []models.RMTranslators
	var userHash string

	isMtr := false

	pageBody, err := tools.GetPageCF(baseServerUrl + "/" + mangaName)
	if err != nil {
		return chaptersList, transList, isMtr, "", err
	}

	var buf bytes.Buffer
	tee := io.TeeReader(pageBody, &buf)

	chaptersPage, err := goquery.NewDocumentFromReader(tee)
	if err != nil {
		return chaptersList, transList, isMtr, "", err
	}

	if chaptersPage.Find(".mtr-message").Length() > 0 {
		isMtr = true
	}

	chaptersPage.Find(".chapters a.chapter-link").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		path := strings.Split(strings.Trim(href, "/"), "/")

		chapter := models.ChaptersList{
			Title: strings.Trim(s.Text(), "\n "),
			Path:  path[1] + "/" + strings.Split(path[2], "?")[0],
		}

		chaptersList = append(chaptersList, chapter)
	})

	// MM
	bufStr := new(strings.Builder)
	_, err = io.Copy(bufStr, &buf)
	if err != nil {
		return chaptersList, transList, isMtr, "", err
	}

	dataUH := regexp.MustCompile(`window.user_hash = '.+';`)
	userHash = strings.Trim(dataUH.FindString(bufStr.String()), "window.user_hash = ';")

	// RM
	chaptersPage.Find("#translation > option").Each(func(i int, s *goquery.Selection) {
		trID, _ := s.Attr("value")
		trName := s.Text()

		trans := models.RMTranslators{
			ID:   trID,
			Name: trName,
		}

		transList = append(transList, trans)
	})

	// MM
	chaptersPage.Find(".translator-selection-item").Each(func(i int, s *goquery.Selection) {
		trID, _ := s.Attr("id")
		trID = strings.Trim(trID, "tr-")
		trName := s.Find(".translator-selection-name").Text()

		trans := models.RMTranslators{
			ID:   trID,
			Name: trName,
		}

		transList = append(transList, trans)
	})

	return tools.ReverseList(chaptersList), transList, isMtr, userHash, nil
}

func DownloadManga(downData models.DownloadOpts) (string, error) {
	var err error
	var chaptersList []models.ChaptersList
	var saveChapters []string
	savedFilesByVol := make(map[string][]string)

	switch downData.Type {
	case "all":
		chaptersList, _, _, _, err = GetChaptersList(downData.MangaURL)
		if err != nil {
			slog.Error(
				"Ошибка при получении списка глав",
				slog.String("Message", err.Error()),
			)
			return "", err
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
			if err.Error() == "noauth" {
				models.WSChan <- models.WSData{
					Cmd: "authErr",
					Payload: map[string]interface{}{
						"type": "err",
						"text": "Для скачивания указанной манги необходимо авторизаваться на сайте!",
					},
				}

				return "", nil
			} else {
				models.WSChan <- models.WSData{
					Cmd: "updateLog",
					Payload: map[string]interface{}{
						"type": "err",
						"text": "-- Ошибка при скачивании главы:" + err.Error(),
					},
				}
			}
		}

		savedFilesByVol[volume] = append(savedFilesByVol[volume], chSavedFiles...)

		chaptersCur++

		saveChapters = append(saveChapters, chapter.Path)

		time.Sleep(time.Duration(configs.Cfg.Readmanga.TimeoutChapter) * time.Millisecond)

		models.WSChan <- models.WSData{
			Cmd: "updateProgress",
			Payload: map[string]interface{}{
				"valNow": chaptersCur,
				"width":  tools.GetPercent(chaptersCur, chaptersTotal),
			},
		}
	}

	chapterPath := path.Join(configs.Cfg.Savepath, downData.SavePath)

	// if downData.PDFvol == "1" {
	// 	models.WSChan <- models.WSData{
	// 		Cmd: "updateLog",
	// 		Payload: map[string]interface{}{
	// 			"type": "std",
	// 			"text": "Создаю PDF для томов",
	// 		},
	// 	}

	// 	pdf.CreateVolPDF(chapterPath, savedFilesByVol, downData.Del)
	// }

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

	models.WSChan <- models.WSData{
		Cmd: "downloadComplete",
		Payload: map[string]interface{}{
			"text": "Скачивание завершено!",
		},
	}

	return chapterPath + ".pdf", nil
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

	chapterURL := baseServerUrl + "/" + downData.MangaURL + "/" + curChapter.Path

	var imageLinks []string

	ptOpt := ""
	mtrOpt := ""
	uhOpt := ""

	if downData.Mtr {
		mtrOpt = "?mtr=1"
	}

	if downData.UserHash != "" {
		if mtrOpt != "" {
			uhOpt = "&"
		} else {
			uhOpt = "?"
		}

		uhOpt = uhOpt + "d=" + downData.UserHash
	}

	if downData.PrefTrans != "" {
		if mtrOpt != "" || uhOpt != "" {
			ptOpt = "&"
		} else {
			ptOpt = "?"
		}

		ptOpt = ptOpt + "tran=" + downData.PrefTrans
	}

	page, err := tools.GetPageCF(chapterURL + mtrOpt + uhOpt + ptOpt)
	if err != nil {
		slog.Error(
			"Ошибка при получении страниц",
			slog.String("Message", err.Error()),
		)
		models.WSChan <- models.WSData{
			Cmd: "updateLog",
			Payload: map[string]interface{}{
				"type": "err",
				"text": "-- Ошибка при получении страниц:" + err.Error(),
			},
		}
		return nil, err
	}

	pageBody, err := io.ReadAll(page)
	if err != nil {
		slog.Error(
			"Ошибка при получении страниц",
			slog.String("Message", err.Error()),
		)
		models.WSChan <- models.WSData{
			Cmd: "updateLog",
			Payload: map[string]interface{}{
				"type": "err",
				"text": "-- Ошибка при получении страниц:" + err.Error(),
			},
		}
		return nil, err
	}

	chapterPage, err := goquery.NewDocumentFromReader(bytes.NewReader(pageBody))
	if err != nil {
		return nil, err
	}

	if chapterPage.Find(".auth-page .alert").Text() != "" {
		return nil, errors.New("noauth")
	}

	r := regexp.MustCompile(`rm_h\.readerInit\(\[\[(.+)\]\],\s(false|true),\s(\[.+\]).+\);`)

	srvList := ServersList{}

	chList := r.FindStringSubmatch(string(pageBody))

	json.Unmarshal([]byte(chList[3]), &srvList)

	imageParts := strings.Split(strings.Trim(chList[1], "[]"), "],[")

	for i := 0; i < len(imageParts); i++ {
		tmpParts := strings.Split(imageParts[i], ",")

		imageLinks = append(imageLinks, strings.Trim(tmpParts[0], "\"'")+strings.Trim(tmpParts[2], "\"'"))
	}

	chapterPath := path.Join(configs.Cfg.Savepath, downData.SavePath, curChapter.Path)

	if _, err := os.Stat(chapterPath); os.IsNotExist(err) {
		os.MkdirAll(chapterPath, 0755)
	}

	var savedFiles []string

	for _, imgURL := range imageLinks {
		fileName, err := DlImage(imgURL, chapterPath, srvList, 0)
		if err != nil {
			models.WSChan <- models.WSData{
				Cmd: "updateLog",
				Payload: map[string]interface{}{
					"type": "err",
					"text": "-- Ошибка при скачивании страницы (" + imgURL + "):" + err.Error(),
				},
			}
			continue
		}

		savedFiles = append(savedFiles, fileName)

		time.Sleep(time.Duration(configs.Cfg.Readmanga.TimeoutImage) * time.Millisecond)
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

func DlImage(imgURL, chapterPath string, srvList ServersList, retry int) (string, error) {
	maxRetry := 5

	if retry > 0 {
		slog.Warn(
			"Повторная попытка",
			slog.Int("Message", retry),
		)
	}

	client := grab.NewClient()
	client.UserAgent = configs.Cfg.UserAgent

	url, _ := urlx.Parse(imgURL)
	host, _, _ := urlx.SplitHostPort(url)

	if host == "one-way.work" {
		imgURL = strings.Split(imgURL, "?")[0]
	} else {
		imgURL = baseServerUrl + imgURL
	}

	req, err := grab.NewRequest(chapterPath, imgURL)
	if err != nil {
		slog.Error(
			"Ошибка при скачивании страницы",
			slog.String("Message", err.Error()),
		)
		if retry == maxRetry {
			return "", err
		} else {
			time.Sleep(time.Duration(configs.Cfg.Readmanga.TimeoutImage) * time.Millisecond)
			return DlImage(imgURL, chapterPath, srvList, retry+1)
		}
	}

	// req.HTTPRequest.Header.Set("Referer", refURL.Scheme+"://"+refURL.Host+"/")

	resp := client.Do(req)
	if resp.Err() != nil {
		if resp.HTTPResponse != nil && resp.HTTPResponse.StatusCode == 404 {
			slog.Error(
				"Ошибка при скачивании страницы",
				slog.String("Message", resp.Err().Error()),
			)
			if retry == maxRetry {
				return "", err
			} else {
				newImgUrl := GetServer(imgURL, srvList)
				time.Sleep(time.Duration(configs.Cfg.Readmanga.TimeoutImage) * time.Millisecond)
				return DlImage(newImgUrl, chapterPath, srvList, retry+1)
			}
		} else {
			slog.Error(
				"Ошибка при скачивании страницы",
				slog.String("Message", resp.Err().Error()),
			)
			if retry == maxRetry {
				return "", err
			} else {
				time.Sleep(time.Duration(configs.Cfg.Readmanga.TimeoutImage) * time.Millisecond)
				return DlImage(imgURL, chapterPath, srvList, retry+1)
			}
		}
	}

	return resp.Filename, nil
}

func GetServer(imgURL string, srvList ServersList) string {
	servers := []string{}

	srcUrl, _ := urlx.Parse(imgURL)

	for _, s := range srvList {
		newUrl, _ := urlx.Parse(s.Path)

		if newUrl.Host != srcUrl.Host {
			servers = append(servers, newUrl.Host)
		}
	}

	rnd := rand.Intn(len(servers) - 1)

	srcUrl.Host = servers[rnd]

	return srcUrl.String()
}
