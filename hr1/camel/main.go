package main

// in ファイルをスキャンして 大文字をカウントする

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// ファイルを開く
	file, err := os.Open("caesar.in")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// ファイルから一行読み込む
	scanner := bufio.NewScanner(file)
	var input string
	// for scanで全行読む  /ｎいれると大文字のカウントが 水増しされる
	for scanner.Scan() {
		input += scanner.Text() // 読み込んだ各行を結合
	}
	fmt.Println(input)
	// ファイル読み込み中のエラーをチェック
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	// 大文字をカウントする
	answer := 0
	for _, ch := range input {
		str := string(ch)
		if strings.ToUpper(str) == str {
			answer++
		}
	}
	fmt.Println(answer)
}
