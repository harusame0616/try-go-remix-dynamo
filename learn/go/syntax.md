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

---

## 構造体（struct）

新しい型を定義する構文。TypeScript の `interface` や `type` に近い。

```go
type Sensor struct {
    ID       string
    Name     string
    Location string
}
```

- `type 名前 struct { }` で定義
- フィールドは `名前 型` の順
- 大文字始まり → 公開、小文字始まり → 非公開

```go
type Repository struct {
    client    *dynamodb.Client  // 小文字 → 非公開
    tableName string            // 小文字 → 非公開
}
```

### 構造体リテラル（インスタンス生成）

```go
s := Sensor{
    ID:       "001",
    Name:     "温度計A",
    Location: "tokyo",
}

// ポインタとして生成
s := &Sensor{
    ID: "001",
}
```

TypeScript の `{ id: "001", name: "温度計A" }` に近いが、型が明示的。

## メソッド（レシーバ付き関数）

構造体に紐づく関数。TypeScript のクラスメソッドに相当。

```go
func (r *Repository) CreateSensor(ctx context.Context, s Sensor) error {
    r.client.PutItem(...)  // r が TypeScript の this に相当
}
```

- `(r *Repository)` がレシーバ
- `*Repository` → ポインタレシーバ（元の構造体を変更可能、コピーなし）
- `Repository` → 値レシーバ（コピーが渡される）
- **実務上はほぼ全てポインタレシーバを使う**

```typescript
// TypeScript だとこう
class Repository {
  createSensor(ctx: Context, s: Sensor): Error {
    this.client.putItem(...)
  }
}
```

### Go でクラス的なものを表現する

Go には `class` キーワードがなく、struct + レシーバーメソッドでクラスを表現する。メソッドは struct の `{}` の外に定義する。

```go
type Sensor struct {
    ID   string
    Name string
}

func (s *Sensor) FullName() string {
    return s.ID + ": " + s.Name
}

func (s *Sensor) Rename(name string) {
    s.Name = name
}
```

### コンストラクタ

Go にはコンストラクタもない。`New~` という命名の関数を使う慣習がある。

```go
func NewSensor(id, name string) *Sensor {
    return &Sensor{
        ID:   id,
        Name: name,
    }
}
```

### 継承の代わり：埋め込み（embedding）

Go に継承はない。埋め込みで似たことができるが、これは継承ではなく委譲（delegation）。is-a ではなく has-a の関係。

```go
type TimestampedSensor struct {
    Sensor              // Sensor を埋め込む
    CreatedAt time.Time
}

ts := TimestampedSensor{
    Sensor:    Sensor{ID: "001", Name: "温度計A"},
    CreatedAt: time.Now(),
}

ts.FullName()  // Sensor のメソッドがそのまま使える
ts.ID          // Sensor のフィールドにも直接アクセスできる
```

### TypeScript のクラスとの対応

| TypeScript | Go |
|---|---|
| `class Sensor` | `type Sensor struct` |
| `constructor()` | `func NewSensor() *Sensor` |
| `this.name` | `s.Name`（レシーバー変数） |
| クラス内メソッド | レシーバーメソッド（外に定義） |
| `extends` (継承) | 埋め込み（has-a、is-a ではない） |
| `implements` | インターフェースを暗黙的に満たす |
| `private` / `public` | 小文字始まり / 大文字始まり |

## context.Context

キャンセル・タイムアウトの伝搬を担う Go 独自の仕組み。

```go
func (r *Repository) CreateSensor(ctx context.Context, s Sensor) error {
```

- Go の慣習として必ず第1引数に置く
- HTTP リクエストのキャンセル（クライアント切断）などを関数チェーン全体に伝える
- TypeScript の `AbortController` が近い概念

```go
// コンテキストの流れ
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
    sensors, err := h.service.ListSensors(r.Context())  // リクエストのコンテキスト
}

func (s *Service) ListSensors(ctx context.Context) ([]Sensor, error) {
    return s.repo.FindAll(ctx)  // そのまま渡す
}
```

### context を渡すべき場合 / 不要な場合

全てのメソッドに必須ではない。「context.Context を受け取るなら第1引数にする」が慣習。

