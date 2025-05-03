package main

import (
	"iter"
	"strings"

	"github.com/ddkwork/golibrary/mylog"
)

func main() {
	memData := []byte{
		0x33,
		0xC9, 0x89, 0x0D, 0xB4, 0x67, // find this pattern in memory
		0x92,
		0x77, 0x89, 0x0D, 0xB8, 0x67, // 这个不符合要求
		0x92, 0x77, 0x88, 0x08, 0x38, 0x48, 0x02,
		0x74, 0x05, 0xE8, 0x94, 0xFF, 0xFF, 0xFF,
		0x33, 0xC0, 0xC3, 0x8B, 0xFF, 0x55, 0x8B,
		0xEC, 0x83, 0xE4, 0xF8,
	}
	pattern := ParsePattern("?9 ?? 0? ?? 67") // 特征码解析,会处理成：F9 FF 0F FF 67 进行匹配

	for offset := range SearchMemoryChunked(memData, pattern, 16) {
		b := memData[offset : offset+len(pattern.Bytes)]
		mylog.HexDump(offset, b)
	}
}

type Pattern struct {
	Bytes      []byte // 特征码字节（已处理通配符）
	Masks      []byte // 通配符掩码（0xFF 表示全通配）
	FirstIndex int    // 第一个非通配符位置
}

func ParsePattern(pattern string) *Pattern {
	clean := strings.ReplaceAll(pattern, " ", "")
	length := len(clean) / 2
	p := &Pattern{
		Bytes:      make([]byte, length),
		Masks:      make([]byte, length),
		FirstIndex: length,
	}

	for i := range length {
		chunk := clean[i*2 : (i+1)*2]
		var mask, value byte

		switch {
		case chunk == "??":
			mask, value = 0xFF, 0x00
		case chunk[0] == '?':
			mask = 0xF0
			value = parseHex(chunk[1])
		case chunk[1] == '?':
			mask = 0x0F
			value = parseHex(chunk[0]) << 4
		default:
			mask = 0x00
			value = (parseHex(chunk[0]) << 4) | parseHex(chunk[1])
		}

		if mask != 0xFF && p.FirstIndex == length {
			p.FirstIndex = i
		}

		p.Bytes[i] = value
		p.Masks[i] = mask
	}

	return p
}

func parseHex(c byte) byte {
	if c >= '0' && c <= '9' {
		return c - '0'
	}
	return c - 'A' + 10
}

func SearchMemoryChunked(mem []byte, pattern *Pattern, chunkSize int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i, chunk := range Chunk(mem, chunkSize) {
			// 为每个块生成独立的迭代器
			for offset := range searchChunk(chunk, pattern) {
				// 计算全局偏移（当前块起始位置 + 块内偏移）
				globalOffset := i*chunkSize + offset
				if !yield(globalOffset) {
					return
				}
			}
		}
	}
}

// 单块内存搜索逻辑
func searchChunk(chunk []byte, pattern *Pattern) iter.Seq[int] {
	length := len(pattern.Bytes)
	return func(yield func(int) bool) {
		for i := 0; i <= len(chunk)-length; i++ {
			// 检查第一个非通配符位置
			firstCheck := i + pattern.FirstIndex
			if firstCheck >= len(chunk) {
				break
			}

			masked := chunk[firstCheck] &^ pattern.Masks[pattern.FirstIndex]
			if masked != pattern.Bytes[pattern.FirstIndex] {
				continue
			}

			// 匹配后续字节
			match := true
			for j := range length {
				pos := i + j
				if pos >= len(chunk) {
					match = false
					break
				}

				if pattern.Masks[j] != 0xFF && (chunk[pos]&^pattern.Masks[j]) != pattern.Bytes[j] {
					match = false
					break
				}
			}

			if match {
				if !yield(i) {
					return
				}
			}
		}
	}
}

func Chunk[Slice ~[]E, E any](s Slice, n int) iter.Seq2[int, Slice] {
	if n < 1 {
		panic("cannot be less than 1")
	}

	return func(yield func(int, Slice) bool) {
		for i := 0; i < len(s); i += n {
			// Clamp the last chunk to the slice bound as necessary.
			end := min(n, len(s[i:]))

			// Set the capacity of each chunk so that appending to a chunk does
			// not modify the original slice.
			if !yield(i, s[i:i+end:i+end]) {
				return
			}
		}
	}
}
