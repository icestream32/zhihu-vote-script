package main

import (
	"bufio"
	"os"
	"zhihu/browser"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	})
}

func execute() {
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
	if err := b.Vote(articleUrls); err != nil {
		log.Errorf("Voting failed: %v", err)
		return
	}
}

func main() {
	// Check if urls.txt exists
	_, err := os.Open("urls.txt")
	if err != nil {
		log.Error("File not found: urls.txt")
		return
	}

	execute()

	log.Info("Press Enter to exit...")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}
