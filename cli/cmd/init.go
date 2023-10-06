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
	configCommons "github.com/envsecrets/envsecrets/cli/config/commons"
	projectConfig "github.com/envsecrets/envsecrets/cli/config/project"
	"github.com/envsecrets/envsecrets/internal/memberships"
	"github.com/envsecrets/envsecrets/internal/organisations"
	"github.com/envsecrets/envsecrets/internal/projects"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	organisationID string
	projectID      string
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

		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {

		//
		// ---	Flow ---
		//
		// 1. Fetch all the organisations the user has access to. And let them choose any one.
		// 2. Fetch all the projects the user has access to in the choosen organisation.
		//	a. Let them choose any one.
		//	b. Let them create a new one.

		//
		//	Call APIs to pull existing entities
		//
		var organisation organisations.Organisation
		var project projects.Project

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
			orgs, err := organisations.GetService().List(commons.DefaultContext, commons.GQLClient)
			if err != nil {
				log.Debug(err)
				log.Fatal("Failed to fetch your organisations")
			}

			var orgsStringList []string
			for _, item := range *orgs {
				orgsStringList = append(orgsStringList, item.Name)
			}

			selection := promptui.Select{
				Label: "Choose Your Organisation",
				Items: orgsStringList,
			}

			index, _, err := selection.Run()
			if err != nil {
				os.Exit(1)
			}

			for itemIndex, item := range *orgs {
				if itemIndex == index {
					organisation = item
					break
				}
			}
		}

		//	Setup project
		if len(projectID) == 0 {

			projectsList, err := projects.List(commons.DefaultContext, commons.GQLClient, &projects.ListOptions{
				OrgID: organisation.ID,
			})
			if err != nil {
				log.Debug(err)
				log.Fatal("Failed to fetch yours projects")
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
				item, err := projects.Create(commons.DefaultContext, commons.GQLClient, &projects.CreateOptions{
					OrgID: organisation.ID,
					Name:  result,
				})
				if err != nil {
					log.Debug(err)
					log.Fatal("Failed to create the project")
				}

				project.ID = item.ID
				project.Name = fmt.Sprint(item.Name)

				//	Wait until default environments are not created.
				log.Info("Creating your default environments. Wait for 5 seconds...")
				time.Sleep(5 * time.Second)
			}
		}

		//	Pull the user's copy of organisation key.
		key, err := memberships.GetKey(commons.DefaultContext, commons.GQLClient, &memberships.GetKeyOptions{
			OrgID:  organisation.ID,
			UserID: commons.AccountConfig.User.ID,
		})
		if err != nil {
			log.Debug(err)
			log.Fatal("Failed to fetch the encryption key")
		}

		//	Write selected entities to project config
		if err := projectConfig.Save(&configCommons.Project{
			//OrgID:     organisation.ID,
			ProjectID: project.ID,
			Key:       key,
			//AutoCapitalize: true,
		}); err != nil {
			log.Debug(err)
			log.Fatal("Failed to save new project configuration locally")
		}
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
	//initCmd.Flags().StringVarP(&environmentID, "environment", "e", "", "Your existing envsecrets environment")
}
