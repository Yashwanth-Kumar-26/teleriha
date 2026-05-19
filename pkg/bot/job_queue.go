package bot

import (
	"fmt"
	"sync"
	"time"
)

// Job represents a task managed by the JobQueue.
type Job struct {
	ID        string
	Name      string
	NextRun   time.Time
	Interval  time.Duration // 0 for run once
	Callback  func(*JobContext) error
	ChatID    int64
	UserID    int64
	Data      interface{}
	removed   bool
	mu        sync.Mutex
}

// Remove marks the job as removed so it won't run again.
func (j *Job) Remove() {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.removed = true
}

// IsRemoved returns whether the job has been removed.
func (j *Job) IsRemoved() bool {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.removed
}

// JobContext holds the environment in which a job runs.
type JobContext struct {
	Bot      *Bot
	Job      *Job
	JobQueue *JobQueue
}

// JobQueue manages and runs scheduled jobs.
type JobQueue struct {
	bot       *Bot
	jobs      []*Job
	mu        sync.Mutex
	stopChan  chan struct{}
	wg        sync.WaitGroup
	running   bool
}

// NewJobQueue creates a new JobQueue instance.
func NewJobQueue(bot *Bot) *JobQueue {
	return &JobQueue{
		bot:      bot,
		stopChan: make(chan struct{}),
	}
}

// JobOption allows configuring additional job parameters.
type JobOption func(*Job)

// WithJobName sets a custom name for the job.
func WithJobName(name string) JobOption {
	return func(j *Job) {
		j.Name = name
	}
}

// WithJobChatID associates a chat ID with the job.
func WithJobChatID(chatID int64) JobOption {
	return func(j *Job) {
		j.ChatID = chatID
	}
}

// WithJobUserID associates a user ID with the job.
func WithJobUserID(userID int64) JobOption {
	return func(j *Job) {
		j.UserID = userID
	}
}

// WithJobData attaches custom data to the job.
func WithJobData(data interface{}) JobOption {
	return func(j *Job) {
		j.Data = data
	}
}

// Start runs the JobQueue background execution loop.
func (jq *JobQueue) Start() {
	jq.mu.Lock()
	if jq.running {
		jq.mu.Unlock()
		return
	}
	jq.running = true
	jq.stopChan = make(chan struct{})
	jq.mu.Unlock()

	jq.wg.Add(1)
	go jq.loop()
}

// Stop gracefully stops the JobQueue.
func (jq *JobQueue) Stop() {
	jq.mu.Lock()
	if !jq.running {
		jq.mu.Unlock()
		return
	}
	jq.running = false
	close(jq.stopChan)
	jq.mu.Unlock()

	jq.wg.Wait()
}

func (jq *JobQueue) loop() {
	defer jq.wg.Done()
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-jq.stopChan:
			return
		case <-ticker.C:
			jq.runPendingJobs()
		}
	}
}

func (jq *JobQueue) runPendingJobs() {
	jq.mu.Lock()
	now := time.Now()
	var pending []*Job

	// Filter out removed jobs and identify pending ones
	var activeJobs []*Job
	for _, job := range jq.jobs {
		if job.IsRemoved() {
			continue
		}
		activeJobs = append(activeJobs, job)

		if now.After(job.NextRun) || now.Equal(job.NextRun) {
			pending = append(pending, job)
		}
	}
	jq.jobs = activeJobs
	jq.mu.Unlock()

	for _, job := range pending {
		jq.wg.Add(1)
		go func(j *Job) {
			defer jq.wg.Done()

			if j.IsRemoved() {
				return
			}

			jobCtx := &JobContext{
				Bot:      jq.bot,
				Job:      j,
				JobQueue: jq,
			}

			// Run job callback
			err := j.Callback(jobCtx)
			if err != nil {
				jq.bot.Logger.Error().Err(err).Str("job_name", j.Name).Msg("Job execution returned error")
			}

			// Update scheduling state
			j.mu.Lock()
			if j.Interval > 0 && !j.removed {
				j.NextRun = time.Now().Add(j.Interval)
			} else {
				j.removed = true
			}
			j.mu.Unlock()
		}(job)
	}
}

// RunOnce schedules a job to run once at a specific time.
func (jq *JobQueue) RunOnce(callback func(*JobContext) error, when time.Time, opts ...JobOption) *Job {
	job := &Job{
		ID:       fmt.Sprintf("job_%d", time.Now().UnixNano()),
		NextRun:  when,
		Callback: callback,
	}
	for _, opt := range opts {
		opt(job)
	}

	jq.mu.Lock()
	jq.jobs = append(jq.jobs, job)
	jq.mu.Unlock()

	return job
}

// RunRepeating schedules a job to run repeatedly with a fixed interval.
func (jq *JobQueue) RunRepeating(callback func(*JobContext) error, interval time.Duration, first time.Time, opts ...JobOption) *Job {
	job := &Job{
		ID:       fmt.Sprintf("job_%d", time.Now().UnixNano()),
		NextRun:  first,
		Interval: interval,
		Callback: callback,
	}
	for _, opt := range opts {
		opt(job)
	}

	jq.mu.Lock()
	jq.jobs = append(jq.jobs, job)
	jq.mu.Unlock()

	return job
}
