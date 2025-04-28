package task

import (
	"sync/atomic"
	"time"
)

type DelayedTask struct {
	timer    *time.Timer
	cancelCh chan struct{}
	done     uint32
}

func NewDelayedTask(delay time.Duration, f func()) *DelayedTask {
	task := &DelayedTask{
		cancelCh: make(chan struct{}),
	}

	task.timer = time.AfterFunc(delay, func() {
		if atomic.CompareAndSwapUint32(&task.done, 0, 1) {
			select {
			case <-task.cancelCh:
				// 任务被取消
			default:
				f()
			}
			// 无论是否执行，都关闭通道避免泄漏
			close(task.cancelCh)
		}
	})

	return task
}

func (t *DelayedTask) Cancel() bool {
	// 先尝试停止定时器
	stopped := t.timer.Stop()
	// 如果停止成功或任务尚未执行，则尝试取消
	if stopped || atomic.LoadUint32(&t.done) == 0 {
		// 使用原子操作确保只取消一次
		if atomic.CompareAndSwapUint32(&t.done, 0, 1) {
			close(t.cancelCh)
			return true
		}
	}
	return false
}
