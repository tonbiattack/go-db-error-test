package usecase

import (
	"context"
	"errors"
	"testing"

	"go-db-error-test/internal/domain"
)

type fakeReader struct {
	individuals []domain.Individual
	corporates  []domain.Corporate
	indErr      error
	corpErr     error
}

func (f *fakeReader) FetchIndividuals(ctx context.Context) ([]domain.Individual, error) {
	if f.indErr != nil {
		return nil, f.indErr
	}
	return f.individuals, nil
}

func (f *fakeReader) FetchCorporates(ctx context.Context) ([]domain.Corporate, error) {
	if f.corpErr != nil {
		return nil, f.corpErr
	}
	return f.corporates, nil
}

func TestBatchService_個人法人とも成功する(t *testing.T) {
	reader := &fakeReader{
		individuals: []domain.Individual{{ID: 1, Name: "個人A"}},
		corporates:  []domain.Corporate{{ID: 10, Name: "法人X"}},
	}
	svc := NewBatchService(reader)

	res, err := svc.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Individuals) != 1 || len(res.Corporates) != 1 {
		t.Fatalf("unexpected result: %#v", res)
	}
}

func TestBatchService_個人だけ失敗しても法人は続行(t *testing.T) {
	reader := &fakeReader{
		indErr: errors.New("個人テーブルの SELECT に失敗"),
		corporates: []domain.Corporate{
			{ID: 10, Name: "法人X"},
		},
	}
	svc := NewBatchService(reader)

	res, err := svc.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Individuals) != 0 {
		t.Fatalf("個人は0件のはず: %#v", res.Individuals)
	}
	if len(res.Corporates) != 1 {
		t.Fatalf("法人は1件のはず: %#v", res.Corporates)
	}
}

func TestBatchService_法人だけ失敗しても個人は続行(t *testing.T) {
	reader := &fakeReader{
		individuals: []domain.Individual{{ID: 1, Name: "個人A"}},
		corpErr:     errors.New("法人テーブルの SELECT に失敗"),
	}
	svc := NewBatchService(reader)

	res, err := svc.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Individuals) != 1 || len(res.Corporates) != 0 {
		t.Fatalf("unexpected result: %#v", res)
	}
}

func TestBatchService_個人法人とも失敗したらバッチ失敗(t *testing.T) {
	reader := &fakeReader{
		indErr:  errors.New("個人テーブルの SELECT に失敗"),
		corpErr: errors.New("法人テーブルの SELECT に失敗"),
	}
	svc := NewBatchService(reader)

	_, err := svc.Run(context.Background())
	if err == nil {
		t.Fatalf("error expected but got nil")
	}
}
