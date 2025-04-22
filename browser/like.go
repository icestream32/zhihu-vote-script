package browser

import (
	"math/rand"
	"time"
	"zhihu/script"

	"github.com/chromedp/chromedp"
)

func (b *Browser) Like(url string) error {
	log.Info("Starting to like...")

	// Random delay between 1-3 seconds
	delay := time.Second * time.Duration(1+rand.Intn(3))
	time.Sleep(delay)

	var isLiked bool

	err := chromedp.Run(b.ctx,
		chromedp.Navigate(url),
		// Wait for article body to load
		chromedp.WaitVisible(`article`, chromedp.ByQuery),
		// Initial scroll to middle to trigger dynamic loading
		chromedp.EvaluateAsDevTools(`window.scrollTo(0, document.body.scrollHeight / 3)`, nil),
		chromedp.Sleep(1*time.Second),
		// Execute like script
		chromedp.EvaluateAsDevTools(script.GetLikeScript(), nil),
		// Wait for like request to complete
		chromedp.Sleep(1*time.Second),
		// Check if liked successfully
		chromedp.EvaluateAsDevTools(script.GetCheckIfLikedScript(), &isLiked),
	)

	if err != nil {
		log.Errorf("Failed to like article: %s, error: %v", url, err)
		return err
	}

	if isLiked {
		log.Infof("Successfully liked article: %s", url)
	} else {
		log.Errorf("Failed to like article: %s, like not detected", url)
	}

	return nil
}
