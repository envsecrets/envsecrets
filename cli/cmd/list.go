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
	"fmt"
	"os"

	"github.com/envsecrets/envsecrets/cli/commons"
	"github.com/envsecrets/envsecrets/cli/internal/secrets"
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List only the keys of your secrets",
	PreRun: func(cmd *cobra.Command, args []string) {

		//	Initialize the common secret.
		InitializeSecret(commons.Log)
	},
	Run: func(cmd *cobra.Command, args []string) {

		options := secrets.ListOptions{
			EnvID: commons.Secret.EnvID,
		}

		if version > -1 {
			options.Version = &version
		}

		secrets, err := secrets.GetService().List(commons.DefaultContext, commons.GQLClient.GQLClient, &options)
		if err != nil {
			commons.Log.Debug(err)

			//	If the dotenv file is not found, skip the error.
			if os.IsNotExist(err) {
				return
			}

			if err.Error() == string(clients.ErrorTypeRecordNotFound) {
				commons.Log.Warn("You haven't set any secrets in this environment")
				commons.Log.Info("Use `envs set --help` for more information")
				os.Exit(1)
			}
			commons.Log.Fatal("Failed to list the secrets")
		}

		for _, key := range secrets.Keys() {
			fmt.Println(key)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	listCmd.Flags().IntVarP(&version, "version", "v", -1, "Version of your secret")
	listCmd.Flags().StringVarP(&environmentName, "env", "e", "", "Remote environment to set the secrets in. Defaults to the local environment.")
}
