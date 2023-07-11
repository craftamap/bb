// Package run contains code from cli/cli
// As cli/cli uses an internal package for internal/run/run.go, we need to
// copy the code from them :(
// https://github.com/cli/cli/blob/34d549e7b61660c7c993181c0be046d6277cad03/internal/run/run.go
package run

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Runnable is typically an exec.Cmd or its stub in tests.
type Runnable interface {
	Output() ([]byte, error)
	Run() error
}

// PrepareCmd extends exec.Cmd with extra error reporting features and provides a
// hook to stub command execution in tests.
var PrepareCmd = func(cmd *exec.Cmd) Runnable {
	return &cmdWithStderr{cmd}
}

// SetPrepareCmd overrides PrepareCmd and returns a func to revert it back.
func SetPrepareCmd(fn func(*exec.Cmd) Runnable) func() {
	origPrepare := PrepareCmd
	PrepareCmd = func(cmd *exec.Cmd) Runnable {
		// normalize git executable name for consistency in tests
		if baseName := filepath.Base(cmd.Args[0]); baseName == "git" || baseName == "git.exe" {
			cmd.Args[0] = "git"
		}
		return fn(cmd)
	}
	return func() {
		PrepareCmd = origPrepare
	}
}

// cmdWithStderr augments exec.Cmd by adding stderr to the error message.
type cmdWithStderr struct {
	*exec.Cmd
}

func (c cmdWithStderr) Output() ([]byte, error) {
	if os.Getenv("DEBUG") != "" {
		_ = printArgs(os.Stderr, c.Cmd.Args)
	}
	if c.Cmd.Stderr != nil {
		return c.Cmd.Output()
	}
	errStream := &bytes.Buffer{}
	c.Cmd.Stderr = errStream
	out, err := c.Cmd.Output()
	if err != nil {
		err = &CmdError{errStream, c.Cmd.Args, err}
	}
	return out, err
}

func (c cmdWithStderr) Run() error {
	if os.Getenv("DEBUG") != "" {
		_ = printArgs(os.Stderr, c.Cmd.Args)
	}
	if c.Cmd.Stderr != nil {
		return c.Cmd.Run()
	}
	errStream := &bytes.Buffer{}
	c.Cmd.Stderr = errStream
	err := c.Cmd.Run()
	if err != nil {
		err = &CmdError{errStream, c.Cmd.Args, err}
	}
	return err
}

// CmdError provides more visibility into why an exec.Cmd had failed.
type CmdError struct {
	Stderr *bytes.Buffer
	Args   []string
	Err    error
}

func (e CmdError) Error() string {
	msg := e.Stderr.String()
	if msg != "" && !strings.HasSuffix(msg, "\n") {
		msg += "\n"
	}
	return fmt.Sprintf("%s%s: %s", msg, e.Args[0], e.Err)
}

func printArgs(w io.Writer, args []string) error {
	if len(args) > 0 {
		// print commands, but omit the full path to an executable
		args = append([]string{filepath.Base(args[0])}, args[1:]...)
	}
	_, err := fmt.Fprintf(w, "%v\n", args)
	return err
}
