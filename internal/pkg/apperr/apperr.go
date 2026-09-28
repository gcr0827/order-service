package apperr

import "errors"

// ============================================================
// 业务错误码
//
// 为什么 code 要和 HTTP status 解耦？
//
//	· code        给【业务】判断用（前端据它决定提示什么、跳哪里）
//	· HTTP status 给【网络层】用（网关 / 负载均衡 / 监控按它统计）
//	· 两者是多对一：同一个 404 可以有多个业务码（商品不存在 / 订单不存在）
//	· 好处：加新业务错误时，HTTP 语义不用变
//
// ============================================================
const (
	CodeInvalidParam  = 40001 // 参数不合法
	CodeNotFound      = 40401 // 资源不存在
	CodeConflict      = 40901 // 业务冲突（库存不足、订单已支付…）
	CodeInternalError = 50001 // 系统内部错误（细节不可暴露给用户）
)

// AppError 携带【业务码】+【原始错误】。
//
// 两条约定：
//
//	① repository / service 只返回 error，【不返回"给用户看的文案"】
//	   —— 文案取决于调用方是谁：HTTP 要给用户看中文、定时任务要写日志、
//	      内部 RPC 要错误码、单测要断言。只有 handler 知道"这是给终端用户的 HTTP 响应"。
//	② 文案由 handler 统一翻译（用 code → 文案 的映射表）
//	   —— 这对应 Laravel 规范里的「Model 返回纯数据」。
type AppError struct {
	Code int
	Err  error
}

func (e *AppError) Error() string { return e.Err.Error() }

// Unwrap 让 errors.Is / errors.As 能【穿过】AppError 继续往下找根因。
//
// ⚠️ 这个方法不能省：没有它，AppError 就是链的终点 ——
//
//	errors.Is(err, apperr.ErrNotFound) 会永远返回 false
//	→ handler 判定不出业务错误 → 业务错误全被当成 500
//	效果等同于用 fmt.Errorf("%v") 包装（断链）。
func (e *AppError) Unwrap() error { return e.Err }

// Wrap 给【已有错误】加上业务码（保留原始错误链，Is/As 仍能穿透）。
// 用在 service 层需要给底层错误补一个业务语义的时候。
func Wrap(code int, err error) *AppError {
	return &AppError{Code: code, Err: err}
}

// ── 预定义业务错误 ────────────────────────────────────────────
//
// ★ 关键设计：它们【本身就是 *AppError】，而不是裸的 errors.New(...)
//
// 这样一个值同时满足两种用法：
//
//	① 当【哨兵】用  → errors.Is(err, apperr.ErrNotFound)     ✅ 能按值匹配
//	② 取【业务码】  → errors.As(err, &appErr) → appErr.Code  ✅ 能拿到 Code
//
// 如果定义成裸哨兵（errors.New），就只能做 ①、做不到 ② ——
// 那 Code 常量和 AppError 类型就成了【没人产生的死代码】。
//
// 为什么没有 New(code, msg) 这个构造函数？
//
//	带 msg 的构造会把【用户文案】塞进错误里，和上面约定 ① 矛盾。
//	需要新的业务错误时，在下面加一个预定义变量即可（还能被 Is 匹配）。
var (
	ErrNotFound          = &AppError{Code: CodeNotFound, Err: errors.New("not found")}
	ErrInsufficientStock = &AppError{Code: CodeConflict, Err: errors.New("insufficient stock")}
)
