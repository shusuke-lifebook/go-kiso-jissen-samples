package main

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
	"unicode/utf8"
)

const maxTitleLen = 20

var (
	ErrEmptyTitle   = errors.New("題が空です")
	ErrTitleTooLong = errors.New("題が長すぎます")
)

// Memo は、メモ1件を表す。
type Memo struct {
	Title     string
	Body      string
	Pinned    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewMemo は、題を確かめてから新しいメモを作る。
func NewMemo(title, body string) (*Memo, error) {
	if err := validateTitle(title); err != nil {
		return nil, err
	}
	now := time.Now()
	return &Memo{
		Title:     title,
		Body:      body,
		Pinned:    false,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Contains は、題か本文に word が含まれるかを返す。
func (m *Memo) Contains(word string) bool {
	return strings.Contains(m.Title, word) || strings.Contains(m.Body, word)
}

// Print は、メモ1件を区切りの線の下に表示する。
func (m *Memo) Print() {
	fmt.Println(strings.Repeat("-", 24))
	fmt.Printf("題：%s\n", m.Title)
	fmt.Printf("本文：%s\n", m.Body)
	if m.Pinned {
		fmt.Println("ピン留め： はい")
	} else {
		fmt.Println("ピン留め： いいえ")
	}
	fmt.Printf("更新: %s\n", m.UpdatedAt.Format("2006-01-02 15:04"))
}

func main() {
	titles := []string{"買い物リスト", "会議のメモ", "読みたい本"}
	bodies := map[string]string{
		"買い物リスト": "牛乳と卵を買う",
		"会議のメモ":  "来週の会議で予算を決める",
		"読みたい本":  "Go の本を2冊",
	}
	pinned := map[string]bool{
		"会議のメモ": true,
	}

	titles = append(titles, "週末の予定")
	bodies["週末の予定"] = "本屋に行く"

	for _, title := range titles {
		if err := validateTitle(title); err != nil {
			fmt.Println("エラー：", err)
			os.Exit(1)
		}
	}

	slices.SortFunc(titles, func(a, b string) int {
		if pinned[a] != pinned[b] {
			if pinned[a] {
				return -1
			}
			return 1
		}
		return strings.Compare(a, b)
	})

	word := ""
	if len(os.Args) > 1 {
		word = os.Args[1]
	}
	found := 0
	for _, title := range titles {
		body := bodies[title]
		if word != "" && !strings.Contains(title, word) && !strings.Contains(body, word) {
			continue
		}
		printMemo(title, body, pinned[title])
		found++
	}
	fmt.Printf("%d 件のメモ (全部で %d 件)\n", found, len(titles))
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
