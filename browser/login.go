package browser

import (
	"context"
	"fmt"
	"time"
	"zhihu/script"

	"github.com/chromedp/chromedp"
)

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
		chromedp.Navigate(loginUrl),
		chromedp.WaitVisible(`div.Qrcode-content`, chromedp.ByQuery),
		// Wait for login success
		chromedp.ActionFunc(func(ctx context.Context) error {
			for {
				var isLoggedIn bool
				err := chromedp.EvaluateAsDevTools(
					script.LoginScript,
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
		return fmt.Errorf("failed to wait for login: %v", err)
	}

	return nil
}
