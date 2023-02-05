/*
Copyright © 2023 Mrinal Wahal mrinalwahal@gmail.com

*/
package cmd

import (
	"fmt"
	"path/filepath"

	configCommons "github.com/envsecrets/envsecrets/config/commons"
	projectConfig "github.com/envsecrets/envsecrets/config/project"
	"github.com/envsecrets/envsecrets/internal/client"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/organisations"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var name string

// organisationsCreateCmd represents the organisationsCreate command
var organisationsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		//	Initialize GQL Client
		client := client.GRAPHQL_CLIENT

		if name == "" {

			//	Validate input
			validate := func(input string) error {
				return nil
			}

			prompt := promptui.Prompt{
				Label:     "Organisation",
				Default:   filepath.Base(filepath.Dir(filepath.Dir(configCommons.EXECUTABLE))),
				AllowEdit: true,
				Validate:  validate,
			}

			result, err := prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			name = result
		}

		//	Create new item
		item, err := organisations.Create(context.DContext, client, &organisations.CreateOptions{
			Name: name,
		})
		if err != nil {
			panic(err)
		}

		// TODO:	Set the new workspace ID in project config
		config, er := projectConfig.Load()
		if er != nil {
			panic(er)
		}

		//	Update the new value
		config.Organisation = item.ID

		if err := projectConfig.Save(config); err != nil {
			panic(err)
		}

	},
	PostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("New workspace created and set in local environment!")
	},
}

func init() {
	organisationsCmd.AddCommand(organisationsCreateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// organisationsCreateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	organisationsCreateCmd.Flags().StringVarP(&name, "name", "n", "", "name of your new workspace")
}
