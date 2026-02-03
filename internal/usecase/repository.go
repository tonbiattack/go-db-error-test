package usecase

import (
	"context"

	"go-db-error-test/internal/domain"
)

type Reader interface {
	FetchIndividuals(ctx context.Context) ([]domain.Individual, error)
	FetchCorporates(ctx context.Context) ([]domain.Corporate, error)
}
