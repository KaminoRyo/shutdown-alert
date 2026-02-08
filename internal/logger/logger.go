//go:build windows

package logger

import (
	"encoding/json"
	"os"
	"path/filepath"
	"shutdown-alert/internal/config"
	"time"
)

// LogEntry はログエントリの構造を表します。
type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Component string                 `json:"component"`
	Message   string                 `json:"message"`
	Error     string                 `json:"error,omitempty"`
	Context   map[string]interface{} `json:"context,omitempty"`
}

// LogError はエラーをJSONファイルに記録します。
// この関数は副作用を持ちます（ファイルへの書き込み）。
func LogError(component, message string, err error, context map[string]interface{}) {
	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     "ERROR",
		Component: component,
		Message:   message,
		Context:   context,
	}

	if err != nil {
		entry.Error = err.Error()
	}

	// ログファイルのパスを取得
	logPath, pathErr := getLogFilePath()
	if pathErr != nil {
		// ログファイルパスの取得に失敗しても、アプリケーションは続行
		return
	}

	// 既存のログを読み込む
	entries := readExistingLogs(logPath)

	// 新しいエントリを追加
	entries = append(entries, entry)

	// 最大エントリ数を超えた場合、古いエントリを削除
	if len(entries) > config.MaxLogEntries {
		entries = entries[len(entries)-config.MaxLogEntries:]
	}

	// JSONファイルに書き込む
	writeLogsToFile(logPath, entries)
}

// getLogFilePath は実行ファイルと同じディレクトリのログファイルパスを返します。
// この関数は副作用を持ちます（ファイルシステムへのアクセス）。
func getLogFilePath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", err
	}

	execDir := filepath.Dir(execPath)
	logPath := filepath.Join(execDir, config.LogFileName)

	return logPath, nil
}

// readExistingLogs は既存のログファイルを読み込みます。
// この関数は副作用を持ちます（ファイルの読み取り）。
func readExistingLogs(logPath string) []LogEntry {
	var entries []LogEntry

	data, err := os.ReadFile(logPath)
	if err != nil {
		// ファイルが存在しない場合は空のスライスを返す
		return entries
	}

	// JSONをパース（失敗しても空のスライスを返す）
	_ = json.Unmarshal(data, &entries)

	return entries
}

// writeLogsToFile はログエントリをJSONファイルに書き込みます。
// この関数は副作用を持ちます（ファイルへの書き込み）。
func writeLogsToFile(logPath string, entries []LogEntry) {
	// JSON形式でエンコード（インデント付き）
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return
	}

	// ファイルに書き込む（失敗しても続行）
	_ = os.WriteFile(logPath, data, 0644)
}
