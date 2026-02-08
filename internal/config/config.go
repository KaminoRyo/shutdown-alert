package config

const (
	// TargetURLは開く対象のURLです。
	TargetURL = "https://www.google.com" // 仕様書で指定

	// 確認ダイアログの最小幅
	DialogWidth = 400
	// 確認ダイアログの最小高さ
	DialogHeight = 150
	//確認ダイアログのタイトル
	DialogTitle = "Shutdown Alert"

	// ShutdownBlockMessageはシャットダウン画面に表示されるメッセージです。
	ShutdownBlockMessage = "確認ダイアログに応答してください"
	// DialogMessageFormatはダイアログメッセージの書式文字列です。
	// URL用のプレースホルダ%sを1つ含む必要があります。
	DialogMessageFormat = "PCをシャットダウンしようとしています。\n%s を開きますか？"

	// OpenButtonLabelは「開く」ボタンのラベルです。
	OpenButtonLabel = "開く(&O)"
	// ExitButtonLabelは終了ボタンのラベルです。
	ExitButtonLabel = "開かない(&E)"

	// TrayIconTooltipはトレイアイコンのツールチップテキストです。
	TrayIconTooltip = "Shutdown Alertが動作しています."
	// TrayMenuTestはテストダイアログメニュー項目のラベルです。
	// &文字はキーボードアクセラレータ（Alt+T）を示します。
	TrayMenuTest = "&ダイアログ表示(&T)"
	// TrayMenuStartupはスタートアップ登録メニュー項目のラベルです。
	// &文字はキーボードアクセラレータ（Alt+S）を示します。
	TrayMenuStartup = "スタートアップに登録(&S)"
	// TrayMenuExitは終了メニュー項目のラベルです。
	// &文字はキーボードアクセラレータ（Alt+E）を示します。
	TrayMenuExit = "&終了(&E)"

	// IconPathはアプリケーションのアイコンファイルパスです。
	IconPath = "internal/icon/icon.ico"

	// IconResourceIDはリソースに埋め込まれたアイコンのIDです。
	// 普通に埋め込んでるのとmanifestで指定してるのと両方使うために定義しています。
	IconResourceID = 2

	// RegistryValueNameはスタートアップ登録に使用するレジストリ値の名前です。
	RegistryValueName = "ShutdownAlert"

	// LogFileNameはエラーログファイルの名前です。
	LogFileName = "error_log.json"

	// MaxLogEntriesはログファイルに保持する最大エントリ数です。
	MaxLogEntries = 100

	// メッセージボックスのタイトル
	MessageBoxTitleError   = "エラー"
	MessageBoxTitleSuccess = "成功"

	// スタートアップ登録成功メッセージ
	StartupRegisterSuccessMessage = "スタートアップに登録しました。\nWindows起動時に自動で起動します。"
	// スタートアップ登録失敗メッセージフォーマット
	StartupRegisterErrorMessageFormat = "スタートアップ登録に失敗しました:\n%v"

	// スタートアップ解除成功メッセージ
	StartupUnregisterSuccessMessage = "スタートアップ登録を解除しました。"
	// スタートアップ解除失敗メッセージフォーマット
	StartupUnregisterErrorMessageFormat = "スタートアップ登録の解除に失敗しました:\n%v"
)
