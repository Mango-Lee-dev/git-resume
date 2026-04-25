package jobs

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/wootaiklee/git-resume/internal/api/dto"
	"github.com/wootaiklee/git-resume/internal/service"
)

// Job represents an async analysis job
type Job struct {
	ID          string
	Status      dto.JobStatus
	Config      service.AnalyzeConfig
	Progress    int
	Phase       string
	Message     string
	CreatedAt   time.Time
	StartedAt   *time.Time
	CompletedAt *time.Time
	ResultCount int
	Error       string

	ctx        context.Context
	cancelFunc context.CancelFunc
}

// ToResponse converts Job to JobResponse DTO
func (j *Job) ToResponse() dto.JobResponse {
	return dto.JobResponse{
		ID:          j.ID,
		Status:      j.Status,
		Progress:    j.Progress,
		Phase:       j.Phase,
		Message:     j.Message,
		CreatedAt:   j.CreatedAt,
		StartedAt:   j.StartedAt,
		CompletedAt: j.CompletedAt,
		ResultCount: j.ResultCount,
		Error:       j.Error,
	}
}

// Manager handles job lifecycle
type Manager struct {
	mu       sync.RWMutex
	jobs     map[string]*Job
	analyzer *service.Analyzer
	workers  *WorkerPool
}

// NewManager creates a job manager
func NewManager(analyzer *service.Analyzer, workerCount int) *Manager {
	m := &Manager{
		jobs:     make(map[string]*Job),
		analyzer: analyzer,
	}

	m.workers = NewWorkerPool(workerCount, m.executeJob)
	return m
}

// Submit creates a new job and queues it for execution
func (m *Manager) Submit(cfg service.AnalyzeConfig) (*Job, error) {
	ctx, cancel := context.WithCancel(context.Background())

	job := &Job{
		ID:         uuid.New().String(),
		Status:     dto.JobStatusPending,
		Config:     cfg,
		CreatedAt:  time.Now(),
		ctx:        ctx,
		cancelFunc: cancel,
	}

	m.mu.Lock()
	m.jobs[job.ID] = job
	m.mu.Unlock()

	m.workers.Submit(job)

	return job, nil
}

// Get retrieves job status
func (m *Manager) Get(id string) (*Job, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	job, ok := m.jobs[id]
	return job, ok
}

// Cancel cancels a running job
func (m *Manager) Cancel(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	job, ok := m.jobs[id]
	if !ok {
		return nil
	}

	if job.Status == dto.JobStatusRunning || job.Status == dto.JobStatusPending {
		job.cancelFunc()
		job.Status = dto.JobStatusCancelled
		now := time.Now()
		job.CompletedAt = &now
		job.Message = "Job cancelled"
	}

	return nil
}

// List returns all jobs
func (m *Manager) List() []*Job {
	m.mu.RLock()
	defer m.mu.RUnlock()

	jobs := make([]*Job, 0, len(m.jobs))
	for _, job := range m.jobs {
		jobs = append(jobs, job)
	}
	return jobs
}

// executeJob runs the analysis (called by worker)
func (m *Manager) executeJob(job *Job) {
	m.mu.Lock()
	job.Status = dto.JobStatusRunning
	now := time.Now()
	job.StartedAt = &now
	m.mu.Unlock()

	progressCh := make(chan service.AnalyzeProgress, 10)

	// Start progress listener
	go func() {
		for p := range progressCh {
			m.mu.Lock()
			job.Progress = p.Progress
			job.Phase = string(p.Phase)
			job.Message = p.Message
			job.ResultCount = p.ResultsCreated
			if p.Error != nil {
				job.Error = p.Error.Error()
			}
			m.mu.Unlock()
		}
	}()

	// Run analysis
	result, err := m.analyzer.Analyze(job.ctx, job.Config, progressCh)

	m.mu.Lock()
	defer m.mu.Unlock()

	completedAt := time.Now()
	job.CompletedAt = &completedAt

	if err != nil {
		if job.Status != dto.JobStatusCancelled {
			job.Status = dto.JobStatusFailed
			job.Error = err.Error()
		}
		return
	}

	job.Status = dto.JobStatusCompleted
	job.Progress = 100
	if result != nil {
		job.ResultCount = len(result.Results)
	}
}

// Cleanup removes old completed/failed jobs
func (m *Manager) Cleanup(olderThan time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	cutoff := time.Now().Add(-olderThan)

	for id, job := range m.jobs {
		if job.CompletedAt != nil && job.CompletedAt.Before(cutoff) {
			delete(m.jobs, id)
		}
	}
}

// Shutdown stops the job manager
func (m *Manager) Shutdown() {
	m.workers.Shutdown()

	m.mu.Lock()
	defer m.mu.Unlock()

	// Cancel all running jobs
	for _, job := range m.jobs {
		if job.Status == dto.JobStatusRunning || job.Status == dto.JobStatusPending {
			job.cancelFunc()
			job.Status = dto.JobStatusCancelled
		}
	}
}
