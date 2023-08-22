package fetch

import (
	"alquran/kedaihelpers"
)

func RunnerFetching(dbs kedaihelpers.DBStruct) {
	go func() {
		// FetchDataSurah(dbs)
		FetchDataAyah(dbs)
	}()

}
