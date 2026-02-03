package usecase

import (
	"context"

	"go-db-error-test/internal/domain"
)

// Reader は個人・法人を取得する読み取りインターフェース。
type Reader interface {
	// FetchIndividuals は個人一覧を取得する。
	FetchIndividuals(ctx context.Context) ([]domain.Individual, error)
	// FetchCorporates は法人一覧を取得する。
	FetchCorporates(ctx context.Context) ([]domain.Corporate, error)
}
