package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestSubscriber(t *testing.T) {
	q := new(SubscriberQueue)
	l := make([]*bytes.Buffer, 0, 10)
	for i := 0; i < 10; i++ {
		var w io.Writer
		b := bytes.NewBuffer(make([]byte, 4096))
		w = b
		if i == 2 {
			a := new(readTwice)
			a.w = b
			w = a
		}
		l = append(l, b)
		q.Add(w)
		fmt.Fprintf(q, "%d%d\n", i, i)
	}
	for _, v := range l {
		fmt.Print(string(v.Bytes()))
	}
}

type readTwice struct {
	n int
	w io.Writer
}

func (r *readTwice) Write(b []byte) (int, error) {
	if !(r.n < 2) {
		return 0, io.ErrClosedPipe
	}
	n, err := r.w.Write(b)
	r.n++
	return n, err
}
