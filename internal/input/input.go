package input

import (
	"syscall"
	"time"
	"unsafe"
)

var (
	user32    = syscall.NewLazyDLL("user32.dll")
	sendInput = user32.NewProc("SendInput")
)

const (
	INPUT_KEYBOARD = 1
	INPUT_MOUSE    = 0

	KEYEVENTF_KEYUP = 0x0002
	EMULATED_FLAG   = 0xBADF00D
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

func SendInput(keys []int) error {
	if len(keys) == 0 {
		return nil
	}

	// Сначала нажимаем все клавиши
	inputs := make([]INPUT, len(keys))
	for i, key := range keys {
		inputs[i].Type = INPUT_KEYBOARD
		inputs[i].Ki.Vk = uint16(key)
		inputs[i].Ki.ExtraInfo = EMULATED_FLAG
	}

	ret, _, err := sendInput.Call(
		uintptr(len(inputs)),
		uintptr(unsafe.Pointer(&inputs[0])),
		unsafe.Sizeof(INPUT{}),
	)

	if ret == 0 {
		return err
	}

	time.Sleep(10 * time.Millisecond)

	for i := range keys {
		inputs[i].Type = INPUT_KEYBOARD
		inputs[i].Ki.Vk = uint16(keys[i])
		inputs[i].Ki.Flags = KEYEVENTF_KEYUP
		inputs[i].Ki.ExtraInfo = EMULATED_FLAG
	}

	ret, _, err = sendInput.Call(
		uintptr(len(inputs)),
		uintptr(unsafe.Pointer(&inputs[0])),
		unsafe.Sizeof(INPUT{}),
	)

	if ret == 0 {
		return err
	}

	return nil
}
