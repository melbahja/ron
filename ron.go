package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

var (
	ver   string
	flags = []cli.Flag{
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
)

func main() {

	app := &cli.App{
		Name:    "Ron",
		Usage:   "The simple task runner.",
		Version: ver,
		Flags:   flags,
		Action:  handle,
		Authors: []*cli.Author{
			{
				Name:  "Mohamed Elbahja",
				Email: "bm9qdW5r@gmail.com",
			},
		},
		BashComplete:         complete,
		EnableBashCompletion: true,
	}

	app.Run(os.Args)
}

func handle(c *cli.Context) error {

	var (
		err     error
		cmds    []string
		ronFile string = fmt.Sprintf("%s/ron.yaml", c.String("dir"))
	)

	if _, err = os.Stat(ronFile); err == nil && c.Args().Len() == 1 {

		cmds, _ = getCmdFromRonFile(ronFile, c.Args().Get(0))
	}

	if cmds == nil {

		dir := fmt.Sprintf("%s/.ron", c.String("dir"))

		_, err = os.Stat(dir)
		check(err)

		cmd, err := getCmd(dir, c.Args().Slice())
		check(err)
		cmds = append(cmds, cmd)
	}

	if c.Bool("dry") {

		fmt.Println("This will run:")
		for cmd := range cmds {
			fmt.Printf(" - %s\n", cmds[cmd])
		}

		return nil
	}

	for cmd := range cmds {
		check(run(c.String("dir"), cmds[cmd]))
	}

	return nil
}

func getCmd(dir string, args []string) (cmd string, err error) {

	var stat os.FileInfo

	cmd = dir

	for i, arg := range args {

		if arg[0] == 0x2d {
			cmd += strings.Join(args[i:], " ")
			break
		}

		cmd += fmt.Sprintf("/%s", arg)
	}

	if stat, err = os.Stat(cmd); err != nil {
		return
	}

	if stat.IsDir() {

		cmd = fmt.Sprintf("%s/.default", cmd)

		if _, err = os.Stat(cmd); err != nil {
			return
		}
	}

	return cmd, nil
}

func getCmdFromRonFile(name, cmd string) ([]string, error) {

	commands, err := parseRonFile(name)

	if err != nil {
		return nil, err
	}

	if cmds, ok := commands[cmd]; ok {

		return cmds, nil
	}

	return nil, fmt.Errorf("Could not find: %s", cmd)
}

func parseRonFile(name string) (map[string][]string, error) {

	data, err := ioutil.ReadFile(name)

	if err != nil {
		return nil, err
	}

	commands := make(map[string][]string)

	if err = yaml.Unmarshal([]byte(data), &commands); err != nil {
		return nil, err
	}

	return commands, nil
}

func complete(c *cli.Context) {

	dir := fmt.Sprintf("%s/.ron/%s", c.String("dir"), strings.Join(c.Args().Slice(), "/"))

	if strings.HasSuffix(dir, "/") == false {
		dir += "/"
	}

	options, err := filepath.Glob(fmt.Sprintf("%s*", dir))
	check(err)

	if c.Args().Len() < 1 {

		ronFile := filepath.Join(c.String("dir"), "ron.yaml")

		if _, err := os.Stat(ronFile); err == nil {

			if commands, err := parseRonFile(ronFile); err == nil {

				cmds := make([]string, 0, len(commands))
				for k := range commands {
					cmds = append(cmds, k)
				}
				options = append(options, cmds...)
			}
		}
	}

	for _, f := range options {
		fmt.Println(strings.Split(strings.Replace(f, dir, "", 1), "/")[0])
	}
}

func run(dir string, c string) error {
	cmd := exec.Command("bash", "-c", c)
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
