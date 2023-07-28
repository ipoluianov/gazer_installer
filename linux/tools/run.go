package tools

import (
	"bufio"
	"errors"
	"fmt"
	"os/exec"
)

func Run(cmd string, path string, args []string) error {
	command := exec.Command(cmd, args...)
	command.Dir = path
	stdout, err := command.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := command.StderrPipe()
	if err != nil {
		return err
	}

	if err := command.Start(); err != nil {
		return err
	}

	stdin := bufio.NewScanner(stdout)
	for stdin.Scan() {
		fmt.Println(stdin.Text())
	}

	stdinerr := bufio.NewScanner(stderr)
	for stdin.Scan() {
		fmt.Println(stdinerr.Text())
	}

	command.Wait()
	if command.ProcessState.ExitCode() != 0 {
		return errors.New("exit_code=" + fmt.Sprint(command.ProcessState.ExitCode()))
	}

	return nil
}
