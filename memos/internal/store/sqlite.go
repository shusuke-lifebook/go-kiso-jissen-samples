package store

import (
	"database/sql"
	"errors"
	"fmt"
	"memos/internal/memo"
	"time"

	_ "modernc.org/sqlite" // database/sqlに"sqlite"のドライバを登録する
)

// SqlStore は、 メモをSQLiteのデータベースに保存する保存先
type SqlStore struct {
	db *sql.DB
}

var _ Store = (*SqlStore)(nil)

const schema = `
CREATE TABLE IF NOT EXISTS memos (
	id TEXT PRIMARY KEY,
	title TEXT NOT NULL,
	body TEXT NOT NULL,
	pinned INTEGER NOT NULL,
	created_at TEXT NOT NULL,
	updated_at TEXT NOT NULL
)`

// NewSqlStore は、 SQLiteのファイルを開き、表が無ければ作る
// busy_timeout と WALは、同時に書き込んだ時に待つための設定。接続ごとに効くようDSNに書く
func NewSqlStore(path string) (*SqlStore, error) {
	dsn := path + "?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("表を作れません：%w", err)
	}
	return &SqlStore{db: db}, nil
}

// Close は、 データベースを閉じる
func (s *SqlStore) Close() error {
	return s.db.Close()
}

// Save は、 メモを保存する。同じIDのメモがあれば書き換える。
func (s *SqlStore) Save(m *memo.Memo) error {
	_, err := s.db.Exec(`
	INSERT INTO memos (id, title, body, pinned, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET
		title = excluded.title,
		body = excluded.body,
		pinned = excluded.pinned,
		created_at = excluded.created_at,
		updated_at = excluded.updated_at
	`, m.ID, m.Title, m.Body, m.Pinned, m.CreatedAt.Format(time.RFC3339Nano), m.UpdatedAt.Format(time.RFC3339Nano))
	return err
}

// Get は、 IDのメモを返す。なければErrNotFouncをラップして返す。
func (s *SqlStore) Get(id string) (*memo.Memo, error) {
	row := s.db.QueryRow(`SELECT id, title, body, pinned, created_at, updated_at FROM memos WHERE id =?`, id)
	m, err := scanMemo(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w：%s", ErrNotFound, id)
	}
	return m, nil
}

// List は、 全てのメモを並べて返す。
func (s *SqlStore) List() ([]*memo.Memo, error) {
	rows, err := s.db.Query(`SELECT id, title, body, pinned, created_at, updated_at FROM memos`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var memos []*memo.Memo
	for rows.Next() {
		m, err := scanMemo(rows)
		if err != nil {
			return nil, err
		}
		memos = append(memos, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	memo.Sort(memos)
	return memos, nil
}

// Delete は、 IDのメモを消す。無ければErrNotFoundをラップして返す。
func (s *SqlStore) Delete(id string) error {
	res, err := s.db.Exec(`DELETE FROM memos WHERE id=?`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("%w：%s", ErrNotFound, id)
	}
	return nil
}

// Scanner は、 *sql.Row と *sql.Rows の両方が持つScanのメソッドだけのインタフェース。
type Scanner interface {
	Scan(dest ...any) error
}

// scanMemo は、 1行を読んでメモにする
func scanMemo(r Scanner) (*memo.Memo, error) {
	var m memo.Memo
	var created, updated string
	if err := r.Scan(&m.ID, &m.Title, &m.Body, &m.Pinned, &created, &updated); err != nil {
		return nil, err
	}
	var err error
	if m.CreatedAt, err = time.Parse(time.RFC3339Nano, created); err != nil {
		return nil, err
	}
	if m.UpdatedAt, err = time.Parse(time.RFC3339Nano, updated); err != nil {
		return nil, err
	}
	return &m, nil
}
