package repository

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/gcr0827/order-service/internal/db"
	"github.com/gcr0827/order-service/internal/model"
)

func TestDeductStock_Concurrent(t *testing.T) {
	gdb, err := db.New()
	if err != nil {
		t.Skipf("跳过继承测试，数据库不可用（%v）--先执行set-a;source .env;set +a", err)
	}

	const (
		testStoreID = 999999
		testSkuID   = 999999
		initStock   = 10
		goroutines  = 100
	)

	// ------------ 准备：清理旧测试数据 + 插入一条一致库存 --------------
	// 必须做Unscoped（）做【物理删除】，因为模型带gorm.DeletedAt
	// 普通 Delete 是软删除 → 行还在 → 唯一键 (store_id, sku_id) 仍被占用 → 后面 Create 会报 Duplicate

	gdb.Unscoped().
		Where("store_id = ? and sku_id = ?", testStoreID, testSkuID).
		Delete(&model.StoreProductInventory{})

	if err = gdb.Create(&model.StoreProductInventory{
		StoreID:        testStoreID,
		SkuID:          testSkuID,
		SpuID:          1,
		TotalStock:     initStock,
		AvailableStock: initStock,
	}).Error; err != nil {
		t.Fatalf("准备测试数据失败：%v", err)
	}

	t.Cleanup(func() {
		gdb.Unscoped().Where("store_id = ? AND sku_id = ?", testStoreID, testSkuID).
			Delete(&model.StoreProductInventory{})
	})

	repo := NewStoreProductInventoryRepository(gdb)
	ctx := context.Background()

	// ---- 并发扣减：100 个 goroutine 各扣 1 ----
	var success, failed atomic.Int64
	var wg sync.WaitGroup

	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()

			if err := repo.DeductStock(ctx, testStoreID, testSkuID, 1); err != nil {
				failed.Add(1)
			} else {
				success.Add(1)
			}
		}()
	}
	wg.Wait()

	// ---- 断言 ----
	t.Logf("成功 = %d 次，失败 = %d 次（期望成功正好 %d）", success.Load(), failed.Load(), initStock)

	if got := success.Load(); got != initStock {
		t.Errorf("❌ 超卖/漏卖！success = %d，期望 %d", got, initStock)
	}

	var final model.StoreProductInventory
	if err := gdb.Where("store_id = ? AND sku_id = ?", testStoreID, testSkuID).
		First(&final).Error; err != nil {
		t.Fatalf("查询最终库存失败: %v", err)
	}
	t.Logf("最终 available_stock = %d（期望 0）", final.AvailableStock)

	if final.AvailableStock != 0 {
		t.Errorf("❌ 最终库存 = %d，期望 0", final.AvailableStock)
	}
}
