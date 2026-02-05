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
	TrayMenuTest = "&Test Dialog"
	// TrayMenuExitは終了メニュー項目のラベルです。
	// &文字はキーボードアクセラレータ（Alt+E）を示します。
	TrayMenuExit = "&Exit"

	// IconPathはアプリケーションのアイコンファイルパスです。
	IconPath = "internal/icon/icon.ico"
)
