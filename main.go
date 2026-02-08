//go:build windows

package main

import (
	"fmt"
	"log"
	"syscall"

	"github.com/lxn/win"

	"shutdown-alert/internal/app"
	"shutdown-alert/internal/config"
	"shutdown-alert/internal/mutex"
)

func main() {
	// 設定ファイルを読み込み
	userConfig, err := config.LoadUserConfig("config.yaml")
	if err != nil {
		// パース失敗時はエラーメッセージを表示
		log.Printf("設定ファイルの読み込みに失敗しました（デフォルト値を使用）: %v", err)
		errorMessage := fmt.Sprintf("設定ファイルの読み込みに失敗しました。\nデフォルト値を使用します。\n\nエラー: %v", err)
		message, _ := syscall.UTF16PtrFromString(errorMessage)
		title, _ := syscall.UTF16PtrFromString("設定エラー")
		win.MessageBox(0, message, title, win.MB_OK|win.MB_ICONWARNING)
	}

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
	a := app.NewApp(userConfig)
	if err := a.Run(); err != nil {
		log.Fatalf("アプリケーションの実行に失敗しました: %v", err)
	}
}
