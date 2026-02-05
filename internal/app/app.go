//go:build windows

package app

import (
	"fmt"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"

	"shutdown-alert/internal/config"
	"shutdown-alert/internal/ui"
	"shutdown-alert/internal/win32"
)

// Appはメインアプリケーションを表します。
type App struct {
	mainWindow *walk.MainWindow
	notifyIcon *walk.NotifyIcon
	urlToOpen  string
}

// NewAppは新しいアプリケーションインスタンスを作成します。
func NewApp() *App {
	return &App{
		urlToOpen: config.TargetURL,
		// mainWindowとnotifyIconはRun内で初期化されます。
	}
}

// Runはアプリケーションを初期化して実行します。
// この関数は副作用（UIの作成、メッセージループの実行）を持ちます。
func (app *App) Run() error {
	// WndProcコールバックのためにappインスタンスを保存します。
	appInstance = app

	err := app.createMainWindow()
	if err != nil {
		return fmt.Errorf("メインウィンドウの作成に失敗しました: %w", err)
	}

	err = app.initNotifyIcon()
	if err != nil {
		return fmt.Errorf("通知アイコンの初期化に失敗しました: %w", err)
	}

	app.installWndProcHook()

	// Windowsのイベント待ちループに入ります。
	// この後の処理はイベントドリブンで行われます。
	app.mainWindow.Run()

	// アプリが終了（=Runが終わる）したら以下をクリーンアップします。
	// finallyブロックのようなものと考えます。
	if app.notifyIcon != nil {
		_ = app.notifyIcon.Dispose()
	}

	// グローバルインスタンスをクリアします。
	appInstance = nil

	return nil
}

// Windowsのアプリは必ずメインウィンドウを持つ必要があるので、非表示のメインウィンドウを作成します。
func (app *App) createMainWindow() error {
	return declarative.MainWindow{
		AssignTo: &app.mainWindow,
		Title:    config.DialogTitle,
		Visible:  false,
		Layout:   declarative.VBox{},
	}.Create()
}

// initNotifyIconは通知アイコンを作成して設定します。
// この関数は副作用（UIの作成）を持ちます。
func (app *App) initNotifyIcon() error {
	var err error
	app.notifyIcon, err = ui.InitNotifyIcon(
		app.mainWindow,
		app.showConfirmationDialog, // テスト用にshowConfirmationDialogを渡す
		func() { walk.App().Exit(0) },
	)
	return err
}

// handleShutdownQueryはシャットダウンが検出されたときに確認ダイアログを表示します。
// この関数は副作用（UIの表示、アプリケーションの終了の可能性）を持ちます。
func (app *App) handleShutdownQuery() {
	err := ui.ShowConfirmationDialog(
		app.mainWindow,
		app.urlToOpen,
		func() {
			app.openURL()
			walk.App().Exit(0)
		},
		func() {
			walk.App().Exit(0)
		},
	)

	if err != nil {
		// ダイアログの表示に失敗した場合は、単純に終了します。
		walk.App().Exit(0)
	}

	// ダイアログをフォアグラウンドに表示します。
	if app.mainWindow != nil {
		win32.SetForegroundWindow(app.mainWindow.Handle())
	}
}

// showConfirmationDialogはシャットダウン確認メッセージを表示します（テスト用）。
// この関数は副作用（UIの表示、アプリケーションの終了の可能性）を持ちます。
func (app *App) showConfirmationDialog() {
	_ = ui.ShowConfirmationDialog(
		app.mainWindow,
		app.urlToOpen,
		func() {
			app.openURL()
			walk.App().Exit(0)
		},
		func() {
			walk.App().Exit(0)
		},
	)

	// ダイアログをフォアグラウンドに表示します。
	if app.mainWindow != nil {
		win32.SetForegroundWindow(app.mainWindow.Handle())
	}
}

// openURLはShellExecuteを使用して対象のURLを開きます。
// この関数は副作用（外部アプリケーションの起動）を持ちます。
func (app *App) openURL() {
	win32.ShellExecute(app.mainWindow.Handle(), app.urlToOpen)
}