```go
// I/O を伴う処理 → 渡す
func (r *Repository) GetSensor(ctx context.Context, id string) (*Sensor, error)

// 純粋な計算やデータ変換 → 不要
func (s *Sensor) FullName() string {
    return s.ID + ": " + s.Name
}
```

| context が必要 | context が不要 |
|---|---|
| DB アクセス | 純粋な計算 |
| HTTP リクエスト | 文字列操作 |
| 外部 API 呼び出し | データ変換 |
| 時間がかかる可能性がある処理 | 即座に完了する処理 |

### context の具体的な使い方

#### 1. HTTP ハンドラから受け取る（最も一般的）

```go
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()  // HTTP サーバーが作った context を受け取る
    sensors, err := h.service.ListSensors(ctx)
    // ...
}
```

ほとんどの場合、自分で context を生成する必要はない。HTTP サーバーが作ったものを受け取って下に渡すだけ。

#### 2. タイムアウトを設定する

```go
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
    defer cancel()  // 必ず呼ぶ（リソースリーク防止）

    sensors, err := h.service.ListSensors(ctx)
    if errors.Is(err, context.DeadlineExceeded) {
        http.Error(w, "timeout", http.StatusGatewayTimeout)
        return
    }
}
```

#### 3. キャンセルの確認

外部ライブラリ（AWS SDK 等）は内部で自動的に context を確認してくれる。自分で確認が必要なのは重い計算ループの場合のみ。

```go
func processLargeData(ctx context.Context, items []SensorData) error {
    for i, item := range items {
        if i%100 == 0 {
            select {
            case <-ctx.Done():
                return ctx.Err()  // context.Canceled or DeadlineExceeded
            default:
            }
        }
        process(item)
    }
    return nil
}
```

### TypeScript との対応

```typescript
// TypeScript（AbortController）
const controller = new AbortController();
setTimeout(() => controller.abort(), 3000);
fetch(url, { signal: controller.signal });
```

`AbortController` + `signal` が、Go の `context.WithTimeout` + `ctx` に対応する。

## map（マップ型）

キーと値のペアを格納するデータ構造。TypeScript の `Record<K, V>` に相当。

```go
m := map[string]string{
    "key1": "value1",
    "key2": "value2",
}

v := m["key1"]        // "value1"
v, ok := m["key3"]    // ok=false（存在しないキー）
```

## インターフェースと型アサーション

インターフェース型から具体的な型を取り出す構文。TypeScript の `as Type` に相当。

```go
// 型アサーション
v := out.Item["name"].(*types.AttributeValueMemberS).Value

// 安全な型アサーション（ok パターン）
v, ok := out.Item["name"].(*types.AttributeValueMemberS)
if !ok {
    // 型が違った場合の処理
}
```

## make / append / len

スライス（可変長配列）を扱う組み込み関数。

```go
// make: スライスの初期化
results := make([]SensorData, 0, len(out.Items))
//              型            長さ  容量

// append: 要素の追加（TypeScript の push に相当、新しいスライスを返す）
results = append(results, newItem)

// len: 長さの取得（TypeScript の .length に相当）
n := len(out.Items)
```

## for range（ループ）

スライスやマップをイテレートする構文。

```go
for i, item := range items {    // i=インデックス, item=値
for _, item := range items {    // インデックスは不要（_ で捨てる）
for i := range items {          // 値は不要（インデックスだけ）
```

`_` はブランク識別子で「この値は使わない」の意味。

```typescript
// TypeScript だとこう
for (const item of items) { ... }
items.forEach((item, i) => { ... })
```

## &（アドレス演算子）

ポインタ `*` の対になる演算子。変数のアドレス（ポインタ）を取得する。

```go
x := "hello"
p := &x         // p は *string 型（x のポインタ）
fmt.Println(*p) // "hello"（ポインタから値を取り出す）
```

- `&変数` → ポインタを取得
- `*ポインタ` → ポインタから値を取り出す（デリファレンス）
- `&Struct{...}` → 構造体をポインタとして生成

SDK の関数がポインタを要求するので `&dynamodb.PutItemInput{...}` のように `&` をつける。

