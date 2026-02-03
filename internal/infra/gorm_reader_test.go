package infra

import (
	"context"
	"os"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func openTestDB(t *testing.T) *gorm.DB {
	if dsn := os.Getenv("TEST_DB_DSN"); dsn != "" {
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			t.Fatalf("failed to connect test db: %v", err)
		}
		return db
	}

	t.Skip("TEST_DB_DSN is not set")
	return nil
}

func migrateTables(t *testing.T, db *gorm.DB) {
	if err := db.AutoMigrate(&IndividualModel{}, &CorporateModel{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
}

func truncateTables(t *testing.T, db *gorm.DB) {
	if err := db.Exec("TRUNCATE TABLE individual_models").Error; err != nil {
		t.Fatalf("failed to truncate individuals: %v", err)
	}
	if err := db.Exec("TRUNCATE TABLE corporate_models").Error; err != nil {
		t.Fatalf("failed to truncate corporates: %v", err)
	}
}

func TestGormReader_FetchIndividualsAndCorporates(t *testing.T) {
	db := openTestDB(t)
	migrateTables(t, db)
	truncateTables(t, db)

	repo := NewGormReader(db)
	ctx := context.Background()

	if err := repo.CreateIndividual(ctx, "個人A", "a@example.com"); err != nil {
		t.Fatalf("insert individual: %v", err)
	}
	if err := repo.CreateCorporate(ctx, "法人X"); err != nil {
		t.Fatalf("insert corporate: %v", err)
	}

	individuals, err := repo.FetchIndividuals(ctx)
	if err != nil {
		t.Fatalf("fetch individuals: %v", err)
	}
	if len(individuals) != 1 {
		t.Fatalf("expected 1 individual but got %d", len(individuals))
	}

	corporates, err := repo.FetchCorporates(ctx)
	if err != nil {
		t.Fatalf("fetch corporates: %v", err)
	}
	if len(corporates) != 1 {
		t.Fatalf("expected 1 corporate but got %d", len(corporates))
	}
}

func TestGormReader_UniqueConstraintError(t *testing.T) {
	// モックが使えないケース: 実 DB でユニーク制約エラーを再現する
	db := openTestDB(t)
	migrateTables(t, db)
	truncateTables(t, db)

	repo := NewGormReader(db)
	ctx := context.Background()

	if err := repo.CreateIndividual(ctx, "個人A", "dup@example.com"); err != nil {
		t.Fatalf("insert individual: %v", err)
	}

	err := repo.CreateIndividual(ctx, "個人B", "dup@example.com")
	if err == nil {
		t.Fatalf("expected unique constraint error but got nil")
	}
}
