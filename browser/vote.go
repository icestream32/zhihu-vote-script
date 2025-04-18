package browser

import (
	"math/rand"
	"time"
	"zhihu/script"

	"github.com/chromedp/chromedp"
)

func (b *Browser) Vote(urls []string) error {
	log.Info("Starting to vote...")

	for _, articleUrl := range urls {
		// Random delay between 1-3 seconds
		delay := time.Second * time.Duration(1+rand.Intn(3))
		time.Sleep(delay)

		var isLiked bool

		err := chromedp.Run(b.ctx,
			chromedp.Navigate(articleUrl),
			// Wait for article body to load
			chromedp.WaitVisible(`article`, chromedp.ByQuery),
			// Initial scroll to middle to trigger dynamic loading
			chromedp.EvaluateAsDevTools(`window.scrollTo(0, document.body.scrollHeight / 3)`, nil),
			chromedp.Sleep(1*time.Second),
			// Execute custom script
			chromedp.EvaluateAsDevTools(script.GetVoteScript(), &isLiked),
			chromedp.Sleep(1*time.Second),
		)

		if err != nil {
			log.Errorf("Failed to vote article: %s, error: %v", articleUrl, err)
			continue
		}

		if isLiked {
			log.Infof("Successfully voted article: %s", articleUrl)
		} else {
			log.Errorf("Failed to vote article: %s, vote not detected", articleUrl)
		}
	}

	return nil
}
