package store

import (
	"context"

	"bigmarket/internal/pkg"
)

// FlushCaches 运营端修改配置后清理策略、活动、规则树缓存
func (s *Store) FlushCaches() {
	ctx := context.Background()
	patterns := []string{
		pkg.RateTableKey + "*", pkg.RateRangeKey + "*", pkg.RuleTreeVOKey + "*",
		pkg.ActivityKey + "*", pkg.SkuKey + "*", pkg.AwardCountKey + "*", pkg.SkuStockKey + "*",
	}
	for _, pattern := range patterns {
		iter := s.RDB.Scan(ctx, 0, pattern, 100).Iterator()
		for iter.Next(ctx) {
			s.RDB.Del(ctx, iter.Val())
		}
	}
}
