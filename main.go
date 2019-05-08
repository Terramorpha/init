package main

import (
	"bufio"
	"fmt"
	"image/color"
	"io/ioutil"
	"os"
	"path"
	"syscall"

	"github.com/Terramorpha/utils"
)

//UnitsFolder is the place where all the unit files will live
const UnitsFolder = "/etc/units"

var SystemLog *Log = NewLog()

func main() {
	SystemLog.AddSubscriber(os.Stdout)

	p := os.Getuid()
	Logf("current uid: %d\n", p)
	Logf("started, line %s\n", getLine())
	var (
		err error
	)
	Logf("STARTED INIT!!!\n\n\n")
	err = syscall.Mount("proc", "/proc", "proc", 0, "")
	if err != nil {
		Logf("err mounting proc:%v\n", err)
	}
	err = syscall.Mount("sys", "/sys", "sysfs", 0, "")
	if err != nil {
		Logf("err mounting sys:%v\n", err)
	}
	s, err := ReadConfigSet(UnitsFolder)
	if err != nil {
		Logf("%s\n", err)
	}
	go Reap(s)

	depUnitSet := make([]depUnit, len(s))
	for i, v := range s {
		depUnitSet[i] = v
	}
	depgraph := MakeDepGraph(depUnitSet)
	depgraph.Start()
	for _, v := range s {
		fmt.Println(v)
	}

	shell(s)
}

func run(name string, args ...string) error {
	trueargs := append([]string{name}, args...)
	p, err := os.StartProcess(name, trueargs[1:], &os.ProcAttr{
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
		Sys: &syscall.SysProcAttr{
			Pgid:    0,
			Setpgid: true,
		},
	})
	if err != nil {
		return err
	}

	_, err = p.Wait()
	return err
}

func shell(svcs []*Unit) {

	b := bufio.NewReader(os.Stdin)
	col := color.RGBA{
		R: 0,
		G: 255,
		B: 0,
		A: 255,
	}
Loop:
	for {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		prompt := utils.SetForeground(fmt.Sprintf("%s $ ", cwd), col)
		fmt.Print(prompt)
		s, _, err := b.ReadLine()
		if err != nil {
			panic(err)
		}
		strs := utils.ParseLine(string(s), ' ')
		//fmt.Println(strs)
		if len(strs) == 0 {
			continue
		}
		switch strs[0] {
		case "svc":
			if len(strs) == 1 {
				fmt.Print(svcs)
				continue

			}
			s := GetUnit(svcs, strs[1])
			if s == nil {
				fmt.Println("no such service")
				continue
			}
			if len(strs) == 2 {
				fmt.Println(s)
				continue
			}
			switch strs[2] {
			case "log":
				fmt.Print(string(s.runtime.log.buffer))
			case "v":
				fmt.Printf("%#+v\n", s)
			}
		case "cd":
			if len(strs) != 2 {
				os.Chdir("/")
				continue
			}
			os.Chdir(strs[1])
		case "exit":
			break Loop
		case "ls":
			s, err := ioutil.ReadDir(cwd)
			if err != nil {
				fmt.Println(err)
			}
			for _, v := range s {
				Logf("%s %v %v\n", v.Name(), v.Mode(), v.Size())
			}
		default:
			if path.Base(strs[0]) == strs[0] { //path
				strs[0] = "/bin/" + strs[0]
			}
			err = run(strs[0], strs[1:]...)
			if err != nil {
				fmt.Println(err)
			}
		}
		continue

	}
}
