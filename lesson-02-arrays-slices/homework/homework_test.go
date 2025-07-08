package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type CircularQueue struct {
	values []int
	front  int
	rear   int
}

func NewCircularQueue(size int) CircularQueue {
	return CircularQueue{
		values: make([]int, size),
		front:  -1,
		rear:   -1,
	}
}

func (q *CircularQueue) Push(value int) bool {
	if q.Full() {
		return false
	}
	if q.Empty() {
		q.front = 0
	}
	if q.rear == cap(q.values)-1 {
		q.rear = 0
	} else {
		q.rear++
	}
	q.values[q.rear] = value
	return true
}

func (q *CircularQueue) Pop() bool {
	return q.pop(false) != -1
}

func (q *CircularQueue) pop(rear bool) int {
	if q.Empty() {
		return -1
	}
	if q.front == q.rear {
		i := q.front
		q.front = -1
		q.rear = -1
		return q.values[i]
	}
	if rear {
		return q.getRear()
	}
	return q.getFront()
}

func (q *CircularQueue) getRear() int {
	i := q.rear
	if q.rear == 0 {
		q.rear = cap(q.values) - 1
	} else {
		q.rear--
	}
	return q.values[i]
}

func (q *CircularQueue) getFront() int {
	i := q.front
	if q.front == cap(q.values)-1 {
		q.front = 0
	} else {
		q.front++
	}
	return q.values[i]
}

func (q *CircularQueue) Front() int {
	return q.pop(false)
}

func (q *CircularQueue) Back() int {
	return q.pop(true)
}

func (q *CircularQueue) Empty() bool {
	return q.front == -1 && q.rear == -1
}

func (q *CircularQueue) Full() bool {
	return (q.front == 0 && q.rear == cap(q.values)-1) || q.front == q.rear+1
}

func TestCircularQueue(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue(queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, -1, queue.Front())
	assert.Equal(t, -1, queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, 1, queue.Front())
	assert.Equal(t, 3, queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{4, 2, 3}, queue.values))

	assert.Equal(t, 4, queue.Front())
	assert.Equal(t, -1, queue.Back())

	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}
