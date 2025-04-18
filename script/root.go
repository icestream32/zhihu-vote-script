package script

import (
	"bufio"
	"os"

	"github.com/sirupsen/logrus"
)

const (
	authority   = "www.zhihu.com"
	origin      = "https://zhuanlan.zhihu.com"
	contentType = "application/json"
	userAgent   = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36"
	baseApiUrl  = "https://www.zhihu.com/api/v4/articles/"
)

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	})
}

func Execute() {
	// read cookie.txt
	cookie := ""
	file, err := os.Open("cookie.txt")
	if err != nil {
		log.Errorf("Failed to open cookie file: %v", err)
		return
	}
	defer file.Close()

	cookieScanner := bufio.NewScanner(file)
	for cookieScanner.Scan() {
		cookie = cookieScanner.Text()
	}

	// read urls.txt
	file, err = os.Open("urls.txt")
	if err != nil {
		log.Errorf("Failed to open zhihu urls file: %v", err)
		return
	}
	defer file.Close()

	// get url by line
	articleUrls := []string{}
	urlScanner := bufio.NewScanner(file)
	for urlScanner.Scan() {
		url := urlScanner.Text()
		articleUrls = append(articleUrls, url)
	}

	go vote(articleUrls, cookie)
	go like(articleUrls, cookie)
}
