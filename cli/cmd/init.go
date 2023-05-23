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
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/envsecrets/envsecrets/cli/auth"
	"github.com/envsecrets/envsecrets/cli/commons"
	"github.com/envsecrets/envsecrets/internal/environments"
	"github.com/envsecrets/envsecrets/internal/memberships"
	"github.com/envsecrets/envsecrets/internal/organisations"
	"github.com/envsecrets/envsecrets/internal/projects"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	configCommons "github.com/envsecrets/envsecrets/cli/config/commons"
	projectConfig "github.com/envsecrets/envsecrets/cli/config/project"
)

var (
	organisationID string
	projectID      string
	environmentID  string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize your current directory for envsecrets",
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {

		//	If the user is not already authenticated,
		//	log them in first.
		if !auth.IsLoggedIn() {
			loginCmd.PreRunE(cmd, args)
			loginCmd.Run(cmd, args)
		}

		//	Re-initialize the commons
		commons.Initialize()

		return nil
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

		//	All names entered by the user must be slugs.
		var re = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
		validate := func(input string) error {
			if len(re.FindAllString(input, -1)) == 0 {
				return errors.New("should be a slug; example: my-new-idea")
			}

			return nil
		}

		//	Setup organisation first
		if len(organisationID) == 0 {

			//	Check whether user has access to at least 1 organisation.
			orgs, er := organisations.List(commons.DefaultContext, commons.GQLClient)
			if er != nil {
				log.Debug(er.Error)
				log.Fatal(er.Message)
			}

			var orgsStringList []string
			for _, item := range *orgs {
				orgsStringList = append(orgsStringList, item.Name)
			}

			selection := promptui.SelectWithAdd{
				Label:    "Choose Your Organisation",
				Items:    orgsStringList,
				AddLabel: "Create New Organisation",
				Validate: validate,
			}

			index, result, err := selection.Run()
			if err != nil {
				os.Exit(1)
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
				item, er := organisations.Create(commons.DefaultContext, commons.GQLClient, &organisations.CreateOptions{
					Name: result,
				})
				if er != nil {
					log.Debug(er.Error)
					log.Fatal(er.Message)
				}

				organisation.ID = item.ID
				organisation.Name = fmt.Sprint(item.Name)
			}
		}

		//	Setup project
		if len(projectID) == 0 {

			projectsList, er := projects.List(commons.DefaultContext, commons.GQLClient, &projects.ListOptions{
				OrgID: organisation.ID,
			})
			if er != nil {
				log.Debug(er.Error)
				log.Fatal(er.Message)
			}

			var projectsStringList []string
			for _, item := range *projectsList {
				projectsStringList = append(projectsStringList, item.Name)
			}

			selection := promptui.SelectWithAdd{
				Label:    "Choose Your Project",
				Items:    projectsStringList,
				AddLabel: "Create New Project",
				Validate: validate,
			}

			index, result, err := selection.Run()
			if err != nil {
				os.Exit(1)
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
				item, er := projects.Create(commons.DefaultContext, commons.GQLClient, &projects.CreateOptions{
					OrgID: organisation.ID,
					Name:  result,
				})
				if er != nil {
					log.Debug(er.Error)
					log.Fatal(er.Message)
				}

				project.ID = item.ID
				project.Name = fmt.Sprint(item.Name)

				//	Wait until default environments are not created.
				log.Info("Creating your default environments. Wait for 5 seconds...")
				time.Sleep(5 * time.Second)
			}
		}

		//	Setup environment
		if len(environmentID) == 0 {

			environmentsList, er := environments.List(commons.DefaultContext, commons.GQLClient, &environments.ListOptions{
				ProjectID: project.ID,
			})
			if er != nil {
				log.Debug(er.Error)
				log.Fatal(er.Message)
			}

			var environmentsStringList []string
			for _, item := range *environmentsList {
				environmentsStringList = append(environmentsStringList, item.Name)
			}

			selection := promptui.SelectWithAdd{
				Label:    "Choose Your Environment",
				Items:    environmentsStringList,
				AddLabel: "Create New Environment",
				Validate: validate,
			}

			index, result, err := selection.Run()
			if err != nil {
				os.Exit(1)
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
				item, er := environments.Create(commons.DefaultContext, commons.GQLClient, &environments.CreateOptions{
					ProjectID: project.ID,
					Name:      result,
				})
				if er != nil {
					log.Debug(er.Error)
					log.Fatal(er.Message)
				}

				environment.ID = item.ID
				environment.Name = fmt.Sprint(item.Name)
			}
		}

		//	Pull the user's copy of organisation key.
		key, err := memberships.GetKey(commons.DefaultContext, commons.GQLClient, &memberships.GetKeyOptions{
			OrgID:  organisation.ID,
			UserID: commons.AccountConfig.User.ID,
		})
		if err != nil {
			log.Debug(err.Error)
			log.Debug(err.Message)
		}

		//	Write selected entities to project config
		if err := projectConfig.Save(&configCommons.Project{
			Version:        1,
			Organisation:   organisation.ID,
			Project:        project.ID,
			Environment:    environment.ID,
			OrgKey:         key,
			AutoCapitalize: true,
		}); err != nil {
			log.Debug(err)
			log.Fatal("Failed to save new project configuration locally")
		}

		/* 		//	Pull environment secrets and populate the contingency file
		   		secret, err := secrets.GetAll(commons.DefaultContext, commons.GQLClient, &secretsCommons.GetSecretOptions{
		   			EnvID: environment.ID,
		   		})
		   		if err != nil {
		   			log.Debug(err.Error)
		   			log.Warn(err.Message)
		   		}

		   		if err := config.GetService().Save(configCommons.Contingency(secret.Data), configCommons.ContingencyConfig); err != nil {
		   			log.Debug(err)
		   			log.Error("Failed to save secrets in Contingency file")
		   		}
		*/
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		log.Info("You can now set your secrets using `envs set`")
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
	initCmd.Flags().StringVarP(&organisationID, "organisation", "w", "", "Your existing envsecrets organisation")
	initCmd.Flags().StringVarP(&projectID, "project", "p", "", "Your existing envsecrets project")
	initCmd.Flags().StringVarP(&environmentID, "environment", "e", "", "Your existing envsecrets environment")
}
