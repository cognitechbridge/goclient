/*
Copyright © 2024 Mohammad Saadatfar
*/

package file

import (
	"ctb-cli/manager"
	"fmt"

	"github.com/spf13/cobra"
)

// DownloadCmd represents the download command
var DownloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a file from cloud",
	Long:  `Download a file from cloud`,
	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("path")
		name, _ := cmd.Flags().GetString("name")

		downloader := manager.Client.NewDownloader(
			path,
			name,
		)
		err := downloader.Download()
		if err != nil {
			fmt.Printf("Error downloading:%v", err)
			return
		}
		fmt.Printf("Download completed. \n")
	},
}

func init() {
	FileCmd.AddCommand(DownloadCmd)

	DownloadCmd.Flags().StringP("name", "n", "", "name on cloud")
	DownloadCmd.Flags().StringP("path", "p", "", "path to download location")
	_ = DownloadCmd.MarkFlagRequired("name")
	_ = DownloadCmd.MarkFlagRequired("path")
}