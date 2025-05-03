package main

import "testing"

// BenchmarkName-8   	  308677	      3719 ns/op
// BenchmarkName-8   	  281802	      3749 ns/op
// BenchmarkName-8   	  269437	      3809 ns/op
// BenchmarkName-8   	  314743	      3585 ns/op
func BenchmarkName(b *testing.B) {
	for b.Loop() {
		memData := []byte{
			0x33, 0xC9, 0x89, 0x0D, 0xB4, 0x67, 0x92, 0x77, 0x89, 0x0D,
			0xB8, 0x67, 0x92, 0x77, 0x88, 0x08, 0x38, 0x48, 0x02, 0x74,
			0x05, 0xE8, 0x94, 0xFF, 0xFF, 0xFF, 0x33, 0xC0, 0xC3, 0x8B,
			0xFF, 0x55, 0x8B, 0xEC, 0x83, 0xE4, 0xF8,
		}

		pattern := ParsePattern("?9 ?? 0? ?? 67") // 特征码解析

		// 使用分块迭代器搜索内存（块大小=16）
		iter := SearchMemoryChunked(memData, pattern, 16)

		// 收集匹配结果
		var matches []int
		for offset := range iter {
			matches = append(matches, offset)
		}
	}
}
