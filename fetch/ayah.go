package fetch

import (
	"alquran/kedaihelpers"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cast"
)

type AyahTypeResponse struct {
	Data []map[string]interface{} `json:"data"`
}

func FetchDataAyah(dbs kedaihelpers.DBStruct) (bool, error) {
	fmt.Println("Start Fetching Data List Ayah")

	idSurah := []int{}
	surahs := dbs.DatabaseQueryRows(`SELECT id from surahs`)
	for _, i := range surahs {
		idSurah = append(idSurah, cast.ToInt(i["id"]))
	}

	// Make GET request to the API
	for _, id := range idSurah {
		resp, err := http.Get("https://web-api.qurankemenag.net/quran-ayah?&surah=" + cast.ToString(id))
		if err != nil {
			return false, err
		}
		defer resp.Body.Close()

		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, err
		}

		// Parse the JSON response
		var dataMap map[string]interface{}
		err = json.Unmarshal(body, &dataMap)
		if err != nil {
			return false, err
		}

		// Extract the "data" array from the response
		dataArray, ok := dataMap["data"].([]interface{})
		if !ok {
			return false, err
		}

		// Convert the data array to []map[string]interface{}
		var result []map[string]interface{}
		for _, item := range dataArray {
			if data, ok := item.(map[string]interface{}); ok {
				result = append(result, data)
			}
		}

		for _, v := range result {
			data := map[string]interface{}{}
			data["id"] = cast.ToInt(v["id"])
			data["surah_id"] = cast.ToInt(v["surah_id"])
			data["ayah"] = cast.ToInt(v["ayah"])
			data["quarter_hizb"] = cast.ToInt(v["quarter_hizb"])
			data["juz"] = cast.ToInt(v["juz"])
			data["manzil"] = cast.ToInt(v["manzil"])
			data["arabic"] = cast.ToString(v["arabic"])
			data["latin"] = cast.ToString(v["latin"])
			arabic_words, _ := json.Marshal(v["arabic_words"])
			data["arabic_words"] = cast.ToString(arabic_words)
			data["translation"] = cast.ToString(v["translation"])
			data["footnotes"] = cast.ToString(v["footnotes"])

			// Mengonversi angka menjadi string dengan format 3 digit
			surahStr := fmt.Sprintf("%03d", cast.ToInt(v["surah_id"]))
			ayahStr := fmt.Sprintf("%03d", cast.ToInt(v["ayah"]))

			// Menggabungkan string surahStr dan ayatStr menjadi id
			idaudio := surahStr + ayahStr
			data["audio"] = "https://media.qurankemenag.net/audio/Abu_Bakr_Ash-Shaatree_aac64/" + idaudio + ".m4a"

			insetNewAyah(data, dbs)

		}
	}
	fmt.Println("Finish Fetching Data List Surah")

	return true, nil
}

func insetNewAyah(newSurah map[string]interface{}, dbs kedaihelpers.DBStruct) {
	_, err := dbs.Dbx.Exec(`
	INSERT
		INTO
	ayahs(
		id,
		surah_id,
		audio,
		ayah,
		quarter_hizb,
		juz,
		manzil,
		arabic,
		latin,
		arabic_words,
		translation,
		footnotes
	) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`, cast.ToInt(newSurah["id"]),
		cast.ToInt(newSurah["surah_id"]), cast.ToString(newSurah["audio"]),
		cast.ToInt(newSurah["ayah"]),
		cast.ToInt(newSurah["quarter_hizb"]), cast.ToInt(newSurah["juz"]), cast.ToInt(newSurah["manzil"]), cast.ToString(newSurah["arabic"]), cast.ToString(newSurah["latin"]), cast.ToString(newSurah["arabic_words"]), cast.ToString(newSurah["translation"]), cast.ToString(newSurah["footnotes"]))
	if err != nil {
		fmt.Println(err.Error())
	}
}
