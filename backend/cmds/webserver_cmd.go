// web 服务

package cmds

import (
	"time"

	"github.com/axiaoxin-com/investool/models"
	"github.com/axiaoxin-com/investool/routes"
	"github.com/axiaoxin-com/investool/webserver"
	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	// ProcessorWebserver web 服务
	ProcessorWebserver = "webserver"
)

// FlagsWebserver cli flags
func FlagsWebserver() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     "config",
			Aliases:  []string{"c"},
			Value:    "./config.yaml",
			Usage:    "配置文件",
			Required: false,
		},
	}
}

// ActionWebserver cli action
func ActionWebserver() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		configFile := c.String("config")

		// 初始化 viper（webserver 需要）
		webserver.InitWithConfigFile(configFile)

		// 加载数据库配置
		if err := models.LoadDatabaseConfig(configFile); err != nil {
			// 加载失败不影响运行
			logrus.Warn("load database config failed:" + err.Error())
		}

		// 初始化数据库连接
		if err := models.InitDatabase(); err != nil {
			logrus.Warn("database initialization failed:" + err.Error())
		}

		// 启动定时任务：每3天执行一次 SyncFund
		startSyncFundScheduler()

		server := webserver.NewGinEngine()
		// 注册路由
		routes.Routes(server)
		// 运行服务
		webserver.Run(server)
		return nil
	}
}

// CommandWebserver 检测器 cli command
func CommandWebserver() *cli.Command {
	flags := FlagsWebserver()
	cmd := &cli.Command{
		Name:   ProcessorWebserver,
		Usage:  "web服务器",
		Flags:  flags,
		Action: ActionWebserver(),
	}
	return cmd
}

// startSyncFundScheduler 启动 SyncFund 定时任务，每3天执行一次
func startSyncFundScheduler() {
	timezone, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		logrus.Errorf("startSyncFundScheduler time LoadLocation error:%v, using Local timezone as default", err.Error())
		timezone, _ = time.LoadLocation("Local")
	}
	sched := gocron.NewScheduler(timezone)

	// 每3天执行一次 SyncFund（从启动时开始，每72小时执行一次）
	sched.Every(48).Hours().Do(func() {
		logrus.Info("Scheduled SyncFund task started")
		UpdateFund()
		logrus.Info("Scheduled SyncFund task completed")
	})

	// 异步启动定时任务
	sched.StartAsync()
	logrus.Info("SyncFund scheduler started: will run every 3 days (72 hours)")
}
