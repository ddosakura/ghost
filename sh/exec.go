package sh

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// Exec command
func Exec(cmd *exec.Cmd, showCmd bool) error {
	if showCmd {
		fmt.Println("cd", cmd.Dir)
		fmt.Println(strings.Join(cmd.Args, " "))
	}

	return execCmd(cmd)
}

func execCmd(cmd *exec.Cmd) error {
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	//cmd.Stdout = os.Stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}
	reader := bufio.NewReader(stdout)
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Print(line)
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}
	if cmd.Process != nil {
		if err = cmd.Process.Kill(); err != nil && err.Error() != "os: process already finished" {
			return err
		}
	}
	return nil
}
