package main

// 文字列を暗号化する

import (
	"fmt"
)

func main() {
	var length, delta int
	var input string

	//scanfは 標準入力から順によみこまれる
	fmt.Scanf("%d\n", &length)
	fmt.Scanf("%s\n", &input)
	fmt.Scanf("%d\n", &delta)

	fmt.Println("lengh", length)
	fmt.Println("input", input)
	fmt.Println("delta", delta)

	// 符号化方式
	// Go言語はUTF-8を採用しているから utf-8の変換をつかえば２桁で1バイトとなる  そして 1byte から 4byteの可変長データで変換
	// code point   Utf-8
	// a : 0061   61
	// あ : 3042  E3 81 82
	// 😨 : 1F628	 F0 9F 98 A8
	// 上記の 4 byteとかになってくると 正しく文字が読めない  �
	// 文字が読めないから  code point が適切  →  rune をつかう

	var ret []rune

	for _, ch := range input { //  上記より string の range は runeが用意されている
		ret = append(ret, cipher(ch, delta))
	}
	fmt.Println(string(ret))
}

func cipher(r rune, delta int) rune {
	if r >= 'A' && r <= 'Z' { // bite が多すぎると化けてしまうから rune でコードポインタで 大文字かどうか判断できる
		return rotate(r, 'A', delta)
	}
	if r >= 'a' && r <= 'z' {
		return rotate(r, 'a', delta)
	}
	return r
}

// 暗号の差分は   大文字or 小文字の  byte を引く
func rotate(r rune, base, delta int) rune {
	tmp := int(r) - base
	tmp = (tmp + delta) % 26
	return rune(tmp + base)
}

//
// func rotate(s rune, delta int, key []rune) rune {
// 	idx := strings.IndexRune(string(key), s)
// 	if idx < 0 {
// 		panic("idx < 0")
// 	}
// 	idx = (idx + delta) % len(key)
// 	return key[idx]
// }
