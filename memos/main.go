package main

import "fmt"

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
}
