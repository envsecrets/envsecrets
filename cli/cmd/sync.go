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
	"encoding/json"
	"net/http"
	"os"

	"github.com/envsecrets/envsecrets/cli/commons"
	"github.com/envsecrets/envsecrets/cli/internal/secrets"
	"github.com/envsecrets/envsecrets/internal/clients"
	environmentCommons "github.com/envsecrets/envsecrets/internal/environments/commons"
	"github.com/envsecrets/envsecrets/internal/events"
	"github.com/envsecrets/envsecrets/internal/integrations"
	"github.com/envsecrets/envsecrets/internal/secrets/pkg/keypayload"
	"github.com/envsecrets/envsecrets/internal/secrets/pkg/payload"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var integrationType string

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync --env [your-remote-environment-name]",
	Short: "Push your secrets to third-party services",
	Long: `This command decrypts your secrets on client side
and pushes them to the chosen third-party service that is activated on that environment.
For example, Github Actions, AWS Secrets Manager, etc.

You can activate your connected integrations on the "integrations" page of your dashboard.`,
	PreRun: func(cmd *cobra.Command, args []string) {

		//	Initialize the common secret.
		InitializeSecret(log)
	},
	Run: func(cmd *cobra.Command, args []string) {

		var err error

		//	Fetch only the required values.
		getOptions := secrets.GetOptions{
			EnvID: commons.Secret.EnvID,
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

		//	Copy the dto.KPMap to keypayload.KPMap
		kpMap := keypayload.KPMap{}
		for key, value := range commons.Secret.Data.GetMapping() {
			kpMap[key] = &payload.Payload{
				Value: value.GetValue(),
			}
		}

		//	Encode all the values before sending them to the server.
		kpMap.Encode()

		options := environmentCommons.SyncRequestOptions{
			Data: &kpMap,
		}

		options.IntegrationType = integrations.Type(integrationType)

		//	Fetch the list of events with their respective type of integrations.
		if options.IntegrationType == "" {

			events, err := events.GetByEnvironment(commons.DefaultContext, commons.GQLClient, commons.Secret.EnvID)
			if err != nil {
				log.Debug(err)
				log.Fatal("failed to fetch active integrations for your environment")
			}

			var types []integrations.Type
			for _, item := range *events {
				types = append(types, item.Integration.Type)
			}

			selection := promptui.Select{
				Label: "Platform to sync your secrets to",
				Items: types,
			}

			index, _, err := selection.Run()
			if err != nil {
				os.Exit(1)
			}

			options.IntegrationType = types[index]
		}

		body, err := json.Marshal(&options)
		if err != nil {
			log.Debug(err)
			log.Fatal("failed to marshal your HTTP request body")
		}

		req, err := http.NewRequestWithContext(commons.DefaultContext, http.MethodPost, commons.API+"/v1/environments/"+commons.Secret.EnvID+"/sync", bytes.NewBuffer(body))
		if err != nil {
			log.Debug(err)
			log.Fatal("failed to create your HTTP request")
		}

		var response clients.APIResponse
		err = commons.HTTPClient.Run(commons.DefaultContext, req, &response)
		if err != nil {
			log.Fatal(err)
		}

		if response.Error != "" {
			log.Fatal(response.Error)
		}

		log.Info("Successfully synced secrets")
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	syncCmd.Flags().IntVarP(&version, "version", "v", -1, "Version of your secret; -1 for latest version")
	syncCmd.Flags().StringVarP(&password, "password", "p", "", "Your envsecrets account password")
	syncCmd.Flags().StringVarP(&integrationType, "type", "t", "", "Type of integration to push secrets to")
	syncCmd.Flags().StringVarP(&environmentName, "env", "e", "", "Remote environment to sync the secrets to.")
	syncCmd.MarkFlagRequired("env")
}