## 値渡し vs ポインタ渡しの使い分け

型によってポインタで渡すべきかどうかが異なる。

### 基本ルール

| 型 | ポインタで渡す？ | 理由 |
|---|---|---|
| **string, int, bool** | No — 値渡し | 軽い（string は内部的にポインタ+長さの16Bだけ） |
| **大きい struct** | Yes — ポインタ | コピーコストを避ける |
| **slice, map** | No — 値渡し | 内部的に既にポインタを持っている（参照型） |

### string がポインタ不要な理由

Go の string は内部的にヘッダ（ポインタ + 長さ = 16B）しかコピーされない。文字列の実体データはコピーされない。

```go
func greet(name string) { }   // Good: 値渡しで十分
func greet(name *string) { }  // Bad: 不要
```

※ AWS SDK で `aws.String("hello")` が出てくるのは「未設定（nil）」と「空文字」を区別するための SDK の都合。Go の一般的な慣習ではない。

### struct のポインタ判断基準

| 判断基準 | 値渡し | ポインタ渡し |
|---|---|---|
| フィールドを変更する？ | — | ポインタ必須 |
| フィールド数が多い（目安5個以上）？ | — | ポインタ推奨 |
| nil を表現したい？ | — | ポインタ必須 |
| 小さくて読み取りだけ？ | 値渡しで十分 | — |

迷ったらポインタにしておけば間違いない。

### slice / map がポインタ不要な理由

slice と map は内部的に既にポインタを持っている（参照型）。値渡ししても実体データはコピーされない。

```go
func process(items []SensorData) { }   // Good: そのまま渡す
func process(items *[]SensorData) { }  // Bad: ほぼ間違い
```

### まとめ

```go
// Good
func CreateSensor(ctx context.Context, name string, s *Sensor) error

// Bad
func CreateSensor(ctx context.Context, name *string, s Sensor) error
```

## ポインタ引数の変更防止

Go には C/C++ の `const` のような「ポインタで渡すが変更を禁止する」キーワードがない。

### 対処法

| アプローチ | 方法 |
|---|---|
| 変更しないなら値渡し | コピーコストが許容できるならポインタにしない |
| インターフェースで読み取り専用メソッドだけ公開 | フィールドに直接触れないようにする |
| 命名で意図を伝える | 関数名やコメントで読み取り専用を明示 |

```go
// 値渡し（最もシンプル）
func printSensor(s Sensor) {  // コピーなので元に影響しない
    s.Name = "changed"  // 元の Sensor は変わらない
}

// インターフェースでフィールドを隠す
type SensorReader interface {
    GetID() string
    GetName() string
}

func process(s SensorReader) {
    name := s.GetName()  // 読み取りだけ可能
}
```

Go のコミュニティでは「ポインタで渡して変更しない」は暗黙の契約として成り立っており、変更が必要ならメソッド名で明示する（`Update~`, `Set~` など）のが一般的。

## 型の比較

### 文字列（string）の比較

`==` でそのまま中身を比較できる。大小比較（辞書順）も可能。

```go
a := "hello"
b := "hello"
c := "world"

a == b   // true（中身が同じ）
a < c    // true（辞書順で比較）
```

辞書順比較のおかげで ISO 8601 形式のタイムスタンプが正しくソートされる。

```go
"2026-01" < "2026-02"  // true
```

JavaScript と違って `===` と `==` の区別はない。`==` だけで常に型安全な比較になる。

### slice の比較

slice は `==` で比較できない。nil との比較だけ例外的に許可。

```go
a := []int{1, 2, 3}
b := []int{1, 2, 3}

a == b     // コンパイルエラー
a == nil   // OK

// 中身の比較
slices.Equal(a, b)        // Go 1.21+（推奨）
reflect.DeepEqual(a, b)   // 古い方法（遅い）
```

### map の比較

map も `==` で比較できない。

```go
m1 := map[string]int{"a": 1}
m2 := map[string]int{"a": 1}

m1 == m2  // コンパイルエラー

maps.Equal(m1, m2)        // Go 1.21+（推奨）
reflect.DeepEqual(m1, m2) // 古い方法
```

