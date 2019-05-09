package sh

import (
	"testing"
)

func TestPWD(t *testing.T) {
	v, err := LoadShell(Shell{
		Name:    "test_pwd.sh",
		Content: `pwd()`,
	})
	if err != nil {
		t.Fatal(err)
	}
	v.Run()
}

func TestEcho(t *testing.T) {
	v, err := LoadShell(Shell{
		Name:    "test_echo.sh",
		Content: `echo("Hello World!")`,
	})
	if err != nil {
		t.Fatal(err)
	}
	v.Run()
}

// TODO: test_cd.sh (cd&pwd)
