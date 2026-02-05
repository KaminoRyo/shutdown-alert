//go:build windows

package main

import (
	"log"

	"shutdown-alert/internal/app"
)

func main() {
	a := app.NewApp()
	if err := a.Run(); err != nil {
		log.Fatalf("アプリケーションの実行に失敗しました: %v", err)
	}
}
