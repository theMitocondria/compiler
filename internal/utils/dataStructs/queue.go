package dataStructs

import (
	"errors"
)

type Queue struct {
	elements []int
}

func (q *Queue) Push(value int) {
	q.elements = append(q.elements, value)
}

func (q *Queue) Pop() (int, error) {
	if len(q.elements) == 0 {
		return 0, errors.New("queue is empty")
	}
	value := q.elements[0]
	q.elements = q.elements[1:]
	return value, nil
}

func (q *Queue) IsEmpty() bool {
	return len(q.elements) == 0
}

func (q *Queue) Size() int {
	return len(q.elements)
}

func InitializeQueue(values []int) Queue {
	q := Queue{}
	for _, value := range values {
		q.Push(value)
	}
	return q
}
