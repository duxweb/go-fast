package monitor

import (
	"context"
	"encoding/json"
	"github.com/duxweb/go-fast/config"
	"github.com/hibiken/asynq"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
)

var (
	logger *log.Logger
)

func Init() {
	logger = log.New(os.Stderr, "", 0)
	logger.SetOutput(getLumberjack("monitor"))
}

func getLumberjack(name string) *lumberjack.Logger {
	path := config.Load("app").GetString("logger.default.path")
	return &lumberjack.Logger{
		Filename:   path + "/service.log",
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}
}

// Control monitor task
func Control(ctx context.Context, t *asynq.Task) error {
	data := GetMonitorData()
	dataJson, err := json.Marshal(data)
	if err == nil {
		logger.Println(string(dataJson))
	}
	return nil
}
