package tests

import (
	"gqlexample/pkg/task"
	"testing"
	"time"
)

func TestDelayedTask(t *testing.T) {
	var count int
	task.NewDelayedTask(2*time.Second, func() {
		count++
	})

	time.Sleep(3 * time.Second)
	if count != 1 {
		t.Errorf("Expected count to be 1, but got %d", count)
	}
}
