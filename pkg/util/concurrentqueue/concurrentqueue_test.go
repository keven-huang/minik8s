package queue

import "testing"

func TestQueue(t *testing.T) {
	q := NewConcurrentQueue()
	q.Push(1)
	q.Push(2)
	q.Push(3)
	if q.Pop() != 1 {
		t.Error("pop error")
	}
	if q.Pop() != 2 {
		t.Error("pop error")
	}
	if q.Pop() != 3 {
		t.Error("pop error")
	}
	if q.Pop() != nil {
		t.Error("pop error")
	}
}
