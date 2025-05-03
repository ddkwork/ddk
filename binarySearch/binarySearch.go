package main

import (
	"fmt"
	"strings"
)

func main() {
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

	fmt.Printf("Matches found at offsets: %v\n", matches)
}

// Pattern struct definition
type Pattern struct {
	Bytes    []byte // 特征码字节（已处理通配符）
	Masks    []byte // 通配符掩码（0xFF 表示全通配）
	FirstIdx int    // 第一个非通配符位置
}

// ParsePattern function
func ParsePattern(pattern string) *Pattern {
	clean := strings.ReplaceAll(pattern, " ", "")
	length := len(clean) / 2
	p := &Pattern{
		Bytes:    make([]byte, length),
		Masks:    make([]byte, length),
		FirstIdx: length,
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

		if mask != 0xFF && p.FirstIdx == length {
			p.FirstIdx = i
		}

		p.Bytes[i] = value
		p.Masks[i] = mask

		// Debug statements
		//fmt.Printf("Chunk: %s, Mask: %02X, Value: %02X\n", chunk, mask, value)
	}

	//fmt.Printf("Parsed Pattern: Bytes=%v, Masks=%v, FirstIndex=%d\n", p.Bytes, p.Masks, p.FirstIndex)
	return p
}

// parseHex function
func parseHex(c byte) byte {
	if c >= '0' && c <= '9' {
		return c - '0'
	}
	return c - 'A' + 10
}

// SearchMemoryChunked function
func SearchMemoryChunked(mem []byte, pattern *Pattern, chunkSize int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for i := 0; i < len(mem); i += chunkSize {
			end := min(i+chunkSize, len(mem))
			chunk := mem[i:end]
			for offset := range searchChunk(chunk, pattern) {
				globalOffset := i + offset
				ch <- globalOffset
			}
		}
	}()
	return ch
}

// searchChunk function
func searchChunk(chunk []byte, pattern *Pattern) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		length := len(pattern.Bytes)
		for i := 0; i <= len(chunk)-length; i++ {
			// 检查第一个非通配符位置
			firstCheck := i + pattern.FirstIdx
			if firstCheck >= len(chunk) {
				break
			}

			masked := chunk[firstCheck] &^ pattern.Masks[pattern.FirstIdx]
			if masked != pattern.Bytes[pattern.FirstIdx] {
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
				// Debug statement
				//fmt.Printf("Match found at chunk offset: %d\n", i)
				ch <- i
			}
		}
	}()
	return ch
}
