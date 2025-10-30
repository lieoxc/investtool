package cmds

import (
	"github.com/axiaoxin-com/investool/cron"
	"github.com/axiaoxin-com/investool/models"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	// ProcessorJSON 导出json数据文件
	ProcessorJSON = "json"
)

// FlagsJSON cli flags
func FlagsJSON() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:    "dump",
			Aliases: []string{"d"},
			Usage:   "导出json数据文件",
		},
		&cli.StringFlag{
			Name:     "config",
			Aliases:  []string{"c"},
			Value:    "./config.yaml",
			Usage:    "配置文件",
			Required: false,
		},
	}
}

// ActionJSON dump json files
func ActionJSON() func(c *cli.Context) error {
	return func(c *cli.Context) error {

		// 加载数据库配置
		configFile := c.String("config")
		if err := models.LoadDatabaseConfig(configFile); err != nil {
			// 加载失败不影响运行
			logrus.Warn("load database config failed:" + err.Error())
			return err
		}

		// 初始化数据库连接
		if err := models.InitDatabase(); err != nil {
			logrus.Warn("database initialization failed:" + err.Error())
			return err
		}

		if c.Bool("d") {
			cron.SyncFund()
			// cron.SyncFundManagers()
			// cron.SyncIndustryList()
			return nil
		}
		return nil
	}
}

// CommandJSON dump json files cmd
func CommandJSON() *cli.Command {
	flags := FlagsJSON()
	cmd := &cli.Command{
		Name:   ProcessorJSON,
		Usage:  "JSON数据",
		Flags:  flags,
		Action: ActionJSON(),
	}
	return cmd
}
