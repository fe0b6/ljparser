package ljparser

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

var (
	dateReg    *regexp.Regexp
	commentReg *regexp.Regexp
)

func init() {
	commentReg = regexp.MustCompile(`([0-9]+) (?:коммента|comment)`)
	dateReg = regexp.MustCompile(`<meta property="og:image" content="[^"]+v=([0-9]+)" />`)
}

func GetComments(lnk string) (pi PostInfo, err error) {
	// Читаем страницу
	resp, err := http.Get(lnk)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		log.Println("[error]", err)
		return
	}

	// Что-то пошло не так
	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
		log.Println("[error]", resp.Status, resp.StatusCode)
		return
	}

	// Читаем ответ
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	var comInt, dateInt int64

	com := commentReg.FindStringSubmatch(string(content))
	if len(com) >= 2 {
		comInt, _ = strconv.ParseInt(com[1], 10, 64)
	}

	date := dateReg.FindStringSubmatch(string(content))
	if len(date) >= 2 {
		dateInt, _ = strconv.ParseInt(date[1], 10, 64)
	}

	pi = PostInfo{
		Comments:  int(comInt),
		Published: int(dateInt),
	}

	return
}
