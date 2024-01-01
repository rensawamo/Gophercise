package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rensawamo/cli-app/cmd"
	"github.com/rensawamo/cli-app/db"
	homedir "github.com/mitchellh/go-homedir"
)

func main() {

	// 違うディレクトリを呼ぶパターン
	home, _ := homedir.Dir()
	dbPath := filepath.Join(home, "tasks.db")
	must(db.Init(dbPath))
	must(cmd.RootCmd.Execute())
}

func must(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
