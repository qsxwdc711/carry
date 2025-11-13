package main

import (
	"log"
	"sx-go/ioc"
)

func main() {
	ioc.InitViper()

	app := InitWebServer()
	ioc.InitLogger(app.mongodb)
	if err := app.server.Run(":8091"); err != nil {
		// 增加错误处理，便于排查启动失败原因
		log.Fatalf("服务启动失败: %v", err)
	}
}
