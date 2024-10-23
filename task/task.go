package task

import (
	"fmt"
	"github.com/RichardKnop/machinery/v2"
	"github.com/RichardKnop/machinery/v2/backends/result"
	brokerseager "github.com/RichardKnop/machinery/v2/brokers/eager"
	machineryConfig "github.com/RichardKnop/machinery/v2/config"
	lockeager "github.com/RichardKnop/machinery/v2/locks/eager"
	"github.com/RichardKnop/machinery/v2/tasks"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"log/slog"
	"time"

	"github.com/duxweb/go-fast/config"
	"github.com/duxweb/go-fast/global"
	"github.com/duxweb/go-fast/logger"
	"github.com/samber/do/v2"
)

type TaskService struct {
	Status bool
	Server *machinery.Server
	Worker map[string]*machinery.Worker
}

func (s *TaskService) Shutdown() error {

	for _, worker := range s.Worker {
		worker.Quit()
	}

	return nil
}

func Init() {
	do.ProvideNamed(global.Injector, "task", NewTask)
}

func NewTask(i do.Injector) (*TaskService, error) {
	err := config.Load("database").ReadInConfig()
	if err != nil {
		return nil, err
	}

	var cnf = &machineryConfig.Config{
		DefaultQueue: "app",
	}

	taskType := config.Load("use").GetString("task.type")
	taskType = lo.Ternary[string](taskType == "", "local", taskType)

	var server *machinery.Server

	switch taskType {
	default:
		broker := brokerseager.New()
		//backend := backendeager.New()
		lock := lockeager.New()
		server = machinery.NewServer(cnf, broker, nil, lock)
		break
	}

	//dbConfig := config.Load("database").GetStringMapString("redis.drivers.default")
	//res := asynq.RedisClientOpt{
	//	Addr:     dbConfig["host"] + ":" + dbConfig["port"],
	//	Password: dbConfig["password"],
	//	DB:       cast.ToInt(dbConfig["db"]),
	//}

	err = server.RegisterTask("ping", func() error {
		fmt.Println("ping")
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &TaskService{
		Server: server,
	}, nil
}

func StartTask() {
	service := do.MustInvokeNamed[*TaskService](global.Injector, "task")

	workers := map[string]*machinery.Worker{}
	workers["app"] = service.Server.NewWorker("app", 10)
	taskWorker := config.Load("use").GetStringMapString("task.worker")

	for name, num := range taskWorker {
		workers[name] = service.Server.NewWorker(name, cast.ToInt(num))
	}
	service.Worker = workers

	//errorHandler := func(err error) {
	//	logger.Log("task").Error("Start task worker", "err", err)
	//}
	//
	//preHandler := func(signature *tasks.Signature) {
	//	logger.Log("task").Info("Pre task worker", "name", signature.Name)
	//}
	//
	//postHandler := func(signature *tasks.Signature) {
	//	logger.Log("task").Info("Post task worker", "name", signature.Name)
	//}

	//for name, worker := range service.Worker {
	//	//worker.SetErrorHandler(errorHandler)
	//	//worker.SetPostTaskHandler(postHandler)
	//	//worker.SetPreTaskHandler(preHandler)
	//	err := worker.Launch()
	//	logger.Log("task").Debug("Start task worker", "name", name)
	//	if err != nil {
	//		logger.Log("task").Error("Queue worker launch", "err", err)
	//	}
	//}

	//Add("ping", []tasks.Arg{})

}

type AddType struct {
	time  *time.Time
	retry int
}

func Add(pattern string, params []tasks.Arg, options ...AddType) *result.AsyncResult {
	return AddTask(pattern, params, options...)
}

func AddDelay(pattern string, params []tasks.Arg, t time.Duration, options ...AddType) *result.AsyncResult {
	option := AddType{}
	if len(options) > 0 {
		option = options[0]
	}
	taskTime := time.Now().Add(t)
	option.time = &taskTime
	return AddTask(pattern, params, option)
}

func AddTime(pattern string, params []tasks.Arg, t time.Time, options ...AddType) *result.AsyncResult {
	option := AddType{}
	if len(options) > 0 {
		option = options[0]
	}
	option.time = &t
	return AddTask(pattern, params, option)
}

func AddTask(pattern string, params []tasks.Arg, options ...AddType) *result.AsyncResult {

	service := do.MustInvokeNamed[*TaskService](global.Injector, "task")

	signature := &tasks.Signature{
		Name:       pattern,
		Args:       params,
		RetryCount: 3,
	}

	var option AddType
	if len(options) > 0 {
		option = options[0]
	}

	if option.retry > 0 {
		signature.RetryCount = option.retry
	}
	if option.time != nil {
		signature.ETA = option.time
	}

	asyncResult, err := service.Server.SendTask(signature)
	if err != nil {
		logger.Log("task").Error("Queue add error", err.Error())
	}

	return asyncResult
}

// AddScheduler registers a task to be executed on a queue
func AddScheduler(cron string, pattern string, params []tasks.Arg) {

	signature := tasks.Signature{
		Name: pattern,
		Args: params,
	}

	service := do.MustInvokeNamed[*TaskService](global.Injector, "task")
	err := service.Server.RegisterPeriodicChain(cron, pattern, &signature)
	if err != nil {
		logger.Log("task").Error("Scheduler add error", err.Error())
	}
}

// ListenerTask registers a task to be executed on a queue
func ListenerTask(pattern string, task any) {
	service := do.MustInvokeNamed[*TaskService](global.Injector, "task")
	err := service.Server.RegisterTask(pattern, task)
	if err != nil {
		panic("Task add error :" + err.Error())
	}
}

type TaskLogger struct {
	Logger *slog.Logger
}

func (t *TaskLogger) Debug(args ...interface{}) {
	t.Logger.Debug(fmt.Sprint(args...))
}

func (t *TaskLogger) Info(args ...interface{}) {
	t.Logger.Info(fmt.Sprint(args...))

}

func (t *TaskLogger) Warn(args ...interface{}) {
	t.Logger.Warn(fmt.Sprint(args...))

}

func (t *TaskLogger) Error(args ...interface{}) {
	t.Logger.Error(fmt.Sprint(args...))

}

func (t *TaskLogger) Fatal(args ...interface{}) {
	t.Logger.Error(fmt.Sprint(args...))
}
