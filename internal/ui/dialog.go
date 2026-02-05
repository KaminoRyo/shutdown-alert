//go:build windows

package ui

import (
	"fmt"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"

	"shutdown-alert/internal/config"
)

// ShowConfirmationDialogはシャットダウン確認ダイアログを表示します。
// この関数は副作用（UIの表示、アプリケーションの終了の可能性）を持ちます。
// ユーザーの選択と発生したエラーを返します。
func ShowConfirmationDialog(owner walk.Form, urlToOpen string, onOpen, onExit func()) error {
	var dlg *walk.Dialog
	var openBtn, exitBtn *walk.PushButton

	_, err := declarative.Dialog{
		AssignTo:      &dlg,
		Title:         config.DialogTitle,
		DefaultButton: &openBtn,
		CancelButton:  &exitBtn,
		MinSize:       declarative.Size{Width: config.DialogWidth, Height: config.DialogHeight},
		Layout:        declarative.VBox{},
		Children: []declarative.Widget{
			declarative.Label{
				Text: fmt.Sprintf(config.DialogMessageFormat, urlToOpen),
			},
			declarative.Composite{
				Layout: declarative.HBox{},
				Children: []declarative.Widget{
					declarative.HSpacer{},
					declarative.PushButton{
						AssignTo: &openBtn,
						Text:     config.OpenButtonLabel,
						OnClicked: func() {
							if onOpen != nil {
								onOpen()
							}
							dlg.Accept()
						},
					},
					declarative.PushButton{
						AssignTo: &exitBtn,
						Text:     config.ExitButtonLabel,
						OnClicked: func() {
							if onExit != nil {
								onExit()
							}
							dlg.Accept()
						},
					},
				},
			},
		},
	}.Run(owner)

	return err
}
