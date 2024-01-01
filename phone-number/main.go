package main

import (
	"bytes"
	"fmt"
)

// func normalize(phone string) string {
// 	re := regexp.MustCompile("\\D")
// 	return re.ReplaceAllString(phone, "")
// }

func normalize(phone string) string {
	var buf bytes.Buffer // この型は 効率的な文字の結合などを援助   + などで 結合すると  メモリの消費がとてもおおい
	for _, ch := range phone {
		if ch >= '0' && ch <= '9' { // rune を用いて数値かどうか判定
			buf.WriteRune(ch)
		}
	}
	fmt.Println(buf)
	return buf.String()  // Buffer 構造体の bufを stringにしている
}
