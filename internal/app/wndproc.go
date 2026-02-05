//go:build windows

// このファイルは Windows のシャットダウン／ログオフイベントを
// Go アプリケーション側で捕捉するための WndProc（Window Procedure）フックを実装している。
//
// 【WndProc（Window Procedure）とは】
//
// Windows の GUI はすべて「メッセージ駆動」で動作している。
// 各ウィンドウは WndProc と呼ばれるコールバック関数を持ち、
//
//   - マウスクリック
//   - キー入力
//   - 再描画
//   - シャットダウン通知
//
// といった OS からのイベントはすべて WndProc に送られる。
//
// つまり WndProc は
//
//	「Windows → アプリケーション の唯一の入口」
//
// になっている低レベル関数である。
//
// 通常 walk が内部で WndProc を管理しているが、
// WM_QUERYENDSESSION（シャットダウン開始通知）は walk の公開APIからは捕捉できないため、
// 本アプリでは Win32 API を使って WndProc を差し替え、独自の処理を割り込ませている。
//
// 【背景】
//
// Windows ではシャットダウンやログオフが開始されると、各ウィンドウに
// WM_QUERYENDSESSION というメッセージが送られる。
// アプリケーションは WndProc 内でこのメッセージを処理し、
// シャットダウンを許可するか、一時的にブロックするかを同期的に返す必要がある。
//
// 【処理フロー】
//
//  1. Windows が WM_QUERYENDSESSION を送信
//  2. wndProcCallback がこれを横取り
//  3. ShutdownBlockReasonCreate でシャットダウン画面に理由を表示
//  4. WM_SHOW_DIALOG を自分自身に PostMessage
//  5. 一旦シャットダウンを拒否（return 0）
//  6. 通常のメッセージループ内で WM_SHOW_DIALOG を受信
//  7. Go 側の確認ダイアログを表示
//
// ※ WM_QUERYENDSESSION のハンドラ内で直接 UI を表示すると不安定になるため、
// PostMessage を使って処理を遅延させている。
//
// 【技術的制約】
//
//   - Win32 のコールバックは Go のクロージャやメソッドを保持できないため、
//     現在の App インスタンスはグローバル変数 appInstance に保持している。
//   - origWndProc には元の walk の WndProc を保存し、未処理メッセージは必ず委譲する。
//
// この構造は Win32 API の制約によるものであり、意図的にこのファイルへ隔離している。
package app

import (
	"syscall"

	"github.com/lxn/win"

	"shutdown-alert/internal/config"
	"shutdown-alert/internal/win32"
)

// appInstanceはWndProcコールバック用の現在のAppインスタンスを保持します。
// これは、WindowsコールバックがGoのクロージャをキャプチャできないために必要です。
// これはWin32 APIの制約であり、回避できません。
var appInstance *App

// origWndProcは元のウィンドウプロシージャを保持します。
var origWndProc uintptr

// wndProcCallbackはWM_QUERYENDSESSIONをインターセプトするためのカスタムウィンドウプロシージャです。
// この関数は副作用（Win32メッセージ処理）を持ちます。
func wndProcCallback(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case win32.WM_QUERYENDSESSION:
		if appInstance != nil {
			// シャットダウン画面に表示されるブロック理由を設定します。
			win32.ShutdownBlockReasonCreate(hwnd, config.ShutdownBlockMessage)

			// このハンドラから戻った後にダイアログを表示するためのメッセージをポストします。
			win32.PostMessage(hwnd, win32.WM_SHOW_DIALOG, 0, 0)

			// シャットダウンを一時的にブロックするためにFALSEを返します。
			return 0
		}
		return 1

	case win32.WM_SHOW_DIALOG:
		if appInstance != nil {
			appInstance.handleShutdownQuery()
			// ダイアログが閉じられた後、ブロック理由をクリアします。
			win32.ShutdownBlockReasonDestroy(hwnd)
		}
		return 0

	case win32.WM_ENDSESSION:
		// セッションが終了しています。必要に応じてクリーンアップします。
		return win.CallWindowProc(origWndProc, hwnd, msg, wParam, lParam)
	}

	// 元のウィンドウプロシージャを呼び出します。
	return win.CallWindowProc(origWndProc, hwnd, msg, wParam, lParam)
}

// installWndProcHookはカスタムウィンドウプロシージャをインストールします。
// この関数は副作用（ウィンドウプロシージャの変更）を持ちます。
func (app *App) installWndProcHook() {
	hwnd := app.mainWindow.Handle()
	origWndProc = win.SetWindowLongPtr(hwnd, win.GWLP_WNDPROC, syscall.NewCallback(wndProcCallback))
}
