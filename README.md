# go-db-error-test

DB エラーが起きても処理を継続できるか、実 DB と GORM で再現・検証するための最小構成サンプルです。
主に「片方の SELECT が失敗しても、もう片方は継続する」「ユニーク制約エラーを実 DB で再現する」ことを目的にしています。

## できること

- 個人・法人の取得処理を分離し、片方の失敗でももう片方を継続するバッチ実行の検証
- 実 DB での SELECT 失敗（テーブル削除）やユニーク制約違反の再現
- GORM を用いた最小限の read/write 実装の動作確認

## ディレクトリ構成

- `internal/domain`: ドメインモデル
- `internal/usecase`: バッチサービスとユニットテスト
- `internal/infra`: GORM 実装 + 実 DB テスト

## テーブル概要

- `individual_models`: `id` (PK), `name`, `email` (unique)
- `corporate_models`: `id` (PK), `name`

※ テスト内で `AutoMigrate` を実行し、必要なテーブルを作成します。

## 事前準備

- Go 1.22+
- MySQL
- `TEST_DB_DSN` の設定

例:

```bash
export TEST_DB_DSN='user:pass@tcp(127.0.0.1:3306)/test_db?parseTime=true&loc=UTC'
```

### Docker で MySQL を起動する場合

```bash
docker compose up -d
```

起動後、以下の DSN を使えます。

```bash
export TEST_DB_DSN='test:test@tcp(127.0.0.1:3306)/test_db?parseTime=true&loc=UTC'
```

## テストの実行

全テスト:

```bash
go test ./...
```

実 DB テストのみ:

```bash
go test ./internal/infra -v
```

ユースケースのユニットテストのみ:

```bash
go test ./internal/usecase -v
```

`TEST_DB_DSN` が未設定の場合、実 DB を使うテストは `skip` されます。

## 実装のポイント

- `internal/infra` のテストでは実 MySQL を使い、個人テーブルを削除して SELECT 失敗を再現します。
- `internal/infra` のテストでは `email` のユニーク制約違反も再現します。
- `internal/usecase` のテストでは、片方の取得が失敗してももう片方を継続できることを確認します。

## よくあるつまずき

- MySQL が起動していない場合は `TEST_DB_DSN` を設定しても接続に失敗します。
- `docker compose up -d` 後は数秒待ってからテストを実行してください。
