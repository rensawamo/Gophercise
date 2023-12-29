package main

//

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
)

func main() {

	// csvの読み込み
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")

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

	correct := 0
	for i, p := range problems {
		fmt.Printf("Problem %d: %s = \n", i+1, p.q)
		var answer string
		fmt.Scanf("%s\n", &answer) // 標準入力まち  入力がanswer に入る
		if answer == p.a {
			correct++
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
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
