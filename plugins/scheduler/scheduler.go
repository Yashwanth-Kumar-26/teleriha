package scheduler

import (
	"context"
	"sync"
	"time"

	"github.com/Yashwanth-Kumar-26/teleriha/pkg/bot"
)

// Task represents a scheduled task.
type Task struct {
	Interval time.Duration
	Func     func()
	LastRun  time.Time
}

// SchedulerPlugin implements the bot.Plugin interface.
type SchedulerPlugin struct {
	bot.BasePlugin
	tasks  []*Task
	mu     sync.Mutex
	stop   chan struct{}
	wg     sync.WaitGroup
}

// New creates a new SchedulerPlugin.
func New() *SchedulerPlugin {
	return &SchedulerPlugin{
		BasePlugin: *bot.NewBasePlugin("scheduler"),
		stop:       make(chan struct{}),
	}
}

// ScheduleEvery schedules a function to run every interval.
func (p *SchedulerPlugin) ScheduleEvery(interval time.Duration, fn func()) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.tasks = append(p.tasks, &Task{
		Interval: interval,
		Func:     fn,
		LastRun:  time.Now(),
	})
}

// Start starts the scheduler background loop.
func (p *SchedulerPlugin) Start(ctx context.Context) error {
	if err := p.BasePlugin.Start(ctx); err != nil {
		return err
	}

	p.wg.Add(1)
	go p.run(ctx)
	return nil
}

func (p *SchedulerPlugin) run(ctx context.Context) {
	defer p.wg.Done()
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-p.stop:
			return
		case <-ticker.C:
			p.checkTasks()
		}
	}
}

func (p *SchedulerPlugin) checkTasks() {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	for _, task := range p.tasks {
		if now.Sub(task.LastRun) >= task.Interval {
			go task.Func()
			task.LastRun = now
		}
	}
}

// Stop stops the scheduler.
func (p *SchedulerPlugin) Stop() error {
	close(p.stop)
	p.wg.Wait()
	return p.BasePlugin.Stop()
}
