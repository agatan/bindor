package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/pkg/errors"
)

const usage = `Usage: bindor COMMAND [arg...]

COMMAND:
	build [PACKAGES...]    Build vendored packages
	exec [args...]         Execute a command with vendored binaries
`

func main() {
	status, err := run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(status)
}

func run(args []string) (int, error) {
	if len(args) < 2 {
		fmt.Println(usage)
		return 1, nil
	}
	switch args[1] {
	case "build":
		return build(args[2:])
	case "exec":
		return execute(args[2:])
	default:
		fmt.Println(usage)
		return 1, nil
	}
}

func binaryName(pack string) string {
	idx := strings.LastIndex(pack, "/")
	if idx == -1 {
		return pack
	}
	return pack[idx+1:]
}

func build(args []string) (int, error) {
	if len(args) < 1 {
		return 1, errors.New("no package name given")
	}
	for _, arg := range args {
		target := fmt.Sprintf("./vendor/%s", arg)
		cmd := exec.Command("go", "build", "-o", filepath.Join(".bindor", binaryName(target)), target)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return 1, errors.Wrapf(err, "failed to build %s", target)
		}
	}
	return 0, nil
}

func execute(args []string) (int, error) {
	if len(args) < 1 {
		return 1, errors.New("exec needs a command to run")
	}
	pwd, err := os.Getwd()
	if err != nil {
		return 1, errors.Wrap(err, "failed to get working directory")
	}
	if err := os.Setenv("PATH", fmt.Sprintf("%s:%s", filepath.Join(pwd, ".bindor"), os.Getenv("PATH"))); err != nil {
		return 1, errors.Wrap(err, "failed to set PATH environment")
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if s, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				return s.ExitStatus(), exitErr
			}
		}
		return 1, err
	}
	return 0, nil
}
