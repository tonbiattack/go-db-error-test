package infra

import (
	"context"

	"go-db-error-test/internal/domain"

	"gorm.io/gorm"
)

// GormReader はGORMでDBから読み書きする実装。
type GormReader struct {
	// db はGORMのDB接続。
	db *gorm.DB
}

// NewGormReader はGormReaderのコンストラクタ。
func NewGormReader(db *gorm.DB) *GormReader {
	return &GormReader{db: db}
}

// IndividualModel は個人テーブルに対応するGORMモデル。
type IndividualModel struct {
	// ID は主キー。
	ID    int64  `gorm:"primaryKey;autoIncrement"`
	// Name は個人名。
	Name  string `gorm:"not null"`
	// Email はユニーク制約のあるメールアドレス。
	Email string `gorm:"uniqueIndex;size:191"`
}

// CorporateModel は法人テーブルに対応するGORMモデル。
type CorporateModel struct {
	// ID は主キー。
	ID   int64  `gorm:"primaryKey;autoIncrement"`
	// Name は法人名。
	Name string `gorm:"not null"`
}

// FetchIndividuals は個人一覧を取得し、ドメイン型へ詰め替える。
func (r *GormReader) FetchIndividuals(ctx context.Context) ([]domain.Individual, error) {
	// DB上のモデルを受け取るためのスライス。
	var models []IndividualModel
	if err := r.db.WithContext(ctx).Order("id").Find(&models).Error; err != nil {
		return nil, err
	}

	// ドメイン型へ変換する。
	individuals := make([]domain.Individual, 0, len(models))
	for _, m := range models {
		individuals = append(individuals, domain.Individual{ID: m.ID, Name: m.Name})
	}
	return individuals, nil
}

// FetchCorporates は法人一覧を取得し、ドメイン型へ詰め替える。
func (r *GormReader) FetchCorporates(ctx context.Context) ([]domain.Corporate, error) {
	// DB上のモデルを受け取るためのスライス。
	var models []CorporateModel
	if err := r.db.WithContext(ctx).Order("id").Find(&models).Error; err != nil {
		return nil, err
	}

	// ドメイン型へ変換する。
	corporates := make([]domain.Corporate, 0, len(models))
	for _, m := range models {
		corporates = append(corporates, domain.Corporate{ID: m.ID, Name: m.Name})
	}
	return corporates, nil
}

// CreateIndividual は個人を新規作成する。
func (r *GormReader) CreateIndividual(ctx context.Context, name string, email string) error {
	// 作成するモデルを組み立てる。
	model := IndividualModel{Name: name, Email: email}
	// DBへINSERTする。
	return r.db.WithContext(ctx).Create(&model).Error
}

// CreateCorporate は法人を新規作成する。
func (r *GormReader) CreateCorporate(ctx context.Context, name string) error {
	// 作成するモデルを組み立てる。
	model := CorporateModel{Name: name}
	// DBへINSERTする。
	return r.db.WithContext(ctx).Create(&model).Error
}
