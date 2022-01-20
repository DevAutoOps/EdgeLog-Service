package demo

import (
	"fmt"
	"github.com/spf13/cobra"
)

//  Define subcommands 
var subCmd = &cobra.Command{
	Use:   "subCmd",
	Short: "subCmd  Brief introduction to command ",
	Long:  ` Command usage details `,
	Args:  cobra.ExactArgs(1), //   Restricted non flag Number of parameters  = 1 , exceed 1 An error will be reported 
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", args[0])
	},
}

// Register subcommand 
func init() {
	Demo1.AddCommand(subCmd)
	//  Subcommands can still be defined  flag  parameter ， For related syntax, see  demo.go  file 
}
