package internal

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func ExecCommand(cmd *exec.Cmd, forwardSignals bool, onExit func()) (int, error) {
	if onExit != nil {
		// ensure the onExit handler is called, regardless of how/when we return
		defer onExit()
	}

	// signal handling logic adapted from aws-vault https://github.com/99designs/aws-vault/
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan)

	if err := cmd.Start(); err != nil {
		return 1, err
	}

	// handle all signals
	go func() {
		for {
			// When running with a TTY, user-generated signals (like SIGINT) are sent to the entire process group.
			// If we forward the signal, the child process will end up receiving the signal twice.
			if forwardSignals {
				// forward to process
				sig := <-sigChan
				cmd.Process.Signal(sig) // #nosec G104
			} else {
				// ignore
				<-sigChan
			}
		}
	}()

	if err := cmd.Wait(); err != nil {
		// ignore errors
		cmd.Process.Signal(os.Kill) // #nosec G104

		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode(), exitError
		}

		return 2, err
	}

	waitStatus, ok := cmd.ProcessState.Sys().(syscall.WaitStatus)
	if !ok {
		return 2, fmt.Errorf("Unexpected ProcessState type, expected syscall.WaitStatus, got %T", waitStatus)
	}
	return waitStatus.ExitStatus(), nil
}
