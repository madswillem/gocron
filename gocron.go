package gocron

import (
	"errors"
	"sync"
	"time"
)

type Registry struct {
	Jobs map[int]Job
}
type Job struct {
	Job       func(*sync.WaitGroup)
	Ticker    *time.Ticker
	ErrorChan chan error
	Wg        *sync.WaitGroup
}

func (r *Registry) Add(j Job) error {
	r.Jobs[len(r.Jobs)] = j
	return nil
}

func (r *Registry) Exec(id int) (*chan error, error) {
	if entry, ok := r.Jobs[id]; ok {
		entry.ErrorChan = make(chan error)
		entry.Wg = &sync.WaitGroup{}
		go entry.Job(entry.Wg)
		go func() {
			entry.Wg.Wait()
			close(entry.ErrorChan)
		}()
		return &entry.ErrorChan, nil
	}
	return nil, errors.New("Job doesnt exist")
}

func New() *Registry {
	return &Registry{}
}
