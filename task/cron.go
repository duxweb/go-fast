package task

import (
	"context"
	"time"

	"github.com/duxweb/go-fast/global"
	"github.com/go-errors/errors"
	"github.com/reugn/go-quartz/job"
	"github.com/reugn/go-quartz/logger"
	"github.com/reugn/go-quartz/quartz"
	"github.com/samber/do/v2"
)

type CronService struct {
	Cron    quartz.Scheduler
	Context context.Context
	Cancel  context.CancelFunc
	Data    map[string]func(ctx context.Context) (any, error)
}

func (s *CronService) Jobs() map[string]func(ctx context.Context) (any, error) {
	return s.Data
}

func (s *CronService) Exists(name string) bool {
	_, ok := s.Data[name]
	return ok
}

func (s *CronService) Check(cron string) bool {
	_, err := quartz.NewCronTrigger(cron)
	return err == nil
}

func (s *CronService) Listener(name string, callback func(ctx context.Context) (any, error)) {
	s.Data[name] = callback
}

func (s *CronService) Scheduler(cron string, name string) error {
	callback, ok := s.Data[name]
	if !ok {
		return errors.New("job not found: " + name)
	}

	cronTrigger, _ := quartz.NewCronTriggerWithLoc(cron, time.Local)
	s.Cron.ScheduleJob(quartz.NewJobDetail(
		job.NewFunctionJob[any](callback),
		quartz.NewJobKey(name),
	), cronTrigger)

	return nil
}

func (s *CronService) Pause(name string) error {
	return s.Cron.PauseJob(quartz.NewJobKey(name))
}

func (s *CronService) Resume(name string) error {
	return s.Cron.ResumeJob(quartz.NewJobKey(name))
}

func (s *CronService) Delete(name string) error {
	return s.Cron.DeleteJob(quartz.NewJobKey(name))
}

func (s *CronService) Clear() error {
	return s.Cron.Clear()
}

func (s *CronService) Start() {

	go func() {
		s.Cron.Start(s.Context)
	}()
}

func (s *CronService) Shutdown() error {
	if s == nil {
		return nil
	}
	s.Cancel()
	s.Cron.Stop()
	s.Cron.Wait(s.Context)
	return nil
}

func CronInit() {
	do.ProvideNamed(global.Injector, "cron", NewCron)
}

func NewCron(i do.Injector) (*CronService, error) {
	logger.SetDefault(logger.NewSimpleLogger(nil, logger.LevelOff))

	ctx, cancel := context.WithCancel(context.Background())
	scheduler := quartz.NewStdScheduler()

	data := &CronService{
		Cron:    scheduler,
		Context: ctx,
		Cancel:  cancel,
		Data:    map[string]func(ctx context.Context) (any, error){},
	}

	return data, nil
}

func Cron() *CronService {
	client := do.MustInvokeNamed[*CronService](global.Injector, "cron")
	return client
}
