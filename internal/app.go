package internal

import (
	"log"
	"repeat-what-shit/internal/hotkeys"
	"repeat-what-shit/internal/input"
	"repeat-what-shit/internal/storage"
	"repeat-what-shit/internal/types"
	"repeat-what-shit/internal/utils"
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
				if stopCh, exists := a.activeMacros[macro.ID]; exists {
					close(stopCh)
					delete(a.activeMacros, macro.ID)
				} else {
					stopCh := make(chan bool)
					a.activeMacros[macro.ID] = stopCh
					go a.executeToggleMacro(macro, stopCh)
				}

			case types.MacroTypeHold:
				if _, exists := a.activeMacros[macro.ID]; !exists {
					stopCh := make(chan bool)
					a.activeMacros[macro.ID] = stopCh
					go a.executeHoldMacro(macro, stopCh)
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

func (a *App) executeHoldMacro(macro types.Macro, stopCh chan bool) {
	for {
		if !hotkeys.IsComboPressed(macro.ActivationKeys) {
			if ch, exists := a.activeMacros[macro.ID]; exists {
				close(ch)
				delete(a.activeMacros, macro.ID)
			}
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

			if !hotkeys.IsComboPressed(macro.ActivationKeys) {
				if ch, exists := a.activeMacros[macro.ID]; exists {
					close(ch)
					delete(a.activeMacros, macro.ID)
				}
				return
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func equalCombos(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
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
