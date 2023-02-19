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
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/envsecrets/envsecrets/config"
	configCommons "github.com/envsecrets/envsecrets/config/commons"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/secrets/commons"
	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		//	Run sanity checks
		if len(args) < 1 || len(args) > 1 {
			panic("invalid key-value pair")
		}

		if !strings.Contains(args[0], "=") {
			panic("invalid key-value pair")
		}

		pair := strings.Split(args[0], "=")

		if len(pair) != 2 {
			panic("invalid key-value pair")
		}

		key := pair[0]
		value := pair[1]

		data := commons.Secret{
			Key:   key,
			Value: value,
		}

		//	Load the project configuration
		projectConfigData, er := config.GetService().Load(configCommons.ProjectConfig)
		if er != nil {
			panic(er.Error())
		}

		projectConfig := projectConfigData.(*configCommons.Project)

		//	Send the secrets to vault
		payload := commons.SetRequest{
			Secret: commons.Secret{
				Key:   key,
				Value: value,
			},
			Path: commons.Path{
				Organisation: projectConfig.Organisation,
				Project:      projectConfig.Project,
				Environment:  projectConfig.Environment,
			},
		}

		reqBody, _ := payload.Marshal()
		req, err := http.NewRequestWithContext(context.DContext, http.MethodPost, os.Getenv("API")+"/api/v1/secrets", bytes.NewBuffer(reqBody))
		if err != nil {
			panic(err)
		}

		//	Load the account configuration
		accountConfigData, er := config.GetService().Load(configCommons.AccountConfig)
		if er != nil {
			panic(er.Error())
		}

		accountConfig := accountConfigData.(*configCommons.Account)

		//	Set Authorization Header
		req.Header.Set("Authorization", "Bearer "+accountConfig.AccessToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(respBody))

		if resp.StatusCode != http.StatusOK {
			panic("failed to set secret")
		}

		//	Set the values in current application
		if err := os.Setenv(key, value); err != nil {
			panic(err)
		}

		//	Export the values in current shell
		if err := exec.Command("sh", "-c", "export", data.String()).Run(); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(setCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//	setCmd.Flags().StringVarP("toggle", "t", false, "Help message for toggle")
}
