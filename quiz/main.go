package main

//

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {

	// csvの読み込み
	// flag.String("オプション名", "初期値", "説明")
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 300, "the time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}

	// problem[]
	problems := parseLines(lines)
	// 時間制限を追加する
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0 // 正解数

problemloop:
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)
		ch := make(chan string)
		go func() {
			var answer string
			fmt.Scan("%s\n", &answer)
			ch <- answer // 入力されたらchに流し込む
		}()

		select {
		// タイムアップでchに流れこんだとき
		case <-timer.C:
			fmt.Println()
			break problemloop
		case answer := <-ch:
			if answer == p.a {
				correct++
			}
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))

	// for i, p := range problems {
	// 	fmt.Printf("Problem %d: %s = \n", i+1, p.q)
	// 	var answer string
	// 	fmt.Scanf("%s\n", &answer) // 標準入力まち  入力がanswer に入る
	// 	if answer == p.a {
	// 		correct++
	// 	}
	// }

	// fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines)) // 問題と解答がセットの連想配列を作成
	for i, line := range lines {
		ret[i] = problem{ // problemの初期化
			q: line[0],
			a: line[1],
		}
	}
	return ret
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
