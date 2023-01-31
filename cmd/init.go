/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/envsecrets/envsecrets/cmd/internal/auth"
	"github.com/envsecrets/envsecrets/config"
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
	workspaceName   string
	projectName     string
	environmentName string
	branchName      string
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
		var workspace workspaces.Workspace
		var project projects.Project
		//	var environment environments.Environment

		//	Initialize GQL Client
		client := client.GRAPHQL_CLIENT

		//	Setup workspace first
		if len(workspaceName) == 0 {

			//	Validate input
			validate := func(input string) error {
				return nil
			}

			prompt := promptui.Prompt{
				Label:     "Workspace",
				Default:   filepath.Base(filepath.Dir(filepath.Dir(config.EXECUTABLE))),
				AllowEdit: true,
				Validate:  validate,
			}

			result, err := prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			//	Create new item
			item, err := workspaces.Create(context.DContext, client, &workspaces.CreateOptions{
				Name: result,
			})
			if err != nil {
				fmt.Println(err)
			}

			workspace.ID = item.ID
			workspace.Name = fmt.Sprint(item.Name)
		}

		//	Setup project
		if len(projectName) == 0 {

			//	Validate input
			validate := func(input string) error {
				return nil
			}

			prompt := promptui.Prompt{
				Label:     "Project",
				Default:   filepath.Base(filepath.Dir(config.EXECUTABLE)),
				AllowEdit: true,
				Validate:  validate,
			}

			result, err := prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			//	Create new item
			item, err := projects.Create(context.DContext, client, &projects.CreateOptions{
				WorkspaceID: workspace.ID,
				Name:        result,
			})
			if err != nil {
				fmt.Println(err)
			}

			project.ID = item.ID
			project.Name = item.Name
			project.WorkspaceID = item.WorkspaceID
		}

		/* 		//	Setup environment
		   		if len(environmentName) == 0 {

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

		   			//	Create new item
		   			item, err := environments.Create(context.DContext, client, &environments.CreateOptions{
		   				ProjectID: project.ID,
		   				Name:      result,
		   			})
		   			if err != nil {
		   				fmt.Println(err)
		   			}

		   			environment.ID = item.ID
		   			environment.Name = item.Name
		   			environment.ProjectID = project.ID
		   		}
		*/
		//	Write selected entities to project config
		if err := projectConfig.Save(&configCommons.Project{
			Version:     1,
			Workspace:   workspace.ID,
			Project:     project.ID,
			Environment: "dev",
			Branch:      "main",
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
	initCmd.Flags().StringVarP(&workspaceName, "workspace", "w", "", "Your existing envsecrets workspace")
	initCmd.Flags().StringVarP(&projectName, "project", "p", "", "Your existing envsecrets project")
	initCmd.Flags().StringVarP(&environmentName, "environment", "e", "dev", "Your existing envsecrets environment")
	initCmd.Flags().StringVarP(&branchName, "branch", "b", "main", "Your existing envsecrets branch")
}
