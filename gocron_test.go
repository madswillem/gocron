package gocron

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestExec(t *testing.T) {
	data := []struct {
		name           string
		job            Job
		execId         string
		expectedError  error
		expectedjoberr error
		wantPanic      bool
	}{
		{
			name: "test",
			job: Job{
				Name: "test",
				Job: func(wg *sync.WaitGroup, done chan bool, err chan error) {
					time.Sleep(2 * time.Second)
					done <- true
					wg.Done()
				},
			},
			expectedError:  nil,
			expectedjoberr: nil,
		},
		{
			name: "test_joberr",
			job: Job{
				Name: "test_joberr",
				Job: func(wg *sync.WaitGroup, done chan bool, err chan error) {
					time.Sleep(2 * time.Second)
					err <- errors.New("test error")
					done <- true
					wg.Done()
				},
			},
			expectedError:  nil,
			expectedjoberr: errors.New("test error"),
		},
		{
			name: "test_wrong_jobid",
			job: Job{
				Name: "test",
				Job: func(wg *sync.WaitGroup, done chan bool, err chan error) {
					time.Sleep(2 * time.Second)
					done <- true
					wg.Done()
				},
			},
			execId:         fmt.Sprint("test"),
			expectedError:  errors.New("Job doesnt exist"),
			expectedjoberr: nil,
		},
		{
			name: "test_jobpanic",
			job: Job{
				Name: "test_joberr",
				Job: func(wg *sync.WaitGroup, done chan bool, err chan error) {
					time.Sleep(2 * time.Second)
					panic("test")
				},
			},
			expectedError:  nil,
			expectedjoberr: nil,
			wantPanic:      true,
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.name, func(t *testing.T) {
			t.Parallel()
			defer func() {
				r := recover()
				if (r != nil) != d.wantPanic {
					t.Errorf("SequenceInt() recover = %v, wantPanic = %v", r, d.wantPanic)
				}
			}()
			entry := New()
			entry.Add(d.job)
			if d.execId != "" {
				d.job.Name = d.execId
			}
			done, joberr, err := entry.Exec(d.job.Name)
			<-done
			if errors.Is(err, d.expectedError) && d.expectedError != nil {
				t.Errorf("Job error was %e but expected %e", d.expectedError, err)
			}
			if e := <-joberr; errors.Is(e, d.expectedjoberr) && d.expectedjoberr != nil {
				t.Errorf("Job error was %e but expected %e", d.expectedjoberr, e)
			}
		})
	}
}
func TestTicker(t *testing.T) {
	r := New()
	executed := false
	err := r.Add(Job{
		Name: "test_ticker",
		Job: func(wg *sync.WaitGroup, done chan bool, err chan error) {
			executed = true
		},
		Ticker: time.NewTicker(5 * time.Millisecond),
	})
	if err != nil {
		t.Errorf("Error: %e", err)
	}
	go func() {
		// Allow some time for the ticker to tick
		time.Sleep(15 * time.Millisecond)
		r.TickerExecuter()
	}()
	time.Sleep(20 * time.Millisecond)
	if !executed {
		t.Errorf("Function was not executed. %v", !executed)
	}
}
