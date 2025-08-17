package core

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/duxweb/go-fast/v2/annotation"
	"github.com/duxweb/go-fast/v2/app"
	"github.com/duxweb/go-fast/v2/cache"
	"github.com/duxweb/go-fast/v2/cmd"
	"github.com/duxweb/go-fast/v2/config"
	"github.com/duxweb/go-fast/v2/database"
	"github.com/duxweb/go-fast/v2/event"
	"github.com/duxweb/go-fast/v2/global"
	"github.com/duxweb/go-fast/v2/hook"
	"github.com/duxweb/go-fast/v2/i18n"
	"github.com/duxweb/go-fast/v2/lock"
	"github.com/duxweb/go-fast/v2/permission"
	"github.com/duxweb/go-fast/v2/queue"
	"github.com/duxweb/go-fast/v2/resources"
	"github.com/duxweb/go-fast/v2/route"
	"github.com/duxweb/go-fast/v2/service"
	"github.com/duxweb/go-fast/v2/views"
	"github.com/duxweb/go-fast/v2/web"
	"github.com/gookit/color"
	"github.com/samber/do/v2"
	"github.com/urfave/cli/v3"
)

// Dux 基础结构
// Dux infrastructure
type App struct {
	Service *service.Service
	Command *cmd.Command
	Ctx     context.Context
}

// New 创建App实例
// New create App instance
func New(ctx ...context.Context) *App {
	app := &App{
		Service: service.New(),
		Command: cmd.New(),
	}
	if len(ctx) > 0 {
		app.Ctx = ctx[0]
	} else {
		app.Ctx = context.Background()
	}
	return app
}

// RegisterApp 应用注册, 注册应用模块
// RegisterApp register application, register application module
func (t *App) RegisterApp(calls ...func() *app.Config) {
	for _, call := range calls {
		app.Register(call())
	}
}

// RegisterService 服务注册，注册框架服务
// RegisterService register service, register framework service
func (t *App) RegisterService(calls ...func() *service.Config) {
	for _, call := range calls {
		service.Register(call())
	}
}

// RegisterCmd 命令注册
// RegisterCmd register command
func (t *App) RegisterCmd(calls ...func() []*cli.Command) {
	for _, call := range calls {
		cmd.Register(call())
	}
}

// RegisterDir 自动创建目录
// RegisterDir register folder creation
func (t *App) RegisterDir(dirs ...string) {
	global.DirList = append(global.DirList, dirs...)
}

// RegisterAnnotations 设置索引文件
// RegisterAnnotations set index file
func (t *App) RegisterAnnotations(data []*annotation.File) {
	annotation.Annotations = data
}

// OnInit 注册初始化钩子
// OnInit register initialization hook
func (t *App) OnInit(hooks ...func() error) {
	for _, v := range hooks {
		hook.Init(v)
	}
}

// OnBoot 注册启动钩子
// OnBoot register startup hook
func (t *App) OnBoot(hooks ...func() error) {
	for _, v := range hooks {
		hook.Boot(v)
	}
}

// OnShutdown 注册关机钩子
// OnShutdown register shutdown hook
func (t *App) OnShutdown(hooks ...func() error) {
	for _, v := range hooks {
		hook.Shutdown(v)
	}
}

// Run 运行命令
// Run command
func (t *App) Run() {

	// 初始化依赖注入
	// Initialize dependency injection
	global.Injector = do.New()
	global.Ctx = t.Ctx

	// 注册核心服务
	t.RegisterService(CoreService)
	t.RegisterService(config.Service)
	t.RegisterService(app.Service)
	t.RegisterService(resources.Service)
	t.RegisterService(web.Service)
	t.RegisterService(route.Service)
	t.RegisterService(database.Service)
	t.RegisterService(permission.Service)
	t.RegisterService(i18n.Service)
	t.RegisterService(views.Service)
	t.RegisterService(queue.Service)
	t.RegisterService(cache.Service)
	t.RegisterService(lock.Service)
	t.RegisterService(event.Service)

	// 初始化框架服务
	// Initialize framework service
	t.Service.Init()

	// 执行初始化 hook
	hook.RunInit()

	// 注册服务命令
	// Register service command
	t.Service.Command()

	// 注册基本命令
	// Register basic commands
	t.RegisterCmd(CoreCommand)

	// 设置优雅关机
	// Setup graceful shutdown
	shutdown(t.Service)

	// 初始化命令行
	// Initialize command line
	appCli := &cli.Command{
		Name:     "dux",
		Version:  global.Version,
		Usage:    "dux framework command line tool",
		Commands: cmd.Commands,
	}

	if err := appCli.Run(context.Background(), os.Args); err != nil {
		color.Red.Println("Error: %s", err)
	}
}

// shutdown 设置优雅关机
// shutdown setup graceful shutdown
func shutdown(svc *service.Service) {
	// 创建信号通道
	sigChan := make(chan os.Signal, 1)

	// 监听中断信号
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		// 等待信号
		sig := <-sigChan
		color.Yellow.Println("⇨ Received signal: %v, shutting down gracefully...\n", sig)

		// 创建超时上下文
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		// 执行关机钩子
		done := make(chan struct{})
		go func() {
			// 1. 先执行用户注册的关机钩子
			hook.RunShutdown()

			// 2. 关闭所有服务
			if err := svc.Shutdown(); err != nil {
				color.Red.Println("⇨ Service shutdown error: %v\n", err)
			}

			// 3. 关闭依赖注入容器
			if err := global.Injector.Shutdown(); err != nil {
				color.Red.Println("⇨ Injector shutdown error: %v\n", err)
			}

			close(done)
		}()

		// 等待关机完成或超时
		select {
		case <-done:
			color.Green.Println("⇨ Shutdown completed successfully")
		case <-ctx.Done():
			color.Red.Println("⇨ Shutdown timeout, forcing exit")
		}

		os.Exit(0)
	}()
}
