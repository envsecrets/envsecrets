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
	"github.com/envsecrets/envsecrets/internal/projects"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// projectsCreateCmd represents the projectsCreate command
var projectsCreateCmd = &cobra.Command{
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
		workspace, er := organisations.Get(context.DContext, client, localConfig.Organisation)
		if er != nil {
			panic(er.Error.Error())
		}

		//	Notify the user about their current workspace
		fmt.Printf("You are about to create a new project in the current '%s' workspace.\n", workspace.Name)

		if name == "" {

			//	Validate input
			validate := func(input string) error {
				return nil
			}

			prompt := promptui.Prompt{
				Label:     "Project",
				Default:   filepath.Base(filepath.Dir(configCommons.EXECUTABLE)),
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
		item, er := projects.Create(context.DContext, client, &projects.CreateOptions{
			OrgID: localConfig.Organisation,
			Name:  name,
		})
		if er != nil {
			panic(er.Error.Error())
		}

		//	Update the new value
		localConfig.Project = item.ID

		if err := projectConfig.Save(localConfig); err != nil {
			panic(err)
		}

	},
	PostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("New project created and set in local environment!")
	},
}

func init() {
	projectsCmd.AddCommand(projectsCreateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// projectsCreateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	projectsCreateCmd.Flags().StringVarP(&name, "name", "n", "", "name of your new project")
}
