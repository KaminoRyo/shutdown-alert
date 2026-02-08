//go:build windows

package main

import (
	"log"
	"syscall"

	"github.com/lxn/win"

	"shutdown-alert/internal/app"
	"shutdown-alert/internal/config"
	"shutdown-alert/internal/mutex"
)

func main() {
	// 2重起動防止: Win32 Mutex APIを使用してアプリケーションミューテックスを取得
	// 名前付きミューテックス（Global\ShutdownAlert-{GUID}）により、
	// プロセス間で排他制御を行い、既に起動している場合はエラーを返す
	appMutex, err := mutex.Acquire()
	if err != nil {
		// 既に起動している場合、GUIモードでもエラーメッセージを表示
		message, _ := syscall.UTF16PtrFromString("アプリケーションは既に起動しています。")
		title, _ := syscall.UTF16PtrFromString(config.DialogTitle)
		win.MessageBox(0, message, title, win.MB_OK|win.MB_ICONINFORMATION)
		return
	}
	// Mutex はOS,リソースハンドルなので必ず解放しないといけない。
	// ハンドルリークするとデバッグ時にゾンビ状態、再起動で内部状態がぐちゃぐちゃな挙動になる
	defer func() {
		if appMutex != nil {
			_ = appMutex.Release()
		}
	}()

	// アプリケーションを実行
	a := app.NewApp()
	if err := a.Run(); err != nil {
		log.Fatalf("アプリケーションの実行に失敗しました: %v", err)
	}
}
