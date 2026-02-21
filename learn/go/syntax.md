# Go 基本構文

## パッケージ宣言

Go のすべてのファイルは必ずどこかのパッケージに属する。

```go
package main
```

- `package main` は特別で「実行可能プログラム」を意味する
- `go run .` や `go build` で実行できるのは `package main` だけ
- ライブラリなら `package mylib` のように別名をつける

## import

他のパッケージを使うための宣言。

```go
import (
    "fmt"
    "net/http"
    "os"
)
```

- 複数パッケージは `()` でまとめる
- **使っていないパッケージを import するとコンパイルエラーになる**

## 関数定義

```go
func 関数名(引数 型) 戻り値の型 {
    return 値
}
```

```go
func NewMux() *http.ServeMux {
    mux := http.NewServeMux()
    return mux
}
```

- `func` キーワードで定義
- 引数は `名前 型` の順（C や Java と逆）
- 戻り値の型は引数リストの後ろに書く

### 多値返却

Go の関数は複数の値を返せる。

```go
func Atoi(s string) (int, error) { ... }

n, err := strconv.Atoi("8080")
```

- エラーを2番目の戻り値で返すのが Go の定番パターン

## 変数宣言

### 短縮変数宣言 `:=`

```go
port := "8080"          // 型を自動推論して宣言＋代入
mux := http.NewServeMux()
```

- **関数内でのみ使える**
- 最も頻繁に使う構文

### var 宣言

```go
var port string = "8080"  // 型を明示
var port = "8080"         // 型推論あり
var port string           // ゼロ値（""）で初期化
```

- 関数外（パッケージレベル）で変数を宣言する場合は `var` を使う

### 代入 `=`

```go
port = "3000"  // 既に宣言済みの変数に代入
```

`:=` は宣言＋代入、`=` は代入のみ。

## ポインタ `*`

```go
var x *http.ServeMux  // ポインタ型：データの場所（アドレス）を持つ
var y http.ServeMux   // 値型：データのコピーを持つ
```

- `*Type` → ポインタ型
- 大きなデータのコピーを避けたいとき、同じデータを共有したいときに使う
- `*http.Request` → Request は大きいのでポインタで渡すのが標準

初学者の理解: 「`*` がついていたらポインタで、大きなデータや共有が必要なときに使う」で十分。

## 無名関数（クロージャ）

```go
mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, `{"status":"ok"}`)
})
```

- `func(引数) { 処理 }` で名前のない関数を定義できる
- JavaScript の `(req, res) => { ... }` と同じ発想

## raw 文字列リテラル

```go
s := `{"status":"ok"}`   // バッククォート：エスケープ不要
s := "{\"status\":\"ok\"}" // ダブルクォート：エスケープ必要
```

- バッククォート `` ` `` で囲むとエスケープが不要
- 改行もそのまま含められる

## if 文

```go
if port == "" {
    port = "8080"
}
```

- 条件に括弧 `()` は不要（書くとエラー）
- `{` は `if` と同じ行に書く（Go のフォーマット規約）

### if with 短縮変数宣言

```go
if err := http.ListenAndServe(addr, mux); err != nil {
    fmt.Println(err)
}
// err はこのブロック外からアクセスできない
```

- セミコロン `;` の前で変数を宣言し、後で条件判定する
- 宣言した変数は `if` ブロック内でのみ有効

## defer

```go
srv := httptest.NewServer(NewMux())
defer srv.Close()  // 関数を抜けるときに実行される
```

- `defer` は「この関数が終了するときに実行する」という予約
- どこで `return` しても、パニックしても、必ず実行される
- JavaScript の `try { ... } finally { cleanup() }` に近い
- リソースの閉じ忘れを防ぐ Go の定番パターン

## nil

Go の null に相当する。

```go
if err != nil {
    // エラーがある場合の処理
}
```

- 「エラーがなければ `nil`、あれば `error` 型の値」が Go の規約

## 公開 / 非公開（exported / unexported）

```go
func NewMux() { }   // 大文字始まり → 他パッケージから参照可能（公開）
func newMux() { }   // 小文字始まり → 同一パッケージ内のみ（非公開）
```

- Go にはアクセス修飾子（`public` / `private`）がない
- **名前の先頭が大文字か小文字かで決まる**
- 関数、変数、型、定数すべてに適用される
