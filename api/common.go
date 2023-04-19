package api

import (
	"encoding/json"
	"net/http"
)

func responseData(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	responseBody := map[string]any{
		"code":    http.StatusOK,
		"message": "success",
		"data":    data,
	}
	responseBytes, _ := json.Marshal(responseBody)
	w.Write(responseBytes)
}

func responseEmpty(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	responseBody := map[string]any{
		"code":    http.StatusOK,
		"message": "success",
	}
	responseBytes, _ := json.Marshal(responseBody)
	w.Write(responseBytes)
}

func responseError(w http.ResponseWriter, errMsg error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	responseBody := map[string]any{
		"code":    http.StatusInternalServerError,
		"message": errMsg.Error(),
	}
	responseBytes, _ := json.Marshal(responseBody)
	w.Write(responseBytes)
}
