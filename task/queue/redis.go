package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/duxweb/go-fast/config"
	"github.com/duxweb/go-fast/logger"
	"github.com/hibiken/asynq"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type Redis struct {
	Context   context.Context
	Cancel    context.CancelFunc
	Server    *asynq.Server
	ServeMuxs map[string]*asynq.ServeMux
	Client    *asynq.Client
	Inspector *asynq.Inspector
}

func NewRedis() *Redis {

	dbConfig := config.Load("database").GetStringMapString("redis.drivers.default")

	res := asynq.RedisClientOpt{
		Addr:     dbConfig["host"] + ":" + dbConfig["port"],
		Password: dbConfig["password"],
		DB:       cast.ToInt(dbConfig["db"]),
	}

	ctx, cancel := context.WithCancel(context.Background())

	server := asynq.NewServer(
		res,
		asynq.Config{
			BaseContext: func() context.Context {
				return ctx
			},
			Logger: &TaskLogger{
				Logger: logger.Log("task"),
			},
			LogLevel:    asynq.WarnLevel,
			Concurrency: 20,
		},
	)

	client := asynq.NewClient(res)
	inspector := asynq.NewInspector(res)

	return &Redis{
		Context:   ctx,
		Cancel:    cancel,
		Server:    server,
		Client:    client,
		Inspector: inspector,
		ServeMuxs: make(map[string]*asynq.ServeMux),
	}
}

func (q *Redis) Worker(queueName string) {
	serveMux := asynq.NewServeMux()
	q.ServeMuxs[queueName] = serveMux
}

func (q *Redis) Start() error {
	for name, serveMux := range q.ServeMuxs {
		queueName := name
		go func() {
			if err := q.Server.Run(serveMux); err != nil {
				logger.Log("task").Error("Queue run", "queue", queueName, "err", err)
			}
		}()
	}
	return nil
}

func (q *Redis) Register(queueName string, name string, callback func(ctx context.Context, params []byte) error) error {
	serveMux, ok := q.ServeMuxs[queueName]
	if !ok {
		return fmt.Errorf("queue %s not found", queueName)
	}
	serveMux.HandleFunc(name, func(ctx context.Context, t *asynq.Task) error {
		return callback(ctx, t.Payload())
	})

	return nil
}

func (q *Redis) Add(queueName string, add QueueAdd) (string, error) {
	return q.AddDelay(queueName, QueueAddDelay{
		QueueAdd: add,
		Delay:    0,
	})
}

func (q *Redis) AddDelay(queueName string, add QueueAddDelay) (string, error) {
	task := asynq.NewTask(add.Name, add.Params)
	opts := []asynq.Option{
		asynq.ProcessIn(add.Delay),
		asynq.MaxRetry(3),
		asynq.Timeout(1 * time.Minute),
		asynq.Retention(24 * time.Hour),
		asynq.Queue(queueName),
	}
	info, err := q.Client.Enqueue(task, opts...)
	if err != nil {
		logger.Log("task").Error("Queue add", "queue", queueName, "err", err.Error())
		return "", err
	}
	return info.ID, nil
}

func (q *Redis) Names() []string {
	return lo.Keys(q.ServeMuxs)
}

func (q *Redis) List(queueName string, page int, limit int) ([]QueueItem, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 50
	}

	queueInfo, err := q.Inspector.GetQueueInfo(queueName)
	if err != nil {
		return nil, 0, err
	}

	var tasks []*asynq.TaskInfo

	// 获取指定队列的任务
	opts := []asynq.ListOption{
		asynq.Page(page),
		asynq.PageSize(limit),
	}

	tasks, err = q.Inspector.ListPendingTasks(queueName, opts...)
	total := int(queueInfo.Pending)

	if err != nil {
		return nil, 0, err
	}

	// 转换为QueueItem格式
	items := make([]QueueItem, 0)
	for _, task := range tasks {
		var params map[string]any
		if err := json.Unmarshal(task.Payload, &params); err != nil {
			continue
		}

		items = append(items, QueueItem{
			ID:        task.ID,
			QueueName: task.Queue,
			Name:      task.Type,
			Params:    params,
			CreatedAt: task.LastFailedAt,
			RunAt:     task.NextProcessAt,
			Retried:   task.Retried,
		})
	}

	return items, int64(total), nil
}

func (q *Redis) Del(queueName string, id string) error {
	err := q.Inspector.DeleteTask(queueName, id)
	if errors.Is(err, asynq.ErrQueueNotFound) {
		return nil
	}
	if errors.Is(err, asynq.ErrTaskNotFound) {
		return nil
	}
	return err
}

func (q *Redis) Close() error {
	q.Cancel()
	q.Server.Shutdown()
	return nil
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
