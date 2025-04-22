package browser

import (
	"context"

	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

const (
	loginUrl  = "https://www.zhihu.com/signin?type=login"
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

func (b *Browser) Run() error {
	return chromedp.Run(b.ctx)
}

func (b *Browser) Close() {
	b.cancel()
}
