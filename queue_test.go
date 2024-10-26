package gocron

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEnqueue(t *testing.T) {
	r := Queue{Size: 5}
	r.Enqueue(&Job{})
	expected := make([]*Job, 1)
	expected[0] = &Job{}
	d := cmp.Diff(expected, r)
	if d == "" {
		t.Error(d)
	}
}
func TestDequeue(t *testing.T) {
	r := Queue{Size: 5}
	r.Elements = append(r.Elements, &Job{})
	e, err := r.Dequeue()
	if err != nil {
		t.Error(err)
	}
	if e == nil {
		t.Error("No element was returned")
	}
}
