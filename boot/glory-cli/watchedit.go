package main

import (
	"github.com/spf13/cobra"
)

var watchEdit = &cobra.Command{
	Use: "watchEdit",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(watchEdit)
}
