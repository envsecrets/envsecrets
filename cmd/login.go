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

	"github.com/envsecrets/envsecrets/cmd/internal/auth"
	accountConfig "github.com/envsecrets/envsecrets/config/account"
	configCommons "github.com/envsecrets/envsecrets/config/commons"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	email    string
	password string
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate your envsecrets account",
	Run: func(cmd *cobra.Command, args []string) {

		if auth.IsLoggedIn() {
			fmt.Println("You are already authenticated!\nLog out with `envsecrets logout`")
			return
		}

		var err error

		if len(email) == 0 {

			//	Take email input
			validate := func(input string) error {
				_, err := mail.ParseAddress(input)
				return err
			}

			emailPrompt := promptui.Prompt{
				Label:    "Email",
				Validate: validate,
			}

			email, err = emailPrompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}
		}

		if len(password) == 0 {

			//	Take password input
			passwordPrompt := promptui.Prompt{
				Label: "Password",
				Mask:  '*',
			}

			password, err = passwordPrompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}
		}

		//	Prepare body
		payload := map[string]interface{}{
			"email":    email,
			"password": password,
		}

		response, err := auth.Login(payload)
		if err != nil {
			panic(err)
		}

		//	Save the account config
		if err := accountConfig.Save(&configCommons.Account{
			AccessToken:  response.Session.AccessToken,
			RefreshToken: response.Session.RefreshToken,
			User:         response.Session.User,
		}); err != nil {
			panic(err)
		}

		fmt.Println("You are logged in!")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	loginCmd.Flags().StringVarP(&email, "email", "e", "", "Your envsecrets account email")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "Your envsecrets account password")
}
