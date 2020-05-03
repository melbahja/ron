package main


import (
	"os"
	"io"
	"os/exec"
)


type Cmd struct {
	Command   string
	Arguments []string
	Stdin     io.Reader
	Stdout    io.Writer
	Stderr    io.Writer
	Environ   []string
}

type CmdResult struct {
	Exit int
	Err  error
}

func run(c Cmd) CmdResult {

	cmd := exec.Command(c.Command, c.Arguments...)
	cmd.Stdin = c.Stdin
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	cmd.Env = c.Environ

	err := cmd.Run()

	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			return CmdResult{
				Exit: e.ExitCode(),
				Err:  err,
			}
		}
	}

	return CmdResult{
		Exit: 0,
		Err:  err,
	}
}



func main() {


	cmd := Cmd{
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		Environ: os.Environ(),
		Command: ".run/" + os.Args[1],
	}

	run(cmd)

}







