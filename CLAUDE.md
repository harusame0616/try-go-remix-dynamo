# Go / Remix / Dynamo を試してみる

## 概要

Go / Remix（React Router v7）/ Dynamo DB を学習することを主目的としたプロジェクトです。
IoT でセンサーからのデータを蓄積し、データの可視化を行うシステムを想定して開発を行います。

## プロジェクト構成

- apps
  - web： センサーの可視化を行うダッシュボードアプリケーション（Remix）
  - api： センサーからのデータを管理する API（GO）
- packages# Go / Remix / Dynamo を試してみる

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
- テスト：testing（標準ライブラリ）
- e2e テスト： testing + net/http/httptest

### 行動指針

- TDD で実装し、RED、GREEN、REFACTOR を厳格に実施する
- 実装の完了条件
  - 全テストが通過する
  - 静的解析が全てパスする
    - リンター, フォーマット, デッドコードチェック, 型チェック
  - テストカバレッジを 70% をキープする
