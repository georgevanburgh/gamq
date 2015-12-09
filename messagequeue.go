package gamq

import (
	"sync"
	"sync/atomic"
)

type queuenode struct {
	data *string
	next *queuenode
}

// A goroutine safe FIFO, based on
// https://github.com/hishboy/gocommons/blob/06389f1595e56cd7c27d9dc9fe48fc771db1b5ef/lang/queue.go
type MessageQueue struct {
	head  *queuenode
	tail  *queuenode
	count uint64
	lock  *sync.Mutex
}

func NewMessageQueue() *MessageQueue {
	q := &MessageQueue{}
	q.lock = &sync.Mutex{}
	return q
}

func (q *MessageQueue) Len() uint64 {
	length := atomic.LoadUint64(&q.count)
	return length
}

func (q *MessageQueue) Push(item *string) {
	q.lock.Lock()
	defer q.lock.Unlock()

	node := &queuenode{data: item}

	if q.tail == nil {
		q.tail = node
		q.head = node
	} else {
		q.tail.next = node
		q.tail = node
	}
	q.count++
}

func (q *MessageQueue) Poll() *string {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.head == nil {
		return nil
	}

	node := q.head
	q.head = q.head.next

	// If q is empty
	if q.head == nil {
		q.tail = nil
	}
	q.count--

	return node.data
}

func (q *MessageQueue) Peek() *string {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.head == nil {
		return nil
	} else {
		return q.head.data
	}
}
