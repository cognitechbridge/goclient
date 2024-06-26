/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// mountCmd represents the mount command
var mountCmd = &cobra.Command{
	Use:   "mount",
	Short: "Mount",
	Long:  `Mount the file system. This command mounts the file system and blocks the terminal.`,
	Run: func(cmd *cobra.Command, args []string) {
		res := ctbApp.PrepareMount(encryptedPrivateKey)
		MarshalOutput(res)
		fmt.Fprint(os.Stdout, "/**********************************\n")
		ctbApp.Mount()
	},
}

func init() {
	rootCmd.AddCommand(mountCmd)
	SetRequiredKeyFlag(mountCmd)
}
