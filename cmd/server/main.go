package main

import (
	_ "Effective-Mobile/docs"
	"Effective-Mobile/internal/app"
	"Effective-Mobile/internal/config"
	"Effective-Mobile/internal/util/logger"
	"context"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic("config: " + err.Error())
	}

	zl, err := logger.New(cfg.Log.Level)
	if err != nil {
		panic("logger: " + err.Error())
	}
	defer zl.Sync()

	sugar := zl.Sugar()

	writer := zap.NewStdLog(zl).Writer()
	gin.DefaultWriter = writer
	gin.DefaultErrorWriter = writer

	application, err := app.New(cfg, sugar)
	if err != nil {
		zl.Fatal("app init", zap.Error(err))
	}
	defer application.DB.Close()

	application.Router.StaticFile(
		"/doc.json",
		"./docs/swagger.json",
	)
	application.Router.GET("/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler),
	)

	application.Router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	srv := &http.Server{
		Addr:         ":" + strconv.Itoa(cfg.Server.Port),
		Handler:      application.Router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zl.Fatal("listen error", zap.Error(err))
		}
	}()
	sugar.Infof("people-enricher up on :%d (%s env)", cfg.Server.Port, cfg.AppEnv)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	sugar.Info("server stopped")
}
