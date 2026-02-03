package usecase

import (
	"context"
	"fmt"
	"log"

	"go-db-error-test/internal/domain"
)

type BatchService struct {
	reader Reader
}

func NewBatchService(reader Reader) *BatchService {
	return &BatchService{reader: reader}
}

type Result struct {
	Individuals []domain.Individual
	Corporates  []domain.Corporate
}

func (s *BatchService) Run(ctx context.Context) (*Result, error) {
	var (
		individuals []domain.Individual
		corporates  []domain.Corporate
	)

	ind, err := s.reader.FetchIndividuals(ctx)
	if err != nil {
		log.Printf("[WARN] 個人データの取得に失敗しました: %v", err)
	} else {
		individuals = ind
	}

	corp, err := s.reader.FetchCorporates(ctx)
	if err != nil {
		log.Printf("[WARN] 法人データの取得に失敗しました: %v", err)
	} else {
		corporates = corp
	}

	if len(individuals) == 0 && len(corporates) == 0 {
		return nil, fmt.Errorf("個人・法人ともに取得に失敗したため終了")
	}

	return &Result{
		Individuals: individuals,
		Corporates:  corporates,
	}, nil
}
