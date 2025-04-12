package main

import (
	"bufio"
	"os"
	"time"
	"zhihu/chrome"
	"zhihu/vote"

	"github.com/djherbis/times"
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

func check() {
	// read urls.txt
	_, err := os.Open("urls.txt")
	if err != nil {
		log.Error("Cannot found file: urls.txt")
	}

	// read cookie.txt
	file, err := os.Open("cookie.txt")
	if err != nil {
		log.Warn("Cannot open cookie file, execute login...")
		chrome.GetCookies()

		// read cookie.txt again
		file, err = os.Open("cookie.txt")
		if err != nil {
			log.Errorf("Failed to open cookie file: %v", err)
			return
		}
	}
	defer file.Close()

	// get file brithtime
	var birthTime time.Time
	if t, err := times.Stat(file.Name()); err == nil && t.HasBirthTime() {
		birthTime = t.BirthTime()
	}

	if birthTime.Before(time.Now().Add(-1 * time.Hour * 24)) {
		log.Warn("Cookie expired, execute login...")
		chrome.GetCookies()
	}
}

func main() {
	check()

	vote.Execute()

	log.Info("Press Enter to continue...")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}
