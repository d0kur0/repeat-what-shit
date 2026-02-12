package internal

import (
	"log"
	"repeat-what-shit/internal/hotkeys"
	"repeat-what-shit/internal/input"
	"repeat-what-shit/internal/storage"
	"repeat-what-shit/internal/types"
	"repeat-what-shit/internal/utils"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type App struct {
	Storage *storage.JsonStorage[types.AppData]
	Version string

	HotkeyService *hotkeys.HotkeyService
	captureMode   bool
	lastCombo     []int
	lastComboTime uint32

	macroMu      sync.Mutex
	activeMacros map[string]chan bool
}

func (a *App) SetupHotkeys() {
	a.activeMacros = make(map[string]chan bool)
	a.HotkeyService = hotkeys.NewHotkeyService()

	a.HotkeyService.Start(func(combo hotkeys.KeyCombo) {
		log.Println(combo.Keys)

		if a.captureMode {
			if len(combo.Keys) == 0 {
				a.lastCombo = nil
				return
			}

			if len(combo.Keys) < len(a.lastCombo) {
				return
			}

			if len(combo.Keys) > len(a.lastCombo) || !equalCombos(combo.Keys, a.lastCombo) {
				a.lastCombo = append([]int(nil), combo.Keys...)
				application.Get().EmitEvent("captured_combo", combo.Keys)
			}
			return
		}

		for _, macro := range a.Storage.GetData().Macros {
			if macro.Disabled {
				continue
			}

			if !equalCombos(combo.Keys, macro.ActivationKeys) {
				continue
			}

			if !utils.IsWindowMatch(utils.GetActiveProcessName(), macro.IncludeTitle) {
				continue
			}

			switch macro.Type {
			case types.MacroTypeSequence:
				go a.executeMacro(macro)

			case types.MacroTypeToggle:
				a.macroMu.Lock()
				if stopCh, exists := a.activeMacros[macro.ID]; exists {
					close(stopCh)
					delete(a.activeMacros, macro.ID)
					a.macroMu.Unlock()
				} else {
					stopCh := make(chan bool)
					a.activeMacros[macro.ID] = stopCh
					a.macroMu.Unlock()
					go a.executeToggleMacro(macro, stopCh)
				}

			case types.MacroTypeHold:
				a.macroMu.Lock()
				if _, exists := a.activeMacros[macro.ID]; !exists {
					stopCh := make(chan bool)
					a.activeMacros[macro.ID] = stopCh
					a.macroMu.Unlock()
					go a.executeHoldMacro(macro, stopCh)
				} else {
					a.macroMu.Unlock()
				}
			}
		}
	})
}

func (a *App) executeMacro(macro types.Macro) {
	for _, action := range macro.Actions {
		input.SendInput(action.Keys)
		if action.Delay > 0 {
			time.Sleep(time.Duration(action.Delay) * time.Millisecond)
		}
	}
}

func (a *App) executeToggleMacro(macro types.Macro, stopCh chan bool) {
	for {
		select {
		case <-stopCh:
			return
		default:
			a.executeMacro(macro)
		}
	}
}

func (a *App) stopMacro(id string) {
	a.macroMu.Lock()
	defer a.macroMu.Unlock()
	if ch, exists := a.activeMacros[id]; exists {
		close(ch)
		delete(a.activeMacros, id)
	}
}

func (a *App) executeHoldMacro(macro types.Macro, stopCh chan bool) {
	defer a.stopMacro(macro.ID)

	for {
		if !hotkeys.IsComboPressed(macro.ActivationKeys) {
			return
		}

		select {
		case <-stopCh:
			return
		default:
			for _, action := range macro.Actions {
				input.SendInput(action.Keys)
				if action.Delay > 0 {
					time.Sleep(time.Duration(action.Delay) * time.Millisecond)
				}
			}
		}
	}
}

func equalCombos(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for _, keyA := range a {
		found := false
		for _, keyB := range b {
			if keyA == keyB {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (a *App) StartCapture() {
	a.captureMode = true
}

func (a *App) StopCapture() {
	a.captureMode = false
	a.lastCombo = nil
	a.lastComboTime = 0
}

func (a *App) ReadAppData() types.AppData {
	return a.Storage.GetData()
}

func (a *App) WriteAppData(data types.AppData) {
	a.Storage.Write(data)
}

func (a *App) GetVersion() string {
	return a.Version
}

func (a *App) GetWindowList() []utils.WindowInfo {
	return utils.GetWindows()
}
