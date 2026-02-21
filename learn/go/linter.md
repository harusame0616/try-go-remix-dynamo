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

#### **【推奨】go tool ディレクティブを使う方法（Go 1.24+）**

Go 1.24 から導入された `go tool` ディレクティブを使うことで、JavaScript の `devDependencies` に相当する仕組みが利用できます。これにより、ツールのバージョンが `go.mod` に記録され、全開発者が同じバージョンを使えるようになります。

```bash
# golangci-lint をプロジェクトに追加（go.mod に記録される）
go get -tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.10.1

# 実行（go.mod に記録されたバージョンが使われる）
go tool golangci-lint run ./...
```

**メリット**:
- go.mod でバージョンが固定され、チーム全員が同じバージョンを使える
- グローバル環境を汚染しない
- JavaScript の `package.json` の `devDependencies` と同様の感覚で管理できる
- goimports など他のツールも同様に管理可能（例：`go get -tool golang.org/x/tools/cmd/goimports@latest`）

**実行方法**:
```bash
# カレントディレクトリとサブディレクトリを全て解析
go tool golangci-lint run ./...

# カレントディレクトリのみ解析
go tool golangci-lint run

# 特定ディレクトリを解析
go tool golangci-lint run ./path/to/...

# 修正可能な問題を自動修正（一部の Linter のみ）
go tool golangci-lint run --fix
```

#### 従来の方法：go install でグローバルインストール

Go 1.24 以前や、グローバルにインストールしたい場合は `go install` を使用します。

**v1（従来版）**

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

**v2（最新版）**

```bash
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
```

v2 では設定ファイルの構造が変更されており、より明確に linter とフォーマッターが分離されています。新規プロジェクトでは v2 の使用を推奨します。

**実行方法**:
```bash
# カレントディレクトリとサブディレクトリを全て解析
golangci-lint run ./...

# カレントディレクトリのみ解析
golangci-lint run

# 特定ディレクトリを解析
golangci-lint run ./path/to/...

# 修正可能な問題を自動修正（一部の Linter のみ）
golangci-lint run --fix
```

**注意点**:
- `go install` はツールを `$GOPATH/bin` または `$GOBIN` にグローバルインストールする
- プロジェクトごとにバージョンを固定できないため、チーム開発では `go tool` ディレクティブの使用を推奨

### デフォルトで有効な Linter

golangci-lint は設定ファイルなしでも以下の Linter がデフォルトで有効になっています：

- **errcheck**: エラー戻り値の未チェック検出
- **govet**: Go の疑わしいコード構造を検出
- **ineffassign**: 無効な代入の検出
- **staticcheck**: 静的解析によるバグ検出
- **unused**: 未使用のコード検出（v1 の deadcode の役割も含む）

### 設定ファイル（.golangci.yml）

プロジェクトのルートに `.golangci.yml` を作成して、Linter を設定します。

#### v1 の設定例

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

#### v2 の設定例（重要な変更点あり）

v2 では `version: "2"` の指定が必須で、linters と formatters が分離されています。

```yaml
# v2 の設定であることを明示（必須）
version: "2"

# Linter の設定
linters:
  enable:
    - errcheck      # error 処理の検出（必須）
    - govet         # Go の疑わしいコード検出（必須）
    - ineffassign   # 無効な代入（推奨）
    - staticcheck   # 一般的なバグ検出（必須）
    - unused        # 未使用コード検出（推奨）
    - revive        # コーディング規約（推奨）
    - gosec         # セキュリティ（セキュリティが重要な場合）
    - unparam       # 未使用パラメータ（推奨）

# フォーマッターの設定（v2 の新機能）
formatters:
  enable:
    - gofmt         # 標準フォーマッター
    - goimports     # import の整理

# Linter ごとの設定
linters-settings:
  revive:
    rules:
      # exported（エクスポート関数のドキュメントコメント）の重要度を下げる
      - name: exported
        severity: warning

  staticcheck:
    checks:
      - all

# 実行設定
run:
  timeout: 5m
```

#### v1 と v2 の主な違い

1. **version フィールドの追加**: v2 では `version: "2"` の明示が必須
2. **formatters セクションの分離**: v1 では `linters.enable` に gofmt や goimports を含めていたが、v2 では `formatters.enable` に分離
3. **unused の役割拡大**: v2 の unused は v1 の deadcode の機能も統合
4. **設定の明確化**: linter（静的解析）とフォーマッター（コード整形）が概念的に分離され、より明確に

### 推奨 Linter（最低限有効にすべき）

| Linter | 目的 | 優先度 | デフォルト有効 |
|--------|------|--------|---------------|
| govet | 基本的なバグ検出 | 必須 | ○ |
| staticcheck | 一般的なバグ・パフォーマンス検出 | 必須 | ○ |
| errcheck | error 処理の検出 | 必須 | ○ |
| unused | 未使用コード検出 | 必須 | ○ |
| ineffassign | 無効な代入 | 推奨 | ○ |
| revive | コーディング規約違反 | 推奨 | × |
| unparam | 未使用パラメータ | 推奨 | × |
| gosec | セキュリティ脆弱性 | セキュリティが重要な場合 | × |

### 主要な Linter の役割と検出例

#### errcheck

エラー戻り値の未チェックを検出します。Go ではエラー処理が重要で、これを見落とすとバグになります。

