package api

import (
	"io"
	"net/http"
	"severless-task-scheduler/db/model"
)

func CreateTask(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		responseError(w, err)
		return
	}

	m := model.Task{
		Parameter: string(bodyBytes),
	}

	err = query.Task.Create(&m)
	if err != nil {
		responseError(w, err)
		return
	}

	responseData(w, m)
}
