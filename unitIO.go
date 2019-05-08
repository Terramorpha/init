package main

import (
	"errors"
	"os"
	"strings"
	"syscall"
)

//dans la partie files du config unit
//il y a les fichiers que le processus va utiliser

const (
	IOTypeUnixSocket = "unix"
	IOTypeUnixConn   = "unixconn"
	IOTypeInetSocket = "inet"
	IOTypeInetConn   = "inetconn"
	IOTypeFile       = "file"
	IOTypeFIFO       = "fifo"
)

var (
	ErrInvalidFlag = errors.New("invalid flag")
	ErrInvalidFmt  = errors.New("file invalid format")
)

func UnitGetIOFiles(f map[string]string) (map[string]*os.File, error) {
	o := make(map[string]*os.File)
	for i, v := range f {

		things := strings.Split(v, " ")
		if len(things) == 0 {
			return nil, nil
		}
		if len(things) < 2 {
			return nil, ErrInvalidFmt
		}
		fileFlag := 0
		switch things[0] {
		case "r":
			fileFlag |= syscall.O_RDONLY
		case "w":
			fileFlag |= syscall.O_WRONLY
			fileFlag |= syscall.O_APPEND
		case "wr":
			fallthrough
		case "rw":
			fileFlag |= syscall.O_RDWR
		case "c":
			fileFlag |= syscall.O_CREAT
		case "cw":
			fallthrough
		case "wc":
			fileFlag |= syscall.O_CREAT
			fileFlag |= syscall.O_WRONLY
		default:
			return nil, ErrInvalidFlag
		}
		switch things[1] {
		case IOTypeUnixSocket:
		case IOTypeUnixConn:
		case IOTypeInetSocket:
		case IOTypeInetConn:
		case IOTypeFile:
			cur, err := os.OpenFile(things[2], fileFlag, 0666)
			if err != nil {
				return nil, err
			}
			o[i] = cur
		}
	}
	return o, nil
}
