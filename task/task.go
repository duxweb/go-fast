package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/demdxx/gocast/v2"
	"github.com/duxweb/go-fast/config"
	"github.com/duxweb/go-fast/global"
	"github.com/duxweb/go-fast/logger"
	"github.com/gookit/color"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/samber/do"
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
	dbConfig := config.Load("database").GetStringMapString("redis")
	res := asynq.RedisClientOpt{
		Addr:     dbConfig["host"] + ":" + dbConfig["port"],
		Password: dbConfig["password"],
		DB:       gocast.Number[int](dbConfig["db"]),
	}

	server := asynq.NewServer(
		res,
		asynq.Config{
			Logger: &TaskLogger{
				Logger: logger.New(
					logger.GetWriter(
						zerolog.LevelDebugValue,
						"task",
						true,
					),
				).With().Timestamp().Logger(),
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
			logger.Log().Error().Msgf("scheduler: ", err.Error())
		},
	})

	do.ProvideValue[*TaskService](global.Injector, &TaskService{
		Server:    server,
		ServeMux:  serveMux,
		Client:    client,
		Inspector: inspector,
		Scheduler: scheduler,
	})
}

type Priority string

const (
	PRIORITY_HIGH    Priority = "high"
	PRIORITY_DEFAULT Priority = "default"
	PRIORITY_LOW     Priority = "low"
)

func StartQueue() {
	if err := do.MustInvoke[*TaskService](global.Injector).Server.Run(do.MustInvoke[*TaskService](global.Injector).ServeMux); err != nil {
		logger.Log().Error().Msgf("Queue service cannot be started: %v", err)
	}
}

func StartScheduler() {
	if err := do.MustInvoke[*TaskService](global.Injector).Scheduler.Run(); err != nil {
		logger.Log().Error().Msgf("Scheduler service cannot be started: %v", err)
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

	info, err := do.MustInvoke[*TaskService](global.Injector).Client.Enqueue(task, opts...)
	if err != nil {
		logger.Log().Error().Msg("Queue add error :" + err.Error())
	}
	return info
}

func DelTask(priority Priority, id string) error {
	err := do.MustInvoke[*TaskService](global.Injector).Inspector.DeleteTask(string(priority), id)
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
	_, err := do.MustInvoke[*TaskService](global.Injector).Scheduler.Register(cron, task, opts...)
	if err != nil {
		panic("Scheduler add error :" + err.Error())
	}
}

// ListenerTask registers a task to be executed on a queue
func ListenerTask(pattern string, handler func(context.Context, *asynq.Task) error) {
	do.MustInvoke[*TaskService](global.Injector).ServeMux.HandleFunc(pattern, handler)
}

type TaskLogger struct {
	Logger zerolog.Logger
}

func (t *TaskLogger) Debug(args ...interface{}) {
	t.Logger.Debug().Msg(fmt.Sprint(args...))
}

func (t *TaskLogger) Info(args ...interface{}) {
	t.Logger.Info().Msg(fmt.Sprint(args...))

}

func (t *TaskLogger) Warn(args ...interface{}) {
	t.Logger.Warn().Msg(fmt.Sprint(args...))

}

func (t *TaskLogger) Error(args ...interface{}) {
	t.Logger.Error().Msg(fmt.Sprint(args...))

}

func (t *TaskLogger) Fatal(args ...interface{}) {
	t.Logger.Fatal().Msg(fmt.Sprint(args...))
}
