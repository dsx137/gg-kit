package lang_test

import (
	"testing"

	"github.com/dsx137/gg-kit/internal/lang"
)

// 添加这个基准测试来验证
func BenchmarkGetGoroutineId(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = lang.GetGoroutineId()
	}
}
