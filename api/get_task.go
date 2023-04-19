package api

import (
	"fmt"
	"net/http"
	"strconv"
)

func GetTask(w http.ResponseWriter, r *http.Request) {
	request := r.URL.Query()
	params, exist := request["id"]
	if !exist {
		responseError(w, fmt.Errorf("id is required"))
		return
	}
	if len(params) == 0 {
		responseError(w, fmt.Errorf("id is required"))
		return
	}
	if params[0] == "" {
		responseError(w, fmt.Errorf("id is required"))
		return
	}
	id, err := strconv.ParseInt(params[0], 10, 64)
	if err != nil {
		responseError(w, err)
		return
	}

	first, err := query.Task.Where(query.Task.ID.Eq(id)).First()
	if err != nil {
		responseError(w, err)
		return
	}

	responseData(w, first)
}
