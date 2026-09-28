package handler

import (
	"errors"
	"net/http"

	"github.com/gcr0827/order-service/internal/pkg/apperr"
	"github.com/gcr0827/order-service/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

// codeMeta 业务码 → (HTTP 状态码, 给用户看的文案)
//
// 集中维护：以后加业务错误，只在这里加一行 —— handler 的流程代码一行不用改。
// 这就是「业务码与 HTTP status 解耦」的落地方式。
// （需要更精确的文案时，再加细分码，如 CodeProductNotFound / CodeOrderNotFound —— 同一个 404 可以对应多个业务码）
var codeMeta = map[int]struct {
	HTTP int
	Msg  string
}{
	apperr.CodeInvalidParam:  {http.StatusBadRequest, "参数不合法"},
	apperr.CodeNotFound:      {http.StatusNotFound, "数据不存在"},
	apperr.CodeConflict:      {http.StatusConflict, "业务冲突"},
	apperr.CodeInternalError: {http.StatusInternalServerError, "服务器内部错误"},
}

// fail 把所有 error 统一翻译成 HTTP 响应 —— 三个 handler 共用，避免重复写 switch。
//
// 两个分支：
//
//	① 是 *apperr.AppError 且码在表里 → 按业务码返回对应状态码 + 文案
//	② 其他（系统错误 / 码未登记）    → 500 + 通用文案（内部细节不外泄）
//
// ⚠️ 为什么用 errors.As 而不是 errors.Is？
//
//	Is 只能回答"是不是它"（拿不到 Code）；要读 Code 必须用 As 把类型取出来。
func fail(c *gin.Context, err error) {
	var appErr *apperr.AppError
	if errors.As(err, &appErr) {
		if meta, ok := codeMeta[appErr.Code]; ok {
			response.Fail(c, meta.HTTP, appErr.Code, meta.Msg)
			return
		}
	}
	// 兜底：不是业务错误，或码没登记 → 系统错误
	response.Fail(c, http.StatusInternalServerError, apperr.CodeInternalError, "服务内部错误")
}
