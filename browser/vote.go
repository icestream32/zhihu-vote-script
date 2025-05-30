package browser

import (
	"math/rand"
	"time"
	"zhihu/script"

	"github.com/chromedp/chromedp"
)

func (b *Browser) Vote(url string) error {
	log.Info("Starting to vote...")

	// Random delay between 1-3 seconds
	delay := time.Second * time.Duration(1+rand.Intn(3))
	time.Sleep(delay)

	var isVoted bool

	err := chromedp.Run(b.ctx,
		chromedp.Navigate(url),
		// Wait for article body to load
		chromedp.WaitVisible(`article`, chromedp.ByQuery),
		// Initial scroll to middle to trigger dynamic loading
		chromedp.EvaluateAsDevTools(`window.scrollTo(0, document.body.scrollHeight / 3)`, nil),
		chromedp.Sleep(1*time.Second),
		// Execute custom script
		chromedp.EvaluateAsDevTools(script.GetVoteScript(), &isVoted),
		chromedp.Sleep(1*time.Second),
	)

	if err != nil {
		log.Errorf("Failed to vote article: %s, error: %v", url, err)
		return err
	}

	if isVoted {
		log.Infof("Successfully voted article: %s", url)
	} else {
		log.Errorf("Failed to vote article: %s, vote not detected", url)
	}

	return nil
}
