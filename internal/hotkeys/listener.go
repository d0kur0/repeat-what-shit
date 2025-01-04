package hotkeys

import (
	"context"
	"log"
	"strings"
	"sync"
	"syscall"
	"unsafe"

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

	LLKHF_INJECTED = 0x00000010 // Флаг для определения эмулированных нажатий
)

var (
	user32   = syscall.NewLazyDLL("user32.dll")
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	psapi    = syscall.NewLazyDLL("psapi.dll")

	getWindowText            = user32.NewProc("GetWindowTextW")
	getClassName             = user32.NewProc("GetClassNameW")
	getForegroundWindow      = user32.NewProc("GetForegroundWindow")
	getWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")

	getModuleFileNameEx = psapi.NewProc("GetModuleFileNameExW")
	openProcess         = kernel32.NewProc("OpenProcess")
	closeHandle         = kernel32.NewProc("CloseHandle")

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
}

func getWindowInfo() (title, class, process string) {
	// Получаем хендл активного окна
	hwnd, _, _ := getForegroundWindow.Call()
	if hwnd == 0 {
		return
	}

	// Получаем заголовок окна
	var titleBuf [256]uint16
	getWindowText.Call(hwnd, uintptr(unsafe.Pointer(&titleBuf[0])), 256)
	title = syscall.UTF16ToString(titleBuf[:])

	// Получаем класс окна
	var classBuf [256]uint16
	getClassName.Call(hwnd, uintptr(unsafe.Pointer(&classBuf[0])), 256)
	class = syscall.UTF16ToString(classBuf[:])

	// Получаем ID процесса
	var processID uint32
	getWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&processID)))

	// Открываем процесс
	const PROCESS_QUERY_INFORMATION = 0x0400
	const PROCESS_VM_READ = 0x0010
	handle, _, _ := openProcess.Call(PROCESS_QUERY_INFORMATION|PROCESS_VM_READ, 0, uintptr(processID))
	if handle != 0 {
		defer closeHandle.Call(handle)

		// Получаем имя исполняемого файла
		var exeBuf [256]uint16
		getModuleFileNameEx.Call(handle, 0, uintptr(unsafe.Pointer(&exeBuf[0])), 256)
		exePath := syscall.UTF16ToString(exeBuf[:])
		if exePath != "" {
			parts := strings.Split(exePath, "\\")
			process = parts[len(parts)-1]
		}
	}

	return
}

// IsWindowMatch проверяет, соответствует ли текущее активное окно заданным условиям
func IsWindowMatch(includeTitles string) bool {
	if includeTitles == "" {
		return true
	}

	title, class, process := getWindowInfo()
	log.Printf("Window info - Title: %s, Class: %s, Process: %s", title, class, process)

	targets := strings.Split(includeTitles, ",")
	log.Printf("Checking against targets: %v", targets)

	for _, target := range targets {
		target = strings.TrimSpace(target)
		if target == "" {
			continue
		}

		log.Printf("Checking target: %s", target)
		// Проверяем совпадение с заголовком, классом или процессом
		if strings.Contains(strings.ToLower(title), strings.ToLower(target)) {
			log.Printf("Match by title")
			return true
		}
		if strings.Contains(strings.ToLower(class), strings.ToLower(target)) {
			log.Printf("Match by class")
			return true
		}
		if strings.Contains(strings.ToLower(process), strings.ToLower(target)) {
			log.Printf("Match by process")
			return true
		}
	}

	log.Printf("No matches found")
	return false
}

func NewHotkeyService() *HotkeyService {
	return &HotkeyService{
		keyboardChannel: make(chan types.KeyboardEvent),
		mouseChannel:    make(chan types.MouseEvent),
		pressedKeys:     make([]int, 0, 4),
		cancelMap:       make(map[int]context.CancelFunc),
	}
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
			// Пропускаем эмулированные нажатия
			if e.Flags&LLKHF_INJECTED != 0 {
				continue
			}

			switch e.Message {
			case types.WM_KEYDOWN, types.WM_SYSKEYDOWN:
				key := int(e.VKCode)
				log.Printf("Key pressed: %X", key)

				// Для новых нажатий
				if !s.isKeyPressed(key) {
					s.pressedKeys = append(s.pressedKeys, key)

					_, cancel := context.WithCancel(context.Background())
					s.cancelMu.Lock()
					s.cancelMap[key] = cancel
					s.cancelMu.Unlock()
				}

				// Всегда вызываем обработчик при нажатии
				if s.handler != nil {
					go s.handler(append([]int(nil), s.pressedKeys...), "")
				}

			case types.WM_KEYUP, types.WM_SYSKEYUP:
				key := int(e.VKCode)
				log.Printf("Key released: %X", key)

				s.cancelEmulation(key)

				for i, k := range s.pressedKeys {
					if k == key {
						s.pressedKeys = append(s.pressedKeys[:i], s.pressedKeys[i+1:]...)
						break
					}
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

// IsComboPressed проверяет, нажата ли указанная комбинация клавиш
func IsComboPressed(combo []int) bool {
	for _, key := range combo {
		if !isKeyPressed(key) {
			return false
		}
	}
	return true
}

// isKeyPressed проверяет, нажата ли указанная клавиша
func isKeyPressed(key int) bool {
	ret, _, _ := getAsyncKeyState.Call(uintptr(key))
	return (ret & 0x8000) != 0
}
