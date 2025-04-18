package browser

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

const (
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36"
)

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	})
}

type Browser struct {
	ctx    context.Context
	cancel context.CancelFunc
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func NewBrowser(headless bool) (*Browser, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", headless),
		chromedp.UserAgent(userAgent),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-notifications", true),
		chromedp.Flag("disable-popup-blocking", true),
		chromedp.Flag("disable-web-security", false),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.Flag("start-maximized", true),
		chromedp.Flag("window-size", "1920,1080"),
	)

	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)

	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(log.Printf),
	)

	return &Browser{
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

func (b *Browser) Close() {
	b.cancel()
}

func (b *Browser) WaitForLogin() error {
	log.Info("Please scan the QR code to login...")

	// Create login success channel
	loginSuccess := make(chan bool, 1)

	// Start a goroutine to check login status
	go func() {
		for {
			time.Sleep(time.Second)
			select {
			case <-loginSuccess:
				return
			default:
				// Continue checking
			}
		}
	}()

	// Use a single Run to handle navigation and login detection
	err := chromedp.Run(b.ctx,
		chromedp.Navigate("https://www.zhihu.com/signin?type=login"),
		chromedp.WaitVisible(`div.Qrcode-content`, chromedp.ByQuery),
		// Wait for login success
		chromedp.ActionFunc(func(ctx context.Context) error {
			for {
				var isLoggedIn bool
				err := chromedp.EvaluateAsDevTools(
					`document.querySelector('#root > div > div:nth-child(2) > header > div.AppHeader-inner.css-11p8nt5 > div.AppHeader-userInfo > div.AppHeader-profile > div.Popover.AppHeader-menu') !== null`,
					&isLoggedIn,
				).Do(ctx)

				if err != nil {
					return err
				}

				if isLoggedIn {
					log.Info("Login successful")
					return nil
				}

				// Wait 1 second before checking again
				time.Sleep(time.Second)
			}
		}),
	)

	if err != nil {
		return fmt.Errorf("Failed to wait for login: %v", err)
	}

	return nil
}

func (b *Browser) VoteArticles(urls []string) error {
	log.Info("Starting to vote...")

	for _, articleUrl := range urls {
		// Random delay between 1-3 seconds
		delay := time.Second * time.Duration(1+rand.Intn(3))
		time.Sleep(delay)

		var isLiked bool
		// Custom JavaScript script
		script := `
			// Find all buttons containing "赞同" followed by space and number
			const buttons = Array.from(document.querySelectorAll('button[aria-label^="赞同 "]'));
			if (buttons.length === 0) {
				throw new Error("No like button found");
			}
			
			// Find the first inactive like button
			const likeButton = buttons.find(button => !button.classList.contains('is-active'));
			if (!likeButton) {
				throw new Error("No clickable like button found");
			}
			
			// Scroll to button position
			likeButton.scrollIntoView({ behavior: 'smooth', block: 'center' });
			
			// Click the button
			likeButton.click();
		`

		err := chromedp.Run(b.ctx,
			chromedp.Navigate(articleUrl),
			// Wait for article body to load
			chromedp.WaitVisible(`article`, chromedp.ByQuery),
			// Initial scroll to middle to trigger dynamic loading
			chromedp.EvaluateAsDevTools(`window.scrollTo(0, document.body.scrollHeight / 3)`, nil),
			chromedp.Sleep(1*time.Second),
			// Execute custom script
			chromedp.EvaluateAsDevTools(script, nil),
			chromedp.Sleep(1*time.Second),
			// Verify if vote was successful
			chromedp.EvaluateAsDevTools(`
				(() => {
					const buttons = Array.from(document.querySelectorAll('button[aria-label^="已赞同 "]'));
					const isLiked = buttons.some(button => button.classList.contains('is-active'));
					return isLiked;
				})();
			`, &isLiked),
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
