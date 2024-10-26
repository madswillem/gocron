package gocron

import (
	"errors"
)

type Queue struct {
	Elements []*Job
	Size     int
}

func (q *Queue) Enqueue(elem *Job) error {
	if len(q.Elements) == q.Size {
		return errors.New("Queue Overflow")
	}
	q.Elements = append(q.Elements, elem)
	return nil
}

func (q *Queue) Dequeue() (*Job, error) {
	if len(q.Elements) == 0 {
		return nil, errors.New("Queue Underflow")
	}
	element := q.Elements[0]
	if len(q.Elements) == 1 {
		q.Elements = nil
		return element, nil
	}
	q.Elements = q.Elements[:1]
	return element, nil
}

func (q *Queue) Peek() (*Job, error) {
	if len(q.Elements) == 0 {
		return nil, errors.New("empty queue")
	}
	return q.Elements[0], nil
}
