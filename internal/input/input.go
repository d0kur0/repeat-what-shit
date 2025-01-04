package input

import (
	"syscall"
	"unsafe"
)

var (
	user32    = syscall.NewLazyDLL("user32.dll")
	sendInput = user32.NewProc("SendInput")
)

const (
	INPUT_KEYBOARD = 1 // keyboard input
	INPUT_MOUSE    = 0 // mouse input

	KEYEVENTF_KEYUP = 0x0002    // key up event
	LLKHF_INJECTED  = 0x10      // injected key event
	EMULATED_FLAG   = 0xBADF00D // unique flag for our emulated key events
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

// SendInput sends an input event (keyboard or mouse)
func SendInput(keys []int) error {
	if len(keys) == 0 {
		return nil
	}

	inputs := make([]INPUT, len(keys)*2)

	for i, key := range keys {
		inputs[i].Type = INPUT_KEYBOARD
		inputs[i].Ki.Vk = uint16(key)
		inputs[i].Ki.ExtraInfo = EMULATED_FLAG // mark as emulated key event
	}

	for i := 0; i < len(keys); i++ {
		j := len(keys)*2 - 1 - i
		inputs[j].Type = INPUT_KEYBOARD
		inputs[j].Ki.Vk = uint16(keys[len(keys)-1-i])
		inputs[j].Ki.Flags = KEYEVENTF_KEYUP
		inputs[j].Ki.ExtraInfo = EMULATED_FLAG // mark as emulated key event
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
