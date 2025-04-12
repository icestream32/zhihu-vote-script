package chrome

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	})
}

func GetCookies() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		if ev, ok := ev.(*network.EventResponseReceived); ok {
			headers := ev.Response.Headers
			if cookie, ok := headers["Set-Cookie"]; ok {
				log.Infof("Found Cookie: %s\n", cookie)
			}
		}
	})

	err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.zhihu.com/signin?type=login"),

		chromedp.WaitVisible(`div.Qrcode-content`, chromedp.ByQuery),
		chromedp.WaitVisible(`div.Topstory-mainColumn`, chromedp.ByQuery),

		chromedp.ActionFunc(func(ctx context.Context) error {
			cookies, err := network.GetCookies().Do(ctx)
			if err != nil {
				return fmt.Errorf("获取cookies失败: %v", err)
			}

			saveCookie(cookies)
			return nil
		}),
	)

	if err != nil {
		log.Errorf("Failed to run: %v", err)
	}

	log.Infof("Login success.")
}

// save cookie to file
func saveCookie(cookies []*network.Cookie) {
	cookieFile, err := os.Create("cookie.txt")
	if err != nil {
		log.Errorf("Failed to create cookie file: %v", err)
		return
	}
	defer cookieFile.Close()

	for _, cookie := range cookies {
		cookieFile.WriteString(fmt.Sprintf("%s=%s;", cookie.Name, cookie.Value))
	}
}
