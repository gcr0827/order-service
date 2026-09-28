package main

// ============================================================
// 服务入口：只做三件事 —— 建依赖 → 建路由 → 跑服务
//
// 为什么入口要"薄"：
//   · 依赖装配   → wire.go
//   · 路由注册   → router.go
//   · 服务生命周期 → 本文件
//   好处：加接口只动 router.go，加依赖只动 wire.go，入口本身很少需要改。
//
// ⚠️ 学习提示：本文件的「优雅关闭」用到 goroutine / 带缓冲 channel /
//    context 超时 / defer —— 对应 P0-⑤ defer、⑥ 并发、⑦ context（W2）。
//    必答 6 问见 docs/main-读码作业.md，学完回来打勾。
// ============================================================

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
	"gorm.io/gorm"

	"github.com/gcr0827/order-service/internal/config"
	"github.com/gcr0827/order-service/internal/db"
)

func main() {
	// ① 读配置（端口等）。DB 凭据走环境变量，不在配置文件里
	cfg := config.Load()

	// ② 建数据库连接。失败直接退出（为什么见 initDB 的注释）
	gdb := initDB()

	// ③ 装配所有依赖 → 得到一个持有全部 handler 的 app
	a := newApp(gdb)

	// ④ 注册路由
	r := setupRouter(a)

	// ⑤ 起服务并阻塞在这里，直到收到退出信号、完成优雅关闭
	if err := run(cfg.Port, r); err != nil {
		log.Fatalf("服务退出: %v", err)
	}
}

// initDB 建数据库连接。
//
// 为什么用 log.Fatalf（fail fast）而不是把 error 返回给 main？
//
//	连不上数据库 = 服务根本没法工作，属于「启动期不可恢复错误」。
//	这种情况应该【立刻退出并打印明确原因】，而不是带着坏掉的依赖继续跑。
//
//	对比：请求处理中的错误是【可恢复】的，必须返回 error 给上层处理，
//	绝不能 Fatal（那会把整个进程干掉）。
//	—— 这正是面试题「什么时候可以用 panic/Fatal」的答案。
func initDB() *gorm.DB {
	g, err := db.New()
	if err != nil {
		log.Fatalf("初始化DB失败: %v", err)
	}

	return g
}

// run 起 HTTP 服务 + 优雅关闭，收到 SIGINT/SIGTERM 后返回。
//
// 这一段是 P0-⑤⑥⑦ 的综合实践，对应 docs/main-读码作业.md 的 6 个必答问题：
//
//	Q1 为什么要 go func() 包起来？
//	Q2 quit 的缓冲为什么是 1？
//	Q3 为什么要排除 http.ErrServerClosed？
//	Q4 Shutdown 和 Close 的区别？
//	Q5 那 5 秒是什么？defer cancel() 为什么必须有？
//	Q6 反复启停会泄漏 goroutine 吗？
//
// （答案分散在下面的注释里，学完 ⑥⑦ 回来把它们串起来）
func run(port string, r *gin.Engine) error {
	srv := &http.Server{Addr: ":" + port, Handler: r}

	// Q1：ListenAndServe 是【阻塞】调用（它内部是个 for 循环在 accept 连接），
	//     所以放进 goroutine；主 goroutine 才能继续往下走去等退出信号。
	go func() {
		log.Printf("server listening on :%s", port)

		// Q3：Shutdown 被调用后，ListenAndServe 会立刻返回 http.ErrServerClosed。
		//     这是"正常关闭"的信号，不是故障 —— 所以要用 errors.Is 排除掉，
		//     否则每次优雅关闭都会走到 log.Fatalf。
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v", err)
		}
	}()

	// Q2：缓冲为什么是 1 —— signal.Notify 往 channel 发信号时【不会阻塞等待接收】，
	//     如果 channel 无缓冲、此刻又没人接收，这个信号就丢了。
	//     缓冲 1 保证"至少能存下一个信号"，一定能被后面 <-quit 取到。
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // Ctrl+C / kill
	<-quit                                               // Q：主 goroutine 阻塞在这里等信号

	// Q4+Q5：
	//   · 用 Shutdown 而不是 Close —— Shutdown 会【等正在处理的请求处理完】
	//     （滚动发布/重启时不掐断用户请求）；Close 是直接断开所有连接。
	//   · 那 5 秒是"最多等多久"：5 秒还没处理完就强制关闭，避免进程永远卡住。
	//   · defer cancel() 即使超时了也必须调用：它会释放 ctx 关联的资源、
	//     停止内部的计时器 goroutine（不调会泄漏 —— P0-⑦ 的考点）。
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server forced to shutdown: %v", err)
	}
	log.Println("server exited")

	// Q6：不会泄漏 —— <-quit 只会走一次；Shutdown 后 ListenAndServe 返回、
	//     那个 goroutine 随之结束；cancel() 也把计时器收掉了。
	return nil
}
