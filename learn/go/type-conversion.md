# 型変換と文字列操作

## Go は暗黙の型変換を一切しない

JavaScript の `":" + 8080` が `":8080"` になるような自動変換は起きず、**コンパイルエラーになる**。

```go
port := 8080
addr := ":" + port  // コンパイルエラー: mismatched types string and int
```

## 数値 → 文字列

### strconv.Itoa（整数 → 文字列）

```go
import "strconv"

port := 8080
s := strconv.Itoa(port)  // "8080"
addr := ":" + s           // ":8080"
```

- "Itoa" は "Integer to ASCII" の略
- 単純な変換のときはこちら

### fmt.Sprintf（フォーマットして文字列を返す）

```go
port := 8080
addr := fmt.Sprintf(":%d", port)  // ":8080"
```

- 複数の値を組み合わせるときはこちら

## 文字列 → 数値

### strconv.Atoi（文字列 → 整数）

```go
s := "8080"
n, err := strconv.Atoi(s)  // "ASCII to Integer"
if err != nil {
    // "abc" のような変換できない文字列の場合ここに来る
}
```

- 戻り値が2つ（`n` と `err`）
- 変換できない文字列の場合は `err` が返る

## []byte ↔ string

```go
b := []byte("hello")   // string → []byte
s := string(b)          // []byte → string
```

- Go では `string` と `[]byte` は別の型なので明示的な変換が必要
- `io.ReadAll` の戻り値は `[]byte` なので、文字列として使うなら `string()` で変換する
