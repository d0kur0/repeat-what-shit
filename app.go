package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"repeat-what-shit/internal/hotkeys"
	"repeat-what-shit/internal/input"
	"repeat-what-shit/internal/logger"
	"sort"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const appName = "repeat what shit"

type MacroType int

const (
	MacroTypeSequence MacroType = iota
	MacroTypeToggle
	MacroTypeHold
)

func readAppData() (AppData, error) {
	appDir, err := logger.GetAppDir()
	if err != nil {
		return AppData{}, fmt.Errorf("failed to get app directory: %w", err)
	}

	dataPath := filepath.Join(appDir, "data.json")
	data, err := os.ReadFile(dataPath)
	if err != nil {
		if os.IsNotExist(err) {
			return AppData{Macros: []Macro{}}, nil
		}
		return AppData{}, fmt.Errorf("failed to read data file: %w", err)
	}

	var appData AppData
	if err := json.Unmarshal(data, &appData); err != nil {
		return AppData{}, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return appData, nil
}

func writeAppData(data AppData) error {
	appDir, err := logger.GetAppDir()
	if err != nil {
		return fmt.Errorf("failed to get app directory: %w", err)
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	dataPath := filepath.Join(appDir, "data.json")
	if err := os.WriteFile(dataPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write data file: %w", err)
	}

	return nil
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
	cancel       context.CancelFunc
	isProduction bool

	hotkeyService  *hotkeys.HotkeyService
	captureHandler hotkeys.KeyComboHandler

	Data AppData

	mu                 sync.RWMutex
	toggledMacros      map[string]bool
	executingMacros    map[string]*string
	holdingMacros      map[string]bool
	activeHoldRoutines map[string]int
	isEmulated         bool
}

func NewApp() *App {
	ctx, cancel := context.WithCancel(context.Background())
	return &App{
		ctx:                ctx,
		cancel:             cancel,
		hotkeyService:      hotkeys.NewHotkeyService(),
		toggledMacros:      make(map[string]bool),
		executingMacros:    make(map[string]*string),
		holdingMacros:      make(map[string]bool),
		activeHoldRoutines: make(map[string]int),
		isEmulated:         false,
	}
}

func (a *App) Shutdown() {
	logger.Info("Завершение работы приложения")
	if a.cancel != nil {
		a.cancel()
	}
	a.mu.Lock()
	for k := range a.toggledMacros {
		a.toggledMacros[k] = false
	}
	a.mu.Unlock()
	logger.Debug("Ожидание завершения горутин")
	time.Sleep(100 * time.Millisecond)
	logger.Info("Приложение завершено")
}

func (a *App) handleHotkey(combo hotkeys.KeyCombo, _ string) {
	defer logger.RecoverWithLog()

	a.mu.RLock()
	isEmulated := a.isEmulated
	a.mu.RUnlock()

	if isEmulated {
		logger.Debug("Пропускаем эмулированное нажатие клавиш")
		return
	}

	logger.Debug("Получена комбинация клавиш: %v", combo)

	if a.captureHandler != nil {
		logger.Debug("Режим захвата активен")
		a.captureHandler(combo, "")
		return
	}

	sortedCombo := make([]int, len(combo))
	copy(sortedCombo, combo)
	sort.Ints(sortedCombo)
	logger.Debug("Отсортированная комбинация: %v", sortedCombo)

	for _, macro := range a.Data.Macros {
		sortedKeys := make([]int, len(macro.ActivationKeys))
		copy(sortedKeys, macro.ActivationKeys)
		sort.Ints(sortedKeys)
		logger.Debug("Проверяем макрос %s (ID: %s), ключи: %v", macro.Name, macro.ID, sortedKeys)

		if reflect.DeepEqual(sortedCombo, sortedKeys) {
			logger.Debug("Найдено совпадение для макроса %s", macro.Name)

			if macro.IncludeTitles != "" && !hotkeys.IsWindowMatch(macro.IncludeTitles) {
				logger.Debug("Пропускаем макрос %s - окно не соответствует", macro.Name)
				continue
			}

			if macro.Type == MacroTypeSequence {
				if _, exists := a.executingMacros[macro.ID]; exists {
					logger.Debug("Макрос %s уже выполняется", macro.Name)
					continue
				}
				a.executingMacros[macro.ID] = new(string)
				defer delete(a.executingMacros, macro.ID)
			}

			switch macro.Type {
			case MacroTypeSequence:
				logger.Debug("Выполняем последовательный макрос %s", macro.Name)
				for _, action := range macro.Actions {
					if err := input.SendInput(action.Keys); err != nil {
						logger.Error("Не удалось отправить комбинацию клавиш: %v", err)
					}
					if action.Delay > 0 {
						time.Sleep(time.Duration(action.Delay) * time.Millisecond)
					}
				}

			case MacroTypeToggle:
				a.mu.Lock()
				isActive := a.toggledMacros[macro.ID]
				a.toggledMacros[macro.ID] = !isActive
				a.mu.Unlock()

				if !isActive {
					logger.Debug("Включаем макрос %s", macro.Name)
					go func(ctx context.Context, macroID string, actions []MacroAction) {
						defer logger.RecoverWithLog()
						logger.Debug("Запущена горутина для макроса-переключателя %s", macro.Name)

						for {
							a.mu.RLock()
							isStillActive := a.toggledMacros[macroID]
							a.mu.RUnlock()

							if !isStillActive {
								logger.Debug("Остановка макроса %s - переключатель выключен", macro.Name)
								return
							}

							select {
							case <-ctx.Done():
								logger.Debug("Остановка макроса %s по контексту", macro.Name)
								return
							default:
								for _, action := range actions {
									a.mu.Lock()
									a.isEmulated = true
									a.mu.Unlock()

									if err := input.SendInput(action.Keys); err != nil {
										logger.Error("Не удалось отправить комбинацию клавиш: %v", err)
									}

									a.mu.Lock()
									a.isEmulated = false
									a.mu.Unlock()

									if action.Delay > 0 {
										time.Sleep(time.Duration(action.Delay) * time.Millisecond)
									}
								}
								time.Sleep(50 * time.Millisecond)
							}
						}
					}(a.ctx, macro.ID, macro.Actions)
				} else {
					logger.Debug("Выключаем макрос %s", macro.Name)
				}

			case MacroTypeHold:
				logger.Debug("Обработка макроса удержания %s", macro.Name)

				a.mu.Lock()
				currentCount := a.activeHoldRoutines[macro.ID]
				if currentCount == 0 {
					a.activeHoldRoutines[macro.ID] = 4
					a.mu.Unlock()

					for i := 1; i <= 4; i++ {
						routineNum := i
						logger.Debug("Запущена горутина %d для макроса %s", routineNum, macro.Name)

						go func(ctx context.Context, macroID string, actions []MacroAction) {
							defer func() {
								a.mu.Lock()
								a.activeHoldRoutines[macroID]--
								if a.activeHoldRoutines[macroID] <= 0 {
									delete(a.activeHoldRoutines, macroID)
								}
								a.mu.Unlock()
								logger.RecoverWithLog()
							}()

							for {
								select {
								case <-ctx.Done():
									logger.Debug("Остановка макроса %s по контексту", macro.Name)
									return
								default:
									if !hotkeys.IsComboPressed(sortedKeys) {
										logger.Debug("Остановка макроса %s - клавиши отпущены", macro.Name)
										return
									}

									for _, action := range actions {
										if err := input.SendInput(action.Keys); err != nil {
											logger.Error("Не удалось отправить комбинацию клавиш: %v", err)
										}
										if action.Delay > 0 {
											time.Sleep(time.Duration(action.Delay) * time.Millisecond)
										}
									}
								}
							}
						}(a.ctx, macro.ID, macro.Actions)
					}
				} else {
					a.mu.Unlock()
					logger.Debug("Пропускаем запуск - макрос %s уже выполняется (%d горутин)", macro.Name, currentCount)
				}
			}
		}
	}
}

func (a *App) StartCapture() error {
	logger.Info("Запуск режима захвата")
	a.captureHandler = func(combo hotkeys.KeyCombo, _ string) {
		if a.ctx != nil {
			runtime.EventsEmit(a.ctx, "combo_captured", combo)
		}
	}
	return a.hotkeyService.Start(a.handleHotkey)
}

func (a *App) StopCapture() {
	logger.Info("Остановка режима захвата")
	a.captureHandler = nil
}

func (a *App) startup(ctx context.Context, isProduction bool) {
	defer logger.RecoverWithLog()
	logger.Info("Запуск приложения (production: %v)", isProduction)

	if a.cancel != nil {
		logger.Debug("Отмена предыдущего контекста")
		a.cancel()
	}
	newCtx, cancel := context.WithCancel(ctx)
	a.ctx = newCtx
	a.cancel = cancel
	a.isProduction = isProduction

	data, err := readAppData()
	if err != nil {
		logger.Error("Failed to read app data: %v", err)
		panic(err)
	}
	logger.Info("Загружено %d макросов", len(data.Macros))
	a.Data = data

	if err := a.hotkeyService.Start(a.handleHotkey); err != nil {
		logger.Error("Failed to start hotkey service: %v", err)
	} else {
		logger.Info("Сервис хоткеев успешно запущен")
	}
}

func (a *App) GetData() AppData {
	return a.Data
}

func (a *App) UpdateData(data AppData) {
	defer logger.RecoverWithLog()
	logger.Info("Обновление данных приложения (%d макросов)", len(data.Macros))

	a.Data = data
	if err := writeAppData(data); err != nil {
		logger.Error("Failed to write app data: %v", err)
	} else {
		logger.Info("Данные успешно сохранены")
	}
}
