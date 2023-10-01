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
	"strings"

	"github.com/envsecrets/envsecrets/cli/commons"
	"github.com/envsecrets/envsecrets/cli/internal/secrets"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [KEY]",
	Short: "Fetch decrypted value corresponding to your secret key",
	PreRun: func(cmd *cobra.Command, args []string) {

		//	If the user has passed a token,
		//	avoid using email+password to authenticate them against the API.
		if XTokenHeader != "" {
			return
		}

		//	Initialize the common secret.
		InitializeSecret(log)
	},
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		key := args[0]

		//	Capitalize the key.
		key = strings.ToUpper(key)

		//	Fetch only the required values.
		getOptions := secrets.GetOptions{
			EnvID: commons.Secret.EnvID,
			Key:   key,
		}

		if version > -1 {
			getOptions.Version = &version
		}

		result, err := secrets.GetService().Get(commons.DefaultContext, commons.GQLClient, &getOptions)
		if err != nil {
			log.Debug(err)
			log.Fatal("Failed to fetch the value")
		}

		commons.Secret = result

		//	Decrypt and decode the common secret.
		DecryptAndDecode()

		if commons.Secret.Get(key) == nil {
			log.Fatal("No value found for this key")
		}

		fmt.Printf("%v", commons.Secret.Get(key).Value)
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	getCmd.Flags().IntVarP(&version, "version", "v", -1, "Version of your secret")
	getCmd.Flags().StringVarP(&environmentName, "env", "e", "", "Remote environment to set the secrets in. Defaults to the local environment.")
	// getCmd.Flags().StringVarP(&XTokenHeader, "token", "t", "", "Environment Token")
}
