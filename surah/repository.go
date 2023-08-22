package surah

import (
	"alquran/kedaihelpers"

	"github.com/spf13/cast"
)

type Repository interface {
	GetSurah(request GetRequest) []map[string]interface{}
	GetDetailSurah(request GetDetailRequest) []map[string]interface{}
}

type repository struct {
	dbs kedaihelpers.DBStruct
}

func NewRepository(dbs kedaihelpers.DBStruct) *repository {
	return &repository{dbs}
}

func (r *repository) GetSurah(request GetRequest) []map[string]interface{} {
	offsets := (request.Offset - 1) * request.Limit
	queryOrder := `ORDER BY ` + cast.ToString(request.OrderColumn) + ` ` + cast.ToString(request.OrderDirection)
	queryLimit := `LIMIT ` + cast.ToString(request.Limit) + ` OFFSET ` + cast.ToString(offsets)

	sql := `SELECT * FROM surahs
	where lower(transliteration) like '%` + request.Search + `%'
	` + queryOrder + `
	` + queryLimit + `
	`
	rows := r.dbs.DatabaseQueryRows(sql)
	return rows
}

func (r *repository) GetDetailSurah(request GetDetailRequest) []map[string]interface{} {
	queryWhere := ""
	if request.Ayah != "" {
		queryWhere += ` and a.ayah = '` + request.Ayah + `'`
	}
	sql := `select a.surah_id,b.location,b.transliteration as surah_name, a.ayah, a.juz, a.arabic,a.latin,a.translation, a.audio 
	from ayahs as a
	left join surahs as b on b.id = a.surah_id
	where a.surah_id = ` + request.SurahID + `
	` + queryWhere + `
	`
	return r.dbs.DatabaseQueryRows(sql)
}
