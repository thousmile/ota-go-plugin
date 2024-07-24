package shared

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func HttpServerQuit(srv *http.Server, quit chan os.Signal, srvName string, callback func()) {
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error(fmt.Sprintf("%s : listen http server", srvName), slog.Any("error", err))
			return
		}
	}()
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	signal.Notify(quit, os.Interrupt)
	<-quit
	slog.Info(fmt.Sprintf("%s : shutdown http server ...", srvName))
	callback()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error(fmt.Sprintf("%s : http server shutdown", srvName), slog.Any("error", err))
		panic(err)
	}
	slog.Info(fmt.Sprintf("%s : http server exit ...", srvName))
}
