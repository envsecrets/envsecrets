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
	"os"

	"github.com/envsecrets/envsecrets/cli/commons"
	"github.com/envsecrets/envsecrets/cli/config"
	configCommons "github.com/envsecrets/envsecrets/cli/config/commons"
	globalCommons "github.com/envsecrets/envsecrets/commons"
	"github.com/envsecrets/envsecrets/internal/auth"
	"github.com/envsecrets/envsecrets/internal/clients"
	userCommons "github.com/envsecrets/envsecrets/internal/users/commons"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	email    string
	password string
)

const (
	KEY_BYTES = 32
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate your envsecrets cloud account",
	Args:  cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {

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
				os.Exit(1)
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
				os.Exit(1)
			}
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		client := clients.NewNhostClient(&clients.NhostConfig{
			BaseURL: configCommons.NHOST_AUTH_URL,
			Logger:  log,
		})

		//	Call the appropriate service handler.
		response, err := auth.GetService().SigninWithPassword(commons.DefaultContext, client, &auth.SigninWithPasswordOptions{
			Email:    email,
			Password: password,
		})
		if err != nil {
			log.Debug(err)
			log.Fatal("Login failed. Recheck your credentials.")
		}

		//	If the user has MFA enabled.
		if response.MFA != nil {

			//	Ask the user for TOTP.
			prompt := promptui.Prompt{
				Label:       "OTP",
				Mask:        '*',
				HideEntered: true,
				Validate: func(input string) error {
					if len(input) != 6 {
						return fmt.Errorf("otp should be 6 digits")
					}
					return nil
				},
			}

			totp, err := prompt.Run()
			if err != nil {
				os.Exit(1)
			}

			response, err = auth.GetService().SigninWithMFA(commons.DefaultContext, client, &auth.SigninWithMFAOptions{
				Ticket: response.MFA["ticket"].(string),
				OTP:    totp,
			})
			if err != nil {
				log.Debug(err)
				log.Fatal("Login failed. Recheck your credentials.")
			}
		}

		var session struct {
			AccessToken          string           `json:"accessToken"`
			AccessTokenExpiresIn int              `json:"accessTokenExpiresIn"`
			RefreshToken         string           `json:"refreshToken"`
			User                 userCommons.User `json:"user"`
		}

		if err := globalCommons.MapToStruct(response.Session, &session); err != nil {
			log.Debug(err)
			log.Fatal("Failed to map the configuration")
		}

		//	Save the account config
		if err := config.GetService().Save(configCommons.Account{
			AccessToken:  session.AccessToken,
			RefreshToken: session.RefreshToken,
			User:         session.User,
		}, configCommons.AccountConfig); err != nil {
			log.Debug(err)
			log.Fatal("Failed to save account configuration locally")
		}

		//	Reload the clients.
		commons.Initialize(log)

		//	Initialize a new GQL client with the user's access token.
		gqlClient := clients.NewGQLClient(&clients.GQLConfig{
			BaseURL:       commons.NHOST_GRAPHQL_URL,
			Authorization: fmt.Sprintf("Bearer %s", response.Session["accessToken"].(string)),
			Logger:        log,
		})

		//	Extract and decrypt keys from user's session.
		pair, err := auth.GetService().DecryptKeysFromSession(commons.DefaultContext, gqlClient, &auth.DecryptKeysFromSessionOptions{
			Session:  response.Session,
			Password: password,
		})
		if err != nil {
			log.Debug(err)
			log.Fatal("Failed to decrypt your keys")
		}

		//	Save the public-private keys locally.
		if err := config.GetService().Save(configCommons.Keys{
			Public:  pair.PublicKey,
			Private: pair.PrivateKey,
		}, configCommons.KeysConfig); err != nil {
			log.Debug(err)
			log.Fatal("Failed to save key configuration locally")
		}
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		log.Info("You are logged in!")
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
