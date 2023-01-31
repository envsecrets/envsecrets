/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/envsecrets/envsecrets/internal/client"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/workspaces"
	"github.com/spf13/cobra"
)

var listJSON bool

// workspacesListCmd represents the workspacesList command
var workspacesListCmd = &cobra.Command{
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

		//	List items
		items, err := workspaces.List(context.DContext, client)
		if err != nil {
			panic(err)
		}

		if listJSON {

			data, err := json.MarshalIndent(items, "", "  ")
			if err != nil {
				panic(err)
			}
			fmt.Println(string(data))

		} else {

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', tabwriter.TabIndent)
			fmt.Fprintf(w, "\t%s\t%s\n", "Name", "ID")
			fmt.Fprintf(w, "\t%s\t%s\n", "----", "----")
			for _, item := range *items {
				fmt.Fprintf(w, "\t%s\t%s\n", item.Name, item.ID)
			}
			w.Flush()
		}
	},
}

func init() {
	workspacesCmd.AddCommand(workspacesListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// workspacesListCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	workspacesListCmd.Flags().BoolVar(&listJSON, "json", false, "Print list in JSON format")
}
