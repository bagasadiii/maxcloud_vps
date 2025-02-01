package utils

import (
	"encoding/json"
	"net/http"
)

type APIResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func JSONResponse(w http.ResponseWriter, statusCode int, data any) {
	response := APIResp{
		Code:    statusCode,
		Message: http.StatusText(statusCode),
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

