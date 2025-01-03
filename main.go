package main

import (
	"context"
	"embed"
	"log"
	"repeat-what-shit/internal/logger"

	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.ico
var icon []byte

func main() {
	defer logger.RecoverWithLog()

	if err := logger.Init(); err != nil {
		log.Fatal(err)
	}
	defer logger.Close()

	app := NewApp()
	defer app.Shutdown()

	go systray.Run(func() {
		systray.SetIcon(icon)
		systray.SetTitle(appName)
		systray.SetTooltip(appName)

		mOpen := systray.AddMenuItem("Открыть", "Открыть главное окно")
		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Выход", "Закрыть приложение")

		go func() {
			defer logger.RecoverWithLog()
			for {
				select {
				case <-app.ctx.Done():
					return
				case <-mOpen.ClickedCh:
					if app.ctx != nil {
						runtime.WindowShow(app.ctx)
					}
				case <-mQuit.ClickedCh:
					systray.Quit()
					if app.ctx != nil {
						runtime.Quit(app.ctx)
					}
					return
				}
			}
		}()
	}, func() {
	})

	err := wails.Run(&options.App{
		Title:     appName,
		Width:     700,
		Height:    600,
		MinWidth:  700,
		MinHeight: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx, runtime.Environment(ctx).BuildType != "development")
		},
		OnShutdown: func(ctx context.Context) {
			app.Shutdown()
		},
		OnBeforeClose: func(ctx context.Context) bool {
			runtime.WindowHide(ctx)
			return false
		},
		Bind: []interface{}{
			app,
		},
		Frameless: true,
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
		},
	})

	if err != nil {
		log.Printf("[ERROR] Application error: %v", err)
	}
}
