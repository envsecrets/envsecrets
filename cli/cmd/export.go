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
	"os"
	"strings"

	"github.com/envsecrets/envsecrets/cli/commons"
	"github.com/envsecrets/envsecrets/cli/internal"
	"github.com/envsecrets/envsecrets/cli/internal/secrets"
	"github.com/envsecrets/envsecrets/dto"
	"github.com/envsecrets/envsecrets/internal/clients"
	"github.com/spf13/cobra"
)

var version int
var exportfile string

var XTokenHeader string

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Prints decrypted list of your environment's (key=value) secret pairs",
	PreRun: func(cmd *cobra.Command, args []string) {

		//	If the user has passed a token,
		//	avoid using email+password to authenticate them against the API.
		if XTokenHeader != "" {
			commons.Secret = &dto.Secret{}
			return
		}

		//	Initialize the common secret.
		InitializeSecret(log)
	},
	Run: func(cmd *cobra.Command, args []string) {

		if XTokenHeader != "" {

			options := &internal.GetValuesOptions{
				Token: XTokenHeader,
			}

			if version > -1 {
				options.Version = &version
			}

			result, err := internal.GetValues(commons.DefaultContext, commons.HTTPClient, options)
			if err != nil {
				log.Debug(err)
				if strings.Compare(err.Error(), string(clients.ErrorTypeRecordNotFound)) == 0 {
					log.Error("You haven't set any secrets in this environment")
					log.Info("Use `envs set --help` for more information")
					os.Exit(1)
				} else {
					log.Fatal("Failed to fetch the secrets")
				}
			}

			for k, v := range result.Data {
				commons.Secret.Set(k, &dto.Payload{
					Value: v.Value,
				})
			}

			commons.Secret.Decode()

		} else {

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
				if strings.Compare(err.Error(), string(clients.ErrorTypeRecordNotFound)) == 0 {
					log.Warn("You haven't set any secrets in this environment")
					log.Info("Use `envs set --help` for more information")
					os.Exit(1)
				} else {
					log.Fatal("Failed to fetch the secrets")
				}
			}

			commons.Secret = result

			//	Decrypt and decode the common secret.
			DecryptAndDecode()
		}

		//	Initialize a new buffer to store key=value lines
		var buffer bytes.Buffer

		buffer.WriteString(strings.Join(commons.Secret.Data.FmtStrings(), "\n"))

		fmt.Println(buffer.String())
		/*
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
						case ".yaml":
						}
					} else {
						fmt.Println(buffer.String())
					}
		*/
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
	exportCmd.Flags().StringVarP(&environmentName, "env", "e", "", "Remote environment to set the secrets in. Defaults to the local environment.")
}
