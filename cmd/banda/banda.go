package main

import (
	"os"

	banda_command "github.com/fikrirnurhidayat/banda-lumaksa/internal/infra/command/banda"
	"github.com/spf13/cobra"
)

var bandaCmd = &cobra.Command{
	Use:   "banda",
	Short: "Financial Tracker",
	Long:  "Financial Tracker",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func main() {
	err := bandaCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	bandaCmd.AddCommand(banda_command.InitCmd)
	bandaCmd.AddCommand(banda_command.ServeCmd)
}
