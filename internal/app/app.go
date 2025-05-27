package app

import (
	"Effective-Mobile/internal/config"
	ghttp "Effective-Mobile/internal/delivery/http"
	"Effective-Mobile/internal/repository/postgres"
	"Effective-Mobile/internal/service"
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type App struct {
	DB     *pgxpool.Pool
	Router *gin.Engine
}

func New(cfg *config.Config, log *zap.SugaredLogger) (*App, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User, cfg.Database.Password,
		cfg.Database.Host, cfg.Database.Port,
		cfg.Database.DBName, cfg.Database.SSLMode,
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}

	repo := postgres.NewPersonRepo(pool)
	enricher := service.NewEnricher(cfg.Enrichment.Timeout)
	peopleSvc := service.NewPeopleService(repo, enricher, log)

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(ghttp.RequestID())

	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Infow("access",
			"req_id", c.GetString("request_id"),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency", time.Since(start),
		)
	})

	handler := ghttp.NewHandler(peopleSvc)
	handler.Register(r)

	return &App{DB: pool, Router: r}, nil
}
