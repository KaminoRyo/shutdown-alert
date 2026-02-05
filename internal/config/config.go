package config

// TargetURL is the URL to be opened.
const TargetURL = "https://www.google.com" // As specified in the requirements document

// Dialog configuration
const (
	// DialogWidth is the minimum width of the confirmation dialog.
	DialogWidth = 400
	// DialogHeight is the minimum height of the confirmation dialog.
	DialogHeight = 150
	// DialogTitle is the title of the confirmation dialog.
	DialogTitle = "Shutdown Alert"
)

// Messages
const (
	// ShutdownBlockMessage is the message displayed on the shutdown screen.
	ShutdownBlockMessage = "確認ダイアログに応答してください"
	// DialogMessageFormat is the format string for the dialog message.
	// It should contain one %s placeholder for the URL.
	DialogMessageFormat = "PCをシャットダウンしようとしています。\n%s を開きますか？"
)

// Button labels
const (
	// OpenButtonLabel is the label for the "Open" button.
	OpenButtonLabel = "開く(&O)"
	// ExitButtonLabel is the label for the "Exit" button.
	ExitButtonLabel = "開かない(&E)"
)

// Tray icon configuration
const (
	// TrayIconTooltip is the tooltip text for the tray icon.
	TrayIconTooltip = "Shutdown Alert is running."
	// TrayMenuTest is the label for the test dialog menu item.
	// The & character indicates the keyboard accelerator (Alt+T).
	TrayMenuTest = "&Test Dialog"
	// TrayMenuExit is the label for the exit menu item.
	// The & character indicates the keyboard accelerator (Alt+E).
	TrayMenuExit = "&Exit"
)
