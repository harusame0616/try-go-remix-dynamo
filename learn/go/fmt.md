# fmt パッケージ

## Print 系の使い分け

出力先が違うだけで、フォーマットの書き方は全部同じ。

| 関数 | 出力先 | 用途 |
|---|---|---|
| `fmt.Printf(format, ...)` | 標準出力 | 画面に表示する |
| `fmt.Sprintf(format, ...)` | 文字列として返す | 変数に格納する |
| `fmt.Fprintf(w, format, ...)` | `io.Writer` に書き込む | ファイル、HTTP レスポンス等 |

```go
// 画面に出す
fmt.Printf("port: %d\n", 8080)

// 文字列として返す（画面に出さない）
s := fmt.Sprintf("port: %d", 8080)

// io.Writer に書き込む
fmt.Fprintf(w, "port: %d", 8080)   // w はファイルや HTTP レスポンス等
```

`f` なし版（`Print`, `Sprint`, `Fprint`）はフォーマット指定子を使わず、引数をそのまま出力する。

```go
fmt.Fprint(w, `{"status":"ok"}`)  // フォーマット不要ならこちら
```

## フォーマット指定子

| 指定子 | 意味 | 例 |
|---|---|---|
| `%s` | 文字列 | `"hello"` |
| `%d` | 整数 | `42` |
| `%f` | 浮動小数点 | `3.140000` |
| `%v` | デフォルト書式（なんでもOK） | どの型でも使える |
| `%q` | ダブルクォート付き文字列 | `"hello"` |
| `%+v` | フィールド名付き構造体表示 | `{Name:Go Version:1.22}` |
| `\n` | 改行 | — |

**迷ったら `%v` を使えばとりあえず動く。** 型に合った出力をしてくれる万能フォーマット。

### テンプレートリテラルはない

Go には JavaScript の `` `Hello ${name}` `` に相当する構文がない。変数を埋め込むには `fmt.Sprintf` を使う。

```go
name := "Go"
version := 1.22

// Go
s := fmt.Sprintf("Hello %s, version %v", name, version)

// JavaScript（比較用）
// const s = `Hello ${name}, version ${version}`
```
