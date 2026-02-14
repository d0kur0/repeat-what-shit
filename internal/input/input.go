package input

import (
	"syscall"
	"time"
	"unsafe"
)

var (
	user32         = syscall.NewLazyDLL("user32.dll")
	sendInputProc  = user32.NewProc("SendInput")
)

const (
	INPUT_MOUSE    = 0
	INPUT_KEYBOARD = 1

	KEYEVENTF_KEYUP = 0x0002
	EMULATED_FLAG   = 0xBADF00D

	WM_LBUTTONDOWN = 0x0201
	WM_RBUTTONDOWN = 0x0204
	WM_MBUTTONDOWN = 0x0207
	WM_XBUTTONDOWN = 0x020B

	MOUSEEVENTF_LEFTDOWN   = 0x0002
	MOUSEEVENTF_LEFTUP     = 0x0004
	MOUSEEVENTF_RIGHTDOWN  = 0x0008
	MOUSEEVENTF_RIGHTUP    = 0x0010
	MOUSEEVENTF_MIDDLEDOWN = 0x0020
	MOUSEEVENTF_MIDDLEUP   = 0x0040
	MOUSEEVENTF_XDOWN      = 0x0080
	MOUSEEVENTF_XUP        = 0x0100
)

type KeyboardInput struct {
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

type MouseInput struct {
	Type      uint32
	Pad0      uint32
	Dx        int32
	Dy        int32
	MouseData uint32
	Flags     uint32
	Time      uint32
	Pad1      uint32
	ExtraInfo uintptr
}

func isMouseKey(key int) bool {
	base := key & 0xFFFF
	return base == WM_LBUTTONDOWN || base == WM_RBUTTONDOWN ||
		base == WM_MBUTTONDOWN || base == WM_XBUTTONDOWN
}

func sendKeyDown(key int) {
	var input KeyboardInput
	input.Type = INPUT_KEYBOARD
	input.Ki.Vk = uint16(key)
	input.Ki.ExtraInfo = EMULATED_FLAG
	sendInputProc.Call(1, uintptr(unsafe.Pointer(&input)), unsafe.Sizeof(input))
}

func sendKeyUp(key int) {
	var input KeyboardInput
	input.Type = INPUT_KEYBOARD
	input.Ki.Vk = uint16(key)
	input.Ki.Flags = KEYEVENTF_KEYUP
	input.Ki.ExtraInfo = EMULATED_FLAG
	sendInputProc.Call(1, uintptr(unsafe.Pointer(&input)), unsafe.Sizeof(input))
}

func sendMouseDown(key int) {
	var input MouseInput
	input.Type = INPUT_MOUSE
	input.ExtraInfo = EMULATED_FLAG

	base := key & 0xFFFF
	switch base {
	case WM_LBUTTONDOWN:
		input.Flags = MOUSEEVENTF_LEFTDOWN
	case WM_RBUTTONDOWN:
		input.Flags = MOUSEEVENTF_RIGHTDOWN
	case WM_MBUTTONDOWN:
		input.Flags = MOUSEEVENTF_MIDDLEDOWN
	case WM_XBUTTONDOWN:
		input.Flags = MOUSEEVENTF_XDOWN
		input.MouseData = uint32(key >> 16)
	}

	sendInputProc.Call(1, uintptr(unsafe.Pointer(&input)), unsafe.Sizeof(input))
}

func sendMouseUp(key int) {
	var input MouseInput
	input.Type = INPUT_MOUSE
	input.ExtraInfo = EMULATED_FLAG

	base := key & 0xFFFF
	switch base {
	case WM_LBUTTONDOWN:
		input.Flags = MOUSEEVENTF_LEFTUP
	case WM_RBUTTONDOWN:
		input.Flags = MOUSEEVENTF_RIGHTUP
	case WM_MBUTTONDOWN:
		input.Flags = MOUSEEVENTF_MIDDLEUP
	case WM_XBUTTONDOWN:
		input.Flags = MOUSEEVENTF_XUP
		input.MouseData = uint32(key >> 16)
	}

	sendInputProc.Call(1, uintptr(unsafe.Pointer(&input)), unsafe.Sizeof(input))
}

func SendInput(keys []int) error {
	if len(keys) == 0 {
		return nil
	}

	for _, key := range keys {
		if isMouseKey(key) {
			sendMouseDown(key)
		} else {
			sendKeyDown(key)
		}
	}

	time.Sleep(10 * time.Millisecond)

	for _, key := range keys {
		if isMouseKey(key) {
			sendMouseUp(key)
		} else {
			sendKeyUp(key)
		}
	}

	return nil
}
