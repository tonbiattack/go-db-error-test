package usecase

import (
	"context"
	"errors"
	"testing"

	"go-db-error-test/internal/domain"
)

// fakeReader はテスト用のReader実装。
type fakeReader struct {
	// individuals は個人取得の正常系返却データ。
	individuals []domain.Individual
	// corporates は法人取得の正常系返却データ。
	corporates  []domain.Corporate
	// indErr は個人取得のエラーを強制するための注入値。
	indErr      error
	// corpErr は法人取得のエラーを強制するための注入値。
	corpErr     error
}

// FetchIndividuals はテスト用の個人取得実装。
func (f *fakeReader) FetchIndividuals(ctx context.Context) ([]domain.Individual, error) {
	// エラーが指定されていればエラーを返す。
	if f.indErr != nil {
		return nil, f.indErr
	}
	// そうでなければ事前に用意したデータを返す。
	return f.individuals, nil
}

// FetchCorporates はテスト用の法人取得実装。
func (f *fakeReader) FetchCorporates(ctx context.Context) ([]domain.Corporate, error) {
	// エラーが指定されていればエラーを返す。
	if f.corpErr != nil {
		return nil, f.corpErr
	}
	// そうでなければ事前に用意したデータを返す。
	return f.corporates, nil
}

func TestBatchService_AllSuccess(t *testing.T) {
	t.Run("個人法人とも成功する", func(t *testing.T) {
		// 両方取得に成功するケース。
		reader := &fakeReader{
			individuals: []domain.Individual{{ID: 1, Name: "個人A"}},
			corporates:  []domain.Corporate{{ID: 10, Name: "法人X"}},
		}
		svc := NewBatchService(reader)

		// 実行。
		res, err := svc.Run(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// 個人・法人がそれぞれ1件取得できていること。
		if len(res.Individuals) != 1 || len(res.Corporates) != 1 {
			t.Fatalf("unexpected result: %#v", res)
		}
	})
}

func TestBatchService_IndividualsFailCorporatesContinue(t *testing.T) {
	t.Run("個人だけ失敗しても法人は続行", func(t *testing.T) {
		// 個人のみ失敗し、法人は成功するケース。
		reader := &fakeReader{
			indErr: errors.New("個人テーブルの SELECT に失敗"),
			corporates: []domain.Corporate{
				{ID: 10, Name: "法人X"},
			},
		}
		svc := NewBatchService(reader)

		// 実行。
		res, err := svc.Run(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// 個人は0件、法人は1件のまま。
		if len(res.Individuals) != 0 {
			t.Fatalf("個人は0件のはず: %#v", res.Individuals)
		}
		if len(res.Corporates) != 1 {
			t.Fatalf("法人は1件のはず: %#v", res.Corporates)
		}
	})
}

func TestBatchService_CorporatesFailIndividualsContinue(t *testing.T) {
	t.Run("法人だけ失敗しても個人は続行", func(t *testing.T) {
		// 法人のみ失敗し、個人は成功するケース。
		reader := &fakeReader{
			individuals: []domain.Individual{{ID: 1, Name: "個人A"}},
			corpErr:     errors.New("法人テーブルの SELECT に失敗"),
		}
		svc := NewBatchService(reader)

		// 実行。
		res, err := svc.Run(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// 個人は1件、法人は0件のまま。
		if len(res.Individuals) != 1 || len(res.Corporates) != 0 {
			t.Fatalf("unexpected result: %#v", res)
		}
	})
}

func TestBatchService_AllFailReturnsError(t *testing.T) {
	t.Run("個人法人とも失敗したらバッチ失敗", func(t *testing.T) {
		// 両方取得に失敗するケース。
		reader := &fakeReader{
			indErr:  errors.New("個人テーブルの SELECT に失敗"),
			corpErr: errors.New("法人テーブルの SELECT に失敗"),
		}
		svc := NewBatchService(reader)

		// 実行してエラーになること。
		_, err := svc.Run(context.Background())
		if err == nil {
			t.Fatalf("error expected but got nil")
		}
	})
}

func TestBatchService_BothSuccessWithEmptyDataReturnsSuccess(t *testing.T) {
	t.Run("個人法人とも成功かつ0件でも正常終了", func(t *testing.T) {
		// 両方の取得処理は成功し、結果のみ空のケース。
		reader := &fakeReader{}
		svc := NewBatchService(reader)

		// 実行してエラーにならないこと。
		res, err := svc.Run(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(res.Individuals) != 0 || len(res.Corporates) != 0 {
			t.Fatalf("unexpected result: %#v", res)
		}
	})
}
