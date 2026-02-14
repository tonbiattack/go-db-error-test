package infra

import (
	"context"
	"os"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func openTestDB(t *testing.T) *gorm.DB {
	// 環境変数でDSNが指定されている場合のみ実DBテストを実行する。
	if dsn := os.Getenv("TEST_DB_DSN"); dsn != "" {
		// MySQLに接続する。
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			t.Fatalf("failed to connect test db: %v", err)
		}
		return db
	}

	// DSNが無い場合はテストをスキップする。
	t.Skip("TEST_DB_DSN is not set")
	return nil
}

func migrateTables(t *testing.T, db *gorm.DB) {
	// テスト用テーブルをマイグレーションする。
	if err := db.AutoMigrate(&IndividualModel{}, &CorporateModel{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
}

func truncateTables(t *testing.T, db *gorm.DB) {
	// 個人テーブルをクリアする。
	if err := db.Exec("TRUNCATE TABLE individual_models").Error; err != nil {
		t.Fatalf("failed to truncate individuals: %v", err)
	}
	// 法人テーブルをクリアする。
	if err := db.Exec("TRUNCATE TABLE corporate_models").Error; err != nil {
		t.Fatalf("failed to truncate corporates: %v", err)
	}
}

func TestGormReader_FetchIndividualsAndCorporates(t *testing.T) {
	// DB接続・準備。
	db := openTestDB(t)
	migrateTables(t, db)
	truncateTables(t, db)

	// リポジトリを構築。
	repo := NewGormReader(db)
	// リクエスト用コンテキスト。
	ctx := context.Background()

	// 事前データ投入（個人）。
	if err := repo.CreateIndividual(ctx, "個人A", "a@example.com"); err != nil {
		t.Fatalf("insert individual: %v", err)
	}
	// 事前データ投入（法人）。
	if err := repo.CreateCorporate(ctx, "法人X"); err != nil {
		t.Fatalf("insert corporate: %v", err)
	}

	// 個人取得が成功すること。
	individuals, err := repo.FetchIndividuals(ctx)
	if err != nil {
		t.Fatalf("fetch individuals: %v", err)
	}
	// 件数が1件であること。
	if len(individuals) != 1 {
		t.Fatalf("expected 1 individual but got %d", len(individuals))
	}

	// 法人取得が成功すること。
	corporates, err := repo.FetchCorporates(ctx)
	if err != nil {
		t.Fatalf("fetch corporates: %v", err)
	}
	// 件数が1件であること。
	if len(corporates) != 1 {
		t.Fatalf("expected 1 corporate but got %d", len(corporates))
	}
}

func TestGormReader_個人テーブルが無いとSELECTが失敗する(t *testing.T) {
	// モックが使えないケース: 実 DB で SELECT の失敗を再現する
	// DB接続・準備。
	db := openTestDB(t)
	migrateTables(t, db)
	truncateTables(t, db)

	// 個人テーブルを削除して SELECT 失敗を再現する。
	if err := db.Migrator().DropTable(&IndividualModel{}); err != nil {
		t.Fatalf("failed to drop individuals: %v", err)
	}

	// リポジトリを構築。
	repo := NewGormReader(db)
	// リクエスト用コンテキスト。
	ctx := context.Background()

	// SELECT が失敗すること。
	_, err := repo.FetchIndividuals(ctx)
	if err == nil {
		t.Fatalf("expected select error but got nil")
	}
}

func TestGormReader_UniqueConstraintError(t *testing.T) {
	// モックが使えないケース: 実 DB でユニーク制約エラーを再現する
	// DB接続・準備。
	db := openTestDB(t)
	migrateTables(t, db)
	truncateTables(t, db)

	// リポジトリを構築。
	repo := NewGormReader(db)
	// リクエスト用コンテキスト。
	ctx := context.Background()

	// 1件目は成功する。
	if err := repo.CreateIndividual(ctx, "個人A", "dup@example.com"); err != nil {
		t.Fatalf("insert individual: %v", err)
	}

	// 同じメールで2件目を投入し、ユニーク制約エラーになること。
	err := repo.CreateIndividual(ctx, "個人B", "dup@example.com")
	if err == nil {
		t.Fatalf("expected unique constraint error but got nil")
	}
}
