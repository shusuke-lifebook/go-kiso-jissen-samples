package main

import (
	"errors"
	"fmt"
	"os"
	"unicode/utf8"
)

const maxTitleLen = 20

var (
	ErrEmptyTitle   = errors.New("題が空です")
	ErrTitleTooLong = errors.New("題が長すぎます")
)

func main() {
	titles := []string{"買い物リスト", "会議のメモ", "読みたい本"}
	title := "買い物リスト"
	body := "牛乳と卵を買う"
	pinned := true

	titles = append(titles, "週末の予定")

	if err := validateTitle(title); err != nil {
		if errors.Is(err, ErrTitleTooLong) {
			fmt.Println("題を短くしてください")
		}
		fmt.Println("エラー：", err)
		os.Exit(1)
	}
	printMemo(title, body, pinned)
}

// validateTitleは、題が空でなく、maxTitleLenの文字以内かを確かめる。
func validateTitle(title string) error {
	n, _ := titleLength(title)
	if n == 0 {
		return ErrEmptyTitle
	}
	if n > maxTitleLen {
		return fmt.Errorf("%w (%d 文字。 %d 文字まで)", ErrTitleTooLong, n, maxTitleLen)
	}
	return nil
}

// titleLengthは、題の文字数とバイト数を返す。
func titleLength(title string) (int, int) {
	return utf8.RuneCountInString(title), len(title)
}

// printMemo は、メモ1件を区切りの線の下に表示する。
func printMemo(title, body string, pinned bool) {
	for i := 0; i < 24; i++ {
		fmt.Print("-")
	}
	fmt.Println()
	fmt.Printf("題：%s\n", title)
	fmt.Printf("本文：%s\n", body)
	if pinned {
		fmt.Println("ピン留め：はい")
	} else {
		fmt.Println("ピン留め：いいえ")
	}

	// n := utf8.RuneCountInString(title)
	n, b := titleLength(title)
	fmt.Printf("題の文字数：%d (バイト数：%d)\n", n, b)
}
