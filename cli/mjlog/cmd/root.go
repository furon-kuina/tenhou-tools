package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "tenhou-log",
	Short: "tenhou-log fetches Tenhou game log and save it to local storage",
	Run: func(cmd *cobra.Command, args []string) {

	},
}
