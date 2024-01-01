package cmd

import (
	"github.com/spf13/cobra"
)

// cobra は cli のインターフェイスを作成できる
var RootCmd = &cobra.Command{
	Use:   "task",
	Short: "Task is a CLI task manager",
}
