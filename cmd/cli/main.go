package main

import (
	_ "edgelog/bootstrap"
	cmd "edgelog/command"
)

//  Development non http Interface class service entry 
func main() {
	//   Set the operation mode to   cli(console)
	cmd.Execute()
}
