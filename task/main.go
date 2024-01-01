package main

import "github.com/rensawamo/task/cmd"

func main() {

	// コブラを参照して  cmdでaddで  go run main.go add とかで 自作したcommandを  呼び込める		
	cmd.RootCmd.Execute()
}
