//go:build windows

package mutex

import (
	"fmt"
	"syscall"

	"shutdown-alert/internal/win32"
)

// アプリケーション固有のミューテックス名
// Global\プレフィックスを使用することで、異なるセッション間でも共有される
const applicationMutexName = "Global\\ShutdownAlert-{E5F8A9C3-4D2B-4A1E-9F3C-8B7D6E5A4C2F}"

// AppMutexはアプリケーションの2重起動を防止するためのミューテックスを表します。
type AppMutex struct {
	handle syscall.Handle
}

// Acquireはアプリケーションミューテックスを取得します。
// 既に別のインスタンスが起動している場合はエラーを返します。
// この関数は副作用（Win32 API呼び出し）を持ちます。
func Acquire() (*AppMutex, error) {
	handle, alreadyExists, err := win32.CreateMutex(applicationMutexName)
	if err != nil {
		return nil, fmt.Errorf("ミューテックスの作成に失敗しました: %w", err)
	}

	if alreadyExists {
		// 既に存在する場合はハンドルを閉じてエラーを返す
		_ = win32.CloseHandle(handle)
		return nil, fmt.Errorf("アプリケーションは既に起動しています")
	}

	return &AppMutex{handle: handle}, nil
}

// Releaseはミューテックスを解放します。
// この関数は副作用（Win32 API呼び出し）を持ちます。
func (mutex *AppMutex) Release() error {
	if mutex.handle == 0 {
		return nil
	}

	// ミューテックスを解放してハンドルを閉じる
	_ = win32.ReleaseMutex(mutex.handle)
	err := win32.CloseHandle(mutex.handle)
	mutex.handle = 0

	return err
}
