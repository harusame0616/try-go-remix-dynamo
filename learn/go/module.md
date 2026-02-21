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
