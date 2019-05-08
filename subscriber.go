package main

import "io"

type SubscriberQueue struct {
	headNode *SubQueueNode
}

type SubQueueNode struct {
	w    io.Writer
	next *SubQueueNode
}

func (s *SubscriberQueue) Add(w io.Writer) {
	lastHead := s.headNode
	s.headNode = new(SubQueueNode)
	s.headNode.next = lastHead
	s.headNode.w = w
}

func (s *SubscriberQueue) Write(bs []byte) (int, error) {
	var (
		last *SubQueueNode
		cur  *SubQueueNode = s.headNode
	)

	for cur != nil {
		_, err := cur.w.Write(bs)
		if err != nil { //couldn't write or whatever
			if last != nil {
				last.next = cur.next
			} else {
				s.headNode = cur.next
			}
			cur = cur.next
			continue
		}
		last = cur
		cur = cur.next
	}
	return len(bs), nil
}
