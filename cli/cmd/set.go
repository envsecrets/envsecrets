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
	"os"
	"path/filepath"
	"strings"

	internalErrors "errors"

	"github.com/envsecrets/envsecrets/cli/commons"
	"github.com/envsecrets/envsecrets/cli/internal/secrets"
	"github.com/envsecrets/envsecrets/dto"
	"github.com/spf13/cobra"
)

var importFile string
var environmentName string

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set KEY=VALUE",
	Short: "Set new key=value pairs in your current environment's secret.",
	Long: `Set new key=value pairs in your current environment's secret.

You can also load your variables directly from files: envs set --file .env

NOTE: This command auto-capitalizes your keys.`,
	PreRun: func(cmd *cobra.Command, args []string) {

		//	Initialize the common secret.
		InitializeSecret(commons.Log)
	},
	Args: func(cmd *cobra.Command, args []string) error {

		if importFile == "" && len(args) != 1 {
			return errors.New("either an import file is required to load variables from or at least 1 key=value pair (of secret) is required")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		if importFile != "" {
			filedata, err := os.ReadFile(importFile)
			if err != nil {
				commons.Log.Debug(err)
				commons.Log.Fatal("Failed to read file: ", importFile)
			}

			switch filepath.Ext(importFile) {
			default:

				lines := strings.Split(string(filedata), "\n")

				for index, item := range lines {

					//	Clean the line.
					item = strings.TrimSpace(item)

					key, payload, err := readPair(item)
					if err != nil {
						commons.Log.Error("Error on line ", index, " of your file")
						commons.Log.Fatal(err)
					}
					commons.Secret.Set(key, payload)
				}

			case ".csv", ".json", ".yaml":
				commons.Log.Error("This file format is not yet supported")
				commons.Log.Info("Use `--help` for more information")
				os.Exit(1)

			}
		} else {

			//	Run sanity checks
			if len(args) < 1 {
				commons.Log.Fatal("Invalid key=value pair")
			}

			key, payload, err := readPair(args[0])
			if err != nil {
				commons.Log.Fatal(err)
			}

			commons.Secret.Set(key, payload)
		}

		//	Encrypt the values.
		Encrypt()

		if err := secrets.GetService().Set(commons.DefaultContext, commons.GQLClient.GQLClient, commons.Secret); err != nil {
			commons.Log.Debug(err)
			commons.Log.Fatal("Failed to set the secrets")
		}

		if commons.Secret.EnvID != "" {
			if commons.Secret.Version != nil {
				commons.Log.Infof("Secrets set! Latest version in remote `%s` is now %d ", environmentName, *commons.Secret.Version)
			}
		} else {
			commons.Log.Info("Secrets set in local environment!")
		}
	},
}

func readPair(data string) (string, *dto.Payload, error) {

	if !strings.Contains(data, "=") {
		return "", nil, internalErrors.New("invalid key=value pair")
	}

	pair := strings.Split(data, "=")

	if len(pair) != 2 {
		return "", nil, internalErrors.New("invalid key=value pair")
	}

	key := pair[0]
	value := pair[1]

	key = strings.ToUpper(key)

	return key, &dto.Payload{
		Value: value,
	}, nil
}

func init() {
	rootCmd.AddCommand(setCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	setCmd.Flags().StringVarP(&importFile, "file", "f", "", "Export secret key-values from a file {.env | .json | .yaml | .txt}")
	setCmd.Flags().StringVarP(&environmentName, "env", "e", "", "Remote environment to set the secrets in. Defaults to the local environment.")
}
