package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"  // 自分のcurrentのファイル構造を探査
	"regexp"
)
// [] の中に特定の文字や数値を忍び込ませる $で終わり
// この正規表現の読み方は 特定の文字からスタートして 4桁の数値 と (数値 of 数値)にして . で拡張子を持つものという意味
var re = regexp.MustCompile("^(.+?) ([0-9]{4}) [(]([0-9]+) of ([0-9]+)[)][.](.+?)$")
var replaceString = "$2 - $1 - $3 of $4.$5"

func main() {
	var dry bool
	flag.BoolVar(&dry, "dry", true, "whether or not this should be a real or dry run")
	flag.Parse()

	walkDir := "sample"
	var toRename []string
	
	filepath.Walk(walkDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {  // 今回はファイルしかみない
			return nil
		}
		if _, err := match(info.Name()); err == nil {
			toRename = append(toRename, path)
		}
		return nil
	})
	for _, oldPath := range toRename {
		dir := filepath.Dir(oldPath)
		fmt.Println("dir",dir)
		filename := filepath.Base(oldPath)
		fmt.Println("filename",filename)
		newFilename, _ := match(filename)
		fmt.Println("newfilename",newFilename)
		newPath := filepath.Join(dir, newFilename)
		fmt.Printf("mv %s => %s\n", oldPath, newPath)
		if !dry {
			err := os.Rename(oldPath, newPath)
			if err != nil {
				fmt.Println("Error renaming:", oldPath, newPath, err.Error())
			}
		}
	}
}

// match returns the new file name, or an error if the file name
// didn't match our pattern.
// 正規表現は 複雑な文字列の解析 ログの解析や  文書処理などにつかうことができる
func match(filename string) (string, error) {
	if !re.MatchString(filename) {
		fmt.Println("正規表現にマッチしていません")
		return "", fmt.Errorf("%s didn't match our pattern", filename)
	}

	fmt.Println("remo",re.ReplaceAllString(filename, replaceString))
	// 順番を変える
	return re.ReplaceAllString(filename, replaceString), nil
}
