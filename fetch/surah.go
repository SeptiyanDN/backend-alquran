package fetch

import (
	"alquran/kedaihelpers"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cast"
)

type SurahTypeResponse struct {
	Data []map[string]interface{} `json:"data"`
}

func FetchDataSurah(dbs kedaihelpers.DBStruct) ([]map[string]interface{}, error) {
	fmt.Println("Start Fetching Data List Surah")
	// Make GET request to the API
	resp, err := http.Get("https://web-api.qurankemenag.net/quran-surah")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the JSON response
	var dataMap map[string]interface{}
	err = json.Unmarshal(body, &dataMap)
	if err != nil {
		return nil, err
	}

	// Extract the "data" array from the response
	dataArray, ok := dataMap["data"].([]interface{})
	if !ok {
		return nil, err
	}

	// Convert the data array to []map[string]interface{}
	var result []map[string]interface{}
	for _, item := range dataArray {
		if data, ok := item.(map[string]interface{}); ok {
			result = append(result, data)
		}
	}

	newResult := []map[string]interface{}{}
	for _, v := range result {
		data := map[string]interface{}{}
		data["id"] = cast.ToInt(v["id"])
		data["arabic"] = cast.ToString(v["arabic"])
		data["latin"] = cast.ToString(v["latin"])
		data["transliteration"] = cast.ToString(v["transliteration"])
		data["translation"] = cast.ToString(v["translation"])
		data["num_ayah"] = cast.ToInt(v["num_ayah"])
		data["location"] = cast.ToString(v["location"])
		newResult = append(newResult, data)
	}
	surahExisted := dbs.DatabaseQueryRows(`SELECT * from surahs`)
	count := 0
	for _, newSurah := range newResult {
		transliteration := cast.ToString(newSurah["transliteration"])
		if !isSurahExisted(transliteration, surahExisted) {
			count++
			fmt.Println(fmt.Sprintf("%s,%d", "Insert Data Ke - ", count))
			insetNewSurah(newSurah, dbs)
		}
	}
	fmt.Println("Finish Fetching Data List Surah")

	return result, nil
}

func isSurahExisted(transliteration string, surahExisted []map[string]interface{}) bool {
	for _, cat := range surahExisted {
		if cast.ToString(cat["transliteration"]) == transliteration {
			return true
		}
	}
	return false
}
func insetNewSurah(newSurah map[string]interface{}, dbs kedaihelpers.DBStruct) {
	_, err := dbs.Dbx.Exec(`
	INSERT
		INTO
	surahs(
		id,
		arabic,
		latin,
		transliteration,
		translation,
		num_ayah,
		location
	) VALUES($1,$2,$3,$4,$5,$6,$7)`, cast.ToInt(newSurah["id"]), cast.ToString(newSurah["arabic"]), cast.ToString(newSurah["latin"]), cast.ToString(newSurah["transliteration"]), cast.ToString(newSurah["translation"]), cast.ToInt(newSurah["num_ayah"]), cast.ToString(newSurah["location"]))
	if err != nil {
		fmt.Println(err.Error())
	}
}
