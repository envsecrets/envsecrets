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
package login

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/envsecrets/envsecrets/cli/clients"
	"github.com/envsecrets/envsecrets/cli/commons"
	"github.com/envsecrets/envsecrets/cli/config"
	configCommons "github.com/envsecrets/envsecrets/cli/config/commons"
	globalCommons "github.com/envsecrets/envsecrets/commons"
	"github.com/envsecrets/envsecrets/internal/auth"
	"github.com/envsecrets/envsecrets/internal/users"
	"github.com/spf13/cobra"
)

const (
	KEY_BYTES = 32
)

// Cmd represents the login command
var Cmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate your envsecrets cloud account",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {

		inputs := []textinput.Model{}

		if !cmd.Flag("email").Changed {
			i := getEmailInput()
			i.Focus()
			inputs = append(inputs, i)
		}

		if !cmd.Flag("password").Changed {
			inputs = append(inputs, getPasswordInput())
		}

		m := model{inputs: inputs}
		if _, err := tea.NewProgram(m).Run(); err != nil {
			commons.Log.Fatal(err)
		}

		email := inputs[0].Value()
		password := inputs[1].Value()

		//	Proceed to login the user.

		client := clients.NewNhostClient(&clients.NhostConfig{
			BaseURL: configCommons.NHOST_AUTH_URL + "/v1",
			Logger:  commons.Log,
		})

		//	Call the appropriate service handler.
		response, err := auth.GetService().SigninWithPassword(commons.DefaultContext, client, &auth.SigninWithPasswordOptions{
			Email:    email,
			Password: password,
		})
		if err != nil {
			commons.Log.Debug(err)
			commons.Log.Fatal("Login failed. Recheck your credentials.")
		}

		//	If the user has MFA enabled.
		if response.MFA != nil {

			i := getOTPInput()
			i.Focus()
			m := model{inputs: []textinput.Model{i}}
			if _, err := tea.NewProgram(m).Run(); err != nil {
				commons.Log.Fatal(err)
			}

			totp := i.Value()

			response, err = auth.GetService().SigninWithMFA(commons.DefaultContext, client, &auth.SigninWithMFAOptions{
				Ticket: response.MFA["ticket"].(string),
				OTP:    totp,
			})
			if err != nil {
				commons.Log.Debug(err)
				commons.Log.Fatal("Login failed. Recheck your credentials.")
			}
		}

		var session struct {
			AccessToken          string     `json:"accessToken"`
			AccessTokenExpiresIn int        `json:"accessTokenExpiresIn"`
			RefreshToken         string     `json:"refreshToken"`
			User                 users.User `json:"user"`
		}

		if err := globalCommons.MapToStruct(response.Session, &session); err != nil {
			commons.Log.Debug(err)
			commons.Log.Fatal("Failed to map the configuration")
		}

		//	Save the account config
		if err := config.GetService().Save(configCommons.Account{
			AccessToken:  session.AccessToken,
			RefreshToken: session.RefreshToken,
			User:         session.User,
		}, configCommons.AccountConfig); err != nil {
			commons.Log.Debug(err)
			commons.Log.Fatal("Failed to save account configuration locally")
		}

		//	Reload the clients.
		commons.Initialize(commons.Log)

		//	Initialize a new GQL client with the user's access token.
		gqlClient := clients.NewGQLClient(&clients.GQLConfig{
			BaseURL:       commons.NHOST_GRAPHQL_URL,
			Authorization: fmt.Sprintf("Bearer %s", response.Session["accessToken"].(string)),
			Logger:        commons.Log,
		})

		//	Extract and decrypt keys from user's session.
		pair, err := auth.GetService().DecryptKeysFromSession(commons.DefaultContext, gqlClient, &auth.DecryptKeysFromSessionOptions{
			Session:  response.Session,
			Password: password,
		})
		if err != nil {
			commons.Log.Debug(err)
			commons.Log.Fatal("Failed to decrypt your keys")
		}

		//	Save the public-private keys locally.
		if err := config.GetService().Save(configCommons.Keys{
			Public:  pair.PublicKey,
			Private: pair.PrivateKey,
		}, configCommons.KeysConfig); err != nil {
			commons.Log.Debug(err)
			commons.Log.Fatal("Failed to save key configuration locally")
		}
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		commons.Log.Info("You are logged in!")
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// login.Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	Cmd.Flags().StringP("email", "e", "", "Your envsecrets account email")
	Cmd.Flags().StringP("password", "p", "", "Your envsecrets account password")
}
