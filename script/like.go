package script

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

func like(urls []string, cookie string) {
	log.Info("Start liking...")

	// create http client
	client := &http.Client{}

	for _, articleUrl := range urls {
		// request slowly
		time.Sleep(time.Second * 1)

		// create http header
		header := http.Header{
			"Cookie":       []string{cookie},
			"Content-Type": []string{contentType},
			"User-Agent":   []string{userAgent},
			"Origin":       []string{origin},
			"Referer":      []string{origin},
			"Authority":    []string{authority},
		}

		// get article id from article url
		articleId := strings.Split(strings.Split(articleUrl, "/")[1], "/")[1]

		// create api url
		apiUrl, err := url.Parse(baseApiUrl + articleId + "/like")
		if err != nil {
			log.Errorf("Failed to parse api url: %v", err)
			continue
		}

		// create http request
		req, err := http.NewRequest("POST", apiUrl.String(), nil)
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

		log.Infof("Successfully liked article: %s", articleUrl)
	}
}
