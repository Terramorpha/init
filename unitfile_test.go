package main

import (
	"os"
	"testing"
)

const testPrefix = "test_configs/"

func TestWriteConfig(t *testing.T) {
	config := Unit{
		Executable: "/bin/echo",
		Name:       "echoServer",
		Needs:      []string{"tools"},
		WorkingDir: "/",
		Files: map[string]string{
			"io": "rw unixconn terramorpha.tech:80",
		},
		Fds: []string{
			"io",
			"io",
			"",
		},
	}

	for _, enc := range []string{"yaml", "json"} {
		f, err := os.Create(testPrefix + "config." + enc)
		if err != nil {
			t.Fatalf("err creating file for %s: %v\n", enc, err)
		}
		err = WriteConfig(f, enc, &config)
		if err != nil {
			t.Fatalf("err encoding %s: %v\n", enc, err)
		}
	}
}
