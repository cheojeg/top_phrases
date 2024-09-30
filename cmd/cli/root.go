package cli

import (
	"github.com/spf13/cobra"
)

//var rootCmd = &cobra.Command{
//	Use:   "app",
//	Short: "A brief description of your application",
//	Long:  `A longer description of your application with examples and usage.`,
//	Run: func(cmd *cobra.Command, args []string) {
//		// Main logic of the application
//		fmt.Println("Running main application logic")
//	},
//}
//
//func Execute() {
//	if err := rootCmd.Execute(); err != nil {
//		fmt.Println(err)
//		os.Exit(1)
//	}
//}

func NewCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "Top Phrases",
		Short: "Root command",
	}

	cmd.AddCommand(newCmdBot())
	return cmd
}
