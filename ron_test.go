package main

import (
	"io/ioutil"
	"os"
	"testing"
)

var (
	err error
	tmp string
	cmd string = `#!/usr/bin/bash
echo "$0";
`
)

func TestRonAll(t *testing.T) {

	tmp, err = ioutil.TempDir("", "ron")
	ck(err)
	defer os.RemoveAll(tmp)

	ck(os.MkdirAll(tmp+"/.ron/sub", 0777))

	mkfile(".default", cmd)
	mkfile("command", cmd)
	mkfile("sub/.default", cmd)
	mkfile("sub/command", cmd)

	t.Run("default", cmdFunc([]string{}))
	t.Run("root command", cmdFunc([]string{"command"}))
	t.Run("sub command default", cmdFunc([]string{"sub"}))
	t.Run("sub command", cmdFunc([]string{"sub", "command"}))

}

func cmdFunc(args []string) func(t *testing.T) {

	return func(t *testing.T) {
		cmd, e := getCmd(tmp+"/.ron", args)
		if e != nil {
			t.Error(e)
		}

		if e = run(tmp, cmd); e != nil {
			t.Error(e)
		}
	}
}

func rootCmd(t *testing.T) {

	cmd, e := getCmd(tmp+"/.ron", []string{"sub"})
	if e != nil {
		t.Error(e)
	}

	if e = run(tmp, cmd); e != nil {
		t.Error(e)
	}
}

func mkfile(f string, content string) {

	ck(ioutil.WriteFile(tmp+"/.ron/"+f, []byte(content), 0777))
}

func ck(err error) {
	if err != nil {
		panic(err)
	}
}
