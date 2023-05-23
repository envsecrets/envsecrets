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
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/envsecrets/envsecrets/cli/commons"
	"github.com/envsecrets/envsecrets/cli/config"
	configCommons "github.com/envsecrets/envsecrets/cli/config/commons"
	"github.com/envsecrets/envsecrets/cli/internal"
	"github.com/envsecrets/envsecrets/internal/keys"
	"github.com/envsecrets/envsecrets/internal/secrets"
	secretsCommons "github.com/envsecrets/envsecrets/internal/secrets/commons"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run -- [command]",
	Short: "Run a command with secrets injected directly into your process",
	Example: `envs run -- YOUR_COMMAND
envs run --command "YOUR_COMMAND && YOUR_OTHER_COMMAND"`,
	PreRun: func(cmd *cobra.Command, args []string) {

		//	If the user has passed a token,
		//	avoid using email+password to authenticate them against the API.
		if XTokenHeader != "" {
			return
		}

		//	Ensure the project configuration is initialized and available.
		if !config.GetService().Exists(configCommons.ProjectConfig) {
			log.Error("Can't read project configuration")
			log.Info("Initialize your current directory with `envs init`")
			os.Exit(1)
		}

		//	If the account configuration doesn't exist,
		//	log-in the user first.
		if !config.GetService().Exists(configCommons.AccountConfig) {
			loginCmd.PreRunE(cmd, args)
			loginCmd.Run(cmd, args)
		}
	},
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

		var orgKey [32]byte
		decryptedOrgKey, err := keys.DecryptAsymmetricallyAnonymous(commons.KeysConfig.Public, commons.KeysConfig.Private, commons.ProjectConfig.OrgKey)
		if err != nil {
			log.Debug(err)
			log.Fatal("Failed to decrypt the organisation's encryption key")
		}
		copy(orgKey[:], decryptedOrgKey)

		//	Get the values from Hasura.
		getOptions := secretsCommons.GetSecretOptions{
			EnvID: commons.ProjectConfig.Environment,
		}

		if version > -1 {
			getOptions.Version = &version
		}

		secret, err := secrets.GetAll(commons.DefaultContext, commons.GQLClient, &getOptions)
		if err != nil {
			log.Debug(err)
			log.Fatal("Failed to fetch the secrets")

		}

		//	Initialize a new buffer to store key=value lines
		var variables []string
		if err := secret.Decrypt(orgKey); err != nil {
			log.Debug(err)
			log.Fatal("Failed to decrypt secrets")
		}

		for key, item := range secret.Data {

			//	Base64 decode the secret value
			if err := item.Decode(); err != nil {
				log.Debug(err)
				log.Fatal("Failed to base64 decode the value for ", key)
			}

			variables = append(variables, fmt.Sprintf("%s=%s", key, item.Value))
		}

		log.Info("Injecting secrets version ", *secret.Version, " into your process")

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

		exitCode, err := internal.ExecCommand(userCmd, false, nil)
		if err != nil {
			log.Debug(err)
			log.Fatal("Command execution failed or completed ungracefully")
		}

		os.Exit(exitCode)
	},
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
	//runCmd.Flags().StringVarP(&XTokenHeader, "token", "t", "", "Environment Token")
}
