package service

import (
	"context"
	"errors"
	"testing"

	"github.com/gcr0827/order-service/internal/model"
	"github.com/gcr0827/order-service/internal/pkg/apperr"
)

// ============================================================
// ① 假实现（mock）：想让它返回什么，就返回什么
// ============================================================

type fakeProductRepo struct {
	spu *model.StoreProductSpu
	err error
}

func (f *fakeProductRepo) GetSpuByID(ctx context.Context, id uint64) (*model.StoreProductSpu, error) {
	return f.spu, f.err
}

func (f *fakeProductRepo) GetSkuByID(ctx context.Context, id uint64) (*model.StoreProductSku, error) {
	return nil, nil
}

func (f *fakeProductRepo) ListSkuBySpuID(ctx context.Context, storeID, spuID uint64) ([]*model.StoreProductSku, error) {
	return nil, nil
}

// 编译期断言：万一 fake 没实现接口，编译就报错（比运行时才发现好）
var _ ProductRepository = (*fakeProductRepo)(nil)

// ============================================================
// ② 表驱动测试：三个 case 写在表里
// ============================================================

func TestProductService_GetSpu(t *testing.T) {
	wantSpu := &model.StoreProductSpu{ID: 1, Name: "测试商品", Code: "P001"}

	tests := []struct {
		name         string
		repoSpu      *model.StoreProductSpu
		repoErr      error
		wantErr      bool // 期望有没有错误
		wantNotFound bool // 期望错误里"是不是 ErrNotFound"
	}{
		{
			name:    "正常找到",
			repoSpu: wantSpu,
			repoErr: nil,
			wantErr: false,
		},
		{
			name:         "记录不存在 → 透传 ErrNotFound",
			repoSpu:      nil,
			repoErr:      apperr.ErrNotFound,
			wantErr:      true,
			wantNotFound: true,
		},
		{
			name:         "数据库故障 → 不能被当成 ErrNotFound",
			repoSpu:      nil,
			repoErr:      errors.New("db down"),
			wantErr:      true,
			wantNotFound: false, // ⭐ 这条锁住了错误分层
		},
	}

	// ③ 循环跑每个 case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 用假实现装配 service（不需要数据库！）
			svc := NewProductService(&fakeProductRepo{spu: tt.repoSpu, err: tt.repoErr})

			got, err := svc.GetSpu(context.Background(), 1)

			// 断言 1：有没有错误
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr = %v", err, tt.wantErr)
			}

			// 断言 2：错误的语义（是不是 ErrNotFound）
			if isNotFound := errors.Is(err, apperr.ErrNotFound); isNotFound != tt.wantNotFound {
				t.Errorf("errors.Is(err, ErrNotFound) = %v, want %v (err=%v)",
					isNotFound, tt.wantNotFound, err)
			}

			// 断言 3：数据对不对（只在期望有数据时检查）
			if tt.repoSpu != nil {
				if got == nil || got.ID != tt.repoSpu.ID || got.Name != tt.repoSpu.Name {
					t.Errorf("got = %+v, want = %+v", got, tt.repoSpu)
				}
			}
		})
	}
}
