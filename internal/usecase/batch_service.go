package usecase

import (
	"context"
	"fmt"
	"log"

	"go-db-error-test/internal/domain"
)

// BatchService は個人・法人取得をまとめて実行するユースケース。
type BatchService struct {
	// reader はデータ取得のための依存。
	reader Reader
}

// NewBatchService は BatchService のコンストラクタ。
func NewBatchService(reader Reader) *BatchService {
	return &BatchService{reader: reader}
}

// Result はバッチ実行結果を表す。
type Result struct {
	// Individuals は取得できた個人一覧。
	Individuals []domain.Individual
	// Corporates は取得できた法人一覧。
	Corporates  []domain.Corporate
}

// Run は個人・法人の取得を試行し、片方が失敗しても継続する。
func (s *BatchService) Run(ctx context.Context) (*Result, error) {
	var (
		// individuals は個人取得の結果。
		individuals []domain.Individual
		// corporates は法人取得の結果。
		corporates  []domain.Corporate
	)

	// 個人取得に失敗しても処理継続する。
	ind, err := s.reader.FetchIndividuals(ctx)
	if err != nil {
		log.Printf("[WARN] 個人データの取得に失敗しました: %v", err)
	} else {
		individuals = ind
	}

	// 法人取得に失敗しても処理継続する。
	corp, err := s.reader.FetchCorporates(ctx)
	if err != nil {
		log.Printf("[WARN] 法人データの取得に失敗しました: %v", err)
	} else {
		corporates = corp
	}

	// 両方失敗（どちらも空）ならバッチ失敗として返す。
	if len(individuals) == 0 && len(corporates) == 0 {
		return nil, fmt.Errorf("個人・法人ともに取得に失敗したため終了")
	}

	// 片方でも成功していれば結果を返す。
	return &Result{
		Individuals: individuals,
		Corporates:  corporates,
	}, nil
}
