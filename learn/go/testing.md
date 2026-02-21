# Go テスト（testing パッケージ）

## testing.T の頻出メソッド

### テスト失敗を報告する

| メソッド | テスト失敗マーク | テスト中断 | 用途 |
|---|---|---|---|
| `t.Error(args...)` | する | しない | 失敗を記録して続行 |
| `t.Errorf(format, args...)` | する | しない | フォーマット付きで失敗を記録して続行 |
| `t.Fatal(args...)` | する | する | 致命的な失敗。即座に中断 |
| `t.Fatalf(format, args...)` | する | する | フォーマット付きで即座に中断 |

**使い分けの基準：**

- **`Errorf`** — 検証の失敗。後続の検証も実行したい場合
- **`Fatalf`** — 前提条件の失敗。続行しても意味がない場合（`err != nil` のチェック等）

```go
// 前提条件 → Fatal（これが失敗したら後の検証は無意味）
resp, err := http.Get(url)
if err != nil {
    t.Fatalf("unexpected error: %v", err)
}

// 検証 → Error（他の検証も見たい）
if resp.StatusCode != 200 {
    t.Errorf("expected 200, got %d", resp.StatusCode)
}
```

### ログ出力

| メソッド | 説明 |
|---|---|
| `t.Log(args...)` | テスト失敗時 or `-v` フラグ時にログ出力 |
| `t.Logf(format, args...)` | フォーマット付きログ出力 |

```go
t.Logf("response body: %s", body) // デバッグ時に便利
```

通常時は表示されず、`go test -v` かテスト失敗時のみ出力される。

### サブテスト（t.Run）

`t.Run` でテストをグループ化し、各検証の意図を明確にできる。

```go
func TestHealthEndpoint(t *testing.T) {
    srv := httptest.NewServer(NewMux())
    defer srv.Close()

    resp, err := http.Get(srv.URL + "/health")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    defer resp.Body.Close()

    t.Run("ステータスコード200を返す", func(t *testing.T) {
        if resp.StatusCode != http.StatusOK {
            t.Errorf("expected status 200, got %d", resp.StatusCode)
        }
    })

    t.Run("JSON形式でstatus:okを返す", func(t *testing.T) {
        b, err := io.ReadAll(resp.Body)
        if err != nil {
            t.Fatalf("failed to read body: %v", err)
        }
        expected := `{"status":"ok"}`
        if string(b) != expected {
            t.Errorf("expected body %q, got %q", expected, string(b))
        }
    })
}
```

サブテスト単体で実行できる：

```bash
go test -run TestHealthEndpoint/ステータスコード200を返す
```

### スキップ・ヘルパー

| メソッド | 説明 |
|---|---|
| `t.Skip(args...)` | テストをスキップする |
| `t.Skipf(format, args...)` | 条件付きスキップ |
| `t.Helper()` | ヘルパー関数としてマークする |

```go
// 環境依存のテストをスキップ
func TestWithDB(t *testing.T) {
    if os.Getenv("DB_URL") == "" {
        t.Skip("DB_URL not set")
    }
}
```

### クリーンアップ・並列実行

| メソッド | 説明 |
|---|---|
| `t.Cleanup(func())` | テスト終了時に実行される後片付け関数を登録 |
| `t.Parallel()` | 他の Parallel テストと並列実行する |
| `t.TempDir()` | テスト終了時に自動削除される一時ディレクトリを返す |

```go
func TestParallel(t *testing.T) {
    t.Parallel()
    dir := t.TempDir() // テスト後に自動で消える
    t.Cleanup(func() {
        db.Close()
    })
}
```

## カスタムアサーション

Go の `testing` パッケージには `assertEqual` などの組み込みアサーション関数がない。毎回 `if` + `t.Errorf` を書く必要があるため、繰り返しを減らすために自作するアサーション関数をカスタムアサーションと呼ぶ。

```go
func assertEqual(t *testing.T, got, want string) {
    t.Helper() // 失敗時に呼び出し元の行番号を報告するために必要
    if got != want {
        t.Errorf("got %q, want %q", got, want)
    }
}

func TestSomething(t *testing.T) {
    assertEqual(t, body, `{"status":"ok"}`)
}
```

### t.Helper() の効果

```
// t.Helper() なし → assertEqual 関数内の行番号が出る（デバッグしづらい）
main_test.go:8: got "ng", want "ok"

// t.Helper() あり → 呼び出し元の行番号が出る（デバッグしやすい）
main_test.go:15: got "ng", want "ok"
```

### Go コミュニティの立場

Go は意図的にアサーション関数を標準に入れていない。「`if` + `t.Errorf` の方が失敗時のメッセージを具体的に書けるから」という思想。

- **小規模プロジェクト** → `if` + `t.Errorf` を素直に書くのが主流
- **大規模プロジェクト** → [testify](https://github.com/stretchr/testify) 等の外部ライブラリを使うことが多い

```go
// testify を使った場合
import "github.com/stretchr/testify/assert"

func TestSomething(t *testing.T) {
    assert.Equal(t, 200, resp.StatusCode)
    assert.JSONEq(t, `{"status":"ok"}`, body)
}
```

## テストファイルの規約

- テストファイルは **`_test.go` で終わる名前にする**（言語仕様）
- `go test` はこのファイルだけをテスト対象として認識する
- `_test.go` のファイルは `go build` ではビルドされない（テストコード専用）
- テスト関数は **`Test` で始まる名前** でなければならない（言語仕様）

### パッケージ名の選択

| パッケージ名 | テスト方式 | 用途 |
|---|---|---|
| `package main` | ホワイトボックス | 非公開関数にもアクセスできる |
| `package main_test` | ブラックボックス | 公開 API だけをテスト |

## テスト実行コマンド

```bash
# 基本: カレントディレクトリのテストを実行
go test

# 再帰的に全パッケージのテストを実行
go test ./...

# 詳細な出力（どのテストが通ったか表示）
go test -v ./...

# カバレッジ付き
go test -cover ./...

# 特定のテスト関数だけ実行（正規表現）
go test -run TestHealthEndpoint

# 組み合わせ
go test -v -cover -run TestHealth ./...
```

### よく使うフラグ

| フラグ | 意味 |
|---|---|
| `-v` | verbose。各テスト関数名と PASS/FAIL を表示 |
| `-cover` | カバレッジ率を表示 |
| `-run <正規表現>` | マッチするテストだけ実行 |
| `-count=1` | キャッシュを無視して必ず再実行 |
| `./...` | カレント以下の全パッケージが対象 |

### テスト結果のキャッシュ

Go はテスト結果をキャッシュする。コードが変わっていなければ再実行せず `(cached)` と表示される。強制的に再実行したい場合は `-count=1` をつける。

```bash
go test -v -count=1 ./...
```

## httptest パッケージ

`net/http/httptest` はテスト用の HTTP サーバーを簡単に立てられる標準パッケージ。

```go
import "net/http/httptest"

srv := httptest.NewServer(handler)  // ランダムなポートでサーバー起動
defer srv.Close()

resp, err := http.Get(srv.URL + "/health")  // srv.URL → "http://127.0.0.1:xxxxx"
```

- `httptest.NewServer` はローカルにテスト用 HTTP サーバーを立ち上げる
- ランダムなポートで起動し、`srv.URL` でアドレスを取得できる
- 実際のネットワーク通信が発生するので、E2E テストに近い形で検証できる
