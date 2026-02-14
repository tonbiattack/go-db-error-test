package usecase

import (
	"context"

	"go-db-error-test/internal/domain"
)

// Reader は個人・法人を取得する読み取りインターフェース。
// ユースケース層はこの境界だけに依存し、永続化手段（GORM/MySQL）を意識しない。
type Reader interface {
	// FetchIndividuals は個人一覧を取得する。
	// 取得順序や絞り込み条件などの詳細は実装側で担保する。
	FetchIndividuals(ctx context.Context) ([]domain.Individual, error)
	// FetchCorporates は法人一覧を取得する。
	// 失敗時はエラーを返し、呼び出し側が継続可否を判断する。
	FetchCorporates(ctx context.Context) ([]domain.Corporate, error)
}
