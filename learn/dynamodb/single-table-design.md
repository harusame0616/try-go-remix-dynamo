# DynamoDB Single Table Design

## Single Table Design とは

DynamoDB では **1つのテーブルに複数のエンティティを混在させる** のが定番パターン。
RDB のように「sensors テーブル」「sensor_data テーブル」と分けるのではなく、PK と SK の命名規則でエンティティを区別する。

## 設計手順

1. **アクセスパターンの洗い出し**（最重要）
2. エンティティの整理
3. PK / SK の命名規則決定
4. GSI の設計

## IoT センサーシステムの設計例

### アクセスパターン

| # | アクセスパターン | 例 |
|---|---|---|
| 1 | センサーのメタ情報を取得 | sensor#001 の名前・設置場所 |
| 2 | センサーデータを記録 | sensor#001 の温度 35.2℃ を保存 |
| 3 | 特定センサーの時系列データ取得 | sensor#001 の直近1時間分 |
| 4 | 特定センサーの最新データ取得 | sensor#001 の最新1件 |
| 5 | 場所別のセンサー一覧取得 | 東京オフィスのセンサー全部 |

### テーブル設計

テーブル名: SensorTable

| PK | SK | name | location | temperature | humidity | GSI1PK | GSI1SK |
|---|---|---|---|---|---|---|---|
| SENSOR#001 | METADATA | 温度計A | tokyo | | | tokyo | 001 |
| SENSOR#001 | DATA#2026-02-22T10:00:00Z | | | 35.2 | 60.1 | | |
| SENSOR#001 | DATA#2026-02-22T10:01:00Z | | | 35.5 | 59.8 | | |
| SENSOR#002 | METADATA | 温度計B | osaka | | | osaka | 002 |
| SENSOR#002 | DATA#2026-02-22T10:00:00Z | | | 28.1 | 55.0 | | |

### GSI 設計

GSI1: 場所別センサー一覧用
- PK: GSI1PK (= location)
- SK: GSI1SK (= sensor_id)

### 各アクセスパターンとキー条件の対応

| # | 操作 | キー条件 |
|---|---|---|
| 1 | GetItem | PK=SENSOR#001, SK=METADATA |
| 2 | PutItem | PK=SENSOR#001, SK=DATA#\<timestamp\> |
| 3 | Query | PK=SENSOR#001, SK begins_with "DATA#" + BETWEEN |
| 4 | Query | PK=SENSOR#001, SK begins_with "DATA#", ScanIndexForward=false, Limit=1 |
| 5 | Query (GSI1) | GSI1PK=tokyo |

### SK の命名規則のポイント

- `METADATA` と `DATA#` のようにプレフィックスでエンティティを区別
- タイムスタンプは ISO 8601（RFC3339）形式にすると辞書順 = 時系列順になる
- `begins_with` でフィルタしやすいようにプレフィックスを設計する
