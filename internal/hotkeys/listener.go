package hotkeys

import (
	"context"
	"sync"
	"syscall"

	"github.com/moutend/go-hook/pkg/keyboard"
	"github.com/moutend/go-hook/pkg/mouse"
	"github.com/moutend/go-hook/pkg/types"
)

const (
	WM_LBUTTONDOWN = 0x0201
	WM_RBUTTONDOWN = 0x0204
	WM_MBUTTONDOWN = 0x0207
	WM_MOUSEWHEEL  = 0x020A
	WM_XBUTTONDOWN = 0x020B

	LLKHF_INJECTED = 0x00000010
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	getAsyncKeyState = user32.NewProc("GetAsyncKeyState")
)

type KeyCombo []int
type KeyComboHandler func(combo KeyCombo, includeTitles string)

type HotkeyService struct {
	keyboardChannel chan types.KeyboardEvent
	mouseChannel    chan types.MouseEvent
	handler         KeyComboHandler
	pressedKeys     []int
	cancelMu        sync.Mutex
	cancelMap       map[int]context.CancelFunc
	emulatedKeys    map[int]bool
}

var globalService *HotkeyService

func NewHotkeyService() *HotkeyService {
	service := &HotkeyService{
		keyboardChannel: make(chan types.KeyboardEvent),
		mouseChannel:    make(chan types.MouseEvent),
		pressedKeys:     make([]int, 0, 4),
		cancelMap:       make(map[int]context.CancelFunc),
		emulatedKeys:    make(map[int]bool),
	}
	globalService = service
	return service
}

func (s *HotkeyService) Start(handler KeyComboHandler) error {
	if err := keyboard.Install(nil, s.keyboardChannel); err != nil {
		return err
	}

	if err := mouse.Install(nil, s.mouseChannel); err != nil {
		keyboard.Uninstall()
		return err
	}

	s.handler = handler
	go s.handleEvents()
	return nil
}

func (s *HotkeyService) Stop() {
	keyboard.Uninstall()
	mouse.Uninstall()
}

func (s *HotkeyService) isKeyPressed(key int) bool {
	for _, k := range s.pressedKeys {
		if k == key {
			return true
		}
	}
	return false
}

func (s *HotkeyService) cancelEmulation(key int) {
	s.cancelMu.Lock()
	defer s.cancelMu.Unlock()
	if cancel, exists := s.cancelMap[key]; exists {
		cancel()
		delete(s.cancelMap, key)
	}
}

func (s *HotkeyService) handleEvents() {
	for {
		select {
		case e := <-s.keyboardChannel:
			if e.Flags&LLKHF_INJECTED != 0 {
				key := int(e.VKCode)
				if e.Message == types.WM_KEYDOWN || e.Message == types.WM_SYSKEYDOWN {
					s.emulatedKeys[key] = true
				} else if e.Message == types.WM_KEYUP || e.Message == types.WM_SYSKEYUP {
					delete(s.emulatedKeys, key)
				}
				continue
			}

			switch e.Message {
			case types.WM_KEYDOWN, types.WM_SYSKEYDOWN:
				key := int(e.VKCode)

				if !s.isKeyPressed(key) {
					s.pressedKeys = append(s.pressedKeys, key)

					_, cancel := context.WithCancel(context.Background())
					s.cancelMu.Lock()
					s.cancelMap[key] = cancel
					s.cancelMu.Unlock()
				}

				if s.handler != nil {
					go s.handler(append([]int(nil), s.pressedKeys...), "")
				}

			case types.WM_KEYUP, types.WM_SYSKEYUP:
				key := int(e.VKCode)
				if s.emulatedKeys[key] {
					continue
				}

				s.cancelEmulation(key)

				for i, k := range s.pressedKeys {
					if k == key {
						s.pressedKeys = append(s.pressedKeys[:i], s.pressedKeys[i+1:]...)
						break
					}
				}

				if s.handler != nil {
					go s.handler(append([]int(nil), s.pressedKeys...), "")
				}
			}

		case e := <-s.mouseChannel:
			switch e.Message {
			case WM_MOUSEWHEEL:
				if s.handler != nil {
					wheelEvent := int(e.Message)
					if int16(e.MouseData>>16) > 0 {
						wheelEvent |= 0x10000
					} else {
						wheelEvent |= 0x20000
					}
					go s.handler([]int{wheelEvent}, "")
				}
			case WM_LBUTTONDOWN, WM_RBUTTONDOWN, WM_MBUTTONDOWN:
				if s.handler != nil {
					go s.handler([]int{int(e.Message)}, "")
				}
			case WM_XBUTTONDOWN:
				if s.handler != nil {
					mouseEvent := int(e.Message) | int(e.MouseData>>16)<<16
					go s.handler([]int{mouseEvent}, "")
				}
			}
		}
	}
}

func (s *HotkeyService) GetPressedKeys() []int {
	return append([]int(nil), s.pressedKeys...)
}

func IsComboPressed(combo []int) bool {
	if globalService == nil {
		return false
	}

	pressed := globalService.GetPressedKeys()
	for _, key := range combo {
		found := false
		for _, p := range pressed {
			if p == key {
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

func isKeyPressed(key int) bool {
	ret, _, _ := getAsyncKeyState.Call(uintptr(key))
	return (ret & 0x8000) != 0
}
