package crawler

import (
	"context"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

const (
	_imageQualityDp = 120
	_width          = 12.0
	_height         = 17.0
	_widthPx        = 1920
	_heightPx       = 1080
)

func UrlToPdf(url string) ([]byte, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	if err := chromedp.Run(
		ctx,
		chromedp.EmulateViewport(_widthPx, _heightPx),
		printToPDF(url, &buf),
	); err != nil {
		return nil, err
	}
	return buf, nil
}

func UrlToScreenshot(url string, outFile string) ([]byte, error) {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()

	var buf []byte
	if err := chromedp.Run(
		ctx,
		chromedp.EmulateViewport(_widthPx, _heightPx),
		fullScreenshot(url, _imageQualityDp, &buf),
	); err != nil {
		return nil, err
	}
	return buf, nil
}

func printToPDF(url string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromeWaitAction(),
		printToPdfAction(res),
	}
}

func printToPdfAction(res *[]byte) chromedp.ActionFunc {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		buf, _, err := page.PrintToPDF().
			WithPrintBackground(true).
			WithPaperWidth(_width).
			WithPaperHeight(_height).
			Do(ctx)
		if err != nil {
			return err
		}
		*res = buf
		return nil
	})
}

func fullScreenshot(urlstr string, quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromeWaitAction(),
		chromedp.FullScreenshot(res, quality),
	}
}

func chromeWaitAction() chromedp.ActionFunc {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		time.Sleep(30 * time.Second)
		return nil
	})
}
