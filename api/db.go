package api

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"severless-task-scheduler/db/api"
)

var db *gorm.DB
var query *api.Query

func init() {
	db, _ := gorm.Open(
		mysql.Open(os.Getenv("DATABASE_DSN")),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		},
	)
	query = api.Use(db)
}
