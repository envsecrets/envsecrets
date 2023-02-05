/*
Copyright © 2023 Mrinal Wahal mrinalwahal@gmail.com

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	projectConfig "github.com/envsecrets/envsecrets/config/project"
	"github.com/envsecrets/envsecrets/internal/client"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/projects"
	"github.com/spf13/cobra"
)

// projectsListCmd represents the projectsList command
var projectsListCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		//	Initialize GQL Client
		client := client.GRAPHQL_CLIENT

		//	Load the project config
		localConfig, err := projectConfig.Load()
		if err != nil {
			panic(err)
		}

		//	List items
		items, er := projects.List(context.DContext, client, &projects.ListOptions{
			OrgID: localConfig.Organisation,
		})
		if er != nil {
			panic(er.Error.Error())
		}

		if listJSON {

			data, err := json.MarshalIndent(items, "", "  ")
			if err != nil {
				panic(err)
			}
			fmt.Println(string(data))

		} else {

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.TabIndent)
			fmt.Fprintf(w, "\t%s\t%s\n", "Name", "ID")
			fmt.Fprintf(w, "\t%s\t%s\n", "----", "----")
			for _, item := range *items {
				if item.ID == localConfig.Project {
					fmt.Fprintf(w, "\t%s\t%s\t(current)\n", item.Name, item.ID)
				} else {
					fmt.Fprintf(w, "\t%s\t%s\n", item.Name, item.ID)
				}
			}
			w.Flush()
		}
	},
}

func init() {
	projectsCmd.AddCommand(projectsListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// projectsListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	projectsListCmd.Flags().BoolVar(&listJSON, "json", false, "Print list in JSON format")
}
