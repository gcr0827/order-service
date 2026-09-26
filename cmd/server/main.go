package main

// 学习提示：本文件的「优雅关闭」用到 goroutine / 带缓冲 channel / context 超时 / defer，
// 对应 P0-⑤ defer（9/24）、⑥ 并发、⑦ context（W2）。
// 必答的 6 个问题见 docs/main-读码作业.md —— 学完后回来打勾。

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gcr0827/order-service/internal/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/gcr0827/order-service/internal/config"
	"github.com/gcr0827/order-service/internal/handler"
	"github.com/gcr0827/order-service/internal/pkg/response"
	"github.com/gcr0827/order-service/internal/repository"
	"github.com/gcr0827/order-service/internal/service"
)

func main() {
	cfg := config.Load()

	// 初始化DB
	gdb := initDB()

	// 依赖装配（手工 DI，W3 起可换 wire）
	repo := repository.NewOrderRepository()
	svc := service.NewOrderService(repo)
	h := handler.NewOrderHandler(svc)

	// 商品
	productRepository := repository.NewMysqlProductRepository(gdb)
	productService := service.NewProductService(productRepository)
	productHandler := handler.NewProductHandler(productService)

	// 库存
	inventoryRpo := repository.NewStoreProductInventoryRepository(gdb)
	inventoryService := service.NewInventoryService(inventoryRpo)
	inventoryHandler := handler.NewInventoryHandler(inventoryService)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		response.OK(c, gin.H{"status": "healthy"})
	})

	v1 := r.Group("/api/v1")
	{
		v1.GET("/orders/:id", h.GetOrder)
		v1.GET("/product/spu/:id", productHandler.GetSpuByID)
		v1.GET("/product/sku/:id", productHandler.GetSkuByID)
		v1.GET("/product/sku/list", productHandler.GetListSkuBySpuID)
		v1.POST("/product/inventory/deduct", inventoryHandler.DeductStock)
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

func initDB() *gorm.DB {
	g, err := db.New()
	if err != nil {
		log.Fatalf("初始化DB失败: %v", err)
	}

	return g
}
