package internal

import (
	"repeat-what-shit/internal/storage"
	"repeat-what-shit/internal/types"
)

type App struct {
	Storage *storage.JsonStorage[types.AppData]
}

func (a *App) ReadAppData() types.AppData {
	return a.Storage.GetData()
}

func (a *App) WriteAppData(data types.AppData) {
	a.Storage.Write(data)
}
