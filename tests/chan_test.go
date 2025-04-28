package tests

import (
	"context"
	"testing"
)

var eventChan chan struct{}

func init() {
	eventChan = make(chan struct{})
}

func calculate(ch chan<- int, ctx context.Context) {
	defer close(ch)
	var res = 1
	calculator := func() {
		res++
	}

	calculator()
	ch <- res

	for {
		select {
		case <-ctx.Done():
			return
		case <-eventChan:
			calculator()
			select {
			case ch <- res:
			case <-ctx.Done(): // 防止写入阻塞
				return
			}
		default:
		}
	}
}

func TestChan(t *testing.T) {
	ch := make(chan int, 1)
	ctx := context.Background()
	go calculate(ch, ctx)

	go func() {
		for range 10 {
			eventChan <- struct{}{}
		}
	}()
	for v := range ch {
		t.Log(v)
	}
}
