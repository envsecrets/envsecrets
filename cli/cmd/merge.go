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

	"github.com/envsecrets/envsecrets/cli/commons"
	"github.com/envsecrets/envsecrets/config"
	configCommons "github.com/envsecrets/envsecrets/config/commons"
	"github.com/envsecrets/envsecrets/internal/auth"
	"github.com/envsecrets/envsecrets/internal/environments"
	secretsCommons "github.com/envsecrets/envsecrets/internal/secrets/commons"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// mergeCmd represents the merge command
var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge secrets from a different environment into current one.",
	Long: `Merge secrets from a different environment into current one.

NOTE: This will overwrite the values of current latest version of secrets.
	
For precaution, this command generates a new version of your secret
with new/overrided values. You can safely rollback to older version
containing original/unedited values.`,
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

		if environmentID == "" {

			//	Load the project config
			projectConfigPayload, err := config.GetService().Load(configCommons.ProjectConfig)
			if err != nil {
				log.Debug(err)
				log.Fatal("Failed to read local project configuration")
			}

			projectConfig := projectConfigPayload.(*configCommons.Project)

			//	Fetch environments
			environmentsList, er := environments.List(commons.DefaultContext, commons.GQLClient, &environments.ListOptions{
				ProjectID: projectConfig.Project,
			})
			if er != nil {
				log.Debug(err)
				log.Fatal("Failed to fetch list of environments")
			}

			//	Remove the existing environment
			var envs []environments.Environment
			for _, item := range *environmentsList {
				if item.ID != projectConfig.Environment {
					envs = append(envs, item)
				}
			}

			if len(envs) == 0 {
				log.Error("You have no other environment in this project to merge from")
				log.Info("First create a new environment using `init` command")
				os.Exit(1)
			}

			var environmentsStringList []string
			for _, item := range envs {
				environmentsStringList = append(environmentsStringList, item.Name)
			}

			selection := promptui.Select{
				Label: "Source Environment To Merge From",
				Items: environmentsStringList,
			}

			index, _, err := selection.Run()
			if err != nil {
				os.Exit(1)
			}

			for itemIndex, item := range envs {
				if itemIndex == index {
					environmentID = item.ID
					break
				}
			}

		}

		//	Load the project configuration
		projectConfigData, er := config.GetService().Load(configCommons.ProjectConfig)
		if er != nil {
			log.Debug(er)
			log.Fatal("Failed to load project configuration")
		}

		projectConfig := projectConfigData.(*configCommons.Project)

		//	Send the secrets to vault
		payload := secretsCommons.MergeRequestOptions{
			OrgID:       projectConfig.Organisation,
			SourceEnvID: environmentID,
			TargetEnvID: projectConfig.Environment,
		}

		if version > -1 {
			payload.SourceVersion = &version
		}

		reqBody, err := payload.Marshal()
		if err != nil {
			log.Debug(err)
			log.Fatal("Failed to prepare request payload")
		}

		req, err := http.NewRequestWithContext(commons.DefaultContext, http.MethodPost, commons.API+"/v1/secrets/merge", bytes.NewBuffer(reqBody))
		if err != nil {
			log.Debug(err)
			log.Fatal("Failed to prepare the request")
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
			log.Fatal("Failed to merge secrets")
		}

		log.Info("Merge Complete! Created version ", response.Data.(map[string]interface{})["version"])
	},
}

func init() {
	rootCmd.AddCommand(mergeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mergeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	mergeCmd.Flags().IntVarP(&version, "version", "v", -1, "Secret version of your source environment")
	mergeCmd.Flags().StringVar(&environmentID, "source-env-id", "", "Environment ID to sync from")
}
