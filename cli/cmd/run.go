/*
Copyright © 2023 Mrinal Wahal <mrinalwahal@gmail.com>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice,
   this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors
   may be used to endorse or promote products derived from this software
   without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
POSSIBILITY OF SUCH DAMAGE.
*/
package cmd

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run -- [command]",
	Short: "Run a command with secrets injected into the environment",
	Example: `envsecrets run -- YOUR_COMMAND --YOUR-FLAG
envsecrets run --command "YOUR_COMMAND && YOUR_OTHER_COMMAND"`,
	Args: func(cmd *cobra.Command, args []string) error {
		// The --command flag and args are mututally exclusive
		usingCommandFlag := cmd.Flags().Changed("command")
		if usingCommandFlag {
			command := cmd.Flag("command").Value.String()
			if command == "" {
				return errors.New("--command flag requires a value")
			}

			if len(args) > 0 {
				return errors.New("arg(s) may not be set when using --command flag")
			}
		} else if len(args) == 0 {
			return errors.New("requires at least 1 arg(s), received 0")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		var variables []string

		//	`envsecrets run -- npm run dev`
		secretPayload, err := export(nil)
		if err != nil {
			log.Debug(err)
			log.Error("Failed to fetch all the secret values")
			return
		}

		for key, item := range secretPayload {
			payload := item.(map[string]interface{})

			//	Base64 decode the secret value
			value, err := base64.StdEncoding.DecodeString(payload["value"].(string))
			if err != nil {
				log.Debug(err)
				log.Error("Failed to base64 decode value for secrets: ", key)
				return
			}

			variables = append(variables, fmt.Sprintf("%s=%s", key, string(value)))
		}

		//	Overwrite reserved keys
		reservedKeys := []string{"PATH", "PS1", "HOME"}
		for _, item := range reservedKeys {
			variables = append(variables, fmt.Sprintf("%s=%s", item, os.Getenv(item)))
		}

		var userCmd *exec.Cmd

		if cmd.Flags().Changed("command") {
			shell := [2]string{"sh", "-c"}
			if runtime.GOOS == "windows" {
				shell = [2]string{"cmd", "/C"}
			} else {
				// these shells all support the same options we use for sh
				shells := []string{"/bash", "/dash", "/fish", "/zsh", "/ksh", "/csh", "/tcsh"}
				envShell := os.Getenv("SHELL")
				for _, s := range shells {
					if strings.HasSuffix(envShell, s) || strings.HasSuffix(envShell, "/bin"+s) {
						shell[0] = envShell
						break
					}
				}
			}
			userCmd = exec.Command(shell[0], shell[1], cmd.Flag("command").Value.String())
		} else {
			userCmd = exec.Command(args[0], args[1:]...)
		}

		userCmd.Env = variables
		userCmd.Stdin = os.Stdin
		userCmd.Stdout = os.Stdout
		userCmd.Stderr = os.Stderr

		exitCode, err := execCommand(userCmd, false, nil)
		if err != nil {
			log.Debugln(err)
			log.Errorln("command execution failed or completed ungracefully")
			os.Exit(1)
		}

		os.Exit(exitCode)
	},
}

func execCommand(cmd *exec.Cmd, forwardSignals bool, onExit func()) (int, error) {
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

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	runCmd.Flags().StringP("command", "c", "", "Command to run. Example: npm run dev")
}
