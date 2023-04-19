package api

import (
	"encoding/json"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"io"
	"net/http"
	"os"
	"severless-task-scheduler/db/api"
	"severless-task-scheduler/db/model"
)

var db *gorm.DB
var query *api.Query

func init() {
	db, _ = gorm.Open(
		mysql.Open(os.Getenv("DATABASE_DSN")),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		},
	)
	query = api.Use(db)
}

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
