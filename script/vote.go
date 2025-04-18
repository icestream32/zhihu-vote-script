package script

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type Vote struct {
	Voting int `json:"voting"`
}

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	})
}

func vote(urls []string, cookie string) {
	log.Info("Start voting...")

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
