/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"ctb-cli/app"

	"github.com/spf13/cobra"
)

// shareCmd represents the share command
var shareCmd = &cobra.Command{
	Use:   "share",
	Short: "Share files with other users",
	Long: `This command shares files with other users.
	The files are shared with the user who has the corresponding private key.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		initKey()
		pattern := args[0]
		recipient, _ := cmd.Flags().GetString("recipient")
		res := app.Share(pattern, recipient)
		MarshalOutput(res)
	},
}

func init() {
	rootCmd.AddCommand(shareCmd)
	SetRequiredKeyFlag(shareCmd)
	shareCmd.PersistentFlags().StringP("recipient", "r", "", "recipient public key. Required.")
	shareCmd.MarkPersistentFlagRequired("recipient")
}
