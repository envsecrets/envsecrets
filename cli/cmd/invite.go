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
	"fmt"
	"net/mail"
	"strings"

	"github.com/envsecrets/envsecrets/config"
	"github.com/envsecrets/envsecrets/config/commons"
	"github.com/envsecrets/envsecrets/internal/client"
	"github.com/envsecrets/envsecrets/internal/context"
	"github.com/envsecrets/envsecrets/internal/invites"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// inviteCmd represents the invite command
var inviteCmd = &cobra.Command{
	Use:   "invite [email]",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {

		//	If email hasn't been passed,
		//	break the flow.
		email, err := mail.ParseAddress(args[0])
		if err != nil {
			panic(err)
		}

		//	Initialize GQL Client
		client := client.GRAPHQL_CLIENT

		var project, environment string

		//	Load the current organisation
		projectConfigData, er := config.GetService().Load(commons.ProjectConfig)
		if er != nil {
			panic(er.Error())
		}

		projectConfig := projectConfigData.(*commons.Project)

		//	Take input for project
		selection := promptui.Select{
			Label: "Choose Project To Invite " + email.String() + " For",
			Items: []string{"All projects in this organisation", "Current Project Only"},
		}

		index, _, err := selection.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		if index == 0 {
			project = "*"
		} else {
			project = projectConfig.Project
		}

		//	Take input for project
		selection = promptui.Select{
			Label: "Choose Environment To Invite " + email.String() + " For",
			Items: []string{"All environments in choosen projects", "Current Environment Only"},
		}

		index, _, err = selection.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		if index == 0 {
			environment = "*"
		} else {
			environment = projectConfig.Environment
		}

		//	Initialize the invitation scope
		scope := strings.Join([]string{project, environment}, "/")

		//	Send the invite
		invite, invitesErr := invites.Create(context.DContext, client, &invites.CreateOptions{
			OrgID:         projectConfig.Organisation,
			Scope:         scope,
			ReceiverEmail: email.Address,
		})
		if invitesErr != nil {
			fmt.Printf("%s either doesn't have an envsecrets account or is already a member of this organisation.\n Ask %s to signup for envsecrets and send this invite again.", email.Address, email.String())
			panic(invitesErr.Error.Error())
		}

		if len(invite.ID) > 0 {
			fmt.Printf("Invitation sent to %s!", email.Address)
		}
	},
}

func init() {
	membersCmd.AddCommand(inviteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// inviteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// inviteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
