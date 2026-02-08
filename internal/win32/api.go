//go:build windows

package win32

import (
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

// Windows API
var (
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	procCreateMutexW = kernel32.NewProc("CreateMutexW")
	procReleaseMutex = kernel32.NewProc("ReleaseMutex")
	procCloseHandle  = kernel32.NewProc("CloseHandle")

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

// CreateMutexは名前付きミューテックスを作成します。
// この関数は副作用（Win32 API呼び出し）を持ちます。
// 戻り値: (ハンドル, 既に存在するか, エラー)
func CreateMutex(name string) (syscall.Handle, bool, error) {
	namePtr, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return 0, false, err
	}

	handle, _, lastErr := procCreateMutexW.Call(
		0, // セキュリティ属性（デフォルト）
		0, // 初期所有者なし
		uintptr(unsafe.Pointer(namePtr)),
	)

	if handle == 0 {
		return 0, false, lastErr
	}

	// ERROR_ALREADY_EXISTSの場合、既に存在する
	alreadyExists := lastErr.(syscall.Errno) == ERROR_ALREADY_EXISTS

	return syscall.Handle(handle), alreadyExists, nil
}

// ReleaseMutexはミューテックスを解放します。
// この関数は副作用（Win32 API呼び出し）を持ちます。
func ReleaseMutex(handle syscall.Handle) error {
	ret, _, err := procReleaseMutex.Call(uintptr(handle))
	if ret == 0 {
		return err
	}
	return nil
}

// CloseHandleはハンドルを閉じます。
// この関数は副作用（Win32 API呼び出し）を持ちます。
func CloseHandle(handle syscall.Handle) error {
	ret, _, err := procCloseHandle.Call(uintptr(handle))
	if ret == 0 {
		return err
	}
	return nil
}
