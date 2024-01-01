package cmd

import (
	"fmt"
	"strings"

	"github.com/rensawamo/cli-app/db"
	"github.com/spf13/cobra"
)

//Golangでの文字列変換 数値変換 (strconvパッケージ)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a task to your task list.",
	Run: func(cmd *cobra.Command, args []string) { // コブラの 標準出力をarg でうけとｆ
		task := strings.Join(args, " ") // スペースでも 空白開けて一つの文字にしてしまう
		_, err := db.CreateTask(task)
		if err != nil {
			fmt.Println("Something went wrong:", err)
			return
		}
		fmt.Printf("Added \"%s\" to your task list.\n", task)
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}
