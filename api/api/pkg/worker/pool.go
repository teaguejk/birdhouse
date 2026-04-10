package worker

import (
	"api/pkg/logging"
	"sync"
	"time"
)

// Task represents a unit of work to be processed by the pool.
type Task struct {
	// Key is used for debouncing — tasks with the same key within the debounce
	// window are merged. If empty, the task is processed immediately without debouncing.
	Key string

	// Fn is the function to execute.
	Fn func()

	// MergeFn is called when a new task merges into an existing pending task.
	// If nil, the newer task's Fn replaces the older one.
	MergeFn func(existing *Task)
}

type pendingTask struct {
	task  *Task
	timer *time.Timer
}

// Pool provides bounded concurrency for async work with optional debouncing.
type Pool struct {
	queue    chan *Task
	workers  int
	debounce time.Duration
	logger   *logging.Logger

	pending map[string]*pendingTask
	mu      sync.Mutex
	wg      sync.WaitGroup
	stop    chan struct{}
	stopped bool
}

// Config for creating a new pool.
type Config struct {
	// Workers is the number of concurrent goroutines processing tasks.
	Workers int

	// QueueSize is the buffer size for the task channel.
	QueueSize int

	// Debounce is the window during which tasks with the same key are merged.
	// Set to 0 to disable debouncing.
	Debounce time.Duration
}

// NewPool creates and starts a worker pool.
func NewPool(logger *logging.Logger, cfg Config) *Pool {
	if cfg.Workers <= 0 {
		cfg.Workers = 5
	}
	if cfg.QueueSize <= 0 {
		cfg.QueueSize = 500
	}

	p := &Pool{
		queue:    make(chan *Task, cfg.QueueSize),
		workers:  cfg.Workers,
		debounce: cfg.Debounce,
		logger:   logger,
		pending:  make(map[string]*pendingTask),
		stop:     make(chan struct{}),
	}

	// Start workers
	for i := 0; i < cfg.Workers; i++ {
		p.wg.Add(1)
		go p.worker()
	}

	logger.Info("worker pool started", "workers", cfg.Workers, "queue_size", cfg.QueueSize, "debounce", cfg.Debounce)
	return p
}

// Submit adds a task to the pool. If the task has a Key and debouncing is enabled,
// it will be held for the debounce window and merged with subsequent same-key tasks.
// Returns false if the pool is stopped or the queue is full.
func (p *Pool) Submit(task *Task) bool {
	p.mu.Lock()
	if p.stopped {
		p.mu.Unlock()
		return false
	}

	// No debounce or no key — send directly to queue
	if p.debounce <= 0 || task.Key == "" {
		p.mu.Unlock()
		select {
		case p.queue <- task:
			return true
		default:
			p.logger.Warn("worker pool queue full, dropping task")
			return false
		}
	}

	// Debounce: check for existing pending task with same key
	if pt, ok := p.pending[task.Key]; ok {
		// Merge into existing
		if task.MergeFn != nil {
			task.MergeFn(pt.task)
		} else {
			pt.task.Fn = task.Fn
		}
		// Reset timer
		pt.timer.Reset(p.debounce)
		p.mu.Unlock()
		return true
	}

	// New key — create pending entry with timer
	pt := &pendingTask{task: task}
	key := task.Key
	pt.timer = time.AfterFunc(p.debounce, func() {
		p.mu.Lock()
		flushed, ok := p.pending[key]
		if ok {
			delete(p.pending, key)
		}
		p.mu.Unlock()

		if ok {
			select {
			case p.queue <- flushed.task:
			default:
				p.logger.Warn("worker pool queue full, dropping debounced task", "key", key)
			}
		}
	})
	p.pending[key] = pt
	p.mu.Unlock()
	return true
}

// Shutdown stops accepting new tasks, flushes all pending debounced tasks,
// and waits for all workers to finish processing.
func (p *Pool) Shutdown() {
	p.mu.Lock()
	p.stopped = true

	// Flush all pending debounced tasks immediately
	for key, pt := range p.pending {
		pt.timer.Stop()
		delete(p.pending, key)
		select {
		case p.queue <- pt.task:
		default:
			p.logger.Warn("worker pool queue full during shutdown, dropping task", "key", key)
		}
	}
	p.mu.Unlock()

	// Close queue to signal workers to drain and exit
	close(p.queue)

	// Wait for all workers to finish
	p.wg.Wait()
	p.logger.Info("worker pool shutdown complete")
}

func (p *Pool) worker() {
	defer p.wg.Done()
	for task := range p.queue {
		func() {
			defer func() {
				if r := recover(); r != nil {
					p.logger.Error("worker pool task panicked", "error", r)
				}
			}()
			task.Fn()
		}()
	}
}
