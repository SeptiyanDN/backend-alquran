package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/cast"
)

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

func APIResponse(message string, code int, status string, data interface{}) Response {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	jsonResponse := Response{
		Meta: meta,
		Data: data,
	}

	jsonResponseBytes, err := json.Marshal(jsonResponse)
	if err != nil {
		return Response{}
	}

	parsedResponse := map[string]interface{}{}

	err = json.Unmarshal(jsonResponseBytes, &parsedResponse)
	if err != nil {
		return Response{}
	}

	parsedResponse["meta"].(map[string]interface{})["message"] = message

	jsonResponseBytes, err = json.Marshal(parsedResponse)
	if err != nil {
		return Response{}
	}

	err = json.Unmarshal(jsonResponseBytes, &jsonResponse)
	if err != nil {
		return Response{}
	}

	return jsonResponse
}

type ErrorDetail struct {
	Key     string `json:"key"`
	Message string `json:"message"`
}

func ErrorEntityRequest(err error) []ErrorDetail {
	var errorDetails []ErrorDetail

	for _, fieldError := range err.(validator.ValidationErrors) {
		var tag string
		if cast.ToString(fieldError.Tag()) == "required" {
			tag = fmt.Sprintf("%s %s! %s", fieldError.Field(), fieldError.Tag(), "Please fill data "+fieldError.Field())
		} else if cast.ToString(fieldError.Tag()) == "email" {
			tag = fmt.Sprintf("%s %s", fieldError.Field(), "must be format email!")
		}
		errorDetail := ErrorDetail{
			Key:     fieldError.Field(),
			Message: tag,
		}
		errorDetails = append(errorDetails, errorDetail)
	}
	return errorDetails
}

func HandleDuplicateKeyError(err error) error {
	errMsg := err.Error()
	switch {
	case strings.Contains(errMsg, "duplicate key value violates unique constraint \"users_username_key\""):
		return errors.New("Username Already Registered! Please Loggin or Forgot Password.")
	case strings.Contains(errMsg, "duplicate key value violates unique constraint \"users_email_key\""):
		return errors.New("Email Already Registered! Please Loggin or Forgot Password.")
	case strings.Contains(errMsg, "duplicate key value violates unique constraint \"roles_role_name_key\""):
		return errors.New("Role Name Already Registered.")
	default:
		return err
	}
}
