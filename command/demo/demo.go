package demo

import (
	"edgelog/app/global/variable"
	"github.com/spf13/cobra"
)

// Demo Sample file ， Let's assume a scenario ：
//  Specified by a command   Search Engines ( Baidu 、 Sogou 、 Google )、 Search type （ text 、 picture ）、 key word   Execute a series of commands 

var (
	// 1. Define a variable ， Receive search engine （ Baidu 、 Sogou 、 Google ）
	SearchEngines string
	// 2. Type of search ( picture 、 written words )
	SearchType string
	// 3. key word 
	KeyWords string
)

var logger = variable.ZapLog.Sugar()

//  Define command 
var Demo1 = &cobra.Command{
	Use:     "sousuo",
	Aliases: []string{"sou", "ss", "s"}, //  Define alias 
	Short:   " This is a Demo， Demonstrate business logic with search content ...",
	Long: ` Call method ：
			1. Enter the project root directory （Ginkeleton）。 
			2. implement  go  run  cmd/cli/main.go sousuo -h  // You can view the user guide 
			3. implement  go  run  cmd/cli/main.go sousuo  Baidu   //  Run one quickly Demo
			4. implement  go  run  cmd/cli/main.go  sousuo  Baidu  -K  key word   -E  baidu -T img    //  Run with specified parameters Demo
		`,
	//Args:    cobra.ExactArgs(2),  //   Restricted non flag parameter （ Also known as position parameters ） The number of must be equal to  2 , Otherwise, an error will be reported 
	// Run Pre functions of commands and subcommands 
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// If you only want to be a callback for subcommands ， It can be judged by relevant parameters ， Only in subcommand execution 
		logger.Infof("Run Pre method of function subcommand ， Position parameters ：%v ，flag parameter ：%s, %s, %s \n", args[0], SearchEngines, SearchType, KeyWords)
	},
	// Run Pre function of command 
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.Infof("Run Pre method of function ， Position parameters ：%v ，flag parameter ：%s, %s, %s \n", args[0], SearchEngines, SearchType, KeyWords)

	},
	// Run  The command is   core   command ， The remaining commands serve this command ， Can delete ， You are free to choose 
	Run: func(cmd *cobra.Command, args []string) {
		//args   Parameter represents non flag（ Also known as position parameters ）， This parameter is stored as an array by default 。
		//fmt.Println(args)
		start(SearchEngines, SearchType, KeyWords)
	},
	// Run Post function of command 
	PostRun: func(cmd *cobra.Command, args []string) {
		logger.Infof("Run Post method of function ， Position parameters ：%v ，flag parameter ：%s, %s, %s \n", args[0], SearchEngines, SearchType, KeyWords)
	},
	// Run Post functions of commands and subcommands 
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		// If you only want to be a callback for subcommands ， It can be judged by relevant parameters ， Only in subcommand execution 
		logger.Infof("Run Post method of function subcommand ， Position parameters ：%v ，flag parameter ：%s, %s, %s \n", args[0], SearchEngines, SearchType, KeyWords)
	},
}

//  Registration command 、 Initialization parameters 
func init() {
	Demo1.AddCommand(subCmd)
	Demo1.Flags().StringVarP(&SearchEngines, "Engines", "E", "baidu", "-E  perhaps  --Engines  Select search engine ， for example ：baidu、sogou")
	Demo1.Flags().StringVarP(&SearchType, "Type", "T", "img", "-T  perhaps  --Type  Select the type of content to search for ， for example ： Picture class ")
	Demo1.Flags().StringVarP(&KeyWords, "KeyWords", "K", " key word ", "-K  perhaps  --KeyWords  Search keywords ")
	//Demo1.Flags().BoolP(1,2,3,5)  // receive bool Type parameter 
	//Demo1.Flags().Int64P()  // receive int type 
}

// Start execution 
func start(SearchEngines, SearchType, KeyWords string) {

	logger.Infof(" Search engine you entered ：%s，  Search type ：%s,  key word ：%s\n", SearchEngines, SearchType, KeyWords)

}
