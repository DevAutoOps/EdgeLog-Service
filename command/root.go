package cmd

import (
	"edgelog/command/demo"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// cli  Command based  https://github.com/spf13/cobra  encapsulation 
// RootCmd represents the base command when called without any subcommands

var RootCmd = &cobra.Command{
	Use:   "Cli",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
			examples and usage of using your application. For example:
			
			Cobra is a CLI library for Go that empowers applications.
			This application is a tool to generate the needed files
			to quickly create a Cobra application.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	//  If the subcommand exists in the subdirectory ， Then it needs to be added uniformly at the entrance ；
	//  If and  root.go  Same directory ， You do not need to add as in the next line 
	RootCmd.AddCommand(demo.Demo1)

}
