package main

import (
	"fmt"
	"log"

	"github.com/gcr0827/order-service/internal/db"
	"github.com/gcr0827/order-service/internal/model"
)

func main() {

	// db.New() 通常是你自己封装的函数：
	// 1. 读取配置（DSN、用户名、密码、地址、库名等）
	// 2. 调用 gorm.Open(...) 建立数据库连接
	// 3. 可能还做了连接池设置、日志设置、AutoMigrate 等
	// 返回的 g 是 *gorm.DB，后续所有数据库操作都通过它完成
	g, err := db.New()
	if err != nil {
		// 初始化失败属于致命错误，直接退出程序
		// log.Fatalf 等价于：先打印日志，再调用 os.Exit(1)
		log.Fatalf("初始化DB失败:%v", err)
	}

	// 构造一个 ProductSpu 结构体实例（指针类型）
	// 这里只设置了部分字段，说明其他字段（如 ID、CreatedAt、UpdatedAt）
	// 要么由数据库自增/默认值生成，要么由 GORM 的约定自动填充
	spu := &model.StoreProductSpu{
		CategoryID: 1,             // 分类 ID
		Name:       "test-1",      // SPU 名称
		Code:       "test-code-1", // SPU 编码（通常要求唯一）
		Status:     0,             // 状态：0 可能表示“禁用/下架”，具体看业务定义
		Sort:       100,           // 排序值，越大越靠前或越靠后取决于业务
	}

	// 新增
	// g.Create(spu) 会把 spu 这条记录 INSERT 到对应的表中
	// 因为传入的是指针 &spu，GORM 会把数据库生成的自增主键回填到 spu.ID
	// 如果表里有 CreatedAt / UpdatedAt 字段，GORM 也会自动填充当前时间
	if err := g.Create(spu).Error; err != nil {
		// 插入失败：可能是唯一索引冲突、字段超长、数据库连接断开等
		log.Fatalf("插入失败：%v", err)
	}

	// 打印插入后的自增 ID
	// 注意：这里能拿到 ID，是因为 Create 时传入了指针，GORM 做了回填
	// 如果传的是值（spu 而不是 &spu），ID 可能不会被写回原变量
	fmt.Printf("插入成功，ID=%d", spu.ID) //// GORM 会把自增 ID 回填

	// 查询
	// 声明一个零值结构体 got，用来接收查询结果
	var got model.StoreProductSpu

	// g.First(&got, spu.ID) 等价于：
	//   SELECT * FROM product_spu WHERE id = ? ORDER BY id LIMIT 1
	// 这里用刚才插入得到的 spu.ID 作为主键条件
	// First 的语义是：按主键升序取第一条，找不到会返回 ErrRecordNotFound
	if err := g.First(&got, spu.ID).Error; err != nil {
		log.Fatalf("查询失败: %v", err)
	}

	// 打印查询结果
	// CreatedAt 是 time.Time 类型，用 Format 转成可读字符串
	// 格式串 "2006-01-02 15:04:05" 是 Go 的固定参考时间，不能随便改
	fmt.Printf("✅ 查到: ID=%d Name=%s Code=%s Status=%d CreatedAt=%s\n",
		got.ID, got.Name, got.Code, got.Status, got.CreatedAt.Format("2006-01-02 15:04:05"))
}
