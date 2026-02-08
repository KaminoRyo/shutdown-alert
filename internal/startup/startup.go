//go:build windows

package startup

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"

	"shutdown-alert/internal/config"
	"shutdown-alert/internal/logger"
)

const (
	// registryKeyPathはスタートアップ登録に使用するレジストリキーのパスです。
	registryKeyPath = `Software\Microsoft\Windows\CurrentVersion\Run`
)

// IsRegistered はアプリケーションがスタートアップに登録されているかを確認します。
// この関数は副作用を持ちます（レジストリの読み取り）。
func IsRegistered() bool {
	key, err := registry.OpenKey(registry.CURRENT_USER, registryKeyPath, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer key.Close()

	registeredPath, _, err := key.GetStringValue(config.RegistryValueName)
	if err != nil {
		return false
	}

	// 現在の実行ファイルのパスを取得
	executablePath, err := getExecutablePath()
	if err != nil {
		return false
	}

	// 引用符付きパスと比較
	quotedPath := fmt.Sprintf(`"%s"`, executablePath)

	// 引用符付きまたは引用符なしのどちらでも一致すればtrueを返す
	// （後方互換性のため）
	return registeredPath == quotedPath || registeredPath == executablePath
}

// Register はアプリケーションをスタートアップに登録します。
// この関数は副作用を持ちます（レジストリへの書き込み）。
func Register() error {
	executablePath, err := getExecutablePath()
	if err != nil {
		return fmt.Errorf("実行ファイルのパス取得に失敗しました: %w", err)
	}

	// パスを引用符で囲む（スペースを含むパスやWindowsのベストプラクティスに対応）
	quotedPath := fmt.Sprintf(`"%s"`, executablePath)

	key, err := registry.OpenKey(registry.CURRENT_USER, registryKeyPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("レジストリキーのオープンに失敗しました: %w", err)
	}
	defer key.Close()

	err = key.SetStringValue(config.RegistryValueName, quotedPath)
	if err != nil {
		return fmt.Errorf("レジストリ値の設定に失敗しました: %w", err)
	}

	return nil
}

// Unregister はアプリケーションをスタートアップから解除します。
// この関数は副作用を持ちます（レジストリからの削除）。
func Unregister() error {
	key, err := registry.OpenKey(registry.CURRENT_USER, registryKeyPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("レジストリキーのオープンに失敗しました: %w", err)
	}
	defer key.Close()

	err = key.DeleteValue(config.RegistryValueName)
	if err != nil && err != registry.ErrNotExist {
		return fmt.Errorf("レジストリ値の削除に失敗しました: %w", err)
	}

	return nil
}

// getExecutablePath は現在の実行ファイルの絶対パスを取得します。
// この関数は純粋関数ではありません（ファイルシステムへのアクセス）。
func getExecutablePath() (string, error) {
	executablePath, err := os.Executable()
	if err != nil {
		return "", err
	}

	// シンボリックリンクを解決して実際のパスを取得
	realPath, err := filepath.EvalSymlinks(executablePath)
	if err != nil {
		return "", err
	}

	return realPath, nil
}

// UpdateIfNeeded はスタートアップ登録パスが現在のパスと異なる場合に自動更新します。
// この関数は副作用を持ちます（レジストリの読み書き、ログファイルへの書き込み）。
func UpdateIfNeeded() error {
	// レジストリに登録されているか確認
	key, err := registry.OpenKey(registry.CURRENT_USER, registryKeyPath, registry.QUERY_VALUE)
	if err != nil {
		// 登録されていない場合は何もしない
		return nil
	}
	defer key.Close()

	registeredPath, _, err := key.GetStringValue(config.RegistryValueName)
	if err != nil {
		// 登録されていない場合は何もしない
		return nil
	}

	// 現在の実行ファイルのパスを取得
	currentPath, err := getExecutablePath()
	if err != nil {
		logger.LogError("startup", "実行ファイルのパス取得に失敗しました", err, nil)
		return err
	}

	// 引用符付きパス
	quotedCurrentPath := fmt.Sprintf(`"%s"`, currentPath)

	// パスが一致している場合は何もしない
	if registeredPath == quotedCurrentPath || registeredPath == currentPath {
		return nil
	}

	// パスが異なる場合は自動更新
	err = Register()
	if err != nil {
		logger.LogError("startup", "スタートアップパスの自動更新に失敗しました", err, map[string]interface{}{
			"old_path": registeredPath,
			"new_path": quotedCurrentPath,
		})
		return err
	}

	// 成功時はログ不要（エラーのみ記録）
	return nil
}
