package main

import (
	"fmt"
	"os"
	"testing"
	"time"
)

var set = []*Unit{
	&Unit{
		Name:     "A",
		Provides: []string{"A"},
		Needs:    []string{
			//nothing
		},
	},
	&Unit{
		Name:     "B",
		Provides: []string{"B"},
		Needs: []string{
			"A",
		},
	},
	&Unit{
		Name:     "C",
		Provides: []string{"C"},
		Needs: []string{
			"B",
		},
	},
	&Unit{
		Name:     "D",
		Provides: []string{"D"},
		Needs: []string{
			"B", "C",
		},
	},
	&Unit{
		Name:     "E",
		Provides: []string{"E"},
		Needs: []string{
			"D", "A",
		},
	},
	&Unit{
		Name:     "F",
		Provides: []string{"F"},
		Needs: []string{
			"D", "C",
		},
	},
}

func TestDepGraph(t *testing.T) {
	MakeDepGraph(UnitSetToDepSlice(set))
}

func TestWriteDepGraph(t *testing.T) {
	prefix := "test_depgraph/"
	for _, v := range set {
		v.Executable = "/bin/sleep 4"
		v.Needs = []string{"tools"}
		v.Type = UnitTypeTask
		f, err := os.Create(prefix + v.Name + ".yaml")
		if err != nil {
			t.Error(err)
			continue
		}
		WriteConfig(f, "yaml", v)
		f.Close()
	}
}

type testDep struct {
	n  []string
	p  []string
	nm string
}

func (t *testDep) ServiceName() string {
	return t.nm
}

func (t *testDep) Start() error {
	time.Sleep(1 * time.Second)
	fmt.Println("started", t.nm)
	return nil
}

func (t *testDep) RequiresServices() []string {
	return t.n
}

func (t *testDep) ProvidesServices() []string {
	return t.p
}

var tGraph = []*testDep{
	&testDep{
		nm: "A",
		p:  []string{"A"},
		n:  []string{
			//nothing
		},
	},
	&testDep{
		nm: "B",
		p:  []string{"B"},
		n: []string{
			"A",
		},
	},
	&testDep{
		nm: "C",
		p:  []string{"C"},
		n: []string{
			"B",
		},
	},
	&testDep{
		nm: "D",
		p:  []string{"D"},
		n: []string{
			"B", "C",
		},
	},
	&testDep{
		nm: "E",
		p:  []string{"E"},
		n: []string{
			"D", "A",
		},
	},
	&testDep{
		nm: "F",
		p:  []string{"F"},
		n: []string{
			"D", "C",
		},
	},
}

func TestOrder(t *testing.T) {
	depL := make([]depUnit, len(tGraph))
	for i, v := range tGraph {
		depL[i] = v
	}
	g := MakeDepGraph(depL)
	g.Start()
}
