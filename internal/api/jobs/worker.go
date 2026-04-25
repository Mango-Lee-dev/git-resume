package jobs

// WorkerPool manages concurrent job execution
type WorkerPool struct {
	workers  int
	jobQueue chan *Job
	quit     chan struct{}
	execute  func(*Job)
}

// NewWorkerPool creates a worker pool
func NewWorkerPool(workers int, execute func(*Job)) *WorkerPool {
	if workers <= 0 {
		workers = 2
	}

	pool := &WorkerPool{
		workers:  workers,
		jobQueue: make(chan *Job, 100),
		quit:     make(chan struct{}),
		execute:  execute,
	}

	pool.start()
	return pool
}

// start launches worker goroutines
func (p *WorkerPool) start() {
	for i := 0; i < p.workers; i++ {
		go p.worker()
	}
}

// worker processes jobs from the queue
func (p *WorkerPool) worker() {
	for {
		select {
		case job := <-p.jobQueue:
			if job != nil {
				p.execute(job)
			}
		case <-p.quit:
			return
		}
	}
}

// Submit adds a job to the queue
func (p *WorkerPool) Submit(job *Job) {
	select {
	case p.jobQueue <- job:
	default:
		// Queue full, job will be dropped
		// In production, consider returning an error or blocking
	}
}

// Shutdown stops all workers
func (p *WorkerPool) Shutdown() {
	close(p.quit)
}

// QueueSize returns the number of pending jobs
func (p *WorkerPool) QueueSize() int {
	return len(p.jobQueue)
}
