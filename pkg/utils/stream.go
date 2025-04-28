package utils

import (
	"bufio"
	"os"

	"go.uber.org/zap"
)

// 将 Stream 改为支持泛型
type Stream[T any] <-chan T

// 修改 NewStream 构造函数
func NewStream[T any](data ...T) Stream[T] {
	ch := make(chan T)
	go func() {
		defer close(ch)
		for _, v := range data {
			ch <- v
		}
	}()
	return ch
}

func (s Stream[T]) ToSlice() []T {
	out := make([]T, 0)
	for v := range s {
		out = append(out, v)
	}
	return out
}

func Map[T, R any](s Stream[T], fn func(T) R) Stream[R] {
	out := make(chan R)
	go func() {
		defer close(out)
		for v := range s {
			out <- fn(v)
		}
	}()
	return out
}

func (s Stream[T]) Filter(fn func(T) bool) Stream[T] {
	out := make(chan T)
	go func() {
		defer close(out)
		for v := range s {
			if fn(v) {
				out <- v
			}
		}
	}()
	return out
}

func ReadLargeFile(filename string) Stream[string] {
	ch := make(chan string)
	go func() {
		file, err := os.Open(filename)
		if err != nil {
			close(ch)
			zap.L().Error("cannot open file", zap.Error(err))
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			ch <- scanner.Text()
		}
		close(ch)
	}()

	return ch
}
