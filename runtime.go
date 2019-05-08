package main

import (
	"os"
	"sync"
)

type UnitRuntime struct {
	lock       sync.Mutex
	Started    bool
	Proc       *os.Process
	IOFiles    map[string]*os.File
	log        *Log
	IsDoneChan chan struct{}
}

func (u *UnitRuntime) GetPid() int {
	u.lock.Lock()
	defer u.lock.Unlock()
	if u.Proc != nil {
		return u.Proc.Pid
	}
	return 0
}
