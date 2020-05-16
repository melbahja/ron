package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type (
	Command struct {
		Cmd  string
		Args []string
	}
)

var flags = []cli.Flag{
	&cli.StringFlag{
		Name:  "dir",
		Usage: "Working directory `path`",
		Value: os.ExpandEnv("$PWD"),
	},
	&cli.BoolFlag{
		Name:  "dry",
		Usage: "Show the file that will be executed",
	},
}

func main() {

	app := &cli.App{
		Name:                 "Ron",
		Usage:                "The simple task runner.",
		Flags:                flags,
		Action:               handle,
		BashComplete:         complete,
		EnableBashCompletion: true,
	}

	app.Run(os.Args)
}

func handle(c *cli.Context) error {

	var (
		err error
		cmd *Command
		dir string = fmt.Sprintf("%s/.ron", c.String("dir"))
	)

	_, err = os.Stat(dir)
	check(err)

	cmd, err = getCmd(dir, c.Args().Slice())
	check(err)

	if c.Bool("dry") {

		fmt.Printf("will run: %s %s\n", cmd.Cmd, strings.Join(cmd.Args, " "))
		return nil
	}

	check(run(c.String("dir"), cmd))

	return nil
}

func getCmd(dir string, args []string) (*Command, error) {

	var (
		err  error
		stat os.FileInfo
		cmd  = Command{
			Cmd: dir,
		}
	)

	for i, arg := range args {

		if arg[0] == 0x2d {
			cmd.Args = args[i:]
			break
		}

		cmd.Cmd += fmt.Sprintf("/%s", arg)
	}

	if stat, err = os.Stat(cmd.Cmd); err != nil {
		return nil, err
	}

	if stat.IsDir() {

		cmd.Cmd = fmt.Sprintf("%s/.default", cmd.Cmd)
		if _, err = os.Stat(cmd.Cmd); err != nil {
			return nil, err
		}
	}

	return &cmd, nil
}

func complete(c *cli.Context) {

	if c.NArg() > 0 {
		return
	}

	dir := fmt.Sprintf("%s/.ron/", c.String("dir"))
	files, err := filepath.Glob(fmt.Sprintf("%s*", dir))
	check(err)

	for _, f := range files {
		fmt.Println(strings.Split(strings.Replace(f, dir, "", 1), "/")[0])
	}
}

func run(dir string, c *Command) error {
	cmd := exec.Command(c.Cmd, c.Args...)
	cmd.Env = os.Environ()
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
