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
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	internalErrors "errors"

	"github.com/envsecrets/envsecrets/cli/commons"
	"github.com/envsecrets/envsecrets/config"
	configCommons "github.com/envsecrets/envsecrets/config/commons"
	"github.com/envsecrets/envsecrets/internal/auth"
	secretsCommons "github.com/envsecrets/envsecrets/internal/secrets/commons"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var encrypt bool
var file string

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set KEY=VALUE",
	Short: "Set new key-value pairs as secrets in your current environment.",
	Long: `Set new key-value pairs as secrets in your current environment.

You can also load your variables directly from files: envsecrets set --file .env

NOTE: This command auto-capitalizes your keys.`,
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
	Args: func(cmd *cobra.Command, args []string) error {

		if file == "" && len(args) != 1 {
			return errors.New("either an import file is required to load variables from or at least 1 key-value pair (of secret) is required")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		data := make(map[string]secretsCommons.Payload)

		if file != "" {
			filedata, err := ioutil.ReadFile(file)
			if err != nil {
				log.Debug(err)
				log.Fatal("Failed to read file: ", file)
			}

			switch filepath.Ext(file) {
			default:

				lines := strings.Split(string(filedata), "\n")

				for index, item := range lines {

					//	Clean the line.
					item = strings.TrimSpace(item)

					key, payload, err := readPair(item)
					if err != nil {
						log.Error("Error on line ", index, " of your file")
						log.Fatal(err)
					}

					data[key] = *payload
				}

			case ".csv":
				log.Error("This file format is not yet supported")
				log.Info("Use `--help` for more information")
				os.Exit(1)

			case ".json":

				var mapping map[string]interface{}
				if err := json.Unmarshal(filedata, &mapping); err != nil {
					log.Debug(err)
					log.Fatal("Failed to read json from file")
				}

				for key, value := range mapping {

					//	Base64 encode the secret value
					value = base64.StdEncoding.EncodeToString([]byte(fmt.Sprint(value)))

					data[key] = secretsCommons.Payload{
						Value: value,
						Type:  secretsCommons.Ciphertext,
					}
				}

			case ".yaml":

				var mapping map[string]interface{}
				if err := yaml.Unmarshal(filedata, &mapping); err != nil {
					log.Debug(err)
					log.Fatal("Failed to read json from file")
				}

				for key, value := range mapping {

					//	Base64 encode the secret value
					value = base64.StdEncoding.EncodeToString([]byte(fmt.Sprint(value)))

					data[key] = secretsCommons.Payload{
						Value: value,
						Type:  secretsCommons.Ciphertext,
					}
				}
			}

		} else {

			//	Run sanity checks
			if len(args) < 1 {
				log.Fatal("Invalid key-value pair")
			}

			key, payload, err := readPair(args[0])
			if err != nil {
				log.Fatal(err)
			}

			data[key] = *payload
		}

		//	Load the project configuration
		projectConfigData, er := config.GetService().Load(configCommons.ProjectConfig)
		if er != nil {
			log.Debug(er)
			log.Fatal("Failed to load project configuration")
		}

		projectConfig := projectConfigData.(*configCommons.Project)

		//	Send the secrets to vault
		payload := secretsCommons.SetRequestOptions{
			Data:  data,
			OrgID: projectConfig.Organisation,
			EnvID: projectConfig.Environment,
		}

		reqBody, err := payload.Marshal()
		if err != nil {
			log.Debug(err)
			log.Fatal("Failed to prepare request payload")
		}

		req, err := http.NewRequestWithContext(commons.DefaultContext, http.MethodPost, commons.API+"/v1/secrets", bytes.NewBuffer(reqBody))
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
			log.Fatal("Failed to set the secrets")
		}

		log.Info("Secrets set! Created version ", response.Data.(map[string]interface{})["version"])
	},
}

func readPair(data string) (string, *secretsCommons.Payload, error) {

	if !strings.Contains(data, "=") {
		return "", nil, internalErrors.New("invalid key-value pair")
	}

	pair := strings.Split(data, "=")

	if len(pair) != 2 {
		return "", nil, internalErrors.New("invalid key-value pair")
	}

	key := pair[0]
	value := pair[1]

	//	Auto-capitalize the key
	key = strings.ToUpper(key)

	//	Whether to encrypt the secret value or not.
	typ := secretsCommons.Ciphertext
	if !encrypt {
		typ = secretsCommons.Plaintext
	}

	//	Base64 encode the secret value
	value = base64.StdEncoding.EncodeToString([]byte(value))

	return key, &secretsCommons.Payload{
		Value: value,
		Type:  typ,
	}, nil
}

func init() {
	rootCmd.AddCommand(setCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	setCmd.Flags().StringVarP(&file, "file", "f", "", "Filepath to import your variables from [.env, .json, .txt, .yaml]")
	setCmd.Flags().BoolVarP(&encrypt, "encrypt", "e", true, "Encrypt the value")
}
