package main

import (
	"bytes"
	"memos/internal/store"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	s := store.NewMemoryStore()
	var out bytes.Buffer

	if err := run(&out, s, []string{"add", "読みたい本", "Goの本2冊"}); err != nil {
		t.Fatalf("addのエラー=%v", err)
	}
	if !strings.Contains(out.String(), "追加しました：") {
		t.Errorf("addの出力%q", out.String())
	}

	out.Reset()
	if err := run(&out, s, []string{"list", "本"}); err != nil {
		t.Fatalf("listの出力=%q", out.String())
	}
	if err := run(&out, s, []string{"remove"}); err == nil {
		t.Error("知らないサブコマンドなのにエラーにならない")
	}
}
