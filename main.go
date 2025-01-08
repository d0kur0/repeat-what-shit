package main

import (
	"embed"
	"fmt"
	"log"
	"log/slog"
	"repeat-what-shit/internal"
	"repeat-what-shit/internal/consts"
	"repeat-what-shit/internal/storage"
	"repeat-what-shit/internal/types"
	"repeat-what-shit/internal/utils"
	"runtime/debug"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed frontend/dist
var uiAssets embed.FS

//go:embed build/appicon.png
var appIcon []byte

func main() {
	appDir, err := utils.GetAppDirPath()
	utils.Catch(err)

	defer func() {
		if r := recover(); r != nil {
			crashLog := fmt.Sprintf("Panic: %v\n\nStack Trace:\n%s", r, debug.Stack())
			s := storage.NewRawStorage(fmt.Sprintf("%s/crash_%s.log", appDir, time.Now().Format("2006-01-02_15-04-05")))
			s.Write([]byte(crashLog))
			log.Println(crashLog)
		}
	}()

	utils.Catch(utils.CreateAppDirIfNotExists())

	appData := storage.NewJsonStorage(fmt.Sprintf("%s/data.json", appDir), types.AppData{})
	utils.Catch(appData.Read())

	a := internal.App{
		Storage: appData,
		Version: consts.Version,
	}

	a.SetupHotkeys()

	app := application.New(application.Options{
		LogLevel:    slog.LevelError,
		Name:        consts.AppName,
		Description: "Repeat what shit",
		Services: []application.Service{
			application.NewService(&a),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(uiAssets),
		},
	})

	createMainWindow(app)
	app.Run()
}

func createMainWindow(app *application.App) {
	mainWindowStartState := application.WindowStateNormal
	if !consts.IsProduction {
		mainWindowStartState = application.WindowStateMinimised
	}

	w := app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		URL:        "/",
		Width:      consts.AppBaseWidth,
		Height:     consts.AppBaseHeight,
		MinWidth:   consts.AppMinWidth,
		MinHeight:  consts.AppMinHeight,
		Title:      consts.AppName,
		StartState: mainWindowStartState,
		Frameless:  true,
		Centered:   true,
	})

	w.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
		app.Quit()
	})

	tray := app.NewSystemTray()
	tray.SetLabel(consts.AppName)
	tray.SetIcon(appIcon)
	tray.SetDarkModeIcon(appIcon)

	trayMenu := app.NewMenu()

	trayMenu.Add("Close").OnClick(func(_ *application.Context) {
		app.Quit()
	})

	tray.SetMenu(trayMenu)

	tray.OnClick(func() {
		w.UnMinimise()
		w.Show()
		w.Focus()
		w.SetAlwaysOnTop(true)
		time.Sleep(100 * time.Millisecond)
		w.SetAlwaysOnTop(false)
	})
}
