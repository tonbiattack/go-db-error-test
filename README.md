# go-db-error-test

DB エラーが起きても継続できるかどうかのテストを行うための最小構成サンプルです。

## 構成

- `internal/usecase`: バッチサービスとユニットテスト（フェイク Reader）
- `internal/infra`: GORM 実装 + 実 DB テスト（SELECT 失敗やユニーク制約違反の再現例）

## 事前準備

- Go 1.22+
- MySQL
- `TEST_DB_DSN` を設定（例）

```
export TEST_DB_DSN='user:pass@tcp(127.0.0.1:3306)/test_db?parseTime=true&loc=UTC'
```

### Docker で MySQL を起動する場合

```
docker compose up -d
```

起動後、以下の DSN を使えます。

```
export TEST_DB_DSN='test:test@tcp(127.0.0.1:3306)/test_db?parseTime=true&loc=UTC'
```

## テスト実行

```
go test ./...
```

`TEST_DB_DSN` が未設定の場合、実 DB を使うテストは `skip` されます。