### struct の比較

フィールドが全て比較可能なら `==` が使える。ただし slice や map フィールドを含むと `==` は使えない。

```go
s1 := Sensor{ID: "001", Name: "温度計A"}
s2 := Sensor{ID: "001", Name: "温度計A"}
s1 == s2  // OK: true

type Data struct {
    Values []int  // slice を含む
}
d1 := Data{Values: []int{1, 2}}
d2 := Data{Values: []int{1, 2}}
d1 == d2  // コンパイルエラー
```

### ポインタ経由の構造体比較

ポインタ同士の `==` はアドレス比較。中身を比較したい場合は `*` でデリファレンスする。

```go
s1 := &Sensor{ID: "001"}
s2 := &Sensor{ID: "001"}
s3 := s1

s1 == s2   // false（別のアドレス）
s1 == s3   // true（同じアドレス）
*s1 == *s2 // true（中身を比較）
```

| 式 | 比較対象 | 意味 |
|---|---|---|
| `s1 == s2` | アドレス | 同じオブジェクトか？ |
| `*s1 == *s2` | 中身 | フィールドの値が同じか？ |

### ポインタの比較

`==` はアドレス比較（同じオブジェクトか）。中身を比較するにはデリファレンスする。

```go
s1 := &Sensor{ID: "001"}
s2 := &Sensor{ID: "001"}
s3 := s1

s1 == s2   // false（別のアドレス）
s1 == s3   // true（同じアドレス）
*s1 == *s2 // true（中身を比較）
```

### まとめ

| 型 | == で比較 | 中身の比較方法 |
|---|---|---|
| string, int, bool | 可能（中身を比較） | そのまま == |
| string の大小比較 | <, >, <=, >= が使える | 辞書順（lexicographic order） |
| struct（比較可能フィールドのみ） | 可能（中身を比較） | そのまま == |
| struct（slice/map フィールドあり） | 不可 | reflect.DeepEqual |
| slice | 不可（nil のみ可） | slices.Equal（Go 1.21+） |
| map | 不可（nil のみ可） | maps.Equal（Go 1.21+） |
| ポインタ | アドレス比較 | *p1 == *p2 でデリファレンス |

## スライスの詳細

### 配列 vs スライス

```go
// 配列 — 固定長（ほぼ使わない）
var a [3]int = [3]int{1, 2, 3}

// スライス — 可変長（こちらを使う）
var s []int = []int{1, 2, 3}
```

Go では**ほぼ常にスライスを使う**。配列は内部的に存在するが、直接使うことは稀。

#### 配列の特徴

- サイズが型の一部（`[3]int` と `[5]int` は別の型）
- 代入や引数渡しで**全要素がコピー**される

```go
a := [3]int{1, 2, 3}
b := a        // 全要素コピー
b[0] = 999
// a[0] は 1 のまま（元は変わらない）
```

#### スライスの特徴

- 可変長。内部構造は「配列へのポインタ + 長さ + 容量」の24B
- 代入や引数渡しではヘッダだけコピー。要素は共有される

```go
s := []int{1, 2, 3}
t := s        // ヘッダだけコピー（実体は共有）
t[0] = 999
// s[0] も 999 になる（元も変わる！）
```

#### 比較表

| | 配列 `[3]int` | スライス `[]int` |
|---|---|---|
| サイズ | 固定（型の一部） | 可変 |
| 代入時 | 全要素コピー | 参照を共有 |
| 関数に渡す | 全要素コピー | ヘッダだけコピー |
| append | 不可 | 可能 |
| == 比較 | 可能 | 不可 |
| 実務で使うか | ほぼ使わない | こちらを使う |

#### 範囲外アクセス

配列もスライスも同じで panic する。

```go
arr := [3]int{10, 20, 30}
_ = arr[5]  // panic: index out of range

s := []int{10, 20, 30}
_ = s[5]    // panic: index out of range
```

安全にアクセスするには自分で長さを確認する。

```go
if i < len(s) {
    v := s[i]
}
```

配列を使う場面は暗号のハッシュ値（`[32]byte`）など固定サイズが保証されてほしい場合のみ。99%のケースでスライスを使う。
