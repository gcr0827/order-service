# main.go 读码作业（W2 完成）

> **为什么有这份文档**：项目骨架是**先写好的**，里面用到了一些你还没学的语言特性
> （goroutine / channel / context / defer）。这是刻意的顺序——**先有可运行的骨架，再逐块学透**。
>
> ⚠️ **但前提是：每一块都要挂上学习时间点。** 否则就变成"自己的代码自己讲不清"——
> 而这个仓库是公开的，面试官**可以直接打开 main.go 问**。
>
> 本文档 = 那些代码的**学习对照表 + 必答问题**。
> 学完 P0-⑤⑥⑦ 后，回来把「能答吗」这列全部打勾。

---

## 1. 骨架代码 ↔ P0 知识点对照

| 代码位置 | 用到的东西 | 对应 P0 | 学的时间 | 现在能讲吗 |
| :--- | :--- | :--- | :--- | :--- |
| `internal/repository/order_repository.go` | 接口 + 实现分离 | **③ interface** | **9/22（今天）** | ⬜ |
| `internal/service/` 依赖接口而非实现 | 面向接口（依赖注入） | **③ interface** | 9/22 | ⬜ |
| `internal/pkg/apperr/` | 错误包装 / 哨兵错误 / `Unwrap` | **④ error** | **9/23** | ⬜ |
| `internal/handler/` 里 `errors.Is` 分支 | `errors.Is` 与错误分层 | **④ error** | 9/23 | ⬜ |
| `main.go` 的 `defer cancel()` | defer 与函数返回的交互 | **⑤ defer** | **9/24** | ⬜ |
| `main.go` 的 `go func(){ srv.ListenAndServe() }()` | goroutine 生命周期 | **⑥ 并发** | **W2** | ⬜ |
| `main.go` 的 `quit := make(chan os.Signal, 1)` | 带缓冲 channel | **⑥ 并发** | W2 | ⬜ |
| `main.go` 的 `context.WithTimeout` / `srv.Shutdown(ctx)` | context 超时与取消传播 | **⑦ context** | **W2** | ⬜ |
| `r.Use(gin.Logger(), gin.Recovery())` | Gin 中间件 | P1（跟着项目学） | W2 | ⬜ |

**结论**：你现在讲不了的只有 **main.go 那 4 处**，而它们全部落在 **9/24 + W2**——
**时间够**：你的首批面试在 **10/1 之后**（9/25–9/30 投 → 筛选 3–7 天），
那时 P0-⑤⑥⑦ 已经学完了。

---

## 2. 必答 6 问（**W2 学完 P0-⑤⑥⑦ 后回来答**）

> **纪律**：先自己想，答不出再去查。**这 6 问就是"能不能讲自己的代码"的标准。**

### Q1 为什么服务要跑在 goroutine 里？

```go
go func() {
    if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
        log.Fatalf("listen: %v", err)
    }
}()
```

- `ListenAndServe` 是**阻塞**的，那**主 goroutine 在干什么**？
- 如果把 `ListenAndServe` 直接写在 main 里（不套 goroutine），会怎样？

### Q2 为什么 `quit` channel 的缓冲是 **1**？

```go
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
```

- `signal.Notify` 往 channel 发信号时**会不会阻塞等待接收**？
- 如果用 `make(chan os.Signal)`（无缓冲），可能出什么问题？

### Q3 为什么要排除 `http.ErrServerClosed`？

- 什么情况下 `ListenAndServe` 会返回这个"错误"？
- 它是**故障**吗？为什么不能当成故障处理（`log.Fatalf`）？

### Q4 `srv.Shutdown(ctx)` 和 `srv.Close()` 的区别？

- 哪个会**等正在处理的请求做完**？哪个**直接掐断**？
- 生产环境该用哪个？为什么？

### Q5 那 5 秒是什么？为什么必须有 `defer cancel()`？

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
if err := srv.Shutdown(ctx); err != nil { ... }
```

- 5 秒到了还没关完，会发生什么？
- `cancel()` 的作用是什么？**即使超时了也要调用吗？**（提示：资源释放 / goroutine 泄漏）

### Q6 自己加一条：这段代码如果**反复启停**，会不会泄漏 goroutine？

- 怎么验证？（提示：`runtime.NumGoroutine()` 或 `go test -race` + 压力测试）

---

## 3. 学会之前的「诚实回复」话术

如果**在 W2 学完之前**就被问到 main.go（可能性不高，但要有准备）：

> 「优雅关闭这块是我搭骨架时按标准写法先落下来的，
> 涉及 goroutine、channel 和 context——我目前的用法是对的，
> 但**细节我还在系统补**（正在按顺序过 defer / 并发 / context），
> 学完我会把这几处的取舍完整讲一遍。」

- ✅ **理直气壮的诚实**：说明你在按计划补，而不是"不懂装懂"
- ❌ 不要硬编「为了资源回收和优雅退出的最佳实践」这种复述式答案——
  面试官下一句就是「那 `Shutdown` 内部怎么实现的？为什么需要 context？」

---

## 4. 下次打开这个文件时

**先看这 6 问能不能答**，答不上就翻 `面试八股/06-并发原语.md` / `07-context.md`（W2 产出）。
