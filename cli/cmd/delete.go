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
	"net/http"
	"os"
	"strings"

	"github.com/envsecrets/envsecrets/cli/commons"
	"github.com/envsecrets/envsecrets/config"
	configCommons "github.com/envsecrets/envsecrets/config/commons"
	"github.com/envsecrets/envsecrets/internal/auth"
	secretsCommons "github.com/envsecrets/envsecrets/internal/secrets/commons"
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [KEY]",
	Short: "Deletes a key-value pair from your current environment's secrets",
	PreRun: func(cmd *cobra.Command, args []string) {

		//	If the user is not already authenticated,
		//	log them in first.
		if !auth.IsLoggedIn() {
			loginCmd.Run(cmd, args)
		}

		//	Ensure the project configuration is initialized and available.
		if !config.GetService().Exists(configCommons.ProjectConfig) {
			log.Error("Can't read project configuration")
			log.Info("Initialize your current directory with `envsecrets init`")
			os.Exit(1)
		}

	},
	Run: func(cmd *cobra.Command, args []string) {

		//	Run sanity checks
		if len(args) != 1 {
			log.Fatal("Invalid key format")
		}
		key := args[0]

		//	Auto-capitalize the key
		key = strings.ToUpper(key)

		var secretVersion *int

		if version > -1 {
			secretVersion = &version
		}

		//	Load the project configuration
		projectConfigData, er := config.GetService().Load(configCommons.ProjectConfig)
		if er != nil {
			log.Debug(er)
			log.Fatal("Failed to fetch project configuration")
		}

		projectConfig := projectConfigData.(*configCommons.Project)

		//	Send the secrets to vault
		payload := secretsCommons.DeleteRequestOptions{
			EnvID:   projectConfig.Environment,
			Key:     key,
			Version: secretVersion,
		}

		reqBody, err := payload.Marshal()
		if err != nil {
			log.Debug(err)
			log.Fatal("Failed to prepare request body")

		}

		req, err := http.NewRequestWithContext(commons.DefaultContext, http.MethodDelete, commons.API+"/v1/secrets", bytes.NewBuffer(reqBody))
		if err != nil {
			log.Debug(err)
			log.Fatal("Failed to prepare request")

		}

		//	Set content-type header
		req.Header.Set("content-type", "application/json")

		var response commons.APIResponse
		if err := commons.HTTPClient.Run(commons.DefaultContext, req, &response); err != nil {
			log.Debug(err)
			log.Fatal("Failed to complete the request")
		}

		if response.Error != "" {
			log.Debug(err)
			log.Fatal("Failed to delete: ", key)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
