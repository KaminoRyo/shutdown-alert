# Shutdown Alert

> [!WARNING]
> 制作中です


## 開発環境のセットアップ

### Goのインストール

このプロジェクトはGo言語でWindowsネイティブアプリケーションを開発するため、ローカル環境へのGoのインストールが必要です。一般的なWebアプリケーション開発でGoをDockerコンテナ内で利用するケースもありますが、本プロジェクトではWindowsのGUIやシャットダウンイベントを直接操作するため、ホストOSにGoを直接インストールする方法を推奨します。

以下の手順でGoをインストールしてください。

1.  **Go公式サイトにアクセス**:
    [https://go.dev/dl/](https://go.dev/dl/)

2.  **Windows向けインストーラーのダウンロード**:
    ページの中から "Microsoft Windows" 向けの `.msi` ファイルをダウンロードします。

3.  **インストールの実行**:
    ダウンロードしたインストーラーを実行し、画面の指示に従って進めます。通常、デフォルト設定で問題ありません。

4.  **インストールの確認**:
    コマンドプロンプトまたはPowerShellを開き、以下のコマンドを実行してGoのバージョンが表示されればインストール完了です。

    ```shell
    go version
    ```

    例: `go version go1.xx.x windows/amd64`

### rsrcツールのインストール

マニフェストとアイコンを実行ファイルに埋め込むために、`rsrc`ツールが必要です。

```shell
go install github.com/akavel/rsrc@latest
```

## ビルド方法

### 1. リソースファイルの生成

マニフェストとアイコンを埋め込むための`.syso`ファイルを生成します。

```shell
rsrc -manifest shutdown-alert.manifest -ico internal/icon/icon.ico -o rsrc.syso
```

### 2. アプリケーションのビルド

GUIモードでビルドします（コンソールウィンドウを非表示）。

```shell
go build -ldflags "-H windowsgui" -o shutdown-alert.exe .
```

### デバッグビルド

コンソール出力を確認したい場合は、`-ldflags`を省略します。

```shell
go build -o shutdown-alert-debug.exe .
```

## 使い方

1. `shutdown-alert.exe`をダブルクリックして起動
2. タスクトレイ（通知領域）にアイコンが表示されます
3. シャットダウン/ログオフ時に確認ダイアログが表示されます
4. 終了するには、タスクトレイのアイコンを右クリックして「Exit」を選択
