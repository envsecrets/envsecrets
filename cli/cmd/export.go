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
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/envsecrets/envsecrets/cli/commons"
	"github.com/envsecrets/envsecrets/cli/config"
	configCommons "github.com/envsecrets/envsecrets/cli/config/commons"
	"github.com/envsecrets/envsecrets/cli/internal"
	"github.com/envsecrets/envsecrets/internal/keys"
	"github.com/envsecrets/envsecrets/internal/secrets"
	secretsCommons "github.com/envsecrets/envsecrets/internal/secrets/commons"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var version int
var exportfile string

var XTokenHeader string

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Prints decrypted list of your environment's (key=value) secret pairs",
	Args:  cobra.NoArgs,
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

		//	Load the user's email.
		accountConfig, err := config.GetService().Load(configCommons.AccountConfig)
		if err != nil {
			loginCmd.PreRunE(cmd, args)
			loginCmd.Run(cmd, args)
		} else {

			accountData := accountConfig.(*configCommons.Account)
			email = accountData.User.Email

			//	Log them in first.
			//	Take password input
			passwordPrompt := promptui.Prompt{
				Label: "Password",
				Mask:  '*',
			}

			password, err = passwordPrompt.Run()
			if err != nil {
				os.Exit(1)
			}

			loginCmd.Run(cmd, args)

			//	Re-initialize the commons
			commons.Initialize()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {

		var secret secretsCommons.GetResponse
		var orgKey [32]byte

		if XTokenHeader != "" {

			options := internal.GetValuesOptions{
				Token: XTokenHeader,
			}

			if version > -1 {
				options.Version = &version
			}

			result, err := internal.GetValues(commons.DefaultContext, commons.HTTPClient, &options)
			if err != nil {
				log.Debug(err)
				log.Fatal("Failed to fetch the secrets")
			}

			secret = *result

		} else {

			decryptedOrgKey, err := keys.DecryptAsymmetricallyAnonymous(commons.KeysConfig.Public, commons.KeysConfig.Private, commons.ProjectConfig.OrgKey)
			if err != nil {
				log.Debug(err)
				log.Fatal("Failed to decrypt the organisation encryption key")
			}
			copy(orgKey[:], decryptedOrgKey)

			//	Get the values from Hasura.
			getOptions := secretsCommons.GetSecretOptions{
				EnvID: commons.ProjectConfig.Environment,
			}

			if version > -1 {
				getOptions.Version = &version
			}

			result, err := secrets.GetAll(commons.DefaultContext, commons.GQLClient, &getOptions)
			if err != nil {
				log.Debug(err)
				log.Fatal("Failed to fetch the secrets")
			}

			secret = *result
		}

		//	Initialize a new buffer to store key=value lines
		var buffer bytes.Buffer
		var variables []string
		for key, item := range secret.Secrets {

			//	Base64 decode the secret value
			decoded, err := base64.StdEncoding.DecodeString(item.Value)
			if err != nil {
				log.Debug(err)
				log.Fatal("Failed to base64 decode the value for ", key)
			}

			if item.Type == secretsCommons.Ciphertext && XTokenHeader == "" {

				//	Decrypt the value using org-key.
				decrypted, err := keys.OpenSymmetrically(decoded, orgKey)
				if err != nil {
					log.Debug(err)
					log.Fatal("Failed to decrypt the secret")
				}

				item.Value = string(decrypted)
			} else {
				item.Value = string(decoded)
			}

			variables = append(variables, fmt.Sprintf("%s=%s", key, item.Value))
		}

		buffer.WriteString(strings.Join(variables, "\n"))

		if exportfile != "" {

			f, err := os.OpenFile(exportfile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
			if err != nil {
				log.Debug(err)
				log.Fatal("Failed to open file: ", exportfile)
			}

			defer f.Close()

			switch filepath.Ext(file) {
			default:
				if _, err := f.WriteString(buffer.String()); err != nil {
					log.Debug(err)
					log.Fatal("Failed to export values to file")
				}

			case ".csv":
				log.Error("This file format is not yet supported")
				log.Info("Use `--help` for more information")
				os.Exit(1)

			case ".json":

				/* 				var mapping map[string]interface{}
				   				for key, value := range mapping {

				   					//	Base64 encode the secret value
				   					value = base64.StdEncoding.EncodeToString([]byte(fmt.Sprint(value)))

				   					data[key] = secretsCommons.Payload{
				   						Value: value,
				   						Type:  secretsCommons.Ciphertext,
				   					}
				   				}
				*/
			case ".yaml":
			}
		} else {
			fmt.Println(buffer.String())
		}
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// exportCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	exportCmd.Flags().IntVarP(&version, "version", "v", -1, "Version of your secret")
	exportCmd.Flags().StringVarP(&exportfile, "file", "f", "", "Export secrets to a file {.json | .yaml | .txt}")
	exportCmd.Flags().StringVarP(&XTokenHeader, "token", "t", "", "Environment Token")
}
