# コードフォーマッター

## gofmt（標準）

Go に同梱されている。別途インストール不要。

```bash
# ファイルを整形して上書き
gofmt -w main.go

# ディレクトリ内を一括整形
gofmt -w .

# 差分を確認（上書きしない）
gofmt -d main.go
```

## goimports（準公式）

`gofmt` の機能に加えて **import 文の自動追加・削除** もやってくれる。実務ではこちらを使うことが多い。

```bash
# インストール（別途必要）
go install golang.org/x/tools/cmd/goimports@latest

# 実行
goimports -w main.go
```

- `go install` は `$(go env GOPATH)/bin`（デフォルト `~/go/bin`）にバイナリを置く
- PATH に `~/go/bin` が含まれていないとコマンドが見つからないので注意

```bash
# PATH が通っていない場合、.zshrc に追加
export PATH=$PATH:$(go env GOPATH)/bin
```

## Go のフォーマットの特徴

- **インデントはタブ**。スペース vs タブ論争は Go では存在しない
- **フォーマットのカスタマイズは一切できない**。Prettier のような設定ファイルはない
- `gofmt` をかけていないコードは受け入れないのが Go コミュニティの文化

設定の余地がないのは Go の設計思想: 「フォーマットの議論に時間を使うな、全員同じにしろ」という割り切り。

## エディタ連携

VSCode + Go 拡張（`golang.go`）を使っていれば、保存時に自動で `goimports` が走る。拡張が `goimports` のパスを自動検出するので、シェルの PATH 設定がなくても動く。
