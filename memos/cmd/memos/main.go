// memos は、 メモをJSONファイルmemos.jsonに保存するコマンド。
//
// memos add <題> <本文> メモを足す
// memos list [言葉] 一覧を表示する。言葉を渡すと、その言葉を含むメモだけ
// memos delete <IDの先頭> メモを消す
// memos import <ディレクトリ> ディレクトリの.txtのファイルを、メモとしてまとめて取り込む
package main

import (
	"errors"
	"fmt"
	"io"
	"memos/internal/importer"
	"memos/internal/memo"
	"memos/internal/store"
	"os"
	"strings"
	"time"
)

const dataFile = "memos.json"

func main() {
	s := store.NewFileStore(dataFile)
	if err := run(os.Stdout, s, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "エラー：", err)
		os.Exit(1)
	}
}

// run は、 argsの最初の言葉で処理を選び、結果をwに書き出す。
func run(w io.Writer, s store.Store, args []string) error {
	if len(args) == 0 {
		return list(w, s, "")
	}
	switch args[0] {
	case "add":
		if len(args) != 3 {
			return errors.New("使い方：memos add <題> <本文>")
		}
		return add(w, s, args[1], args[2])
	case "list":
		word := ""
		if len(args) > 1 {
			word = args[1]
		}
		return list(w, s, word)
	case "delete":
		if len(args) != 2 {
			return errors.New("使い方：memos delete <IDの先頭>")
		}
		return remove(w, s, args[1])
	case "import":
		if len(args) != 2 {
			return errors.New("使い方：memos import <ディレクトリ>")
		}
		return importDir(w, s, args[1])
	default:
		return fmt.Errorf("知らないコマンドです。:%s (add・list・delete・importのどれか)", args[0])
	}
}

// add は、 新しいメモを作って保存する。
func add(w io.Writer, s store.Store, title, body string) error {
	m, err := memo.New(title, body)
	if err != nil {
		return err
	}
	if err := s.Save(m); err != nil {
		return err
	}
	fmt.Fprintf(w, "追加しました：%s %s\n", m.ID[:8], m.Title)
	return nil
}

// list は、 wordを含むメモを1件ずつ書き出す。wordが空なら全部
func list(w io.Writer, s store.Store, word string) error {
	memos, err := s.List()
	if err != nil {
		return err
	}
	found := 0
	for _, m := range memos {
		if word != "" && !m.Contains(word) {
			continue
		}
		fmt.Fprintln(w, m)
		found++
	}
	fmt.Fprintf(w, "%d件のメモ（全部で%d件）\n", found, len(memos))
	return nil
}

// remove は、 IDがprefixで始まるメモを1件消す。
func remove(w io.Writer, s store.Store, prefix string) error {
	memos, err := s.List()
	if err != nil {
		return err
	}
	var hit []*memo.Memo
	for _, m := range memos {
		if strings.HasPrefix(m.ID, prefix) {
			hit = append(hit, m)
		}
	}
	switch len(hit) {
	case 0:
		return fmt.Errorf("%w:%s", store.ErrNotFound, prefix)
	case 1:
		if err := s.Delete(hit[0].ID); err != nil {
			return err
		}
		fmt.Fprintf(w, "消しました：%s %s\n", hit[0].ID[:8], hit[0].Title)
		return nil
	default:
		return fmt.Errorf("%s で始まるメモが%d件あります。IDをもっと長く書いてください", prefix, len(hit))
	}
}

// importDir は、dirの.txtのファイルを同時に読み込んでメモにし、保存する。
func importDir(w io.Writer, s store.Store, dir string) error {
	start := time.Now()
	memos, err := importer.ReadDir(dir, 8)
	if err != nil {
		return err
	}
	for _, m := range memos {
		if err := s.Save(m); err != nil {
			return err
		}
	}
	fmt.Fprintf(w, "%d件を取り込みました(%v) \n", len(memos), time.Since(start).Round(time.Millisecond))
	return nil
}
