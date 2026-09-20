package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/gcr0827/order-service/internal/config"
	"github.com/gcr0827/order-service/internal/handler"
	"github.com/gcr0827/order-service/internal/pkg/response"
	"github.com/gcr0827/order-service/internal/repository"
	"github.com/gcr0827/order-service/internal/service"
)

func main() {
	cfg := config.Load()

	// 依赖装配（手工 DI，W3 起可换 wire）
	repo := repository.NewOrderRepository()
	svc := service.NewOrderService(repo)
	h := handler.NewOrderHandler(svc)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		response.OK(c, gin.H{"status": "healthy"})
	})

	v1 := r.Group("/api/v1")
	{
		v1.GET("/orders/:id", h.GetOrder)
	}

	srv := &http.Server{Addr: ":" + cfg.Port, Handler: r}

	go func() {
		log.Printf("server listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v", err)
		}
	}()

	// 优雅关闭：收到信号后给 5 秒处理存量请求
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server forced to shutdown: %v", err)
	}
	log.Println("server exited")
}
