package sh

import (
	"os/exec"
	"testing"

	"github.com/mitchellh/go-homedir"
)

func TestExec(t *testing.T) {
	home, err := homedir.Dir()
	if err != nil {
		t.Fatal(err)
	}
	c := exec.Command("cd .config")
	c.Dir = home
	execCmd(c)

	c = exec.Command("pwd")
	c.Dir = RootDirTmp
	execCmd(c)
}
