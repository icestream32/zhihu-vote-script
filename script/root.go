package script

import (
	"bufio"
	"os"

	"zhihu/browser"

	"github.com/sirupsen/logrus"
)

// const (
// 	authority   = "www.zhihu.com"
// 	origin      = "https://zhuanlan.zhihu.com"
// 	contentType = "application/json"
// 	userAgent   = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36"
// 	baseApiUrl  = "https://www.zhihu.com/api/v4/articles/"
// )

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	})
}

func Execute() {
	// Create browser instance
	b, err := browser.NewBrowser(false) // false means show browser window
	if err != nil {
		log.Errorf("Failed to create browser: %v", err)
		return
	}
	defer b.Close()

	// Wait for user login
	if err := b.WaitForLogin(); err != nil {
		log.Errorf("Login failed: %v", err)
		return
	}

	// Read article URLs
	file, err := os.Open("urls.txt")
	if err != nil {
		log.Errorf("Failed to open article URL file: %v", err)
		return
	}
	defer file.Close()

	// Get URL list
	articleUrls := []string{}
	urlScanner := bufio.NewScanner(file)
	for urlScanner.Scan() {
		url := urlScanner.Text()
		articleUrls = append(articleUrls, url)
	}

	// Execute voting
	if err := b.VoteArticles(articleUrls); err != nil {
		log.Errorf("Voting failed: %v", err)
		return
	}
}
