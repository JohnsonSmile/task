package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"twitter_task/infra/config"
	"twitter_task/infra/database"
	"twitter_task/infra/logger"
	"twitter_task/service/router"

	"go.uber.org/zap"
)

/*
	[x] 1. 进入官网
	[x] 2. 授权推特
	[] 3. 关注推特
	[] 4. 发布指定推特文案
	[] 5. 成功获取白单/测试资格
*/

func main() {

	// is debug
	var isDebug = flag.Bool("debug", false, "当前是否在debug下开启服务")

	flag.Parse()

	// initialize logger
	logger.Initialize(*isDebug)
	// initialize config
	config.Initialize(*isDebug)
	// initialize database
	database.Initialize(*isDebug)

	// 初始化router
	eg := router.Initialize(*isDebug)
	go func() {
		eg.Run(fmt.Sprintf(":%d", config.GetServerConfig().Port))
	}()

	// 优雅退出
	// kill -9 是捕捉不到的
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	zap.L().Info("Shutdown Server.")
	// 同步logger信息
	zap.L().Sync()
}
