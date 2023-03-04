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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/envsecrets/envsecrets/cli/commons"

	"github.com/envsecrets/envsecrets/config"
	configCommons "github.com/envsecrets/envsecrets/config/commons"
	secretsCommons "github.com/envsecrets/envsecrets/internal/secrets/commons"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		//	Run sanity checks
		if len(args) != 1 {
			panic("invalid key-value pair")
		}
		key := args[0]

		var secretVersion *int

		if version > -1 {
			secretVersion = &version
		}

		//	Load the project config
		projectConfigPayload, err := config.GetService().Load(configCommons.ProjectConfig)
		if err != nil {
			panic(err)
		}

		projectConfig := projectConfigPayload.(*configCommons.Project)

		//	Get the secret service
		payload := &secretsCommons.GetRequestOptions{
			Path: secretsCommons.Path{
				Organisation: projectConfig.Organisation,
				Environment:  projectConfig.Environment,
				Project:      projectConfig.Project,
			},
			Key:     key,
			Version: secretVersion,
		}

		reqBody, _ := payload.Marshal()
		req, err := http.NewRequestWithContext(commons.DefaultContext, http.MethodGet, commons.API+"/v1/secrets", bytes.NewBuffer(reqBody))
		if err != nil {
			panic(err)
		}

		//	Set content-type header
		req.Header.Set("content-type", "application/json")

		resp, er := commons.HTTPClient.Run(commons.DefaultContext, req)
		if er != nil {
			panic(er)
		}

		if resp.StatusCode != http.StatusOK {
			panic("failed to set secret")
		}

		defer resp.Body.Close()

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		var response secretsCommons.APIResponse
		if err := json.Unmarshal(respBody, &response); err != nil {
			panic(err)
		}

		secret := response.Data.(map[string]interface{})

		//	Base64 decode the secret value
		data, err := base64.StdEncoding.DecodeString(secret["value"].(string))
		if err != nil {
			panic(err)
		}

		fmt.Println(string(data))

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
}
