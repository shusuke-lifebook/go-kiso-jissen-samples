package store

import (
	"errors"
	"fmt"
	"memos/internal/memo"
	"path/filepath"
	"sync"
	"testing"
)

// TestStores は、同じテストをMemoryStore と FileStoreの両方で動かす。
func TestStores(t *testing.T) {
	stores := []struct {
		name string
		s    Store
	}{
		{"MemoryStore", NewMemoryStore()},
		{"FileStore", NewFileStore(filepath.Join(t.TempDir(), "memos.json"))},
		{"SqlStore", newTestSqlStore(t)},
	}
	for _, tt := range stores {
		t.Run(tt.name, func(t *testing.T) {
			testStore(t, tt.s)
		})
	}
}

func testStore(t *testing.T, s Store) {
	t.Helper()
	list, err := s.List()
	if err != nil || len(list) != 0 {
		t.Fatalf("最初のList=%d件、%v、期待したのは0件", len(list), err)
	}

	m, err := memo.New("買い物リスト", "牛乳と卵を買う")
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Save(m); err != nil {
		t.Fatalf("Saveのエラー=%v", err)
	}
	got, err := s.Get(m.ID)
	if err != nil {
		t.Fatalf("Getのエラー=%v", err)
	}
	if got.Title != "買い物リスト" {
		t.Errorf("Getした題%q、期待したのは%q", got.Title, "買い物リスト")
	}
	if err := s.Delete(m.ID); err != nil {
		t.Fatalf("Deleteのエラー=%v", err)
	}
	if _, err := s.Get(m.ID); !errors.Is(err, ErrNotFound) {
		t.Errorf("消した後のGetのエラー=%v、期待したのはErrNotFound", err)
	}
	if err := s.Delete(m.ID); !errors.Is(err, ErrNotFound) {
		t.Errorf("2回目のDeleteのエラー=%v、期待したのはErrNotFound", err)
	}
}

// TestFileStoreKeepsMemoは、FileStoreに保存したメモが、別のFileStoreからも読めることを確かめる。
func TestFileStoreKeepsMemos(t *testing.T) {
	path := filepath.Join(t.TempDir(), "memos.json")
	m, err := memo.New("会議のメモ", "来週の会議で予算を決める")
	if err != nil {
		t.Fatal(err)
	}
	if err := NewFileStore(path).Save(m); err != nil {
		t.Fatal(err)
	}
	got, err := NewFileStore(path).Get(m.ID)
	if err != nil {
		t.Fatalf("別のFileStoreからのGetのエラー=%v", err)
	}
	if got.Body != m.Body || !got.CreatedAt.Equal(m.CreatedAt) {
		t.Errorf("読み直したメモ=%q%v、期待したのは%q%v", got.Body, got.CreatedAt, m.Body, m.CreatedAt)
	}
}

// newTestSqlStore は、 テスト用のSQLiteのファイルを一時ディレクトリに作る。テストが終わると閉じる
func newTestSqlStore(t *testing.T) *SqlStore {
	t.Helper()
	s, err := NewSqlStore(filepath.Join(t.TempDir(), "memos.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

// TestSqlStoreConcurrentSave は、 50個のgroutineから同時に保存しても、50件とも残ることを確かめる。
func TestSqlStoreConcurrentSave(t *testing.T) {
	s := newTestSqlStore(t)
	var wg sync.WaitGroup
	for i := range 50 {
		wg.Go(func() {
			m, err := memo.New(fmt.Sprintf("メモ%d", i), "同時に保存する")
			if err != nil {
				t.Error(err)
				return
			}
			if err := s.Save(m); err != nil {
				t.Errorf("Saveのエラー=%v", err)
			}
		})
	}
	wg.Wait()
	list, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 50 {
		t.Errorf("保存できたメモ=%d件、期待したのは50件", len(list))
	}
}
