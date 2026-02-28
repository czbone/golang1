# Hello World - Go

Go言語によるシンプルな Hello World プログラムです。

## 必要条件

- Go 1.26 以上

## 実行方法

```bash
go run main.go
```

## 出力

```
Hello World
```

## PowerShell からの実行

```powershell
# ソースを直接実行
go run main.go

# ビルドして実行
go build -o hello.exe .
.\hello.exe
```

## VS Code での実行とデバッグ

VS Code の「実行とデバッグ」パネル（`Ctrl+Shift+D`）から起動できます。

### 前提条件

- [Go 拡張機能](https://marketplace.visualstudio.com/items?itemName=golang.go) のインストール

### 起動方法

1. サイドバーの「実行とデバッグ」アイコンをクリックするか、`Ctrl+Shift+D` を押す
2. 上部のドロップダウンで **「Launch Package」** を選択
3. `F5` を押して実行（またはデバッグ）を開始する

ブレークポイントを設定してステップ実行することも可能です。

### launch.json 設定

`.vscode/launch.json` に以下の設定が定義されています。

```json
{
    "name": "Launch Package",
    "type": "go",
    "request": "launch",
    "mode": "auto",
    "program": "${workspaceFolder}"
}
```

## ビルド

```bash
go build -o hello .
```

ビルド後、生成されたバイナリを実行します。

```bash
./hello
```

## プロジェクト構成

```
.
├── go.mod    # モジュール定義
├── main.go   # エントリーポイント
└── README.md # このファイル
```
