package apperr

import "errors"

// 业务错误码（与 HTTP status 分离）
const (
	CodeInvalidParam  = 40001
	CodeNotFound      = 40401
	CodeConflict      = 40901
	CodeInternalError = 50001
)

// AppError 携带业务错误码与原始错误。
// 约定：repository / service 只返回 error，不返回 message 字符串；
// 面向用户的文案由 handler 层统一翻译。
// 这对应 Laravel 规范里的「Model 返回纯数据」。
type AppError struct {
	Code int
	Err  error
}

func (e *AppError) Error() string { return e.Err.Error() }
func (e *AppError) Unwrap() error { return e.Err }

// New 用文案构造业务错误
func New(code int, msg string) *AppError {
	return &AppError{Code: code, Err: errors.New(msg)}
}

// Wrap 包装已有错误并附上业务错误码
func Wrap(code int, err error) *AppError {
	return &AppError{Code: code, Err: err}
}

// 哨兵错误：用 errors.Is 判断业务条件
var (
	ErrNotFound = errors.New("order not found")
)
