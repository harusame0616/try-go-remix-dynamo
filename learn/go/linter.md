# Go の Linter

Linter は静的解析ツールで、実行せずにコードの問題を検出します。Go にはコンパイラが型チェックを厳格に行うため、他の言語ほどは Linter に頼らなくても大丈夫ですが、バグの早期発見やコード品質の向上に役立ちます。

## 標準ツール

### go vet

Go に標準搭載の静的解析ツール。`go build` の時点でも一部のチェックが走ります。

```bash
go vet ./...
```

以下のような怪しいコードを検出：
- 未使用の変数
- printf のフォーマット不一致
- 常に false の条件式
- 型の不整合

**特徴**
- インストール不要（Go に付属）
- 比較的高速
- Go の設計思想に沿ったチェックに特化

## サードパーティ Linter

### staticcheck

最も信頼性の高い静的解析ツール。golangci-lint のコアメンバーが開発しており、バグ検出の正確性が高い。

```bash
go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck ./...
```

以下を検出：
- バグになりうるコードパターン
- パフォーマンス問題（例：不要なアロケーション）
- 非推奨 API の使用
- シンプルコードへのリファクタリング提案

### errcheck

戻り値の error を無視しているコードを検出。Go では error 処理が重要で、これを見落とすとバグになります。

```bash
go install github.com/kisielk/errcheck@latest
errcheck ./...
```

例：
```go
// 以下は errcheck で検出される
file.Close()  // Close() の error 返り値を無視している
```

### revive

非推奨になった golint の後継。コーディング規約違反を検出し、ルールのカスタマイズが柔軟。

```bash
go install github.com/mgechev/revive@latest
revive ./...
```

検出例：
- エクスポートされた関数にドキュメントコメントがない
- interface が 1 メソッドしかない（Reader, Writer など）
- 名前付き戻り値を使いすぎている
- 関数が長すぎる

### gosec

セキュリティ脆弱性に特化した Linter。

```bash
go install github.com/securego/gosec/v2/cmd/gosec@latest
gosec ./...
```

検出例：
- SQL インジェクション
- ハードコードされた認証情報
- ユーザー入力の不適切な使用
- 弱い暗号化

### ineffassign

無効な代入（後で上書きされるだけの変数への代入）を検出。

```bash
go install github.com/gordonklaus/ineffassign@latest
ineffassign ./...
```

### deadcode

到達不能コード（デッドコード）を検出。削除できるコードを見つけられます。

```bash
go install golang.org/x/tools/cmd/deadcode@latest
deadcode ./...
```

### unparam

使われていない関数パラメータを検出。

```bash
go install mvdan.cc/unparam@latest
unparam ./...
```

## メタ Linter: golangci-lint

100 以上の Linter をまとめて実行できるメタ Linter。Go のデファクトスタンダード。個別に Linter を入れるのではなく、これを使うのが実務的です。

### インストール

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### 実行

```bash
# カレントディレクトリを解析
golangci-lint run

# 特定ディレクトリを解析
golangci-lint run ./path/to/...

# 修正可能な問題を自動修正（一部の Linter のみ）
golangci-lint run --fix
```

### 設定ファイル（.golangci.yml）

プロジェクトのルートに `.golangci.yml` を作成して、Linter を設定：

```yaml
# 有効にする Linter を指定
linters:
  enable:
    - errcheck      # error 処理の検出（必須）
    - staticcheck   # 一般的なバグ検出（必須）
    - revive        # コーディング規約（推奨）
    - gosec         # セキュリティ（セキュリティが重要な場合）
    - ineffassign   # 無効な代入（推奨）
    - unparam       # 未使用パラメータ（推奨）

# Linter ごとの設定
linters-settings:
  revive:
    rules:
      # exported（エクスポート関数のドキュメントコメント）の重要度を下げる
      - name: exported
        severity: warning

  staticcheck:
    # SA チェック（バグ）は全て有効
    # ST チェック（スタイル）は一部無効にする場合
    checks:
      - all
      - -SA3000  # 例：特定チェックを除外

# 実行設定
run:
  timeout: 5m
  deadline: 10m

# 出力フォーマット
output:
  format: github-actions
```

### 推奨 Linter（最低限有効にすべき）

| Linter | 目的 | 優先度 |
|--------|------|--------|
| go vet | 基本的なバグ検出 | 必須（デフォルト有効） |
| staticcheck | 一般的なバグ・パフォーマンス検出 | 必須 |
| errcheck | error 処理の検出 | 必須 |
| revive | コーディング規約違反 | 推奨 |
| ineffassign | 無効な代入 | 推奨 |
| unparam | 未使用パラメータ | 推奨 |
| gosec | セキュリティ脆弱性 | セキュリティが重要な場合 |

## 実務でのベストプラクティス

### 1. 個別 Linter を手動で管理しない

複数の Linter を個別に入れるのは管理が大変です。必ず golangci-lint を使って一括管理します。

```bash
# これはやらない
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/kisielk/errcheck@latest
# ... 以下、各 Linter をインストール

# これをやる
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
# .golangci.yml で全て設定
```

### 2. CI/CD パイプラインに組み込む

GitHub Actions の例：

```yaml
name: Lint
on: [push, pull_request]

jobs:
  golangci:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
        with:
          go-version: '1.22'
      - uses: golangci/golangci-lint-action@v3
        with:
          version: latest
          args: --timeout=5m
```

### 3. 新規プロジェクトは最初から Linter を導入する

既存のコードに Linter を後から導入すると、大量のエラーが出て対応が大変になります。プロジェクト開始時から導入するのが最適です。

### 4. チーム内で Linter ルールを統一する

`.golangci.yml` をバージョン管理に含めて、チーム全員が同じルールで実装するようにします。

### 5. コミットフックで自動実行（Optional）

```bash
# .git/hooks/pre-commit
#!/bin/bash
golangci-lint run ./...
```

これにより、問題のあるコードをコミットする前に検出できます。

## go vet と golangci-lint の関係

- `go vet` は Go に標準搭載で、比較的高速
- `golangci-lint run` は `go vet` を含めて複数の Linter を実行
- 実務では golangci-lint を使えば、go vet は別途実行する必要がない

```bash
# go vet だけ
go vet ./...

# go vet + 他の Linter（推奨）
golangci-lint run
```
