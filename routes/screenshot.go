package routes

import (
	"context"
	"encoding/json"
	"github.com/asphaltbot/screenshotter/util"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/kataras/iris/v12"
	"github.com/mafredri/cdp/devtool"
	"io/ioutil"
	"log"
	"os"
)

type ScreenshotRequest struct {
	URL string `json:"url"`
}

func RegisterScreenshotRoute(app *iris.Application) {
	app.Post("/screenshot", TakeScreenshot)
}

func TakeScreenshot(irisContext iris.Context) {
	bgContext := context.Background()
	devTools := devtool.New("http://127.0.0.1:9222") // headless chrome running on docker
	pt, err := devTools.Get(bgContext, devtool.Page)

	if err != nil {
		pt, err = devTools.Create(bgContext)
		if err != nil {
			panic(err)
		}
	}

	allocCtx, cancel := chromedp.NewRemoteAllocator(context.Background(), pt.WebSocketDebuggerURL)
	defer cancel()

	bodyBytes, err := ioutil.ReadAll(irisContext.Request().Body)
	util.CheckError(err)

	var screenshotRequest ScreenshotRequest
	err = json.Unmarshal(bodyBytes, &screenshotRequest)
	util.CheckError(err)

	var imageBuf []byte
	if err := chromedp.Run(allocCtx, ScreenshotTasks(screenshotRequest.URL, &imageBuf)); err != nil {
		log.Fatal(err)
	}

	randomFileName := util.RandStringRunes(8)
	if err := ioutil.WriteFile(randomFileName+".png", imageBuf, 0644); err != nil {
		log.Fatal(err)
	} else {
		fileBytes, err := ioutil.ReadFile(randomFileName + ".png")
		util.CheckError(err)

		_ = os.Remove(randomFileName + ".png")

		irisContext.ContentType("image/jpeg")
		_, _ = irisContext.Write(fileBytes)
	}

}

func ScreenshotTasks(url string, imageBuf *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		emulation.SetDeviceMetricsOverride(1920, 1080, 1.0, false),
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) (err error) {
			*imageBuf, err = page.CaptureScreenshot().WithQuality(100).Do(ctx)
			return err
		}),
	}
}
