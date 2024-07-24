package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"os"
	"ota-go-plugin/shared"
	"time"
)

var Impl MainProgramImpl

type MainProgramImpl struct {
	cfg  shared.AppConfig
	quit chan os.Signal
}

func (m *MainProgramImpl) Start(cfg shared.AppConfig) {
	m.cfg = cfg
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "version v2",
			"time":    time.Now().Format(time.DateTime),
		})
	})
	addr := fmt.Sprintf(":%d", m.cfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}
	slog.Info("v2 listen http server", slog.Any("addr", srv.Addr))
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	m.quit = make(chan os.Signal, 1)
	shared.HttpServerQuit(srv, m.quit, "v2", func() {})
}

func (m *MainProgramImpl) Stop() {
	m.quit <- os.Kill
}
