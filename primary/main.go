package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"os"
	"ota-go-plugin/shared"
	"plugin"
	"time"
)

func main() {
	var pluginPath string
	flag.StringVar(&pluginPath, "plugin", "v1.so", "path to plugin")
	flag.Parse()
	defSym, err := OpenPlugin(pluginPath)
	if err != nil {
		slog.Error("open plugin: ", slog.Any("error", err))
		return
	}
	cfg := shared.AppConfig{
		Port: 18081,
		Mqtt: shared.MqttConfig{
			Broker:   "127.0.0.1",
			Port:     1883,
			Username: "admin",
			Password: "admin123",
		},
	}
	// 使用反射 将插件 转成 接口
	defImpl := defSym.(shared.MainProgram)
	// 调用 接口 的方法
	go defImpl.Start(cfg)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/stop", func(c *gin.Context) {
		defImpl.Stop()
		c.JSON(200, gin.H{
			"message": "stop success",
			"time":    time.Now().Format(time.RFC3339),
		})
	})
	router.GET("/restart", func(c *gin.Context) {
		query := c.DefaultQuery("pluginPath", pluginPath)
		newSym, err := OpenPlugin(query)
		if err != nil {
			c.JSON(200, gin.H{
				"message": fmt.Sprintf("open plugin fail: %s", err.Error()),
				"time":    time.Now().Format(time.RFC3339),
			})
			return
		}
		// 使用反射 将插件 转成 接口
		newImpl := newSym.(shared.MainProgram)
		// 暂停默认的
		defImpl.Stop()
		// 调用 接口 的方法
		go newImpl.Start(cfg)
		defImpl = newImpl
		// 暂停3秒，等待 插件http服务停止
		time.Sleep(3 * time.Second)
		c.JSON(200, gin.H{
			"message": "restart success",
			"time":    time.Now().Format(time.RFC3339),
		})
	})
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", 18086),
		Handler: router,
	}
	slog.Info("launcher listen http server", slog.Any("addr", srv.Addr))
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	shared.HttpServerQuit(
		srv, quit, "launcher",
		func() {
			defImpl.Stop()
		},
	)
}

func OpenPlugin(pluginPath string) (plugin.Symbol, error) {
	// 打开 .so 文件
	p, err := plugin.Open(pluginPath)
	if err != nil {
		slog.Error("open plugin: ", slog.Any("error", err))
		return nil, err
	}
	// 符号解析，获取插件中的 Impl 变量。必须是 插件中 定义 的变量名称
	sym, err := p.Lookup("Impl")
	if err != nil {
		slog.Error("looking up plugin: ", slog.Any("error", err))
		return nil, err
	} else {
		return sym, nil
	}
}
