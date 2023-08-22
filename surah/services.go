package surah

type Services interface {
	GetRequest(request GetRequest) []map[string]interface{}
	GetDetailSurah(request GetDetailRequest) map[string]interface{}
}

type services struct {
	repository Repository
}

func NewServices(repository Repository) *services {
	return &services{repository}
}

func (s *services) GetRequest(request GetRequest) []map[string]interface{} {
	return s.repository.GetSurah(request)
}

func (s *services) GetDetailSurah(request GetDetailRequest) map[string]interface{} {
	rows := s.repository.GetDetailSurah(request)
	obj := map[string]interface{}{}
	obj["surah_name"] = rows[0]["surah_name"]
	obj["location"] = rows[0]["location"]
	obj["detail"] = rows

	return obj
}
