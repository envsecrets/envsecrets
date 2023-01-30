/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/envsecrets/envsecrets/cmd/internal/auth"
	"github.com/envsecrets/envsecrets/internal/client"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/projects"
	"github.com/envsecrets/envsecrets/internal/workspaces"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	configCommons "github.com/envsecrets/envsecrets/config/commons"
	projectConfig "github.com/envsecrets/envsecrets/config/project"
)

var (
	workspace   string
	project     string
	environment string
	branch      string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize your project for envsecrets",
	Run: func(cmd *cobra.Command, args []string) {

		//	If the user is not already authenticated,
		//	log them in first.
		if !auth.IsLoggedIn() {
			loginCmd.Run(cmd, args)
		}

		//
		//	Call APIs to pull existing entities
		//

		//	Initialize GQL Client
		client := client.GRAPHQL_CLIENT

		//	Setup workspace first
		if len(workspace) == 0 {

			//	Fetch users workspaces
			workspaces, err := workspaces.List(context.DContext, client)
			if err != nil {
				panic(err)
			}

			//	[TODO] If the user doesn't have any workspace,
			//	then initiate the flow to create one.

			prompt := promptui.Select{
				Label: "Workspace",
				Items: *workspaces,
				Templates: &promptui.SelectTemplates{
					Active:   "\U0001F336 {{ .Name }}",
					Selected: "\U0001F336 {{ .Name }}",
				},
			}

			index, _, err := prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			for item, value := range *workspaces {
				if item == index {
					workspace = value.Name
				}
			}
		}

		//	Setup project
		if len(project) == 0 {

			//	Fetch users projects
			projects, err := projects.List(context.DContext, client)
			if err != nil {
				panic(err)
			}

			//	[TODO] If the user doesn't have any project,
			//	then initiate the flow to create one.

			prompt := promptui.Select{
				Label: "Project",
				Items: *projects,
				Templates: &promptui.SelectTemplates{
					Active:   "\U0001F336 {{ .Name }}",
					Selected: "\U0001F336 {{ .Name }}",
				},
			}

			index, _, err := prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			for item, value := range *projects {
				if item == index {
					project = value.Name
				}
			}
		}

		//	Write selected entities to project config
		if err := projectConfig.Save(&configCommons.Project{
			Version:   1,
			Workspace: workspace,
			Project:   project,
		}); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	initCmd.Flags().StringVarP(&workspace, "workspace", "w", "", "Your existing envsecrets workspace")
	initCmd.Flags().StringVarP(&project, "project", "p", "", "Your existing envsecrets project")
	initCmd.Flags().StringVarP(&environment, "environment", "e", "dev", "Your existing envsecrets environment")
	initCmd.Flags().StringVarP(&branch, "branch", "b", "main", "Your existing envsecrets branch")
}
