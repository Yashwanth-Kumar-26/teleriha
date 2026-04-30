package scheduler

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/Yashwanth-Kumar-26/teleriha/pkg/bot"
)

func TestSchedulerPlugin(t *testing.T) {
	b := bot.New("test-token")
	
	p := New()
	err := p.Init(b)
	assert.NoError(t, err)
	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	err = p.Start(ctx)
	assert.NoError(t, err)
	
	var mu sync.Mutex
	called := 0
	
	p.ScheduleEvery(50*time.Millisecond, func() {
		mu.Lock()
		called++
		mu.Unlock()
	})
	
	time.Sleep(120 * time.Millisecond)
	
	mu.Lock()
	assert.GreaterOrEqual(t, called, 2)
	mu.Unlock()
	
	p.Stop()
}
