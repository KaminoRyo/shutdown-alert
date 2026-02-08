//go:build windows

package win32

// Win32メッセージ定数
const (
	WM_QUERYENDSESSION = 0x0011      // セッション終了の問い合わせ
	WM_ENDSESSION      = 0x0016      // セッション終了
	WM_USER            = 0x0400      // ユーザー定義メッセージの開始
	WM_SHOW_DIALOG     = WM_USER + 1 // ダイアログ表示用のカスタムメッセージ
)

// ShellExecute定数
const (
	SW_SHOWNORMAL = 1 // ウィンドウを通常表示
)

// Mutex定数
const (
	ERROR_ALREADY_EXISTS = 183 // 既に存在するエラーコード
)
