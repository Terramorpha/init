package main

import "fmt"

type UnitSet []*Unit

func (s *UnitSet) GetProvider(provider string) *Unit {
	for _, v := range *s {
		if v.DoesProvide(provider) {
			return v
		}
	}
	return nil
}

func (s *UnitSet) StartProviders(ps []string) {
	for _, v := range ps {

		p := s.GetProvider(v)

		if p == nil {
			continue
		}
		if p.runtime.Started {
			continue
		}
		s.StartProviders(p.Needs)
		err := p.Start()
		if err != nil {
			Logf("error starting %s: %v\n", p.Name, err)
		}
	}
}

func (s *UnitSet) StartAll() {
	for _, v := range *s {
		if v.runtime.Started {
			continue
		}
		s.StartProviders(v.Needs)
		err := v.Start()
		if err != nil {
			Logf("error starting %s: %v\n", v.Name, err)
		}
	}
}

func (u *UnitSet) String() string {
	o := ""
	for _, v := range *u {
		o += fmt.Sprintln(v)
	}
	return o
}

func (u *UnitSet) GetUnit(x string) *Unit {
	for _, v := range *u {
		if v.Name == x {
			return v
		}
	}
	return nil
}
