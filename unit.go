package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Terramorpha/utils"
	"golang.org/x/xerrors"
)

const (
	UnitTypeTask    = "task"
	UnitTypeService = "service"
)

//Unit represents a config for a Unit
type Unit struct {
	Name       string            `json:"name" yaml:"name"`
	Needs      []string          `json:"needs" yaml:"needs"`
	Provides   []string          `json:"provides" yaml:"provides"`
	Executable string            `json:"executable" yaml:"executable"`
	Type       string            `json:"type" yaml:"type"`
	WorkingDir string            `json:"workingdir" yaml:"workingdir"`
	Files      map[string]string `json:"files" yaml:"files"`
	Fds        []string          `json:"fds" yaml:"fds"`
	runtime    UnitRuntime
}

func (u *Unit) String() string {
	return fmt.Sprintf("%s needs %v, provides %v and uses %v", u.Name, u.Needs, u.Provides, u.Executable)
}

func (u *Unit) DoesProvide(ser string) bool {
	for _, v := range u.Provides {
		if ser == v {
			return true
		}
	}
	return false
}

func (u *Unit) ProvidesAny(needs []string) bool {
	for _, need := range needs {
		for _, service := range u.Provides {
			if need == service {
				return true
			}
		}
	}
	return false
}

func (u *Unit) Start() error {
	u.runtime.lock.Lock()
	u.runtime.Started = true
	u.runtime.log = NewLog()
	fs, err := UnitGetIOFiles(u.Files)
	if err != nil {
		return err
	}

	u.runtime.IOFiles = fs
	if u.Type == UnitTypeTask {
		u.runtime.IsDoneChan = make(chan struct{})
	}

	files, err := u.getFilesList(fs)
	if err != nil {
		return err
	}

	attrs := &os.ProcAttr{
		Files: files,
		Dir:   u.WorkingDir,
	}
	args := utils.ParseLine(u.Executable, ' ')
	p, err := os.StartProcess(args[0], args, attrs)
	if err != nil {
		return xerrors.Errorf("error starting %s: %w", u.Name, err)
	}
	u.runtime.Proc = p

	u.runtime.lock.Unlock()

	fmt.Fprintf(u.runtime.log, "[%v]: started\n", time.Now())
	if u.Type == UnitTypeTask {
		<-u.runtime.IsDoneChan
	}
	Logf("%s started: (%d)\n", u.Name, u.runtime.GetPid())
	return nil
}

func (u *Unit) Restart() {
	args := utils.ParseLine(u.Executable, ' ')

	fmt.Fprintf(u.runtime.log, "[%v]: stopped, restarting...\n", time.Now())
	files, err := u.getFilesList(u.runtime.IOFiles)
	if err != nil {
		Logf("error restarting %s: %v\n", u.Name, err)
		return
	}
	process, err := os.StartProcess(args[0], args, &os.ProcAttr{
		Dir:   u.WorkingDir,
		Files: files,
	})
	if err != nil {
		Logf("error starting service %s: %v\n", u.Name, err)
		return
	}
	u.runtime.Proc = process
}

func (u *Unit) RequiresServices() []string {
	return u.Needs
}

func (u *Unit) ProvidesServices() []string {
	return u.Provides
}
func (u *Unit) ServiceName() string {
	return u.Name
}

func GetUnit(list []*Unit, name string) *Unit {
	for _, v := range list {
		if v.Name == name {
			return v
		}
	}
	return nil
}

func (u *Unit) getFilesList(Files map[string]*os.File) ([]*os.File, error) {

	var files []*os.File
	if len(u.Fds) < 3 { //on crée une slice de fichiers de longueur 3 minimum
		files = make([]*os.File, 3)
	} else {
		files = make([]*os.File, len(u.Fds))
	}

	/*
		on fait une copie de fds (les nom de file handle)
		d'une grosseur égale au nombre de fds
	*/
	fileNames := make([]string, len(files))

	copy(fileNames, u.Fds) //on copie les noms
	//Logf("%s: %v\n", getLine(), fs)
	for i := range fileNames {
		//Logf("%s: filename %s\n", getLine(), fileNames[i])
		fileDescriptor, ok := Files[fileNames[i]]
		if (fileNames[i] == "") || !ok { //from null or to sys log
			if i == 0 {
				files[0] = nil
				continue
			}
			r, w, err := os.Pipe()
			if err != nil {
				return nil, err
			}
			files[i] = w
			go io.Copy(u.runtime.log, r)
			continue
		}
		files[i] = fileDescriptor
	}
	return files, nil
}
