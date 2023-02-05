/*
Copyright © 2023 Mrinal Wahal mrinalwahal@gmail.com

*/
package cmd

import (
	"fmt"

	projectConfig "github.com/envsecrets/envsecrets/config/project"
	"github.com/envsecrets/envsecrets/internal/client"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/environments"
	"github.com/envsecrets/envsecrets/internal/projects"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// environmentsCreateCmd represents the environmentsCreate command
var environmentsCreateCmd = &cobra.Command{
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

		//	Load the project config
		localConfig, err := projectConfig.Load()
		if err != nil {
			panic(err)
		}

		//	Fetch the current workspace
		project, er := projects.Get(context.DContext, client, localConfig.Project)
		if er != nil {
			panic(er.Error.Error())
		}

		//	Notify the user about their current workspace
		fmt.Printf("You are about to create a new environment in the current '%s' project.\n", project.Name)

		if name == "" {

			//	Validate input
			validate := func(input string) error {
				return nil
			}

			prompt := promptui.Prompt{
				Label:     "Environment",
				Default:   "dev",
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
		item, er := environments.Create(context.DContext, client, &environments.CreateOptions{
			ProjectID: localConfig.Project,
			Name:      name,
		})
		if er != nil {
			panic(er.Error.Error())
		}

		//	Update the new value
		localConfig.Environment = item.Name

		if err := projectConfig.Save(localConfig); err != nil {
			panic(err)
		}

	},
	PostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("New environment created and set in project configuration!")
	},
}

func init() {
	environmentsCmd.AddCommand(environmentsCreateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// environmentsCreateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// environmentsCreateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
