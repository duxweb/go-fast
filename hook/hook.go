package hook

var (
	Boots     []func() error
	Inits     []func() error
	Shutdowns []func() error
)

// Boot 启动钩子
// Boot startup hook
func Boot(f func() error) {
	Boots = append(Boots, f)
}

// Init 初始化钩子
// Init initialization hook
func Init(f func() error) {
	Inits = append(Inits, f)
}

// RunInit 执行初始化钩子
// RunInit execute initialization hook
func RunInit() {
	var err error
	for _, f := range Inits {
		err = f()
		if err != nil {
			panic(err)
		}
	}
}

// RunBoot 执行启动钩子
// RunBoot execute startup hook
func RunBoot() {
	var err error
	for _, f := range Boots {
		err = f()
		if err != nil {
			panic(err)
		}
	}
}

// Shutdown 关机钩子
// Shutdown shutdown hook
func Shutdown(f func() error) {
	Shutdowns = append(Shutdowns, f)
}

// RunShutdown 执行关机钩子
// RunShutdown execute shutdown hook
func RunShutdown() {
	var err error
	for _, f := range Shutdowns {
		err = f()
		if err != nil {
			// 关机时不要 panic，只记录错误
			println("Shutdown error:", err.Error())
		}
	}
}
