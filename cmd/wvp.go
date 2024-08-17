package cmd

import (
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"plugin_wvp/services"
	"time"
)

func StartInit() {
	c := cron.New()

	// 添加一个任务，使用 cmd 表达式 "*/5 * * * *" 表示每 5 分钟执行一次
	c.AddFunc("*/5 * * * *", func() {
		logrus.Debug("每 5 分钟执行一次任务:", time.Now())
		services.NewWvpService().DeviceMqttPublish()
	})

	// 启动 cmd 调度器
	c.Start()
}
