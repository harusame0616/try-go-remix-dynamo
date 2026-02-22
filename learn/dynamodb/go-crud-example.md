# DynamoDB Go CRUD 実装例（AWS SDK for Go v2）

## 使用パッケージ

```go
import (
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)
```

## 構造体定義

```go
type Sensor struct {
    ID       string
    Name     string
    Location string
}

type SensorData struct {
    SensorID    string
    Timestamp   time.Time
    Temperature float64
    Humidity    float64
}

type Repository struct {
    client    *dynamodb.Client
    tableName string
}
```

## Create（PutItem）

### センサーメタ情報の登録

```go
func (r *Repository) CreateSensor(ctx context.Context, s Sensor) error {
    _, err := r.client.PutItem(ctx, &dynamodb.PutItemInput{
        TableName: &r.tableName,
        Item: map[string]types.AttributeValue{
            "PK":       &types.AttributeValueMemberS{Value: "SENSOR#" + s.ID},
            "SK":       &types.AttributeValueMemberS{Value: "METADATA"},
            "name":     &types.AttributeValueMemberS{Value: s.Name},
            "location": &types.AttributeValueMemberS{Value: s.Location},
            "GSI1PK":   &types.AttributeValueMemberS{Value: s.Location},
            "GSI1SK":   &types.AttributeValueMemberS{Value: s.ID},
        },
        ConditionExpression: aws.String("attribute_not_exists(PK)"),
    })
    return err
}
```

### センサーデータの記録

```go
func (r *Repository) PutSensorData(ctx context.Context, d SensorData) error {
    ts := d.Timestamp.Format(time.RFC3339)
    _, err := r.client.PutItem(ctx, &dynamodb.PutItemInput{
        TableName: &r.tableName,
        Item: map[string]types.AttributeValue{
            "PK":          &types.AttributeValueMemberS{Value: "SENSOR#" + d.SensorID},
            "SK":          &types.AttributeValueMemberS{Value: "DATA#" + ts},
            "temperature": &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", d.Temperature)},
            "humidity":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", d.Humidity)},
        },
    })
    return err
}
```

## Read（GetItem / Query）

### メタ情報の取得（GetItem）

```go
func (r *Repository) GetSensor(ctx context.Context, sensorID string) (*Sensor, error) {
    out, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
        TableName: &r.tableName,
        Key: map[string]types.AttributeValue{
            "PK": &types.AttributeValueMemberS{Value: "SENSOR#" + sensorID},
            "SK": &types.AttributeValueMemberS{Value: "METADATA"},
        },
    })
    if err != nil {
        return nil, err
    }
    if out.Item == nil {
        return nil, fmt.Errorf("sensor %s not found", sensorID)
    }
    return &Sensor{
        ID:       sensorID,
        Name:     out.Item["name"].(*types.AttributeValueMemberS).Value,
        Location: out.Item["location"].(*types.AttributeValueMemberS).Value,
    }, nil
}
```

### 時系列データの範囲取得（Query + BETWEEN）

```go
func (r *Repository) QuerySensorData(ctx context.Context, sensorID string, from, to time.Time) ([]SensorData, error) {
    out, err := r.client.Query(ctx, &dynamodb.QueryInput{
        TableName:              &r.tableName,
        KeyConditionExpression: aws.String("PK = :pk AND SK BETWEEN :from AND :to"),
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":pk":   &types.AttributeValueMemberS{Value: "SENSOR#" + sensorID},
            ":from": &types.AttributeValueMemberS{Value: "DATA#" + from.Format(time.RFC3339)},
            ":to":   &types.AttributeValueMemberS{Value: "DATA#" + to.Format(time.RFC3339)},
        },
    })
    if err != nil {
        return nil, err
    }

    results := make([]SensorData, 0, len(out.Items))
    for _, item := range out.Items {
        results = append(results, unmarshalSensorData(sensorID, item))
    }
    return results, nil
}
```

### 最新データ1件取得（降順 + Limit）

```go
func (r *Repository) GetLatestSensorData(ctx context.Context, sensorID string) (*SensorData, error) {
    out, err := r.client.Query(ctx, &dynamodb.QueryInput{
        TableName:              &r.tableName,
        KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :prefix)"),
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":pk":     &types.AttributeValueMemberS{Value: "SENSOR#" + sensorID},
            ":prefix": &types.AttributeValueMemberS{Value: "DATA#"},
        },
        ScanIndexForward: aws.Bool(false),
        Limit:            aws.Int32(1),
    })
    if err != nil {
        return nil, err
    }
    if len(out.Items) == 0 {
        return nil, fmt.Errorf("no data for sensor %s", sensorID)
    }
    d := unmarshalSensorData(sensorID, out.Items[0])
    return &d, nil
}
```

## Update（UpdateItem）

```go
func (r *Repository) UpdateSensorName(ctx context.Context, sensorID, newName string) error {
    _, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
        TableName: &r.tableName,
        Key: map[string]types.AttributeValue{
            "PK": &types.AttributeValueMemberS{Value: "SENSOR#" + sensorID},
            "SK": &types.AttributeValueMemberS{Value: "METADATA"},
        },
        UpdateExpression: aws.String("SET #name = :name"),
        ExpressionAttributeNames: map[string]string{
            "#name": "name", // "name" は予約語
        },
        ExpressionAttributeValues: map[string]types.AttributeValue{
            ":name": &types.AttributeValueMemberS{Value: newName},
        },
        ConditionExpression: aws.String("attribute_exists(PK)"),
    })
    return err
}
```

## Delete（DeleteItem）

```go
func (r *Repository) DeleteSensor(ctx context.Context, sensorID string) error {
    _, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
        TableName: &r.tableName,
        Key: map[string]types.AttributeValue{
            "PK": &types.AttributeValueMemberS{Value: "SENSOR#" + sensorID},
            "SK": &types.AttributeValueMemberS{Value: "METADATA"},
        },
    })
    return err
}
```

## 注意点

| ポイント | 説明 |
|---|---|
| **予約語** | `name`, `status`, `data` などは DynamoDB の予約語。`ExpressionAttributeNames` でエスケープが必要 |
| **数値型** | DynamoDB の数値は `AttributeValueMemberN` で、値は文字列で渡す（`"35.2"` であって `35.2` ではない） |
| **PutItem の上書き** | デフォルトで同じキーがあれば上書き。防ぎたければ `ConditionExpression` を使う |
| **DeleteItem の冪等性** | 存在しないアイテムの削除はエラーにならない（成功扱い） |
| **Query のページング** | 1回の Query で返るデータは最大 1MB。`LastEvaluatedKey` が返ったら続きを取得する必要がある |
