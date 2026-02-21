# Go モジュール（go.mod）

## go mod init の基本

`go mod init <モジュール名>` でモジュールを初期化する。

```bash
go mod init github.com/harusame0616/try-go-remix-dynamo/apps/api
```

慣例として **リポジトリのパスをモジュール名に使う**。これにより import パスとディレクトリ構造が一致し、`go get` でも正しく取得できる。

## モジュールパスの選び方

モノレポで Go と非 Go（TypeScript 等）が混在する場合、リポジトリルートではなく **Go コードが存在するディレクトリのパス** を使う。

```
# 適切（Go のスコープを apps/api に閉じる）
github.com/harusame0616/try-go-remix-dynamo/apps/api

# 不適切（リポジトリ全体を Go モジュールとして扱ってしまう）
github.com/harusame0616/try-go-remix-dynamo
```

## モジュールパスが影響する 4 つの要素

### 1. import パス

モジュール名がそのまま `import` 文のプレフィックスになる。

```go
// モジュール名が .../apps/api の場合 → ディレクトリ構造と一致する
import "github.com/harusame0616/try-go-remix-dynamo/apps/api/handler"
import "github.com/harusame0616/try-go-remix-dynamo/apps/api/model"

// モジュール名がリポジトリルートの場合 → apps/api/ の中にあるのに import パスに含まれず乖離する
import "github.com/harusame0616/try-go-remix-dynamo/handler"
import "github.com/harusame0616/try-go-remix-dynamo/model"
```

### 2. go.mod の配置場所とスコープ

| パターン | go.mod の場所 | Go が管理するスコープ |
|---|---|---|
| `.../apps/api` | `apps/api/go.mod` | `apps/api/` 以下のみ |
| `.../try-go-remix-dynamo` | リポジトリルートの `go.mod` | リポジトリ全体 |

リポジトリルートに `go.mod` を置くと、Go のツール（`go build`, `go test` 等）がリポジトリ全体を Go プロジェクトとして扱おうとする。TypeScript の `apps/web` があるリポジトリでこれをやると不要な混乱が起きる。

### 3. go get / go install

他の人がモジュールを取得する場合、モジュール名とリポジトリ内の実際のパスが一致していないと `go get` が正しく動かない。

```bash
go get github.com/harusame0616/try-go-remix-dynamo/apps/api
```

### 4. マルチモジュール構成（go.work）

将来 `packages/` に Go の共有ライブラリを追加したくなったとき：

- **モジュール名が `.../apps/api` の場合** → `packages/shared` にも独立した `go.mod` を置き、`go.work` で束ねられる
- **モジュール名がルートの場合** → 全部が 1 つのモジュールに含まれ、分離が難しくなる

## モノレポ（Go + 非 Go 混在）でのベストプラクティス

Go と TypeScript が混在するリポジトリでは、**Go コードのディレクトリにスコープを閉じる**のが正解。

1. **Go のスコープを閉じる** — リポジトリ全体を Go モジュールにすると、TypeScript プロジェクトにまで Go ツールの影響が及ぶ
2. **パスの一致** — モジュール名と実際のディレクトリ構造が一致し、import パスの可読性が上がる
3. **将来の拡張** — 複数の Go モジュールが必要になった場合、Go Workspaces（`go.work`）で管理できる

## go.sum とは

go.sum は **依存パッケージのチェックサム（ハッシュ値）を記録するファイル** である。

JavaScript でいう pnpm-lock.yaml や package-lock.json のハッシュ検証部分に相当する。

## go.sum の役割

go.sum は go.mod に記載された各依存パッケージの **正確なハッシュ値を記録** する。

```
github.com/example/package v1.2.3 h1:abc123...
github.com/example/package v1.2.3/go.mod h1:def456...
```

`go mod download` や `go build` 時に、ダウンロードしたパッケージのハッシュが go.sum の記録と一致するか検証される。一致しなければ **ビルドを拒否** する（改ざん検知、サプライチェーン攻撃対策）。

## なぜ go.sum が必要か

**go.mod だけではバージョン番号しかわからない** という問題がある。

