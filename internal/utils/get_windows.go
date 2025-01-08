package utils

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"path/filepath"
	"syscall"
	"unsafe"
)

type ICONINFO struct {
	FIcon    int32
	XHotspot int32
	YHotspot int32
	HbmMask  syscall.Handle
	HbmColor syscall.Handle
}

type BITMAP struct {
	Type       int32
	Width      int32
	Height     int32
	WidthBytes int32
	Planes     uint16
	BitsPixel  uint16
	Bits       *byte
}

type WindowInfo struct {
	Handle     syscall.Handle `json:"handle"`
	Process    string         `json:"process"`
	IconBase64 string         `json:"iconBase64"`
}

var (
	user32                   = syscall.NewLazyDLL("user32.dll")
	gdi32                    = syscall.NewLazyDLL("gdi32.dll")
	kernel32                 = syscall.NewLazyDLL("kernel32.dll")
	psapi                    = syscall.NewLazyDLL("psapi.dll")
	shell32                  = syscall.NewLazyDLL("shell32.dll")
	enumWindows              = user32.NewProc("EnumWindows")
	isWindowVisible          = user32.NewProc("IsWindowVisible")
	getIconInfo              = user32.NewProc("GetIconInfo")
	getDC                    = user32.NewProc("GetDC")
	createCompatibleDC       = gdi32.NewProc("CreateCompatibleDC")
	createCompatibleBitmap   = gdi32.NewProc("CreateCompatibleBitmap")
	selectObject             = gdi32.NewProc("SelectObject")
	drawIcon                 = user32.NewProc("DrawIcon")
	getBitmapBits            = gdi32.NewProc("GetBitmapBits")
	deleteDC                 = gdi32.NewProc("DeleteDC")
	releaseDC                = user32.NewProc("ReleaseDC")
	deleteObject             = gdi32.NewProc("DeleteObject")
	getWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	openProcess              = kernel32.NewProc("OpenProcess")
	getModuleFileNameEx      = psapi.NewProc("GetModuleFileNameExW")
	closeHandle              = kernel32.NewProc("CloseHandle")
	extractIconEx            = shell32.NewProc("ExtractIconExW")
	destroyIcon              = user32.NewProc("DestroyIcon")
)

func GetWindows() []WindowInfo {
	var windows []WindowInfo
	seenProcesses := make(map[string]bool)

	cb := syscall.NewCallback(func(hwnd syscall.Handle, lparam uintptr) uintptr {
		if !isWindowVisible_(hwnd) {
			return 1
		}

		process := getProcessName(hwnd)
		if process == "" {
			return 1
		}

		if seenProcesses[process] {
			return 1
		}
		seenProcesses[process] = true

		iconBase64 := getWindowIcon(hwnd)

		windows = append(windows, WindowInfo{
			Handle:     hwnd,
			Process:    process,
			IconBase64: iconBase64,
		})

		return 1
	})

	enumWindows.Call(cb, 0)
	return windows
}

func isWindowVisible_(hwnd syscall.Handle) bool {
	ret, _, _ := isWindowVisible.Call(uintptr(hwnd))
	return ret != 0
}

func getProcessName(hwnd syscall.Handle) string {
	var processID uint32
	getWindowThreadProcessId.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&processID)),
	)

	const PROCESS_QUERY_INFORMATION = 0x0400
	const PROCESS_VM_READ = 0x0010

	handle, _, _ := openProcess.Call(
		PROCESS_QUERY_INFORMATION|PROCESS_VM_READ,
		0,
		uintptr(processID),
	)

	if handle == 0 {
		return ""
	}
	defer closeHandle.Call(handle)

	var buf [260]uint16
	ret, _, _ := getModuleFileNameEx.Call(
		handle,
		0,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(len(buf)),
	)
	if ret == 0 {
		return ""
	}

	return filepath.Base(syscall.UTF16ToString(buf[:]))
}

