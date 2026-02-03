package infra

import (
	"context"

	"go-db-error-test/internal/domain"

	"gorm.io/gorm"
)

type GormReader struct {
	db *gorm.DB
}

func NewGormReader(db *gorm.DB) *GormReader {
	return &GormReader{db: db}
}

type IndividualModel struct {
	ID    int64  `gorm:"primaryKey;autoIncrement"`
	Name  string `gorm:"not null"`
	Email string `gorm:"uniqueIndex;size:255"`
}

type CorporateModel struct {
	ID   int64  `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"not null"`
}

func (r *GormReader) FetchIndividuals(ctx context.Context) ([]domain.Individual, error) {
	var models []IndividualModel
	if err := r.db.WithContext(ctx).Order("id").Find(&models).Error; err != nil {
		return nil, err
	}

	individuals := make([]domain.Individual, 0, len(models))
	for _, m := range models {
		individuals = append(individuals, domain.Individual{ID: m.ID, Name: m.Name})
	}
	return individuals, nil
}

func (r *GormReader) FetchCorporates(ctx context.Context) ([]domain.Corporate, error) {
	var models []CorporateModel
	if err := r.db.WithContext(ctx).Order("id").Find(&models).Error; err != nil {
		return nil, err
	}

	corporates := make([]domain.Corporate, 0, len(models))
	for _, m := range models {
		corporates = append(corporates, domain.Corporate{ID: m.ID, Name: m.Name})
	}
	return corporates, nil
}

func (r *GormReader) CreateIndividual(ctx context.Context, name string, email string) error {
	model := IndividualModel{Name: name, Email: email}
	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *GormReader) CreateCorporate(ctx context.Context, name string) error {
	model := CorporateModel{Name: name}
	return r.db.WithContext(ctx).Create(&model).Error
}
