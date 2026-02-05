//go:build windows

package ui

import (
	"fmt"

	"github.com/lxn/walk"

	"shutdown-alert/internal/config"
)

// InitNotifyIconは通知アイコンを作成して設定します。
// この関数は副作用（UI要素の作成）を持ちます。
func InitNotifyIcon(mainWindow *walk.MainWindow, onTest, onExit func()) (*walk.NotifyIcon, error) {
	icon, err := walk.NewIconFromFile(config.IconPath)
	if err != nil {
		return nil, fmt.Errorf("アイコンの読み込みに失敗しました: %w", err)
	}

	notifyIcon, err := walk.NewNotifyIcon(mainWindow)
	if err != nil {
		return nil, fmt.Errorf("通知アイコンの作成に失敗しました: %w", err)
	}

	if err := notifyIcon.SetIcon(icon); err != nil {
		return nil, fmt.Errorf("アイコンの設定に失敗しました: %w", err)
	}
	if err := notifyIcon.SetToolTip(config.TrayIconTooltip); err != nil {
		return nil, fmt.Errorf("ツールチップの設定に失敗しました: %w", err)
	}

	// テストアクションを作成します（確認ダイアログのテスト用）。
	testAction := walk.NewAction()
	if err := testAction.SetText(config.TrayMenuTest); err != nil {
		return nil, fmt.Errorf("テストテキストの設定に失敗しました: %w", err)
	}
	testAction.Triggered().Attach(func() {
		if onTest != nil {
			onTest()
		}
	})

	// 終了アクションを作成します。
	exitAction := walk.NewAction()
	if err := exitAction.SetText(config.TrayMenuExit); err != nil {
		return nil, fmt.Errorf("終了テキストの設定に失敗しました: %w", err)
	}
	exitAction.Triggered().Attach(func() {
		if onExit != nil {
			onExit()
		}
	})

	// コンテキストメニューにアクションを追加します。
	if err := notifyIcon.ContextMenu().Actions().Add(testAction); err != nil {
		return nil, fmt.Errorf("テストアクションの追加に失敗しました: %w", err)
	}
	if err := notifyIcon.ContextMenu().Actions().Add(exitAction); err != nil {
		return nil, fmt.Errorf("終了アクションの追加に失敗しました: %w", err)
	}

	if err := notifyIcon.SetVisible(true); err != nil {
		return nil, fmt.Errorf("表示設定に失敗しました: %w", err)
	}

	return notifyIcon, nil
}
