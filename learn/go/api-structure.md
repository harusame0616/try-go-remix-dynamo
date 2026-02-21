# Go Web API のコード分離パターンとルーティング

## API コードの分離パターン

Go の Web API では、ハンドラ定義・ビジネスロジック・エンドポイント定義を分離することで、テスタビリティと保守性を高める。一般的なパターンは以下の 3 つ。

### 1. レイヤー型（最も一般的）

handler / service / repository / model をパッケージで分離するパターン。

```
apps/api/
├── main.go
├── handler/
│   └── sensor.go
├── service/
│   └── sensor.go
├── repository/
│   └── sensor.go
└── model/
    └── sensor.go
```

各レイヤーの責務：

| レイヤー | 責務 |
|---|---|
| **handler** | HTTP の関心事（リクエストのパース、レスポンスの組み立て、ステータスコード） |
| **service** | ビジネスルール（HTTP に依存しない純粋なロジック） |
| **repository** | DB 操作の抽象化 |
| **model** | データ構造の定義 |

ルーティング定義は main.go に残すか router/ パッケージに切り出す。

```go
// handler/sensor.go
package handler

type SensorHandler struct {
    service *service.SensorService
}

func (h *SensorHandler) List(w http.ResponseWriter, r *http.Request) {
    sensors, err := h.service.ListSensors(r.Context())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(sensors)
}
```

```go
// service/sensor.go
package service

type SensorService struct {
    repo *repository.SensorRepository
}

func (s *SensorService) ListSensors(ctx context.Context) ([]model.Sensor, error) {
    return s.repo.FindAll(ctx)
}
```

### 2. ドメイン型（機能単位）

機能・ドメインごとにパッケージを分けるパターン。

```
apps/api/
├── main.go
├── sensor/
│   ├── handler.go
│   ├── service.go
│   ├── repository.go
│   └── model.go
└── device/
    ├── handler.go
    ├── service.go
    ├── repository.go
    └── model.go
```

- 関連コードが 1 箇所にまとまりナビゲーションしやすい
- Go では **循環参照が禁止** なのでドメイン間の依存に注意が必要
- ドメイン間で共有する型は別パッケージ（例: `shared/`）に切り出す

```go
// sensor/handler.go
package sensor

type Handler struct {
    svc *Service
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
    sensors, err := h.svc.ListAll(r.Context())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(sensors)
}
```

### 3. ハイブリッド型（小〜中規模向け）

handler / service / model の 3 層構成。repository は service に内包し、規模に応じて後から分離する。

```
apps/api/
├── main.go
├── handler/
│   ├── health.go
│   └── sensor.go
├── service/
│   └── sensor.go
└── model/
    └── sensor.go
```

- 小規模なうちは repository 層を省略して service が直接 DB を操作する
- プロジェクトが成長したら repository を切り出せばよい
- 最初からレイヤーを増やしすぎない実用的なアプローチ

```go
// service/sensor.go
package service

type SensorService struct {
    db *dynamodb.Client
}

// repository 層を持たず、service が直接 DB を操作する
func (s *SensorService) ListSensors(ctx context.Context) ([]model.Sensor, error) {
    out, err := s.db.Scan(ctx, &dynamodb.ScanInput{
        TableName: aws.String("sensors"),
    })
    if err != nil {
        return nil, err
    }
    // ... 結果をマッピング
}
```

### 共通するポイント

- ルーティング定義は main.go か専用の router.go に置く
- ハンドラは HTTP の関心事だけ（パース・バリデーション・レスポンス）
- ビジネスロジックはハンドラから分離し、HTTP に依存しない形にする（テストしやすくなる）
- 依存性の注入はハンドラの構造体フィールドや関数引数で行う

```go
// 依存性の注入の例
type SensorHandler struct {
    service SensorServicer // インターフェースで受け取るとテスト時にモック可能
}

func NewSensorHandler(svc SensorServicer) *SensorHandler {
    return &SensorHandler{service: svc}
}
```

## ルーティングのモジュール化

Express（Node.js）ではルーターを分離してプレフィックス付きでマウントする仕組みが標準で備わっている。Go の標準ライブラリではどうするか。

### Express の例

```js
// routes/sensor.js
const router = express.Router()
router.get('/', listSensors)
router.post('/', createSensor)
module.exports = router

// app.js
app.use('/sensors', sensorRouter)  // プレフィックス付きでマウント
```

### Go でのパターン

#### パターン1: 関数でルートをグループ化（最も一般的・デファクト）

```go
// handler/sensor.go
package handler

func RegisterSensorRoutes(mux *http.ServeMux) {
    mux.HandleFunc("GET /sensors", listSensors)
    mux.HandleFunc("POST /sensors", createSensor)
}
```

```go
// main.go
func NewMux() *http.ServeMux {
    mux := http.NewServeMux()
    handler.RegisterSensorRoutes(mux)
    handler.RegisterDeviceRoutes(mux)
    return mux
}
```

- Express と違いプレフィックスの自動付与はないので各ルートにフルパスを書く
- Go 1.22+ のパターンマッチが強力なので実用上はこれで十分
- シンプルで理解しやすく、最も広く使われている

#### パターン2: ServeMux のネスト（サブルーター的）

```go
sensorMux := http.NewServeMux()
sensorMux.HandleFunc("GET /sensors/", listSensors)
sensorMux.HandleFunc("POST /sensors/", createSensor)

mux := http.NewServeMux()
mux.Handle("/sensors/", sensorMux)
```

- Express ほど綺麗にプレフィックスを切り出せない（ハンドラ側でもフルパスを意識する必要がある）
- あまり使われない

#### パターン3: サードパーティルーター（chi 等）

```go
r := chi.NewRouter()
r.Route("/sensors", func(r chi.Router) {
    r.Get("/", listSensors)
    r.Post("/", createSensor)
})
```

- chi は net/http と互換性があり、Express の `Router()` に最も近い体験を提供する
- ミドルウェアのチェーンも Express ライクに書ける

### Express との比較

| Express (Node.js) | Go 標準ライブラリ | 備考 |
|---|---|---|
| `express.Router()` | `http.NewServeMux()` | Go はサブルーター機能が限定的 |
| `app.use('/prefix', router)` | `mux.Handle("/prefix/", subMux)` | Go はプレフィックス自動付与なし |
| ルートにプレフィックスが自動付与 | 各ルートにフルパスを記述 | Go の方がシンプルだが冗長 |
| ミドルウェアチェーン標準対応 | 自前で実装 or chi 等を使用 | 標準ライブラリには組み込みなし |

### まとめ

Go の標準ライブラリだけで開発する場合、**パターン1（関数でルートをグループ化）** が最も一般的で推奨される。Express のようなプレフィックス自動付与はないが、Go 1.22+ のメソッドベースルーティング（`"GET /sensors"` 形式）と組み合わせれば十分に実用的である。Express に近い体験が欲しい場合は chi の導入を検討するとよい。
