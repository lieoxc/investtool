// web 服务

package cmds

import (
	"github.com/axiaoxin-com/investool/routes"
	"github.com/axiaoxin-com/investool/webserver"
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
			Value:    "./config.toml",
			Usage:    "配置文件",
			Required: false,
		},
	}
}

// ActionWebserver cli action
func ActionWebserver() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		configFile := c.String("config")
		webserver.InitWithConfigFile(configFile)

		// 启动定时任务
		//cron.RunCronJobs(true)

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
