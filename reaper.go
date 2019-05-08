package main

import (
	"fmt"
	"syscall"
	"time"
)

func Reap(units []*Unit) {
	var (
		rusage    = syscall.Rusage{}
		ws        syscall.WaitStatus
		pid       int
		err       error
		lastError error
	)
	for {
		pid, err = syscall.Wait4(0, &ws, syscall.WALL, &rusage)
		if err != nil {
			if err != lastError {
				Logf("error while waiting: %v\n", err)
			}
			lastError = err
			if err == syscall.ECHILD {
				time.Sleep(1 * time.Second)
			}
			continue
		}

		for _, v := range units {
			if !v.runtime.Started { //skip if process hasn't yet been started
				continue //protects again segfault: dereference of v.runtime.Proc
			}
			if v.runtime.GetPid() == pid {
				fmt.Fprintf(v.runtime.log, "[%v]: stopped\n", time.Now())
				switch v.Type {
				case UnitTypeService:
					v.Restart()
				case UnitTypeTask:
					fmt.Println()
					v.runtime.IsDoneChan <- struct{}{}
				}
				break
			}
		}

		Logf("process %d stopped with exit status %d\n", pid, ws.ExitStatus())
	}
}
