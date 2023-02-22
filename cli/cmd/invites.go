/*
Copyright © 2023 Mrinal Wahal mrinalwahal@gmail.com

*/
package cmd

import (
	"github.com/spf13/cobra"
)

var (
	accepted bool
)

// invitesCmd represents the invites command
var invitesCmd = &cobra.Command{
	Use:   "invites",
	Short: "A brief description of your command",
	/* 	Run: func(cmd *cobra.Command, args []string) {

	   		//	Fetch the existing invites
	   		items, err := invites.List(commons.DefaultContext, commons.GQLClient, &invites.ListOptions{
	   			Accepted: false,
	   		})
	   		if err != nil {
	   			panic(err.Error)
	   		}

	   		//	Offer acceptance selection
	   		//	Take input for project
	   		selection := promptui.Select{
	   			Label: "Choose an invite to accept",
	   			Items: *items,
	   			Templates: &promptui.SelectTemplates{
	   				Active: fmt.Sprintf("%s {{ .Organisation.Name | underline }}", promptui.IconSelect),
	   			},
	   		}

	   		_, result, er := selection.Run()
	   		if er != nil {
	   			fmt.Printf("Prompt failed %v\n", er)
	   			return
	   		}

	   		fmt.Println("Selected: ", result)
	   	},
	*/}

func init() {
	rootCmd.AddCommand(invitesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// invitesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// invitesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
