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
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/envsecrets/envsecrets/cli/commons"
	"github.com/envsecrets/envsecrets/config"
	configCommons "github.com/envsecrets/envsecrets/config/commons"
	"github.com/envsecrets/envsecrets/internal/memberships"
	"github.com/spf13/cobra"
)

var listJSON bool

// membersCmd represents the members command
var membersCmd = &cobra.Command{
	Use:   "members",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PreRun: func(cmd *cobra.Command, args []string) {

		//	Ensure organisation is saved in project config.
		if !config.GetService().Exists(configCommons.ProjectConfig) {

			// TODO: Run the init flow
			fmt.Println("Organisation doesn't exist in local project config")
		}
	},
	Run: func(cmd *cobra.Command, args []string) {

		//	Load the project config
		localConfigData, er := config.GetService().Load(configCommons.ProjectConfig)
		if er != nil {
			panic("project config not found: " + er.Error())
		}

		localConfig, ok := localConfigData.(*configCommons.Project)
		if !ok {
			panic("failed type assertion for project config")
		}

		//	Pull members of the current organisation
		items, err := memberships.List(commons.DefaultContext, commons.GQLClient, &memberships.ListOptions{
			OrgID: localConfig.Organisation,
		})
		if err != nil {
			panic(err.Error.Error())
		}

		if listJSON {

			data, err := json.MarshalIndent(items, "", "  ")
			if err != nil {
				panic(err)
			}
			fmt.Println(string(data))

		} else {

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
			fmt.Fprintf(w, "\t%s\t%s\n", "Email", "Added On")
			fmt.Fprintf(w, "\t%s\t%s\n", "----", "----")
			for _, item := range *items {
				fmt.Fprintf(w, "\t%s\t%s\n", item.User.Email, item.CreatedAt)
			}
			w.Flush()
		}
	},
}

func init() {
	rootCmd.AddCommand(membersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// membersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	membersCmd.Flags().BoolVar(&listJSON, "json", false, "Print list in JSON format")
}
