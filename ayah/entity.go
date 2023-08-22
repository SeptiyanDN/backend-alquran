package ayah

import (
	"github.com/jinzhu/gorm/dialects/postgres"
)

type Ayah struct {
	ID          int64          `json:"id"`
	SurahID     int64          `json:"surah_id"`
	Ayah        int64          `json:"ayah"`
	Audio       string         `json:"audio"`
	QuarterHizb int64          `json:"quarter_hizb"`
	Juz         int64          `json:"juz"`
	Manzil      int64          `json:"manzil"`
	Arabic      string         `json:"arabic"`
	Latin       string         `json:"latin"`
	ArabicWords postgres.Jsonb `gorm:"type:jsonb" json:"arabic_words"`
	Translation string         `json:"translation"`
	Footnotes   string         `json:"footnotes"`
}
