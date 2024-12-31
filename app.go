package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"repeat-what-shit/internal/hotkeys"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type MacroType int

const (
	MacroTypeSequence MacroType = iota
	MacroTypeWhilePressed
	MacroTypeToggle
)

var configFilePath string

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	configFilePath = filepath.Join(homeDir, ".repeat-what-shit", "data.json")
}

type MacroAction struct {
	Keys  []int `json:"keys"`
	Delay int   `json:"delay"`
}

type Macro struct {
	ID             string        `json:"id"`
	Name           string        `json:"name"`
	ActivationKeys []int         `json:"activation_keys"`
	Type           MacroType     `json:"type"`
	Actions        []MacroAction `json:"actions"`
}

type AppData struct {
	Macros []Macro `json:"macros"`
}

type App struct {
	ctx          context.Context
	isProduction bool

	hotkeyService *hotkeys.HotkeyService

	Data AppData
}

func NewApp() *App {
	return &App{
		hotkeyService: hotkeys.NewHotkeyService(),
	}
}

func (a *App) startup(ctx context.Context, isProduction bool) {
	a.ctx = ctx
	a.isProduction = isProduction

	// Создаем директорию для конфига, если её нет
	configDir := filepath.Dir(configFilePath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		panic(err)
	}

	// Проверяем существование файла
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// Файл не существует, создаем новый с пустой структурой
		emptyData := AppData{Macros: []Macro{}}
		jsonData, err := json.MarshalIndent(emptyData, "", "  ")
		if err != nil {
			panic(err)
		}
		if err := os.WriteFile(configFilePath, jsonData, 0644); err != nil {
			panic(err)
		}
		a.Data = emptyData
	} else {
		// Читаем существующий файл
		jsonData, err := os.ReadFile(configFilePath)
		if err != nil {
			panic(err)
		}
		if err := json.Unmarshal(jsonData, &a.Data); err != nil {
			panic(err)
		}
	}
}

func (a *App) GetData() AppData {
	return a.Data
}

func (a *App) UpdateData(data AppData) {
	a.Data = data

	// Сериализуем данные в JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return
	}

	// Записываем в файл
	if err := os.WriteFile(configFilePath, jsonData, 0644); err != nil {
		return
	}
}

func (a *App) StartCapture() error {
	log.Println("StartCapture")
	return a.hotkeyService.Start(func(combo hotkeys.KeyCombo) {
		log.Println("combo_captured", combo)
		runtime.EventsEmit(a.ctx, "combo_captured", combo)
	})
}

func (a *App) StopCapture() {
	log.Println("StopCapture")
	a.hotkeyService.Stop()
}
