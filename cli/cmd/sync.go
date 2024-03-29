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

	"github.com/envsecrets/envsecrets/cli/clients"
	"github.com/envsecrets/envsecrets/cli/commons"
	"github.com/envsecrets/envsecrets/cli/internal/secrets"
	"github.com/envsecrets/envsecrets/internal/environments"
	"github.com/envsecrets/envsecrets/internal/events"
	"github.com/envsecrets/envsecrets/internal/integrations"
	"github.com/envsecrets/envsecrets/internal/secrets/pkg/keypayload"
	"github.com/envsecrets/envsecrets/internal/secrets/pkg/payload"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var all bool

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
		InitializeSecret(commons.Log)
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

		result, err := secrets.GetService().Get(commons.DefaultContext, commons.GQLClient.GQLClient, &getOptions)
		if err != nil {
			commons.Log.Debug(err)
			commons.Log.Fatal("Failed to fetch the value")
		}

		commons.Secret = result

		//	Decrypt and decode the common secret.
		DecryptAndDecode()

		//	Encrypt the secrets with the sync key.
		var syncKey [32]byte
		copy(syncKey[:], commons.KeysConfig.Sync)
		if err := commons.Secret.Encrypt(syncKey); err != nil {
			commons.Log.Debug(err)
			commons.Log.Fatal("Failed to decrypt the secret")
		}

		//	Copy the dto.KPMap to keypayload.KPMap
		kpMap := keypayload.KPMap{}
		for key, value := range commons.Secret.Data.GetMapping() {
			kpMap[key] = &payload.Payload{
				Value: value.GetValue(),
			}
		}
		kpMap.MarkAllEncoded()

		options := environments.SyncOptions{
			Pairs: &kpMap,
		}

		//	Fetch the list of events with their respective type of integrations.
		if !all {

			events, err := events.GetService().GetByEnvironment(commons.DefaultContext, commons.GQLClient.GQLClient, commons.Secret.EnvID)
			if err != nil {
				commons.Log.Debug(err)
				commons.Log.Fatal("failed to fetch active integrations for your environment")
			}

			type item struct {
				ID    string
				Title string
				Type  integrations.Type
			}

			var items []item
			for _, event := range *events {
				items = append(items, item{
					ID:    event.ID,
					Title: event.GetEntityTitle(),
					Type:  event.Integration.Type,
				})
			}

			selection := promptui.Select{
				Label: "Which platform do you want to sync your secrets to?",
				Items: items,
				Templates: &promptui.SelectTemplates{
					Active:   `{{ ">" | blue }} [{{ .Type }}] {{ .Title }}`,
					Inactive: `[{{ .Type }}] {{ .Title }}`,
					Selected: `{{ "✔" | green }} [{{ .Type }}] {{ .Title }}`,
				},
			}

			index, _, err := selection.Run()
			if err != nil {
				os.Exit(1)
			}

			options.EventIDs = []string{(*events)[index].ID}
		}

		body, err := json.Marshal(&options)
		if err != nil {
			commons.Log.Debug(err)
			commons.Log.Fatal("failed to marshal your HTTP request body")
		}

		req, err := http.NewRequestWithContext(commons.DefaultContext, http.MethodPost, clients.API+"/v1/environments/"+commons.Secret.EnvID+"/sync", bytes.NewBuffer(body))
		if err != nil {
			commons.Log.Debug(err)
			commons.Log.Fatal("failed to create your HTTP request")
		}

		var response clients.APIResponse
		err = commons.HTTPClient.Run(commons.DefaultContext, req, &response)
		if err != nil {
			commons.Log.Fatal(err)
		}

		if response.Error != "" {
			commons.Log.Fatal(response.Error)
		}

		commons.Log.Info("Successfully synced secrets to connected services")
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
	syncCmd.Flags().BoolVarP(&all, "all", "a", false, "Bypass selection and sync to all integrations connected to the environment")
	syncCmd.Flags().IntVarP(&version, "version", "v", -1, "Version of your secret; -1 for latest version")
	syncCmd.Flags().StringVarP(&environmentName, "env", "e", "", "Remote environment to sync the secrets to.")
	syncCmd.MarkFlagRequired("env")
}
