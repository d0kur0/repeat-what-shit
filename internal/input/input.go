package input

import (
	"sync/atomic"
	"syscall"
	"unsafe"
)

var (
	isEmulating atomic.Bool
	user32      = syscall.NewLazyDLL("user32.dll")
	sendInput   = user32.NewProc("SendInput")
)

const (
	INPUT_KEYBOARD = 1
	INPUT_MOUSE    = 0

	KEYEVENTF_KEYUP = 0x0002
)

type INPUT struct {
	Type uint32
	Ki   struct {
		Vk        uint16
		Scan      uint16
		Flags     uint32
		Time      uint32
		ExtraInfo uintptr
		Padding1  uint32
		Padding2  uint32
	}
}

// IsEmulating возвращает true если сейчас идет эмуляция ввода
func IsEmulating() bool {
	return isEmulating.Load()
}

// SendInput отправляет событие ввода (клавиатура или мышь)
func SendInput(keys []int) error {
	if len(keys) == 0 {
		return nil
	}

	isEmulating.Store(true)
	defer func() {
		isEmulating.Store(false)
	}()

	inputs := make([]INPUT, len(keys)*2)

	for i, key := range keys {
		inputs[i].Type = INPUT_KEYBOARD
		inputs[i].Ki.Vk = uint16(key)
	}

	for i := 0; i < len(keys); i++ {
		j := len(keys)*2 - 1 - i
		inputs[j].Type = INPUT_KEYBOARD
		inputs[j].Ki.Vk = uint16(keys[len(keys)-1-i])
		inputs[j].Ki.Flags = KEYEVENTF_KEYUP
	}

	ret, _, err := sendInput.Call(
		uintptr(len(inputs)),
		uintptr(unsafe.Pointer(&inputs[0])),
		unsafe.Sizeof(INPUT{}),
	)
	if ret == 0 {
		return err
	}
	return nil
}
