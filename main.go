package main

import (
	"go.uber.org/zap"
	"log"
	"sx-go/ioc"
)

func main() {
	ioc.InitViper()
	db := ioc.InitMongodb()
	loggers := ioc.InitLogger(db)

	defer loggers.Logg.Sync()
	zap.ReplaceGlobals(loggers.Logg) // 关键：让全局 zap.L() 使用初始化后的日志器

	// 4. 再初始化 Web 服务（此时中间件和 handler 中的 zap.L() 已生效）
	app := InitWebServer()
	app.DB = db // 复用已初始化的 DB
	if err := app.server.Run(":8888"); err != nil {
		// 增加错误处理，便于排查启动失败原因
		log.Fatalf("服务启动失败: %v", err)
	}
}
