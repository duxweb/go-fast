package queue

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"time"

	"github.com/maragudk/goqite"
	"github.com/maragudk/goqite/jobs"
	"github.com/samber/lo"
)

type BaseQueue struct {
	QueueName string
	Queue     *goqite.Queue
	Job       *jobs.Runner
}

type Base struct {
	db      *sql.DB
	queues  map[string]*BaseQueue
	Context context.Context
	Cancel  context.CancelFunc
}

func NewBase() *Base {

	ctx, cancel := context.WithCancel(context.Background())

	db, err := sql.Open("sqlite3", "file:/queue.db?vfs=memdb&_journal=WAL&_timeout=5000&_fk=true")
	if err != nil {
		log.Fatalln(err)
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if err := goqite.Setup(ctx, db); err != nil {
		slog.Error("Error in setup", "error", err)
	}

	return &Base{db: db, Context: ctx, Cancel: cancel, queues: make(map[string]*BaseQueue)}
}

func (q *Base) Worker(queueName string) {
	queue := goqite.New(goqite.NewOpts{
		DB:   q.db,
		Name: queueName,
	})
	job := jobs.NewRunner(jobs.NewRunnerOpts{
		Limit:        1,
		Log:          nil,
		PollInterval: 10 * time.Millisecond,
		Queue:        queue,
	})

	q.queues[queueName] = &BaseQueue{
		QueueName: queueName,
		Queue:     queue,
		Job:       job,
	}
}

func (q *Base) Start() error {
	for _, queue := range q.queues {
		go queue.Job.Start(q.Context)
	}
	return nil
}

func (q *Base) Register(queueName string, name string, callback func(ctx context.Context, params []byte) error) error {
	queue, ok := q.queues[queueName]
	if !ok {
		return fmt.Errorf("queue %s not found", queueName)
	}
	queue.Job.Register(name, callback)
	return nil
}
func (q *Base) Add(queueName string, add QueueAdd) (string, error) {
	return q.AddDelay(queueName, QueueAddDelay{
		QueueAdd: add,
		Delay:    0,
	})
}

func (q *Base) AddDelay(queueName string, add QueueAddDelay) (string, error) {
	queue, ok := q.queues[queueName]
	if !ok {
		return "", fmt.Errorf("queue %s not found", queueName)
	}
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(queueMessage{Name: add.Name, Message: add.Params}); err != nil {
		return "", err
	}

	id, err := queue.Queue.SendAndGetID(q.Context, goqite.Message{Body: buf.Bytes(), Delay: add.Delay})
	if err != nil {
		return "", err
	}
	return string(id), nil
}

func (q *Base) Names() []string {
	return lo.Keys(q.queues)
}

func (q *Base) List(queueName string, page int, limit int) ([]QueueItem, int64, error) {
	var rows *sql.Rows
	var err error

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 50
	}

	offset := (page - 1) * limit
	var count int64
	var countQuery string
	var listQuery string
	args := []any{}
	countArgs := []any{}

	countQuery = "SELECT COUNT(*) FROM goqite"
	listQuery = "SELECT id, queue, body, timeout, received, created, datetime(created) as date FROM goqite"

	var conditions []string
	if queueName != "" {
		conditions = append(conditions, "queue = ?")
		args = append(args, queueName)
		countArgs = append(countArgs, queueName)
	}

	if len(conditions) > 0 {
		whereClause := " WHERE " + strings.Join(conditions, " AND ")
		countQuery += whereClause
		listQuery += whereClause
	}

	listQuery += " ORDER BY date DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)
	// 获取总数
	if queueName == "" {
		err = q.db.QueryRowContext(q.Context, countQuery, countArgs...).Scan(&count)
	} else {
		err = q.db.QueryRowContext(q.Context, countQuery, countArgs...).Scan(&count)
	}
	if err != nil {
		return nil, 0, err
	}

	// 获取列表数据
	fmt.Println(listQuery, args)
	rows, err = q.db.QueryContext(q.Context, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []QueueItem
	for rows.Next() {
		var id, queue string
		var body []byte
		var timeout string
		var received uint
		var created string
		var date string
		err = rows.Scan(&id, &queue, &body, &timeout, &received, &created, &date)
		if err != nil {
			return nil, 0, err
		}

		var msg queueMessage
		if err = gob.NewDecoder(bytes.NewReader(body)).Decode(&msg); err != nil {
			return nil, 0, err
		}

		var params map[string]any
		if err = json.Unmarshal(msg.Message, &params); err != nil {
			// Handle invalid JSON by returning empty map instead of error
			params = make(map[string]any)
		}
		// 将UTC时间转换为当前时区时间
		createdTime, err := time.Parse(time.RFC3339, created)
		if err != nil {
			return nil, 0, fmt.Errorf("解析时间失败: %v", err)
		}

		runAt, err := time.Parse(time.RFC3339, timeout)
		if err != nil {
			return nil, 0, fmt.Errorf("解析时间失败: %v", err)
		}

		items = append(items, QueueItem{
			ID:        id,
			QueueName: queue,
			Name:      msg.Name,
			Params:    params,
			Retried:   int(received),
			CreatedAt: createdTime.In(time.Local),
			RunAt:     runAt.In(time.Local),
		})
	}

	return items, count, nil

}

func (q *Base) Del(queueName string, id string) error {
	queue, ok := q.queues[queueName]
	if !ok {
		return fmt.Errorf("queue %s not found", queueName)
	}
	return queue.Queue.Delete(q.Context, goqite.ID(id))
}

func (q *Base) Close() error {
	q.Cancel()
	return q.db.Close()
}

type queueMessage struct {
	Name    string
	Message []byte
}

type BaseLogger struct{}

func (l *BaseLogger) Info(msg string, args ...any) {
	slog.Info(msg, args...)
}
