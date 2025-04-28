package timewheel

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestTimeWheel_New(t *testing.T) {
	tests := []struct {
		name            string
		intervalSeconds int
		slotNum         int
		job             Job
		wantNil         bool
	}{
		{
			name:            "valid parameters",
			intervalSeconds: 1,
			slotNum:         60,
			job:             func(any) {},
			wantNil:         false,
		},
		{
			name:            "invalid interval",
			intervalSeconds: -1,
			slotNum:         60,
			job:             func(any) {},
			wantNil:         true,
		},
		{
			name:            "invalid slot number",
			intervalSeconds: 1,
			slotNum:         0,
			job:             func(any) {},
			wantNil:         true,
		},
		{
			name:            "nil job",
			intervalSeconds: 1,
			slotNum:         60,
			job:             nil,
			wantNil:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tw := New(int64(tt.intervalSeconds), tt.slotNum, tt.job)
			if (tw == nil) != tt.wantNil {
				t.Errorf("New() = %v, want nil: %v", tw, tt.wantNil)
			}
		})
	}
}

func TestTimeWheel_AddAndExecuteTimer(t *testing.T) {
	var wg sync.WaitGroup
	executed := make(map[string]bool)
	var mu sync.Mutex

	job := func(data any) {
		mu.Lock()
		executed[data.(string)] = true
		mu.Unlock()
		wg.Done()
	}

	tw := New(2, 10, job)
	if tw == nil {
		t.Fatal("Failed to create TimeWheel")
	}

	tw.Start()
	defer tw.Stop()

	// 添加多个定时器任务
	testCases := []struct {
		key   string
		delay time.Duration
	}{
		{"task1", 200 * time.Millisecond},
		{"task2", 300 * time.Millisecond},
		{"task3", 400 * time.Millisecond},
	}

	wg.Add(len(testCases))
	for _, tc := range testCases {
		tw.AddTimer(tc.delay, tc.key, tc.key)
	}

	// 等待所有任务执行完成
	wg.Wait()

	// 验证所有任务是否都已执行
	for _, tc := range testCases {
		if !executed[tc.key] {
			t.Errorf("Task %s was not executed", tc.key)
		}
	}
}

func TestTimeWheel_RemoveTimer(t *testing.T) {
	var wg sync.WaitGroup
	executed := make(map[string]bool)
	var mu sync.Mutex

	job := func(data any) {
		mu.Lock()
		executed[data.(string)] = true
		mu.Unlock()
		wg.Done()
	}

	tw := New(2, 10, job)
	if tw == nil {
		t.Fatal("Failed to create TimeWheel")
	}

	tw.Start()
	defer tw.Stop()

	// 添加两个定时器任务
	wg.Add(1) // 只等待一个任务完成
	tw.AddTimer(3, "task1", "task1")
	tw.AddTimer(3, "task2", "task2")

	// 立即删除task2
	tw.RemoveTimer("task2")

	// 等待task1执行完成
	wg.Wait()

	// 验证task1执行，task2未执行
	if !executed["task1"] {
		t.Error("Task1 should have been executed")
	}
	if executed["task2"] {
		t.Error("Task2 should not have been executed")
	}
}

func TestTimeWheel_MultipleCircles(t *testing.T) {
	var wg sync.WaitGroup
	executed := false
	var mu sync.Mutex

	job := func(data any) {
		mu.Lock()
		executed = true
		mu.Unlock()
		wg.Done()
	}

	tw := New(1, 10, job)
	if tw == nil {
		t.Fatal("Failed to create TimeWheel")
	}

	tw.Start()
	defer tw.Stop()

	wg.Add(1)
	// 添加一个需要多圈才能执行的任务
	tw.AddTimer(4, "long-task", "long-task")

	// 等待任务执行完成
	wg.Wait()

	if !executed {
		t.Error("Long duration task was not executed")
	}
}

func TestTimeWheel_StressTest(t *testing.T) {
	const (
		numTasks      = 100000 // 总任务数
		numGoroutines = 100   // 并发goroutine数
		maxDelay      = 5     // 最大延迟时间(秒)
	)

	var wg sync.WaitGroup
	executed := make(map[string]bool)
	var mu sync.Mutex

	// 创建任务执行回调函数
	job := func(data any) {
		taskID := data.(string)
		mu.Lock()
		executed[taskID] = true
		mu.Unlock()
		wg.Done()
	}

	// 创建时间轮，使用较小的时间间隔和槽数以增加压力
	tw := New(1, 60, job)
	if tw == nil {
		t.Fatal("Failed to create TimeWheel")
	}

	tw.Start()
	defer tw.Stop()

	// 启动多个goroutine并发添加任务
	tasksPerGoroutine := numTasks / numGoroutines
	wg.Add(numTasks)

	for g := range numGoroutines {
		go func(goroutineID int) {
			for i := range tasksPerGoroutine {
				taskID := fmt.Sprintf("task-%d-%d", goroutineID, i)
				delay := time.Duration(rand.Intn(maxDelay)) * time.Second
				tw.AddTimer(delay, taskID, taskID)

				// 随机删除一些任务
				if rand.Float32() < 0.1 { // 10%的概率删除任务
					tw.RemoveTimer(taskID)
					mu.Lock()
					delete(executed, taskID)
					mu.Unlock()
					wg.Done() // 如果任务被删除，减少等待计数
				}
			}
		}(g)
	}

	// 等待所有未删除的任务执行完成
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// 设置超时时间
	timeout := time.After(time.Duration(maxDelay+20) * time.Second)

	// 等待任务完成或超时
	select {
	case <-done:
		// 成功完成
	case <-timeout:
		t.Fatal("Stress test timed out")
	}

	// 统计执行结果
	mu.Lock()
	executedCount := len(executed)
	mu.Unlock()

	t.Logf("Successfully executed %d tasks", executedCount)

	// 验证内存使用情况
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	t.Logf("Memory usage: Alloc = %v MiB", m.Alloc/1024/1024)
}
