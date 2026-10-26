package main

import (
	"fmt"
	"unicode/utf8"
)

const maxTitleLen = 20

func main() {
	title := "買い物リスト"
	body := "牛乳と卵を買う"
	pinned := true
	var views int

	for i := 0; i < 24; i++ {
		fmt.Print("-")
	}
	fmt.Println("")

	fmt.Printf("題：%s\n", title)
	fmt.Printf("本文：%s\n", body)
	fmt.Printf("ピン留め：%t\n", pinned)
	fmt.Printf("表示回数：%d\n", views)
	fmt.Printf("題の長さの上限：%d\n", maxTitleLen)

	if pinned {
		fmt.Println("ピン留め：はい")
	} else {
		fmt.Println("ピン留め：いいえ")
	}

	n := utf8.RuneCountInString(title)
	switch {
	case n == 0:
		fmt.Println("題が空です")
	case n > maxTitleLen:
		fmt.Printf("題が長すぎます (%d文字まで)\n", maxTitleLen)
	default:
		fmt.Printf("題の文字数：%d (バイト数：%d)\n", n, len(title))
	}
}
