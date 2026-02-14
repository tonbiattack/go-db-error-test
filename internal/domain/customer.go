package domain

// Individual は個人顧客の最小ドメインモデルを表す。
// 永続化層の都合（テーブル名/タグ等）を持ち込まないことで、
// ユースケース層の関心を業務データへ限定する。
type Individual struct {
	// ID はDB上で割り振られる識別子。
	ID   int64
	// Name は表示名。
	Name string
}

// Corporate は法人顧客の最小ドメインモデルを表す。
// Individual と同じく、ドメインでは最小限の属性のみ保持する。
type Corporate struct {
	// ID はDB上で割り振られる識別子。
	ID   int64
	// Name は表示名。
	Name string
}
