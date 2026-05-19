package bot

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJobQueue_RunOnce(t *testing.T) {
	bot := New("token")
	jq := NewJobQueue(bot)
	jq.Start()
	defer jq.Stop()

	var wg sync.WaitGroup
	wg.Add(1)

	called := false
	var receivedData interface{}

	jq.RunOnce(func(ctx *JobContext) error {
		called = true
		receivedData = ctx.Job.Data
		wg.Done()
		return nil
	}, time.Now().Add(10*time.Millisecond), WithJobData("test-payload"), WithJobName("once-job"))

	// Wait with timeout
	c := make(chan struct{})
	go func() {
		wg.Wait()
		close(c)
	}()

	select {
	case <-c:
		assert.True(t, called)
		assert.Equal(t, "test-payload", receivedData)
	case <-time.After(1 * time.Second):
		t.Fatal("Job once-job did not execute in time")
	}
}

func TestJobQueue_RunRepeating(t *testing.T) {
	bot := New("token")
	jq := NewJobQueue(bot)
	jq.Start()
	defer jq.Stop()

	var mu sync.Mutex
	runs := 0

	job := jq.RunRepeating(func(ctx *JobContext) error {
		mu.Lock()
		runs++
		mu.Unlock()
		return nil
	}, 20*time.Millisecond, time.Now())

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	assert.GreaterOrEqual(t, runs, 3)
	mu.Unlock()

	job.Remove()
	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	lastRuns := runs
	mu.Unlock()

	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	assert.Equal(t, lastRuns, runs)
	mu.Unlock()
}
