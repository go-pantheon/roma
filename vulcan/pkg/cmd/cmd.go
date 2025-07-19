package cmd

import (
	"bytes"
	"io"
	"os"
	"os/exec"

	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/fabrica-util/xsync"
)

func CmdExecute(dir, name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	cmd.Dir = dir

	stdoutIn, err := cmd.StdoutPipe()
	if err != nil {
		return "", errors.Wrapf(err, "command stdout pipe failed. name:%s, args:%v", name, arg)
	}

	stderrIn, err := cmd.StderrPipe()
	if err != nil {
		return "", errors.Wrapf(err, "command stderr pipe failed. name:%s, args:%v", name, arg)
	}

	var (
		errStdout, errStderr error
		stdoutBuf, stderrBuf bytes.Buffer
	)

	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)

	if err := cmd.Start(); err != nil {
		return "", errors.Wrapf(err, "command start failed. name:%s, args:%v", name, arg)
	}

	xsync.Go("cmd.Execute.stdout", func() error {
		_, errStdout = io.Copy(stdout, stdoutIn)
		return nil
	})

	xsync.Go("cmd.Execute.stderr", func() error {
		_, errStderr = io.Copy(stderr, stderrIn)
		return nil
	})

	if err := cmd.Wait(); err != nil {
		return "", errors.Wrapf(err, "command wait failed. name:%s, args:%v", name, arg)
	}

	if errStdout != nil {
		return "", errors.Wrapf(errStdout, "command output. name:%s, args:%v", name, arg)
	}

	if errStderr != nil {
		return "", errors.Wrapf(errStderr, "command output failed. name:%s, args:%v", name, arg)
	}

	outStr, errStr := stdoutBuf.String(), stderrBuf.String()
	if errStr != "" {
		return "", errors.Errorf("command output failed. name:%s, args:%v %s", name, arg, errStr)
	}

	// slog.Info("Command executed.", "name", name, "args", arg, "out", outStr)
	return outStr, nil
}
