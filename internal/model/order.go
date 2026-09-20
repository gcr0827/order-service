package model

// Order 订单模型。
// 约定：model 只承载数据结构，不做校验、不拼接 message。
type Order struct {
	ID       int64  `json:"id"`
	UserID   int64  `json:"user_id"`
	Amount   int64  `json:"amount"` // 单位：分
	Status   int8   `json:"status"` // 1=待支付 2=已支付 3=已取消
	CreateAt string `json:"create_at"`
}
