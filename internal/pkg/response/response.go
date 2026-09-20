package response

import "github.com/gin-gonic/gin"

// Body 统一响应结构。
// code=0 表示成功，与 HTTP status 解耦，便于前端区分网络错误与业务错误。
type Body struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

// OK 成功响应
func OK(c *gin.Context, data any) {
	c.JSON(200, Body{Code: 0, Msg: "ok", Data: data})
}

// Fail 失败响应：httpStatus 给网络层看，code 给业务层看
func Fail(c *gin.Context, httpStatus int, code int, msg string) {
	c.JSON(httpStatus, Body{Code: code, Msg: msg})
}
