//go:build windows

package ui

import (
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"

	"shutdown-alert/internal/config"
)

// ShowConfirmationDialogはシャットダウン確認ダイアログを表示します。
// この関数は副作用（UIの表示、アプリケーションの終了の可能性）を持ちます。
// ユーザーの選択と発生したエラーを返します。
func ShowConfirmationDialog(owner walk.Form, urlToOpen string, dialogWidth, dialogHeight int, dialogMessage string, onOpen, onExit func()) error {
	var dlg *walk.Dialog
	var openBtn, exitBtn *walk.PushButton

	// URLの有無によってボタンの構成を決定します。
	var buttons []declarative.Widget
	buttons = append(buttons, declarative.HSpacer{})

	if urlToOpen != "" {
		// URLがある場合：「開く」と「閉じる」両方のボタンを表示
		buttons = append(buttons, declarative.PushButton{
			AssignTo: &openBtn,
			Text:     config.OpenButtonLabel,
			OnClicked: func() {
				if onOpen != nil {
					onOpen()
				}
				dlg.Accept()
			},
		})
	}

	// 「閉じる」ボタンは常に表示
	buttons = append(buttons, declarative.PushButton{
		AssignTo: &exitBtn,
		Text:     config.ExitButtonLabel,
		OnClicked: func() {
			if onExit != nil {
				onExit()
			}
			dlg.Accept()
		},
	})

	// デフォルトボタンとキャンセルボタンの設定
	var defaultButton **walk.PushButton
	if urlToOpen != "" {
		defaultButton = &openBtn
	} else {
		defaultButton = &exitBtn
	}

	_, err := declarative.Dialog{
		AssignTo:      &dlg,
		Title:         config.DialogTitle,
		DefaultButton: defaultButton,
		CancelButton:  &exitBtn,
		MinSize:       declarative.Size{Width: dialogWidth, Height: dialogHeight},
		Layout:        declarative.VBox{},
		Children: []declarative.Widget{
			declarative.Label{
				Text: dialogMessage,
			},
			declarative.Composite{
				Layout:   declarative.HBox{},
				Children: buttons,
			},
		},
	}.Run(owner)

	return err
}
