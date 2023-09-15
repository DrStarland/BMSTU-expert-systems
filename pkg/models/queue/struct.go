package queue

import (
	"errors"
	"sync"
)

// Using Generics to define Type in Stake to Use Structs, too.
type Queue[T any] struct {
	lock *sync.Mutex // Mutex for Thread safety
	S    []T         // Slice
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{lock: &sync.Mutex{}, S: []T{}}
}

func (st Queue[T]) Len() int {
	st.lock.Lock()
	defer st.lock.Unlock()
	return len(st.S)
}

func (queue *Queue[T]) Push(element T) {
	queue.lock.Lock()
	defer queue.lock.Unlock()
	queue.S = append(queue.S, element)
}

func (queue *Queue[T]) Pop() (T, error) {
	queue.lock.Lock()
	defer queue.lock.Unlock()
	l := len(queue.S)
	if l == 0 {
		var empty T
		return empty, errors.New("empty Queue")
	}
	element := queue.S[0]
	// log.Println("завис здесь")
	queue.S = queue.S[1:]
	return element, nil
}

func (queue Queue[T]) Peek() (T, error) {
	queue.lock.Lock()
	defer queue.lock.Unlock()
	l := len(queue.S)
	if l == 0 {
		var empty T
		return empty, errors.New("empty Queue")
	}
	return queue.S[0], nil
}