func getWindowIcon(hwnd syscall.Handle) string {
	var processID uint32
	getWindowThreadProcessId.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&processID)),
	)

	const PROCESS_QUERY_INFORMATION = 0x0400
	const PROCESS_VM_READ = 0x0010

	handle, _, _ := openProcess.Call(
		PROCESS_QUERY_INFORMATION|PROCESS_VM_READ,
		0,
		uintptr(processID),
	)

	if handle == 0 {
		return ""
	}
	defer closeHandle.Call(handle)

	var exePath [260]uint16
	ret, _, _ := getModuleFileNameEx.Call(
		handle,
		0,
		uintptr(unsafe.Pointer(&exePath[0])),
		uintptr(len(exePath)),
	)
	if ret == 0 {
		return ""
	}

	var iconSmall, iconLarge syscall.Handle
	ret, _, _ = extractIconEx.Call(
		uintptr(unsafe.Pointer(&exePath[0])),
		0,
		uintptr(unsafe.Pointer(&iconLarge)),
		uintptr(unsafe.Pointer(&iconSmall)),
		1,
	)

	if ret == 0 {
		return ""
	}

	icon := iconSmall
	if icon == 0 {
		icon = iconLarge
	}
	if iconSmall != 0 {
		defer destroyIcon.Call(uintptr(iconSmall))
	}
	if iconLarge != 0 {
		defer destroyIcon.Call(uintptr(iconLarge))
	}

	if icon == 0 {
		return ""
	}

	var iconInfo ICONINFO
	ret, _, _ = getIconInfo.Call(
		uintptr(icon),
		uintptr(unsafe.Pointer(&iconInfo)),
	)
	if ret == 0 {
		return ""
	}

	hdcScreen, _, _ := getDC.Call(0)
	hdcMem, _, _ := createCompatibleDC.Call(hdcScreen)
	hbmp, _, _ := createCompatibleBitmap.Call(hdcScreen, 32, 32)

	selectObject.Call(hdcMem, hbmp)
	drawIcon.Call(hdcMem, 0, 0, uintptr(icon))

	img := image.NewRGBA(image.Rect(0, 0, 32, 32))
	bitsSize := 32 * 32 * 4
	getBitmapBits.Call(hbmp, uintptr(bitsSize), uintptr(unsafe.Pointer(&img.Pix[0])))

	for i := 0; i < len(img.Pix); i += 4 {
		img.Pix[i], img.Pix[i+2] = img.Pix[i+2], img.Pix[i]
	}

	deleteDC.Call(hdcMem)
	releaseDC.Call(0, hdcScreen)
	deleteObject.Call(hbmp)
	if iconInfo.HbmColor != 0 {
		deleteObject.Call(uintptr(iconInfo.HbmColor))
	}
	if iconInfo.HbmMask != 0 {
		deleteObject.Call(uintptr(iconInfo.HbmMask))
	}

	var buf bytes.Buffer
	png.Encode(&buf, img)
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

func GetActiveProcessName() string {
	var processID uint32
	user32 := syscall.NewLazyDLL("user32.dll")
	getForegroundWindow := user32.NewProc("GetForegroundWindow")
	getWindowThreadProcessId := user32.NewProc("GetWindowThreadProcessId")

	hwnd, _, _ := getForegroundWindow.Call()
	if hwnd == 0 {
		return ""
	}

	getWindowThreadProcessId.Call(
		hwnd,
		uintptr(unsafe.Pointer(&processID)),
	)

	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	psapi := syscall.NewLazyDLL("psapi.dll")
	openProcess := kernel32.NewProc("OpenProcess")
	getModuleFileNameEx := psapi.NewProc("GetModuleFileNameExW")
	closeHandle := kernel32.NewProc("CloseHandle")

	const PROCESS_QUERY_INFORMATION = 0x0400
	const PROCESS_VM_READ = 0x0010

	handle, _, _ := openProcess.Call(
		PROCESS_QUERY_INFORMATION|PROCESS_VM_READ,
		0,
		uintptr(processID),
	)

	if handle == 0 {
		return ""
	}
	defer closeHandle.Call(handle)

	var buf [260]uint16
	ret, _, _ := getModuleFileNameEx.Call(
		handle,
		0,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(len(buf)),
	)
	if ret == 0 {
		return ""
	}

	return filepath.Base(syscall.UTF16ToString(buf[:]))
}
