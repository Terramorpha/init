package main

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

type Log struct {
	buffer      []byte
	lock        sync.Mutex
	subscribers *SubscriberQueue
}

func NewLog() *Log {
	o := &Log{
		buffer:      make([]byte, 0, 4096),
		subscribers: &SubscriberQueue{},
	}
	return o
}

func (l *Log) Write(in []byte) (int, error) {
	if l == nil { //by default is a mock writer
		return len(in), nil
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	l.buffer = append(l.buffer, in...)
	l.subscribers.Write(in)
	return len(in), nil
}

func (l *Log) AddSubscriber(w io.Writer) {
	l.lock.Lock()
	defer l.lock.Unlock()
	io.Copy(w, bytes.NewReader(l.buffer))
	l.subscribers.Add(w)
}

func Logf(format string, args ...interface{}) (int, error) {
	return fmt.Fprintf(SystemLog, format, args...)
}
