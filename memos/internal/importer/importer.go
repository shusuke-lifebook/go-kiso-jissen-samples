// Package importerは、テキストファイルからメモを取り込む
package importer

import (
	"errors"
	"fmt"
	"memos/internal/memo"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// ReadDir は、 dir の中の.txtのファイルを読んでメモにする。
// ファイルの1行目を題、2行目から後ろを本文にする。同時に読みファイルはlimit個まで。
func ReadDir(dir string, limit int) ([]*memo.Memo, error) {
	if _, err := os.Stat(dir); err != nil {
		return nil, err
	}
	paths, err := filepath.Glob(filepath.Join(dir, "*.txt"))
	if err != nil {
		return nil, err
	}
	return ParallelMap(paths, limit, readFile)
}

// readFile は、1つのファイルを読み込んでメモにする。
func readFile(path string) (*memo.Memo, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	title, body, _ := strings.Cut(string(data), "\n")
	m, err := memo.New(strings.TrimSpace(title), strings.TrimSpace(body))
	if err != nil {
		return nil, fmt.Errorf("%s:%w", filepath.Base(path), err)
	}
	return m, nil
}

// ParallelMap は、itemsの値ごとにfをgoroutineで呼び、itemsと同じ順の結果を返す。
// 同時に動かすfはlimit個まで。エラーがあれば、すべてのエラーをまとめて返す。
func ParallelMap[T, U any](items []T, limit int, f func(T) (U, error)) ([]U, error) {
	results := make([]U, len(items))
	errs := make([]error, len(items))
	sem := make(chan struct{}, limit)
	var wg sync.WaitGroup
	for i, item := range items {
		wg.Go(func() {
			sem <- struct{}{}
			defer func() { <-sem }()
			results[i], errs[i] = f(item)
		})
	}
	wg.Wait()
	if err := errors.Join(errs...); err != nil {
		return nil, err
	}
	return results, nil
}
