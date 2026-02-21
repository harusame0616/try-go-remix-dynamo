# Go のコード実行方法

## コード実行方法

### go run — コンパイル＋即実行

- `go run .` でカレントディレクトリの main パッケージを実行
- `go run main.go` でファイル指定でも可
- コンパイルして一時ファイルとして実行し、終了後に削除される
- 開発中はこれが基本
- **常に `package main` + `main()` 関数の両方が必要**。どちらが欠けてもエラー

### go build — バイナリを生成

- `go build -o api .` でバイナリを生成
- 本番デプロイ用。生成されたバイナリは単体で動く（ランタイム不要）

### JS との比較表

| JS | Go | 用途 |
|----|----|------|
| `node app.js` / `tsx app.ts` | `go run .` | 開発中の実行 |
| `tsc && node dist/app.js` | `go build -o app . && ./app` | ビルド＋実行 |
| `nodemon` / `vite dev` | `air` など外部ツール | ホットリロード |

## main パッケージと実行の制約

### go run の実行条件

- `package main` であること（必須）
- `main()` 関数があること（必須）
- ファイル指定（`go run main.go`）でもパッケージ指定（`go run .`）でもこの条件は同じ

### 特定の関数だけ実行する方法

Go には直接的な方法がない（REPL もない）。

#### 方法1: main() を一時的に書き換えて試す

```go
func main() {
    result := myFunction()
    fmt.Println(result)
}
```

#### 方法2: テストとして書く（Go らしいやり方）

```go
func TestMyFunction(t *testing.T) {
    result := myFunction()
    fmt.Println(result)
}
```

```bash
go test -run TestMyFunction ./...
```

## グローバルスコープの制約

### 関数の外に書けるもの（宣言のみ）

```go
var x = 10                    // OK: 変数宣言
const y = "hello"             // OK: 定数宣言
type MyStruct struct{}        // OK: 型宣言
```

### 関数の外に書けないもの

```go
fmt.Println("hello")          // NG: 文（statement）は書けない
```

### init() 関数

- main() より前に自動実行される特殊関数
- パッケージの初期化処理に使う
- 濫用すると依存関係が見えにくくなるので注意

```go
func init() {
    // パッケージのインポート時に自動実行される
    fmt.Println("初期化処理")
}

func main() {
    fmt.Println("main 関数")
}
```

実行結果:
```
初期化処理
main 関数
```

### JS との比較

| やりたいこと | Go | JS |
|-------------|----|----|
| エントリポイント | main() のみ | どのファイルでも node file.js |
| 特定の関数を実行 | テストで書くか main() を書き換える | node -e "require('./foo').bar()" |
| グローバルにコード実行 | 不可（宣言のみ） | 可能 |
| REPL | なし（標準） | node |

Go は「明示的であること」を重視する言語なので、実行の起点が常に main() に限定されている。自由度は低いが、コードの追跡がしやすい。
