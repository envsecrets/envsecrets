/*
Copyright © 2023 Mrinal Wahal mrinalwahal@gmail.com

*/
package cmd

import (
	"fmt"

	"github.com/envsecrets/envsecrets/internal/auth"
	"github.com/envsecrets/envsecrets/internal/client"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/environments"
	"github.com/envsecrets/envsecrets/internal/organisations"
	"github.com/envsecrets/envsecrets/internal/projects"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	configCommons "github.com/envsecrets/envsecrets/config/commons"
	projectConfig "github.com/envsecrets/envsecrets/config/project"
)

var (
	organisationName string
	projectName      string
	environmentName  string
	branchName       string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize your project for envsecrets",
	PreRun: func(cmd *cobra.Command, args []string) {

		//	If the user is not already authenticated,
		//	log them in first.
		if !auth.IsLoggedIn() {
			loginCmd.Run(cmd, args)
		}

	},
	Run: func(cmd *cobra.Command, args []string) {

		//
		//	---	Flow ---
		//
		//	1. Check whether user is part of at least 1 organisation.
		//		-> Yes = Show option to choose from existing organisations or create a new one.
		//		-> No = Start the flow to create a new organisation.
		//	2. Check whether user has access to at least 1 project in the choosen organisation.
		//		-> Yes = Show option to choose from existing projects or create a new one.
		//		-> No = Start the flow to create a new project.
		//	3. Check whether user has access to at least 1 environment in the choosen project.
		//		-> Yes = Show option to choose from existing environments or create a new one.
		//		-> No = Start the flow to create a new environment.
		//

		//
		//	Call APIs to pull existing entities
		//
		var organisation organisations.Organisation
		var project projects.Project
		var environment environments.Environment

		//	Initialize GQL Client
		client := client.GRAPHQL_CLIENT

		//	Setup organisation first
		if len(organisationName) == 0 {

			//	Check whether user has access to at least 1 organisation.
			orgs, er := organisations.List(context.DContext, client)
			if er != nil {
				panic(er.Error)
			}

			var orgsStringList []string
			for _, item := range *orgs {
				orgsStringList = append(orgsStringList, item.Name)
			}

			selection := promptui.SelectWithAdd{
				Label:    "Choose Your Organisation",
				Items:    orgsStringList,
				AddLabel: "Create New Organisation",
			}

			index, result, err := selection.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			if index > -1 {

				for itemIndex, item := range *orgs {
					if itemIndex == index {
						organisation = item
						break
					}
				}

			} else {

				//	Create new item
				item, er := organisations.Create(context.DContext, client, &organisations.CreateOptions{
					Name: result,
				})
				if er != nil {
					panic(er.Error.Error())
				}

				organisation.ID = item.ID
				organisation.Name = fmt.Sprint(item.Name)
			}
		}

		//	Setup project
		if len(projectName) == 0 {

			projectsList, er := projects.List(context.DContext, client, &projects.ListOptions{
				OrgID: organisation.ID,
			})
			if er != nil {
				panic(er.Error)
			}

			var projectsStringList []string
			for _, item := range *projectsList {
				projectsStringList = append(projectsStringList, item.Name)
			}

			selection := promptui.SelectWithAdd{
				Label:    "Choose Your Project",
				Items:    projectsStringList,
				AddLabel: "Create New Project",
			}

			index, result, err := selection.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			if index > -1 {

				for itemIndex, item := range *projectsList {
					if itemIndex == index {
						project = item
						break
					}
				}

			} else {

				//	Create new item
				item, er := projects.Create(context.DContext, client, &projects.CreateOptions{
					OrgID: organisation.ID,
					Name:  result,
				})
				if er != nil {
					panic(er.Error.Error())
				}

				project.ID = item.ID
				project.Name = fmt.Sprint(item.Name)
			}
		}

		//	Setup project
		if len(environmentName) == 0 {

			environmentsList, er := environments.List(context.DContext, client, &environments.ListOptions{
				ProjectID: project.ID,
			})
			if er != nil {
				panic(er.Error)
			}

			var environmentsStringList []string
			for _, item := range *environmentsList {
				environmentsStringList = append(environmentsStringList, item.Name)
			}

			selection := promptui.SelectWithAdd{
				Label:    "Choose Your Environment",
				Items:    environmentsStringList,
				AddLabel: "Create New Environment",
			}

			index, result, err := selection.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			if index > -1 {

				for itemIndex, item := range *environmentsList {
					if itemIndex == index {
						environment = item
						break
					}
				}

			} else {

				//	Create new item
				item, er := environments.Create(context.DContext, client, &environments.CreateOptions{
					ProjectID: project.ID,
					Name:      result,
				})
				if er != nil {
					panic(er.Error.Error())
				}

				environment.ID = item.ID
				environment.Name = fmt.Sprint(item.Name)
			}
		}

		//	Write selected entities to project config
		if err := projectConfig.Save(&configCommons.Project{
			Version:      1,
			Organisation: organisation.ID,
			Project:      project.ID,
			Environment:  environment.ID,
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
	initCmd.Flags().StringVarP(&organisationName, "organisation", "w", "", "Your existing envsecrets organisation")
	initCmd.Flags().StringVarP(&projectName, "project", "p", "", "Your existing envsecrets project")
	initCmd.Flags().StringVarP(&environmentName, "environment", "e", "", "Your existing envsecrets environment")
}
