package utils

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewStream(t *testing.T) {
	// 测试空数据
	stream := NewStream[int]()
	count := 0
	for range stream {
		count++
	}
	assert.Equal(t, 0, count)

	stream = NewStream(1, 2, 3)
	values := make([]int, 0)
	for v := range stream {
		values = append(values, v)
	}
	assert.Equal(t, []int{1, 2, 3}, values)
}

func TestStream_Map(t *testing.T) {
	stream := NewStream[int](1, 2, 3)
	doubled := Map[int, int](stream, func(v int) int {
		return v * 2
	})

	results := make([]int, 0)
	for v := range doubled {
		results = append(results, v)
	}
	assert.Equal(t, []int{2, 4, 6}, results)
}

func TestStream_Filter(t *testing.T) {
	stream := NewStream[int](1, 2, 3, 4, 5)
	evenNums := stream.Filter(func(v int) bool {
		return v%2 == 0
	})

	results := make([]int, 0)
	for v := range evenNums {
		results = append(results, v)
	}
	assert.Equal(t, []int{2, 4}, results)
}

func TestReadLargeFile(t *testing.T) {
	// 创建临时测试文件
	content := []byte("line1\nline2\nline3")
	tmpfile, err := os.CreateTemp("", "test*.txt")
	assert.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.Write(content)
	assert.NoError(t, err)
	tmpfile.Close()

	// 测试读取文件
	stream := ReadLargeFile(tmpfile.Name())
	lines := make([]string, 0)

	// 使用 select 和 timeout 来避免可能的死锁
	timeout := time.After(2 * time.Second)
	for {
		select {
		case line, ok := <-stream:
			if !ok {
				assert.Equal(t, []string{"line1", "line2", "line3"}, lines)
				return
			}
			lines = append(lines, line)
		case <-timeout:
			t.Fatal("test timed out")
			return
		}
	}
}
