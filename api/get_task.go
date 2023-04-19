package api

import (
	"encoding/json"
	"io"
	"net/http"
)

type GetTaskRequest struct {
	ID int64 `json:"id"`
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		responseError(w, err)
		return
	}

	request := GetTaskRequest{}
	err = json.Unmarshal(bodyBytes, &request)
	if err != nil {
		responseError(w, err)
		return
	}

	first, err := query.Task.Where(query.Task.ID.Eq(request.ID)).First()
	if err != nil {
		responseError(w, err)
		return
	}

	responseData(w, first)
}
