package gocron

import (
	"errors"
	"sync"
	"time"
)

type Registry struct {
	Jobs map[string]Job
}
type Job struct {
	Name      string
	Job       func(*sync.WaitGroup, chan bool, chan error)
	Ticker    *time.Ticker
	ErrorChan chan error
	Done      chan bool
	Wg        *sync.WaitGroup
}

func (r *Registry) Add(j Job) error {
	if _, ok := r.Jobs[j.Name]; ok {
		return errors.New("Job already exists")
	}
	r.Jobs[j.Name] = j
	return nil
}

func (r *Registry) Exec(name string) (chan bool, chan error, error) {
	if entry, ok := r.Jobs[name]; ok {
		entry.ErrorChan = make(chan error, 1)
		entry.Done = make(chan bool, 1)
		entry.Wg = &sync.WaitGroup{}
		entry.Wg.Add(1)
		go entry.Job(entry.Wg, entry.Done, entry.ErrorChan)
		go func() {
			entry.Wg.Wait()
			close(entry.ErrorChan)
			close(entry.Done)
		}()
		return entry.Done, entry.ErrorChan, nil
	}
	return nil, nil, errors.New("Job doesnt exist")
}

func New() *Registry {
	return &Registry{Jobs: make(map[string]Job)}
}
