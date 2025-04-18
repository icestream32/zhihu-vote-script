package vote

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type Vote struct {
	Voting int `json:"voting"`
}

var log = logrus.New()

const (
	authority   = "www.zhihu.com"
	origin      = "https://zhuanlan.zhihu.com"
	contentType = "application/json"
	userAgent   = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36"
	baseApiUrl  = "https://www.zhihu.com/api/v4/articles/"
)

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

	handleUrls(articleUrls, cookie)
}

func handleUrls(urls []string, cookie string) {
	// create http client
	client := &http.Client{}

	for _, articleUrl := range urls {
		// request slowly
		time.Sleep(time.Second * 1)

		// create http header
		header := http.Header{
			"Cookie":       []string{cookie},
			"Content-Type": []string{contentType},
			"Authority":    []string{authority},
			"Origin":       []string{origin},
			"User-Agent":   []string{userAgent},
			"Referer":      []string{articleUrl},
		}

		// create http body
		body, err := json.Marshal(Vote{
			Voting: 1,
		})
		if err != nil {
			log.Errorf("Failed to marshal vote: %v", err)
			continue
		}

		// get article id from article url
		articleId := strings.Split(strings.Split(articleUrl, origin+"/")[1], "/")[1]

		// create api url
		apiUrl, err := url.Parse(baseApiUrl + articleId + "/voters")
		if err != nil {
			log.Errorf("Failed to parse api url: %v", err)
			continue
		}

		// create http request
		req, err := http.NewRequest("POST", apiUrl.String(), bytes.NewBuffer(body))
		if err != nil {
			log.Errorf("Failed to create http request: %v", err)
			continue
		}

		// set http header
		req.Header = header

		// send http request
		_, err = client.Do(req)
		if err != nil {
			log.Errorf("Failed to send http request: %v", err)
			continue
		}

		log.Infof("Successfully voted for article: %s", articleUrl)
	}
}
