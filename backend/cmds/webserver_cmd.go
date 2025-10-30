// web 服务

package cmds

import (
	"github.com/axiaoxin-com/investool/models"
	"github.com/axiaoxin-com/investool/routes"
	"github.com/axiaoxin-com/investool/webserver"
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