同じバージョン番号（例: `v1.2.3`）でも、パッケージの中身がすり替えられている可能性がある。攻撃者がパッケージレジストリを改ざんした場合や、中間者攻撃でパッケージが書き換えられた場合、バージョン番号だけでは検知できない。

**go.sum のハッシュ値** によって、パッケージの内容が意図したものと完全に一致するかを検証できる。ハッシュ値が 1 ビットでも変われば検証に失敗するため、改ざんを確実に検知できる。

## go.sum の運用

### バージョン管理に含める

go.sum は **必ず git にコミットする**。

- チーム全員が同じハッシュ値を共有することで、全員が同じパッケージを使っていることを保証する
- CI/CD でも go.sum を使って検証することで、本番環境での改ざんを防ぐ

### 手動で編集しない

go.sum は **自動生成されるファイル** であり、手動で編集してはいけない。

`go mod tidy` を実行すると、不要なエントリが削除され、必要なエントリが追加される。

```bash
go mod tidy
```

### go mod tidy で整理される

- 依存パッケージを追加した場合 → 自動的に go.sum にハッシュが記録される
- 依存パッケージを削除した場合 → `go mod tidy` で不要なエントリが削除される

## go mod tidy の詳細

`go mod tidy` は **go.mod と go.sum を実際のコードの使用状況に合わせて整理するコマンド** である。

### やること

1. **不要な依存を削除** — コードで使っていないのに go.mod に残っている依存を削除する
2. **足りない依存を追加** — コードで使っているのに go.mod に書かれていない依存を追加する
3. **go.sum を同期** — 上記に合わせて go.sum のエントリも追加・削除する

### いつ使うか

- 依存パッケージを追加・削除した後
- go.mod を手動で編集した後
- CI でビルドする前（整合性チェックとして）

## 外部パッケージの導入方法

### 方法1: コードに import を書いてから go mod tidy（一般的）

```go
import "github.com/labstack/echo/v4"
```

```bash
go mod tidy
```

### 方法2: go get で明示的に追加

```bash
go get github.com/labstack/echo/v4@latest    # 最新版
go get github.com/labstack/echo/v4@v4.11.0   # バージョン指定
```

### 使い分け

- **import + go mod tidy**: 通常の開発フロー。コードを書きながら自然に追加
- **go get**: バージョンを明示的に指定したい時、まだ import を書く前に先に入れたい時

### 削除する場合

コードから import を消してから `go mod tidy` で不要な依存が削除される。

### JavaScript との比較

| JavaScript（pnpm） | Go | 説明 |
|---|---|---|
| `pnpm add express` | `go get github.com/.../express` | パッケージを追加 |
| `pnpm install` | `go mod download` | lockfile から依存を復元 |
| `package.json` | `go.mod` | 依存パッケージの宣言 |
| `pnpm-lock.yaml` | `go.sum` | ハッシュ値の記録 |

**Go の特徴**: `go mod tidy` がインストールと整理を兼ねている。Go はビルド時に自動でダウンロードするので、明示的な「インストール」ステップが不要な場面が多い。

## go.mod と go.sum の関係

| ファイル | 役割 | 何を記録するか |
|---|---|---|
| **go.mod** | 依存パッケージとバージョンの宣言 | **何を使うか**（パッケージ名とバージョン） |
| **go.sum** | 依存パッケージの整合性検証 | **本物かどうか**（ハッシュ値） |

go.mod はパッケージのバージョンを指定するが、go.sum はそのバージョンの内容が正しいかを検証する。両方が揃って初めて、安全かつ再現可能なビルドが実現できる。

### 例

go.mod:

```go
module github.com/harusame0616/try-go-remix-dynamo/apps/api

go 1.23.4

require github.com/example/package v1.2.3
```

go.sum:

```
github.com/example/package v1.2.3 h1:abc123...
github.com/example/package v1.2.3/go.mod h1:def456...
```

`go build` 時：

1. go.mod を読んで `github.com/example/package v1.2.3` が必要だとわかる
2. パッケージをダウンロードしてハッシュ値を計算する
3. go.sum に記録されたハッシュ値と比較する
4. 一致すればビルド続行、不一致ならエラーで停止
