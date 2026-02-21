# Go / Remix / Dynamo を試してみる

## 概要

Go / Remix（React Router v7）/ Dynamo DB を学習することを主目的としたプロジェクトです。
IoT でセンサーからのデータを蓄積し、データの可視化を行うシステムを想定して開発を行います。

## プロジェクト構成

- apps
  - web： センサーの可視化を行うダッシュボードアプリケーション（Remix）
  - api： センサーからのデータを管理する API（GO）
- packages
- learn：学習した内容のドキュメント置き場

## 概要

Go / Remix（React Router v7）/ Dynamo DB を学習することを主目的としたプロジェクトです。
IoT でセンサーからのデータを蓄積し、データの可視化を行うシステムを想定して開発を行います。

## 技術スタック

### フロントエンド

- 言語：TypeScript
- フロントエンドフレームワーク： Remix（React Router v7）
- パッケージマネージャー：pnpm
- テスト：vitest
- コンポーネント：vitest browser mode
- e2e テスト： playwright

### バックエンド

- 言語：Go
- http フレームワーク：net/http（標準ライブラリ）
- データベース：Dynamo DB
- リンター：golangci-lint v2
- フォーマッター：golangci-lint v2（gofmt, goimports）
- デッドコードチェック：golangci-lint v2（unused）
- 型チェック：golangci-lint v2（typecheck）
- テスト：testing（標準ライブラリ）
- e2e テスト： testing + net/http/httptest

## コマンド

### バックエンド（apps/api）

```bash
# リンター + フォーマッター
go -C apps/api tool golangci-lint run ./...

# テスト
go -C apps/api test ./...
```

## 行動指針

- TDD で実装し、RED、GREEN、REFACTOR を厳格に実施する
- 実装の完了条件
  - 全テストが通過する
  - 静的解析が全てパスする
    - リンター, フォーマット, デッドコードチェック, 型チェック
  - テストカバレッジを 70% をキープする

## 学習内容の管理

本プロジェクトは学習目的のため、学習した内容について整理して管理する。
そのためユーザーが、言語、フレームワーク、ライブラリ、DB などについて質問した場合、回答した内容は整理して必ず learn フォルダ配下に保存すること

### 出力形式

markdown 形式

### フォルダ構成

- learn
  - go
  - remix
  - dynamodb

## 注意点

- 学習内容はトピックごとにまとめる
- 学習内容の管理には必ずサブエージェントを使う
  - ファイルの作成、既存のドキュメントの確認、修正が必要かの判断と修正内用の検討、修正など
