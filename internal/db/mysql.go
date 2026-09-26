package db

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// New 创建并返回一个 GORM 数据库连接实例
// 返回值:
//   - *gorm.DB: GORM 数据库操作对象，用于后续的增删改查操作
//   - error: 如果连接或初始化过程中出现错误，返回相应的错误信息
func New() (*gorm.DB, error) {
	// 构建 MySQL 的 DSN (Data Source Name) 连接字符串
	// 格式: username:password@tcp(host:port)/dbname?charset=xxx&parseTime=True&loc=Local
	// 使用 getEnv 函数从环境变量中读取配置，若未设置则使用默认值
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		getEnv("DB_USER", "root"),         // 数据库用户名，默认 root
		getEnv("DB_PASSWORD", "root"),     // 数据库密码，默认 root
		getEnv("DB_HOST", "127.0.0.1"),    // 数据库主机地址，默认本地
		getEnv("DB_PORT", "3306"),         // 数据库端口，默认 3306
		getEnv("DB_DATABASE", "order_db"), // 数据库名称，默认 order_db
	)

	// 使用 GORM 打开 MySQL 连接
	// mysql.Open(dsn): 使用 MySQL 驱动解析 DSN
	// gorm.Config 中配置 NamingStrategy:
	//   SingularTable: true 表示表名使用单数形式（默认 GORM 会使用复数形式）
	//   例如: 结构体 User 默认映射到表 users，设置为 true 后映射到表 user
	gdb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true}, // 将表名设置为单数
	})

	// 如果打开数据库连接失败，返回带有上下文信息的错误
	// 使用 %w 包装原始错误，便于上层使用 errors.Is / errors.As 进行错误判断
	if err != nil {
		return nil, fmt.Errorf("mysql连接失败: %w", err)
	}

	// 从 GORM 实例中获取底层的 *sql.DB 对象
	// 用于配置连接池的相关参数
	sqlDB, err := gdb.DB()

	// 如果获取底层 sql.DB 失败，返回错误
	if err != nil {
		return nil, fmt.Errorf("获取sql.DB失败: %w", err)
	}

	// 设置数据库的最大打开连接数（包括正在使用和空闲的连接）
	// 这里设置为 10，表示同一时刻最多有 10 个连接
	// 注意: 该值通常应 >= SetMaxIdleConns，否则空闲连接数会被限制
	sqlDB.SetMaxOpenConns(25)

	// 设置连接池中最大空闲连接数
	// 空闲连接是指当前没有被使用但保留在池中的连接
	// 这里设置为 25，表示最多保留 25 个空闲连接备用
	sqlDB.SetMaxIdleConns(10)

	// 设置连接的最大存活时间
	// 超过该时间的连接会被关闭并从连接池中移除
	// 这里设置为 1 小时，有助于避免使用过期连接（如 MySQL 的 wait_timeout 限制）
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 返回初始化完成的 GORM 数据库实例
	return gdb, nil
}

// 获取配置信息
func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
