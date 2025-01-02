package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"repeat-what-shit/internal/hotkeys"
	"repeat-what-shit/internal/input"
	"sort"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type MacroType int

const (
	MacroTypeSequence MacroType = iota
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
	IncludeTitles  string        `json:"include_titles"`
}

type AppData struct {
	Macros []Macro `json:"macros"`
}

type App struct {
	ctx          context.Context
	isProduction bool

	hotkeyService  *hotkeys.HotkeyService  // единый сервис для всех операций с клавишами
	captureHandler hotkeys.KeyComboHandler // колбек для режима захвата

	Data AppData

	// Состояния макросов
	toggledMacros   map[string]bool    // для макросов-переключателей
	executingMacros map[string]*string // для отслеживания выполняющихся макросов
}

func NewApp() *App {
	return &App{
		hotkeyService:   hotkeys.NewHotkeyService(),
		toggledMacros:   make(map[string]bool),
		executingMacros: make(map[string]*string),
	}
}

func (a *App) handleHotkey(combo hotkeys.KeyCombo, _ string) {
	// Если есть колбек захвата, вызываем его
	if a.captureHandler != nil {
		a.captureHandler(combo, "")
		return
	}

	// Сортируем полученную комбинацию
	sortedCombo := make([]int, len(combo))
	copy(sortedCombo, combo)
	sort.Ints(sortedCombo)

	// Проверяем каждый макрос
	for _, macro := range a.Data.Macros {
		// Сортируем ожидаемую комбинацию
		sortedKeys := make([]int, len(macro.ActivationKeys))
		copy(sortedKeys, macro.ActivationKeys)
		sort.Ints(sortedKeys)

		// Сравниваем комбинации клавиш
		if reflect.DeepEqual(sortedCombo, sortedKeys) {
			// Проверяем окно, если указаны ограничения
			if macro.IncludeTitles != "" && !hotkeys.IsWindowMatch(macro.IncludeTitles) {
				continue
			}

			switch macro.Type {
			case MacroTypeSequence:
				for _, action := range macro.Actions {
					if err := input.SendInput(action.Keys); err != nil {
						log.Printf("[ERROR] Не удалось отправить комбинацию клавиш: %v", err)
					}
					if action.Delay > 0 {
						time.Sleep(time.Duration(action.Delay) * time.Millisecond)
					}
				}

			case MacroTypeToggle:
				a.toggledMacros[macro.ID] = !a.toggledMacros[macro.ID]
				if a.toggledMacros[macro.ID] {
					go func(macroID string, actions []MacroAction) {
						for a.toggledMacros[macroID] {
							for _, action := range actions {
								if err := input.SendInput(action.Keys); err != nil {
									log.Printf("[ERROR] Не удалось отправить комбинацию клавиш: %v", err)
								}
								if action.Delay > 0 {
									time.Sleep(time.Duration(action.Delay) * time.Millisecond)
								}
							}
						}
					}(macro.ID, macro.Actions)
				}
			}
		}
	}
}

func (a *App) StartCapture() error {
	a.captureHandler = func(combo hotkeys.KeyCombo, _ string) {
		runtime.EventsEmit(a.ctx, "combo_captured", combo)
	}
	return a.hotkeyService.Start(a.handleHotkey)
}

func (a *App) StopCapture() {
	a.captureHandler = nil
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

	// Запускаем отслеживание с общим обработчиком
	if err := a.hotkeyService.Start(a.handleHotkey); err != nil {
		log.Printf("Ошибка запуска отслеживания: %v", err)
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
