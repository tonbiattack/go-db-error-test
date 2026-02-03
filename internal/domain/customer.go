package domain

// Individual は個人顧客の最小ドメインモデルを表す。
type Individual struct {
	// ID はDB上で割り振られる識別子。
	ID   int64
	// Name は表示名。
	Name string
}

// Corporate は法人顧客の最小ドメインモデルを表す。
type Corporate struct {
	// ID はDB上で割り振られる識別子。
	ID   int64
	// Name は表示名。
	Name string
}
