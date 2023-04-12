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
	"net/http"
	"os"

	"github.com/envsecrets/envsecrets/cli/commons"
	"github.com/envsecrets/envsecrets/config"
	configCommons "github.com/envsecrets/envsecrets/config/commons"
	"github.com/envsecrets/envsecrets/internal/auth"
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/spf13/cobra"
)

var version int
var exportfile string

var XTokenHeader string
var XOrgIDHeader string

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Prints decrypted list of your environment's (key-value) secret pairs",
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
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		secretPayload, err := export(nil)
		if err != nil {
			log.Fatal(err)
		}

		log.Debug("Fetched secret version ", secretPayload["version"])

		if secretPayload["data"] != nil {

			for key, item := range secretPayload["data"].(map[string]interface{}) {
				payload := item.(map[string]interface{})

				//	If the value is empty/nil,
				//	then it either doesn't exist or wasn't fetched.
				if payload["value"] == nil {
					log.Fatal("Values not found for key: ", key)
				}

				//	Base64 decode the secret value
				value, err := base64.StdEncoding.DecodeString(payload["value"].(string))
				if err != nil {
					log.Debugf("key: %s; value %v", key, payload["value"])
					log.Debug(err)
					log.Fatal("Failed to base64 decode the secret value")
				}

				fmt.Printf("%s=%v", key, string(value))
				fmt.Println()
			}

		}
	},
}

func export(key *string) (map[string]interface{}, error) {

	var secretVersion *int

	if version > -1 {
		secretVersion = &version
	}

	//	Load the project config
	projectConfigPayload, err := config.GetService().Load(configCommons.ProjectConfig)
	if err != nil {
		return nil, err
	}

	projectConfig := projectConfigPayload.(*configCommons.Project)

	req, err := http.NewRequestWithContext(commons.DefaultContext, http.MethodGet, commons.API+"/v1/secrets", nil)
	if err != nil {
		return nil, err
	}

	//	Initialize the query params.
	query := req.URL.Query()
	if key != nil {
		query.Set("key", *key)
	}
	if secretVersion != nil {
		query.Set("version", fmt.Sprint(*secretVersion))
	}

	//	Initialize the HTTP Client
	client := commons.HTTPClient

	//	If the environment secret is passed,
	//	create a new HTTP client and attach it in the header.
	if XTokenHeader != "" {

		//	Validate the mandatory `x-org-id` header.
		if XOrgIDHeader == "" {
			return nil, errors.New("Passing the --org-id flag is mandatory with environment token")
		}

		client = clients.NewHTTPClient(&clients.HTTPConfig{
			BaseURL: commons.API + "/v1",
			Logger:  commons.Logger,
			CustomHeaders: []clients.CustomHeader{
				{
					Key:   string(clients.TokenHeader),
					Value: XTokenHeader,
				},
				{
					Key:   string(clients.OrgIDHeader),
					Value: XOrgIDHeader,
				},
			},
		})

	} else {
		query.Set("org_id", projectConfig.Organisation)
		query.Set("env_id", projectConfig.Environment)
	}

	req.URL.RawQuery = query.Encode()

	var response commons.APIResponse
	if err := client.Run(commons.DefaultContext, req, &response); err != nil {
		return nil, errors.New(err.Message)
	}

	if response.Error != "" {
		log.Debug(response.Error)
		return nil, errors.New(response.Message)
	}

	return response.Data.(map[string]interface{}), nil
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
	exportCmd.Flags().StringVar(&XOrgIDHeader, "org-id", "", "Organisation ID; mandatory with environment token")
}
