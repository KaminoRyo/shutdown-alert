package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	// ---ユーザーが設定ファイルで上書き可能な設定---

	// TargetURLは開く対象のURLです。
	TargetURL = "https://www.google.com" // 仕様書で指定

	// 確認ダイアログの最小幅
	DialogWidth = 600
	// 確認ダイアログの最小高さ
	DialogHeight = 400

	// DialogMessageFormatはダイアログメッセージの書式文字列です。
	DialogMessageFormat = `PCをシャットダウンしようとしています。
	https://www.google.com を開きますか？`

	// ---上書き不可能な設定---
	//確認ダイアログのタイトル
	DialogTitle = "Shutdown Alert"

	// ShutdownBlockMessageはシャットダウン画面に表示されるメッセージです。
	ShutdownBlockMessage = "確認ダイアログに応答してください"

	// OpenButtonLabelは「開く」ボタンのラベルです。
	OpenButtonLabel = "開く(&O)"
	// ExitButtonLabelは終了ボタンのラベルです。
	ExitButtonLabel = "閉じる(&E)"

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

// UserConfig はユーザーが設定ファイルで指定可能な設定を保持します。
type UserConfig struct {
	TargetURL     string `yaml:"target_url"`
	DialogWidth   int    `yaml:"dialog_width"`
	DialogHeight  int    `yaml:"dialog_height"`
	DialogMessage string `yaml:"dialog_message"`
}

// LoadUserConfig は設定ファイルを読み込み、デフォルト値とマージした設定を返します。
// この関数は副作用（ファイル読み込み）を持ちます。
func LoadUserConfig(configPath string) (UserConfig, error) {
	// デフォルト値で初期化
	config := UserConfig{
		TargetURL:     TargetURL,
		DialogWidth:   DialogWidth,
		DialogHeight:  DialogHeight,
		DialogMessage: DialogMessageFormat,
	}

	// ファイルが存在すれば読み込んで上書き
	data, err := os.ReadFile(configPath)
	if err != nil {
		// ファイルが存在しない場合もエラーとして返す
		return config, err
	}

	// YAMLをパース（ポインタ型を使用してフィールドの存在を判定）
	var userConfig struct {
		TargetURL     *string `yaml:"target_url,omitempty"`
		DialogWidth   *int    `yaml:"dialog_width,omitempty"`
		DialogHeight  *int    `yaml:"dialog_height,omitempty"`
		DialogMessage *string `yaml:"dialog_message,omitempty"`
	}

	if err := yaml.Unmarshal(data, &userConfig); err != nil {
		return config, err
	}

	// 設定値が指定されていれば上書き（nilチェックでフィールドの存在を判定）
	if userConfig.TargetURL != nil {
		// URLのバリデーション
		if err := validateURL(*userConfig.TargetURL); err != nil {
			return config, fmt.Errorf("target_url のバリデーションエラー: %w", err)
		}
		config.TargetURL = *userConfig.TargetURL
	}
	if userConfig.DialogWidth != nil {
		// ダイアログ幅のバリデーション
		if *userConfig.DialogWidth < 0 || *userConfig.DialogWidth > 10000 {
			return config, fmt.Errorf("dialog_width は 0 から 10000 の範囲で指定してください: %d", *userConfig.DialogWidth)
		}
		config.DialogWidth = *userConfig.DialogWidth
	}
	if userConfig.DialogHeight != nil {
		// ダイアログ高さのバリデーション
		if *userConfig.DialogHeight < 0 || *userConfig.DialogHeight > 10000 {
			return config, fmt.Errorf("dialog_height は 0 から 10000 の範囲で指定してください: %d", *userConfig.DialogHeight)
		}
		config.DialogHeight = *userConfig.DialogHeight
	}
	if userConfig.DialogMessage != nil {
		// メッセージのバリデーション
		if err := validateDialogMessage(*userConfig.DialogMessage); err != nil {
			return config, fmt.Errorf("dialog_message のバリデーションエラー: %w", err)
		}
		config.DialogMessage = *userConfig.DialogMessage
	}

	return config, nil
}

// validateDialogMessage はダイアログメッセージの妥当性を検証します。
// この関数は純粋関数です。
func validateDialogMessage(message string) error {
	// 空文字列は許可しない
	if message == "" {
		return fmt.Errorf("メッセージが空です")
	}

	return nil
}

// validateURL はURLの安全性を検証します。
// この関数は純粋関数です。
func validateURL(targetURL string) error {
	// 空文字列は許可（アラートモード）
	if targetURL == "" {
		return nil
	}

	// URLのパース
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return fmt.Errorf("無効なURL形式です: %w", err)
	}

	// スキームの検証（http, httpsのみ許可）
	scheme := strings.ToLower(parsedURL.Scheme)
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("URLスキームは http または https のみ許可されています: %s", parsedURL.Scheme)
	}

	// ホスト名の検証（空でないこと）
	if parsedURL.Host == "" {
		return fmt.Errorf("URLにホスト名が指定されていません")
	}

	return nil
}
