package main

import (
	"context"
	"embed"

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
	app := NewApp()

	go systray.Run(func() {
		systray.SetIcon(icon)
		systray.SetTitle("repeat-what-shit")
		systray.SetTooltip("repeat-what-shit")

		mOpen := systray.AddMenuItem("Открыть", "Открыть главное окно")
		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Выход", "Закрыть приложение")

		go func() {
			for {
				select {
				case <-mOpen.ClickedCh:
					runtime.WindowShow(app.ctx)
				case <-mQuit.ClickedCh:
					systray.Quit()
					runtime.Quit(app.ctx)
				}
			}
		}()
	}, func() {
		// Cleanup
	})

	err := wails.Run(&options.App{
		Title:     "repeat-what-shit",
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
		println("Error:", err.Error())
	}
}
