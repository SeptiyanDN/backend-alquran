package handler

import (
	"alquran/helpers"
	"alquran/surah"
	"net/http"

	"github.com/gin-gonic/gin"
)

type handlerSurah struct {
	services surah.Services
}

func NewSurahHandler(services surah.Services) *handlerSurah {
	return &handlerSurah{services}
}

func (h *handlerSurah) GetSurah(c *gin.Context) {
	var request surah.GetRequest
	err := c.ShouldBind(&request)
	if err != nil {
		response := helpers.APIResponse("Invalid Request! Please Check Your Request", http.StatusBadRequest, "Success", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	rows := h.services.GetRequest(request)
	response := helpers.APIResponse("Success To Get All Surah", http.StatusOK, "Successfully", rows)
	c.JSON(http.StatusOK, response)
}

func (h *handlerSurah) GetDetailSurah(c *gin.Context) {
	var request surah.GetDetailRequest
	err := c.ShouldBind(&request)
	if err != nil {
		response := helpers.APIResponse("Invalid Request! Please Check Your Request", http.StatusBadRequest, "Success", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	rows := h.services.GetDetailSurah(request)
	response := helpers.APIResponse("Success To Get Detail Surah", http.StatusOK, "Successfully", rows)
	c.JSON(http.StatusOK, response)
}
