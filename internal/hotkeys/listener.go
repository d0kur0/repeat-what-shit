package hotkeys

import (
	"sync"

	"github.com/moutend/go-hook/pkg/keyboard"
	"github.com/moutend/go-hook/pkg/mouse"
	"github.com/moutend/go-hook/pkg/types"
)

const (
	WM_MOUSEMOVE   = 0x0200
	WM_LBUTTONDOWN = 0x0201
	WM_LBUTTONUP   = 0x0202
	WM_RBUTTONDOWN = 0x0204
	WM_RBUTTONUP   = 0x0205
	WM_MBUTTONDOWN = 0x0207
	WM_MBUTTONUP   = 0x0208
	WM_MOUSEWHEEL  = 0x020A
	WM_XBUTTONDOWN = 0x020B
	WM_XBUTTONUP   = 0x020C
)

type KeyCombo []int
type KeyComboHandler func(combo KeyCombo)

type HotkeyService struct {
	isRunning       bool
	currentKeys     sync.Map
	keyboardChannel chan types.KeyboardEvent
	mouseChannel    chan types.MouseEvent
	handler         KeyComboHandler
}

func NewHotkeyService() *HotkeyService {
	return &HotkeyService{
		keyboardChannel: make(chan types.KeyboardEvent),
		mouseChannel:    make(chan types.MouseEvent),
	}
}

func (s *HotkeyService) Start(handler KeyComboHandler) error {
	if s.isRunning {
		return nil
	}

	if err := keyboard.Install(nil, s.keyboardChannel); err != nil {
		return err
	}

	if err := mouse.Install(nil, s.mouseChannel); err != nil {
		keyboard.Uninstall()
		return err
	}

	s.handler = handler
	s.isRunning = true
	go s.handleEvents()
	return nil
}

func (s *HotkeyService) Stop() {
	if !s.isRunning {
		return
	}
	keyboard.Uninstall()
	mouse.Uninstall()
	s.isRunning = false
}

func (s *HotkeyService) getCurrentCombo() KeyCombo {
	var keys []int
	s.currentKeys.Range(func(key, value interface{}) bool {
		if k, ok := key.(int); ok {
			keys = append(keys, k)
		}
		return true
	})
	return keys
}

func (s *HotkeyService) handleEvents() {
	for {
		select {
		case e := <-s.keyboardChannel:
			if !s.isRunning {
				return
			}

			switch e.Message {
			case types.WM_KEYDOWN, types.WM_SYSKEYDOWN:
				s.currentKeys.Store(int(e.VKCode), true)
			case types.WM_KEYUP, types.WM_SYSKEYUP:
				combo := s.getCurrentCombo()
				if len(combo) > 0 && s.handler != nil {
					s.handler(combo)
				}
				s.currentKeys = sync.Map{} // Очищаем после обработки
			}

		case e := <-s.mouseChannel:
			if !s.isRunning {
				return
			}

			switch e.Message {
			case WM_MOUSEMOVE:
				continue
			case WM_MOUSEWHEEL:
				// Для колеса мыши сразу отправляем одиночное событие
				delta := int16(e.MouseData >> 16)
				var wheelEvent int
				if delta > 0 {
					wheelEvent = int(e.Message) | 0x10000 // Up flag
				} else {
					wheelEvent = int(e.Message) | 0x20000 // Down flag
				}
				if s.handler != nil {
					s.handler([]int{wheelEvent})
				}
			case WM_LBUTTONDOWN, WM_RBUTTONDOWN, WM_MBUTTONDOWN:
				// Для обычных кнопок мыши отправляем событие сразу
				if s.handler != nil {
					s.handler([]int{int(e.Message)})
				}
			case WM_XBUTTONDOWN:
				// Для X-кнопок добавляем информацию о номере кнопки
				buttonNum := uint16(e.MouseData >> 16)
				mouseEvent := int(e.Message) | int(buttonNum)<<16
				if s.handler != nil {
					s.handler([]int{mouseEvent})
				}
			}
		}
	}
}
