package hotkeys

import (
	"context"
	"repeat-what-shit/internal/input"
	"sync"

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
)

type KeyCombo struct {
	Keys []int
	Time uint32
}

type KeyComboHandler func(combo KeyCombo)

type HotkeyService struct {
	keyboardChannel chan types.KeyboardEvent
	mouseChannel    chan types.MouseEvent
	handler         KeyComboHandler
	pressedKeys     map[int]struct{}
	cancelMu        sync.Mutex
	cancelMap       map[int]context.CancelFunc
	lastEventTime   uint32
}

var globalService *HotkeyService

func NewHotkeyService() *HotkeyService {
	service := &HotkeyService{
		keyboardChannel: make(chan types.KeyboardEvent),
		mouseChannel:    make(chan types.MouseEvent),
		pressedKeys:     make(map[int]struct{}),
		cancelMap:       make(map[int]context.CancelFunc),
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
	s.cancelMu.Lock()
	defer s.cancelMu.Unlock()
	_, exists := s.pressedKeys[key]
	return exists
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
			if uint32(e.DWExtraInfo) == input.EMULATED_FLAG {
				continue
			}

			switch e.Message {
			case types.WM_KEYDOWN, types.WM_SYSKEYDOWN:
				key := int(e.VKCode)

				if s.isKeyPressed(key) {
					continue
				}

				s.cancelMu.Lock()
				s.pressedKeys[key] = struct{}{}
				s.lastEventTime = e.Time
				s.cancelMu.Unlock()

				_, cancel := context.WithCancel(context.Background())
				s.cancelMu.Lock()
				s.cancelMap[key] = cancel
				s.cancelMu.Unlock()

				if s.handler != nil {
					combo := KeyCombo{
						Keys: append([]int(nil), s.GetPressedKeys()...),
						Time: e.Time,
					}
					go s.handler(combo)
				}

			case types.WM_KEYUP, types.WM_SYSKEYUP:
				key := int(e.VKCode)

				s.cancelEmulation(key)

				s.cancelMu.Lock()
				delete(s.pressedKeys, key)
				s.lastEventTime = e.Time
				s.cancelMu.Unlock()

				if s.handler != nil {
					combo := KeyCombo{
						Keys: append([]int(nil), s.GetPressedKeys()...),
						Time: e.Time,
					}
					go s.handler(combo)
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
					combo := KeyCombo{
						Keys: []int{wheelEvent},
						Time: e.Time,
					}
					go s.handler(combo)
				}
			case WM_LBUTTONDOWN, WM_RBUTTONDOWN, WM_MBUTTONDOWN:
				if s.handler != nil {
					combo := KeyCombo{
						Keys: []int{int(e.Message)},
						Time: e.Time,
					}
					go s.handler(combo)
				}
			case WM_XBUTTONDOWN:
				if s.handler != nil {
					mouseEvent := int(e.Message) | int(e.MouseData>>16)<<16
					combo := KeyCombo{
						Keys: []int{mouseEvent},
						Time: e.Time,
					}
					go s.handler(combo)
				}
			}
		}
	}
}

func (s *HotkeyService) GetPressedKeys() []int {
	s.cancelMu.Lock()
	defer s.cancelMu.Unlock()

	pressed := make([]int, 0, len(s.pressedKeys))
	for key := range s.pressedKeys {
		pressed = append(pressed, key)
	}
	return pressed
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
