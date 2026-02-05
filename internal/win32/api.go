//go:build windows

package win32

import (
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

// Windows API
var (
	user32                         = syscall.NewLazyDLL("user32.dll")
	procShutdownBlockReasonCreate  = user32.NewProc("ShutdownBlockReasonCreate")
	procShutdownBlockReasonDestroy = user32.NewProc("ShutdownBlockReasonDestroy")
	procPostMessageW               = user32.NewProc("PostMessageW")
	procSetForegroundWindow        = user32.NewProc("SetForegroundWindow")
)

// ShutdownBlockReasonCreateはシャットダウンをブロックする理由を設定します。
// この関数は副作用（Win32 API呼び出し）を持ちます。
func ShutdownBlockReasonCreate(hwnd win.HWND, reason string) {
	reasonPtr, _ := syscall.UTF16PtrFromString(reason)
	_, _, _ = procShutdownBlockReasonCreate.Call(uintptr(hwnd), uintptr(unsafe.Pointer(reasonPtr)))
}

// ShutdownBlockReasonDestroyはシャットダウンブロックの理由をクリアします。
// この関数は副作用（Win32 API呼び出し）を持ちます。
func ShutdownBlockReasonDestroy(hwnd win.HWND) {
	_, _, _ = procShutdownBlockReasonDestroy.Call(uintptr(hwnd))
}

// PostMessageはウィンドウのメッセージキューにメッセージをポストします。
// この関数は副作用（Win32 API呼び出し）を持ちます。
func PostMessage(hwnd win.HWND, msg uint32, wParam, lParam uintptr) {
	_, _, _ = procPostMessageW.Call(uintptr(hwnd), uintptr(msg), wParam, lParam)
}

// SetForegroundWindowはウィンドウをフォアグラウンドに表示します。
// この関数は副作用（Win32 API呼び出し）を持ちます。
func SetForegroundWindow(hwnd win.HWND) {
	_, _, _ = procSetForegroundWindow.Call(uintptr(hwnd))
}

// ShellExecuteはデフォルトのアプリケーションを使用してファイルまたはURLを開きます。
// この関数は副作用（Win32 API呼び出し）を持ちます。
func ShellExecute(hwnd win.HWND, url string) {
	verb, _ := syscall.UTF16PtrFromString("open")
	file, _ := syscall.UTF16PtrFromString(url)
	win.ShellExecute(hwnd, verb, file, nil, nil, SW_SHOWNORMAL)
}