検出例：
```go
// NG: Close() の戻り値を無視
resp.Body.Close()

// OK: 明示的に無視（戻り値があることを認識している）
_, _ = resp.Body.Close()

// OK: defer で適切に処理
defer func() {
    _ = resp.Body.Close()
}()
```

#### govet

Go の疑わしいコード構造を検出します。

検出例：
- printf のフォーマット不一致
- 常に false の条件式
- 型の不整合
- 構造体のコピーによる問題

#### ineffassign

無効な代入（後で上書きされるだけの変数への代入）を検出します。

検出例：
```go
// NG: value は後で上書きされるので、最初の代入は無駄
value := 10
value = 20
fmt.Println(value)

// OK
value := 20
fmt.Println(value)
```

#### staticcheck

静的解析によるバグ検出を行います。golangci-lint のコアメンバーが開発している信頼性の高いツールです。

検出例：
- バグになりうるコードパターン
- パフォーマンス問題（不要なアロケーション）
- 非推奨 API の使用
- シンプルコードへのリファクタリング提案

#### unused

未使用のコード（変数、関数、型、定数など）を検出します。v2 では v1 の deadcode の機能も統合されています。

検出例：
```go
// NG: 未使用の変数
func example() {
    unused := 10  // 使われていない
    used := 20
    fmt.Println(used)
}

// NG: 未使用の関数
func neverCalled() {  // どこからも呼ばれていない
    // ...
}
```

#### revive

golint の後継で、コーディングスタイルチェックを行います。

検出例：
```go
// NG: パッケージコメントがない
package main

// OK: パッケージコメントあり
// Package main provides the entry point for the application.
package main

// NG: エクスポートされた関数にコメントがない
func PublicFunction() {
}

// OK: エクスポートされた関数にコメントあり
// PublicFunction does something important.
func PublicFunction() {
}

// NG: 未使用のパラメータ
func handler(w http.ResponseWriter, r *http.Request) {
    // r を使っていない
    fmt.Fprint(w, "hello")
}

// OK: 未使用のパラメータは _ にリネーム
func handler(w http.ResponseWriter, _ *http.Request) {
    fmt.Fprint(w, "hello")
}
```

### 実際の修正例

以下は、golangci-lint v2 で検出された問題と修正内容の実例です。

#### 1. errcheck: fmt.Fprint の戻り値未チェック

```go
// 検出された問題
func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Health check OK")  // Error return value is not checked
}

// 修正後
func handler(w http.ResponseWriter, _ *http.Request) {
    _, _ = fmt.Fprint(w, "Health check OK")  // 明示的に無視
}
```

#### 2. errcheck: resp.Body.Close() の戻り値未チェック

```go
// 検出された問題
resp, err := http.Get("http://localhost:8080/health")
if err != nil {
    t.Fatal(err)
}
defer resp.Body.Close()  // Error return value is not checked

// 修正後
resp, err := http.Get("http://localhost:8080/health")
if err != nil {
    t.Fatal(err)
}
defer func() {
    _ = resp.Body.Close()  // defer 内で明示的に無視
}()
```

#### 3. revive/package-comments: パッケージコメントがない

```go
// 検出された問題
package main

// 修正後
// Package main provides the entry point for the application.
package main
```

#### 4. revive/unused-parameter: 未使用パラメータ

```go
// 検出された問題
func handler(w http.ResponseWriter, r *http.Request) {
    // r を使っていない
    fmt.Fprint(w, "hello")
}

// 修正後
func handler(w http.ResponseWriter, _ *http.Request) {
    // 未使用のパラメータは _ にリネーム
    _, _ = fmt.Fprint(w, "hello")
}
```

## 実務でのベストプラクティス

### 1. 個別 Linter を手動で管理しない

複数の Linter を個別に入れるのは管理が大変です。必ず golangci-lint を使って一括管理します。

```bash
# これはやらない（個別インストール）
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/kisielk/errcheck@latest
# ... 以下、各 Linter をインストール

# これをやる（推奨：go tool ディレクティブで管理）
go get -tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.10.1
# .golangci.yml で全て設定

# または（従来の方法）
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
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

v2 を使う場合は、golangci-lint のバージョンを明示的に指定します：

```yaml
- uses: golangci/golangci-lint-action@v3
  with:
    version: v2.0.0  # v2 のバージョンを指定
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
go tool golangci-lint run ./...  # go tool ディレクティブを使う場合
# または
# golangci-lint run ./...  # go install でグローバルインストールした場合
```

これにより、問題のあるコードをコミットする前に検出できます。

### 6. Makefile を使ったタスク管理（Optional）

Go プロジェクトでは、`make lint` のようなタスクランナーとして Makefile を使うことがあります。特に CI/CD パイプラインで複数のコマンドを実行する場合や、チーム内でコマンドを統一したい場合に便利です。

```makefile
# Makefile
.PHONY: lint
lint:
	go tool golangci-lint run ./...

.PHONY: lint-fix
lint-fix:
	go tool golangci-lint run --fix ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: fmt
fmt:
	go tool goimports -w .
	go fmt ./...
```

実行例：
```bash
make lint      # Linter を実行
make lint-fix  # Linter を実行して自動修正
make test      # テストを実行
make fmt       # フォーマットを実行
```

ただし、**小規模プロジェクトでは Makefile を使わず、直接コマンドを実行する方がシンプルで十分**です。Makefile は複雑なビルドプロセスや複数のコマンドを組み合わせる必要がある場合に検討しましょう。

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
