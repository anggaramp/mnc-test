package main

import (
	"database/sql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"mnc-test/shared/utils"
	"mnc-test/transport"
	"net/http"
)

func main() {
	var dsn string
	var log *zap.Logger
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	dsn = "host=0.0.0.0 port=10090 user=mnc_test password=.mnc_test! dbname=mnc_test sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	sqlDB, err := db.DB()
	if err != nil {
		return
	}

	defer func(sqlDb *sql.DB) {
		_ = sqlDb.Close()
	}(sqlDB)

	log, err = utils.InitLog(`{
		"level": "debug",
		"encoding": "json",
		"outputPaths": ["./logs/app.log"],
		"errorOutputPaths": ["./logs/app.log"],
		"encoderConfig": {
			"messageKey": "message",
			"levelKey": "level",
			"levelEncoder": "lowercase",
			"callerEncoder":"zapcore.ShortCallerEncoder",
			"callerKey":"caller",
			"timeEncoder":"nanos",
			"timeKey":"time"
		}
	}`)
	if err != nil {
		return
	}

	transport.Setup(e, db, log)

	e.Logger.Fatal(e.Start(":8889"))

}
