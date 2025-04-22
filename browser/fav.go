package browser

import (
	"math/rand"
	"time"
	"zhihu/script"

	"github.com/chromedp/chromedp"
)

func (b *Browser) Fav(url string) error {
	log.Info("Starting to fav...")

	// Random delay between 1-3 seconds
	delay := time.Second * time.Duration(1+rand.Intn(3))
	time.Sleep(delay)

	var isFaved bool

	err := chromedp.Run(b.ctx,
		// Wait for article body to load
		chromedp.WaitVisible(`article`, chromedp.ByQuery),
		// Initial scroll to middle to trigger dynamic loading
		chromedp.EvaluateAsDevTools(`window.scrollTo(0, document.body.scrollHeight / 3)`, nil),
		chromedp.Sleep(1*time.Second),
		// Click fav entry button
		chromedp.EvaluateAsDevTools(script.GetFavEntryButtonScript(), nil),
		// Wait for fav dialog to appear
		chromedp.Sleep(1*time.Second),
		// Click fav button
		chromedp.EvaluateAsDevTools(script.GetFavButtonScript(), nil),
		// Wait for fav request to complete
		chromedp.Sleep(1*time.Second),
		// Check if faved successfully
		chromedp.EvaluateAsDevTools(script.GetCheckIfFavScript(), &isFaved),
	)

	if err != nil {
		log.Errorf("Failed to fav article: %s, error: %v", url, err)
		return err
	}

	if isFaved {
		log.Infof("Successfully faved article: %s", url)
	} else {
		log.Errorf("Failed to fav article: %s, fav not detected", url)
	}

	return nil
}
