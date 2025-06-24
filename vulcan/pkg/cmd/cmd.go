package cmd

import (
	"bytes"
	"io"
	"os"
	"os/exec"

	"github.com/go-pantheon/fabrica-util/errors"
)

func CmdExecute(dir, name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	cmd.Dir = dir

	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	var errStdout, errStderr error
	var stdoutBuf, stderrBuf bytes.Buffer
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)

	err := cmd.Start()
	if err != nil {
		return "", errors.Wrapf(err, "command start failed. name:%s, args:%v", name, arg)
	}

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()
	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
	}()

	err = cmd.Wait()
	if err != nil {
		return "", errors.Wrapf(err, "command wait failed. name:%s, args:%v", name, arg)
	}
	if errStdout != nil || errStderr != nil {
		return "", errors.Wrapf(err, "command output failed. name:%s, args:%v", name, arg)
	}

	outStr, errStr := stdoutBuf.String(), stderrBuf.String()
	if errStr != "" {
		return "", errors.Wrapf(err, "command output failed. name:%s, args:%v %s", name, arg, errStr)
	}
	// slog.Info("Command executed.", "name", name, "args", arg, "out", outStr)
	return outStr, nil
}
