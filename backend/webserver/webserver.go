package webserver

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/axiaoxin-com/goutils"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// InitWithConfigFile 根据 webserver 配置文件初始化 webserver
func InitWithConfigFile(configFile string) {
	// 加载配置文件内容到 viper 中以便使用
	configPath, file := path.Split(configFile)
	if configPath == "" {
		configPath = "./"
	}
	ext := path.Ext(file)
	configType := strings.Trim(ext, ".")
	configName := strings.TrimSuffix(file, ext)
	logrus.Infof("load %s type config file %s from %s", configType, configName, configPath)

	if err := goutils.InitViper(configFile, func(e fsnotify.Event) {
		logrus.Warn("Config file changed:" + e.Name)
		logrus.SetLevel(logrus.Level(viper.GetInt("logging.level")))
	}); err != nil {
		// 文件不存在时 1 使用默认配置，其他 err 直接 panic
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(err)
		}
		logrus.Error("Init viper error:" + err.Error())
	}

	// 设置 viper 中 webserver 配置项默认值
	viper.SetDefault("env", "localhost")
	viper.SetDefault("server.addr", ":4869")
	viper.SetDefault("server.pprof", true)
}

// 注意：这里依赖 viper ，必须在外部先对 viper 配置进行加载
func Run(app http.Handler) {
	// 判断是否加载 viper 配置
	if !goutils.IsInitedViper() {
		panic("Running server must init viper by config file first!")
	}

	// 创建 server
	addr := viper.GetString("server.addr")
	srv := &http.Server{
		Addr:         addr,
		Handler:      app,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 10 * time.Minute,
	}
	// Shutdown 时需要调用的方法
	srv.RegisterOnShutdown(func() {
		// TODO
	})

	// 启动 http server
	go func() {
		var ln net.Listener
		var err error
		if strings.ToLower(strings.Split(addr, ":")[0]) == "unix" {
			ln, err = net.Listen("unix", strings.Split(addr, ":")[1])
			if err != nil {
				panic(err)
			}
		} else {
			ln, err = net.Listen("tcp", addr)
			if err != nil {
				panic(err)
			}
		}
		if err := srv.Serve(ln); err != nil {
			logrus.Error("Serve error:" + err.Error())
		}
	}()
	logrus.Infof("Server is running on %s", srv.Addr)

	// 监听中断信号， WriteTimeout 时间后优雅关闭服务
	// syscall.SIGTERM 不带参数的 kill 命令
	// syscall.SIGINT ctrl-c kill -2
	// syscall.SIGKILL 是 kill -9 无法捕获这个信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Infof("Server is shutting down.")

	// 创建一个 context 用于通知 server 3 秒后结束当前正在处理的请求
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Error("Server shutdown with error: " + err.Error())
	}
	logrus.Info("Server exit.")
}
