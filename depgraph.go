package main

import (
	"fmt"
	"sync"
)

type depUnit interface {
	ServiceName() string
	Start() error
	RequiresServices() []string
	ProvidesServices() []string
}

type DepGraph struct {
	Nodes []*depNode
}

type depNode struct {
	parents  []*depNode
	children []*depNode
	donechan chan interface{}
	unit     depUnit
}

//MakeDepGraph
//assumes no recursion
//be careful when creating  complex
//dependency graphes
func MakeDepGraph(set []depUnit) *DepGraph {
	nds := make([]*depNode, len(set))
	dg := &DepGraph{
		Nodes: nds,
	}

	for i, unit := range set {
		nds[i] = &depNode{
			unit: unit,
		}
	}

	for _, n := range nds { //for each depNode, get all it's needs
		for _, v := range nds { //and update links
			if ShareStuff(v.unit.ProvidesServices(), n.unit.RequiresServices()) {
				n.parents = append(n.parents, v)
				v.children = append(v.children, n)
			}
		}

	}
	//created tree,
	//need buffered chans

	for _, v := range nds {
		v.donechan = make(chan interface{}, len(v.children))
	}

	fmt.Println("depgraph:")
	for _, v := range nds {
		fmt.Printf("%v\n", v)
	}
	return dg
}

func (d *depNode) String() string {
	pn := []string{}
	cn := []string{}
	for i := range d.parents {
		pn = append(pn, d.parents[i].unit.ServiceName())
	}
	for i := range d.children {
		cn = append(cn, d.children[i].unit.ServiceName())
	}
	o := ""
	o = fmt.Sprintf("%s: p: %v c: %v", d.unit.ServiceName(), pn, cn)
	return o
}

func (d *depNode) GoString() string {
	return fmt.Sprintf("%s: p: %v c: %v", d.unit.ProvidesServices()[0], d.parents, d.children)
}

func (d *DepGraph) Start() {
	var err error
	wg := sync.WaitGroup{}
	wg.Add(len(d.Nodes))
	for i := range d.Nodes {
		go func(nds *depNode) {
			//wait for dependencies to finish
			for _, parent := range nds.parents {
				<-parent.donechan
			}
			err = nds.unit.Start()
			if err != nil {
				Logf("error starting process %s: %v\n", nds.unit.ServiceName(), err)
			}
			for range nds.children {
				nds.donechan <- "done"
			}
			wg.Done()
			Logf("started %s\n", nds.unit.ServiceName())
		}(d.Nodes[i])
	}
	wg.Wait()
	Logf("everything was started\n")
}

func ShareStuff(slice []string, slice2 []string) bool {
	for _, v := range slice {
		for _, b := range slice2 {
			if v == b {
				return true
			}
		}
	}
	return false

}

func UnitSetToDepSlice(tGraph []*Unit) []depUnit {
	depL := make([]depUnit, len(tGraph))
	for i, v := range tGraph {
		depL[i] = v
	}
	return depL
}
