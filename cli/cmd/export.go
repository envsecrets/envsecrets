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

	"github.com/envsecrets/envsecrets/cli/commons"
	"github.com/envsecrets/envsecrets/cli/internal"
	"github.com/envsecrets/envsecrets/config"
	configCommons "github.com/envsecrets/envsecrets/config/commons"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var version int
var exportfile string

var XTokenHeader string

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Prints decrypted list of your environment's (key-value) secret pairs",
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
			log.Info("Initialize your current directory with `envsecrets init`")
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

		options := internal.GetValuesOptions{}

		if version > -1 {
			options.Version = &version
		}

		if XTokenHeader == "" {

			//	Load the project config
			projectConfigPayload, err := config.GetService().Load(configCommons.ProjectConfig)
			if err != nil {
				log.Debug(err)
				log.Error("Can't read project configuration")
				log.Info("Initialize your current directory with `envsecrets init`")
				os.Exit(1)
			}

			projectConfig := projectConfigPayload.(*configCommons.Project)
			options.EnvID = projectConfig.Environment

		} else {
			options.Token = XTokenHeader
		}

		secrets, err := internal.GetValues(commons.DefaultContext, commons.HTTPClient, &options)
		if err != nil {
			log.Debug(err.Error)
			log.Fatal(err.Message)
		}

		for key, item := range secrets.Data {

			//	Base64 decode the secret value
			value, err := base64.StdEncoding.DecodeString(item.Value.(string))
			if err != nil {
				log.Debug(err)
				log.Fatal("Failed to base64 decode the value for %s", key)
			}

			fmt.Printf("%s=%s", key, string(value))
			fmt.Println()
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
