package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

func ReadConfig(r io.Reader, enc string) (*Unit, error) {
	enc = enc[1:]
	var (
		err error
		s   *Unit = new(Unit)
		bs  []byte
	)

	switch enc {
	default:
		return nil, fmt.Errorf("invalid encoding: %s", enc)
	case "json":
		err = json.NewDecoder(r).Decode(s)
		if err != nil {
			return nil, err
		}
	case "yaml":
		bs, err = ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(bs, s)
		if err != nil {
			return nil, err
		}
	}
	return s, nil
}

func WriteConfig(w io.Writer, enc string, s *Unit) error {
	var (
		err error
		bs  []byte
	)

	switch enc {
	case "json":
		err = json.NewEncoder(w).Encode(s)
	case "yaml":
		bs, err = yaml.Marshal(s)
		if err != nil {
			return err
		}
		l := len(bs)
		Wn := 0
		n := 0
		for Wn < l {
			n, err = w.Write(bs[Wn:])
			if err != nil {
				return err
			}
			Wn += n
		}
	}
	return err
}

func ReadConfigSet(dir string) ([]*Unit, error) {
	s, err := readConfigs(dir)

	return s, err
}

func readConfigs(dir string) ([]*Unit, error) {
	o := make([]*Unit, 0, 2)
	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, f := range fs {
		if f.IsDir() {
			subUnits, err := readConfigs(dir + "/" + f.Name())
			if err != nil {
				fmt.Println(err)
			}
			o = append(o, subUnits...)
			continue
		}
		Logf("got %s\n", dir+"/"+f.Name())
		file, err := os.Open(dir + "/" + f.Name())
		if err != nil {
			fmt.Println(err)
			continue
		}

		s, err := ReadConfig(file, path.Ext(file.Name()))
		if err != nil {
			fmt.Println(err)
			continue
		}
		o = append(o, s)
	}
	return o, nil
}
