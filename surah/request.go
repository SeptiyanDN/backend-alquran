package surah

type GetRequest struct {
	Search         string `json:"search"`
	Limit          int    `json:"limit"`
	Offset         int    `json:"offset"`
	OrderColumn    string `json:"order_column"`
	OrderDirection string `json:"order_direction"`
}

type GetDetailRequest struct {
	SurahID string `json:"surah_id"`
	Ayah    string `json:"ayah"`
}
