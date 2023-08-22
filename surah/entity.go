package surah

type Surah struct {
	ID              int64  `json:"id"`
	Arabic          string `json:"arabic"`
	Latin           string `json:"latin"`
	Transliteration string `json:"transliteration"`
	Translation     string `json:"translation"`
	NumAyah         int64  `json:"num_ayah"`
	Location        string `json:"location"`
}
