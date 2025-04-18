package main

import (
	"bufio"
	"os"
	"zhihu/script"

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

func main() {
	// Check if urls.txt exists
	_, err := os.Open("urls.txt")
	if err != nil {
		log.Error("File not found: urls.txt")
		return
	}

	script.Execute()

	log.Info("Press Enter to exit...")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}
