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
	"fmt"
	"os"
	"strings"

	"github.com/envsecrets/envsecrets/cli/commons"
	"github.com/envsecrets/envsecrets/cli/config"
	configCommons "github.com/envsecrets/envsecrets/cli/config/commons"
	"github.com/envsecrets/envsecrets/internal/keys"
	"github.com/envsecrets/envsecrets/internal/secrets"
	secretsCommons "github.com/envsecrets/envsecrets/internal/secrets/commons"
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
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		key := args[0]

		//	Auto-capitalize the key
		key = strings.ToUpper(key)

		decryptedOrgKey, err := keys.DecryptOrganisationKey(commons.KeysConfig.Public, commons.KeysConfig.Private, commons.ProjectConfig.OrgKey)
		if err != nil {
			log.Debug(err.Error)
			log.Fatal(err.Message)
		}

		//	Get the values from Hasura.
		getOptions := secretsCommons.GetSecretOptions{
			EnvID: commons.ProjectConfig.Environment,
			Key:   key,
		}

		if version > -1 {
			getOptions.Version = &version
		}

		secret, err := secrets.Get(commons.DefaultContext, commons.GQLClient, &getOptions)
		if err != nil {
			log.Debug(err.Error)
			log.Fatal(err.Message)
		}

		for _, item := range secret.Data {

			if item.Value != nil {

				if item.Type == secretsCommons.Ciphertext {

					//	Base64 decode the secret value
					decoded, er := base64.StdEncoding.DecodeString(item.Value.(string))
					if er != nil {
						log.Debug(er)
						log.Fatal("Failed to base64 decode the value for ", key)
					}

					//	Decrypt the value using org-key.
					decrypted, err := keys.OpenSymmetrically(decoded, decryptedOrgKey)
					if err != nil {
						log.Debug(err.Error)
						log.Fatal(err.Message)
					}

					item.Value = string(decrypted)
				}

				fmt.Printf("%v", item.Value)
				fmt.Println()
			}
		}
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
	getCmd.Flags().StringVarP(&XTokenHeader, "token", "t", "", "Environment Token")
}
