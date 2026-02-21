# Go の命名規約

## スコープが狭いほど短い名前

Go はスコープが狭い変数ほど短い名前を推奨する文化がある。

| スコープ | 例 | 説明 |
|---|---|---|
| ハンドラー引数 | `w`, `r` | 全員が知っている慣習 |
| ループ変数 | `i`, `k`, `v` | `for i, v := range items` |
| レシーバー | `s`, `h` | `func (s *Server) Start()` |
| パッケージレベル | 長い名前 | `DefaultTimeout`, `ErrNotFound` |

**「スコープが狭い → 短く」「スコープが広い → 説明的に」** がルール。

## 慣習的な省略

以下は「省略」ではなく「Go を書く人なら誰でも知っている語彙」。

```go
// HTTP ハンドラーの引数
func handler(w http.ResponseWriter, r *http.Request) {}

// context は ctx
func handler(ctx context.Context) {}

// error は err
if err := doSomething(); err != nil {}

// testing.T は t、testing.B は b
func TestFoo(t *testing.T) {}
func BenchmarkFoo(b *testing.B) {}
```

これらを `response` や `request` と書くと「Go に慣れていない人が書いたコード」という印象を与える。

## 公開 / 非公開と命名

- **大文字始まり** → 公開（exported）: `NewMux`, `StatusOK`, `ErrNotFound`
- **小文字始まり** → 非公開（unexported）: `newMux`, `port`, `err`

`New〜` で始まる関数はコンストラクタ的な役割を持つ慣習:

```go
func NewMux() *http.ServeMux { ... }     // ServeMux を生成して返す
func NewServer(port int) *Server { ... } // Server を生成して返す
```
