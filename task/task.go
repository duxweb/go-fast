package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/duxweb/go-fast/config"
	"github.com/duxweb/go-fast/global"
	"github.com/duxweb/go-fast/logger"
	"github.com/gookit/color"
	"github.com/hibiken/asynq"
	"github.com/samber/do/v2"
	"github.com/spf13/cast"
	"log/slog"
	"time"
)

type TaskService struct {
	Server    *asynq.Server
	ServeMux  *asynq.ServeMux
	Client    *asynq.Client
	Inspector *asynq.Inspector
	Scheduler *asynq.Scheduler
}

func (s *TaskService) Shutdown() error {
	s.Server.Shutdown()
	s.Scheduler.Shutdown()
	return nil
}

func Init() {
	do.ProvideNamed(global.Injector, "task", NewTask)
}

func NewTask(i do.Injector) (*TaskService, error) {
	dbConfig := config.Load("database").GetStringMapString("redis.drivers.default")
	res := asynq.RedisClientOpt{
		Addr:     dbConfig["host"] + ":" + dbConfig["port"],
		Password: dbConfig["password"],
		DB:       cast.ToInt(dbConfig["db"]),
	}
	server := asynq.NewServer(
		res,
		asynq.Config{
			Logger: &TaskLogger{
				Logger: logger.Log("task"),
			},
			LogLevel:    asynq.WarnLevel,
			Concurrency: 20,
			Queues: map[string]int{
				"high":    10,
				"default": 7,
				"low":     3,
			},
		},
	)

	serveMux := asynq.NewServeMux()
	client := asynq.NewClient(res)
	inspector := asynq.NewInspector(res)

	serveMux.HandleFunc("ping", func(ctx context.Context, t *asynq.Task) error {
		color.Print("⇨ <green>Task server start</>\n")
		return nil
	})

	scheduler := asynq.NewScheduler(res, &asynq.SchedulerOpts{
		LogLevel: asynq.ErrorLevel,
		Location: global.TimeLocation,
		PostEnqueueFunc: func(info *asynq.TaskInfo, err error) {
			if err == nil {
				return
			}
			logger.Log("task").Error("scheduler", err)
		},
	})

	return &TaskService{
		Server:    server,
		ServeMux:  serveMux,
		Client:    client,
		Inspector: inspector,
		Scheduler: scheduler,
	}, nil
}

type Priority string

const (
	PRIORITY_HIGH    Priority = "high"
	PRIORITY_DEFAULT Priority = "default"
	PRIORITY_LOW     Priority = "low"
)

func StartQueue() {
	service := do.MustInvokeNamed[*TaskService](global.Injector, "task")
	if err := service.Server.Run(service.ServeMux); err != nil {
		logger.Log("task").Error("Queue run", "err", err)
	}

}

func StartScheduler() {
	service := do.MustInvokeNamed[*TaskService](global.Injector, "task")
	if err := service.Scheduler.Run(); err != nil {
		logger.Log().Error("Scheduler run", err)
	}
}

func Add(typename string, params any, priority ...Priority) *asynq.TaskInfo {
	group := PRIORITY_DEFAULT
	if len(priority) > 0 {
		group = priority[0]
	}
	return AddTask(typename, params, asynq.Queue(string(group)))
}

func AddDelay(typename string, params any, t time.Duration, priority ...Priority) *asynq.TaskInfo {
	group := PRIORITY_DEFAULT
	if len(priority) > 0 {
		group = priority[0]
	}
	return AddTask(typename, params, asynq.ProcessIn(t), asynq.Queue(string(group)))
}

func AddTime(typename string, params any, t time.Time, priority ...Priority) *asynq.TaskInfo {
	group := PRIORITY_DEFAULT
	if len(priority) > 0 {
		group = priority[0]
	}
	return AddTask(typename, params, asynq.ProcessAt(t), asynq.Queue(string(group)))
}

func AddTask(typename string, params any, opts ...asynq.Option) *asynq.TaskInfo {
	payload, _ := json.Marshal(params)
	task := asynq.NewTask(typename, payload)
	opts = append(opts, asynq.MaxRetry(3))            // Retry count
	opts = append(opts, asynq.Timeout(1*time.Minute)) // Timeout period
	opts = append(opts, asynq.Retention(2*time.Hour)) // Retention time

	info, err := do.MustInvokeNamed[*TaskService](global.Injector, "task").Client.Enqueue(task, opts...)
	if err != nil {
		logger.Log("task").Error("Queue add error", err.Error())
	}
	return info
}

func DelTask(priority Priority, id string) error {
	err := do.MustInvokeNamed[*TaskService](global.Injector, "task").Inspector.DeleteTask(string(priority), id)
	if errors.Is(err, asynq.ErrQueueNotFound) {
		return nil
	}
	if errors.Is(err, asynq.ErrTaskNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

// ListenerScheduler registers a task to be executed on a schedule
// cron: the schedule for the task
// typename: the name of the task type
// params: parameters for the task (can be of any type)
// priority: (optional) the priority group for the task
func ListenerScheduler(cron string, typename string, params any, priority ...Priority) {
	payload, _ := json.Marshal(params)
	task := asynq.NewTask(typename, payload)
	var opts []asynq.Option
	opts = append(opts, asynq.MaxRetry(3))
	opts = append(opts, asynq.Timeout(30*time.Minute))
	opts = append(opts, asynq.Retention(2*time.Hour))
	group := PRIORITY_DEFAULT
	if len(priority) > 0 {
		group = priority[0]
	}
	opts = append(opts, asynq.Queue(string(group)))
	_, err := do.MustInvokeNamed[*TaskService](global.Injector, "task").Scheduler.Register(cron, task, opts...)
	if err != nil {
		panic("Scheduler add error :" + err.Error())
	}
}

// ListenerTask registers a task to be executed on a queue
func ListenerTask(pattern string, handler func(context.Context, *asynq.Task) error) {
	do.MustInvokeNamed[*TaskService](global.Injector, "task").ServeMux.HandleFunc(pattern, handler)
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
