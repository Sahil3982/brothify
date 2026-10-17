package cron

import (
	"context"
	"fmt"
	"log"
	"sync"

	robfigcron "github.com/robfig/cron/v3"
)

type Job interface {
	Name() string
	Spec() string
	Run(ctx context.Context) error
}

type Scheduler struct {
	cron *robfigcron.Cron
	mu   sync.Mutex
	jobs map[string]struct{}
}

func NewScheduler() *Scheduler {   //dependency injection
	return &Scheduler{
		cron: robfigcron.New(robfigcron.WithChain(
			robfigcron.SkipIfStillRunning(robfigcron.DefaultLogger),
		)),
		jobs: make(map[string]struct{}),
	}
}

func (s *Scheduler) Register(job Job) error {
	if job == nil {
		return fmt.Errorf("cannot register a nil cron job")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.jobs[job.Name()]; exists {
		return fmt.Errorf("cron job %q is already registered", job.Name())
	}

	_, err := s.cron.AddFunc(job.Spec(), func() {
		if err := job.Run(context.Background()); err != nil {
			log.Printf("Cron job %q failed: %v", job.Name(), err)
		}
	})
	if err != nil {
		return fmt.Errorf("register cron job %q: %w", job.Name(), err)
	}
	s.jobs[job.Name()] = struct{}{}
	return nil
}

func (s *Scheduler) Start() {
	s.cron.Start()
}

func (s *Scheduler) Stop() {
	<-s.cron.Stop().Done()
}
