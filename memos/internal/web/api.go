// Package web は、 memosのWeb APIと画面を作る。
package web

import (
	"encoding/json/v2"
	"errors"
	"log"
	"memos/internal/memo"
	"memos/internal/store"
	"net/http"
	"time"
)

// Server は、 メモの保存先を持ち、HTTPリクエストに答える。
type Server struct {
	store store.Store
}

// NewServer は、 sを保存先にするServerを作る
func NewServer(s store.Store) *Server {
	return &Server{store: s}
}

// Handler は、 URLとハンドラの対応を登録したhttp.Handlerを返す。
func (srv *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/memos", srv.listMemos)
	mux.HandleFunc("POST /api/memos", srv.createMemo)
	mux.HandleFunc("GET /api/memos/{id}", srv.getMemo)
	mux.HandleFunc("PUT /api/memos/{id}", srv.updateMemo)
	mux.HandleFunc("DELETE /api/memos/{id}", srv.deleteMemo)
	return logRequest(mux)
}

// memoInput は、追加と更新のリクエストのJSONの形。
type memoInput struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	Pinned bool   `json:"pinned"`
}

func (srv *Server) listMemos(w http.ResponseWriter, r *http.Request) {
	memos, err := srv.store.List()
	if err != nil {
		writeError(w, err)
		return
	}
	word := r.URL.Query().Get("q")
	found := []*memo.Memo{}
	for _, m := range memos {
		if word == "" || m.Contains(word) {
			found = append(found, m)
		}
	}
	writeJSON(w, http.StatusOK, found)
}

func (srv *Server) createMemo(w http.ResponseWriter, r *http.Request) {
	var in memoInput
	if err := json.UnmarshalRead(r.Body, &in); err != nil {
		writeJSON(w, http.StatusBadRequest, errorBody{Error: "JSONを読めません：" + err.Error()})
		return
	}
	m, err := memo.New(in.Title, in.Body)
	if err != nil {
		writeError(w, err)
		return
	}
	m.Pinned = in.Pinned
	if err := srv.store.Save(m); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, m)
}

func (srv *Server) getMemo(w http.ResponseWriter, r *http.Request) {
	m, err := srv.store.Get(r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, m)
}

func (srv *Server) updateMemo(w http.ResponseWriter, r *http.Request) {
	m, err := srv.store.Get(r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return
	}
	var in memoInput
	if err := json.UnmarshalRead(r.Body, &in); err != nil {
		writeJSON(w, http.StatusBadRequest, errorBody{Error: "JSONを読めません：" + err.Error()})
		return
	}
	m.Pinned = in.Pinned
	if err := srv.store.Save(m); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, m)
}

func (srv *Server) deleteMemo(w http.ResponseWriter, r *http.Request) {
	if err := srv.store.Delete(r.PathValue("id")); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// errorBody は、エラーの時に返すJSONの型
type errorBody struct {
	Error string `json:"error"`
}

// writeJSON は、 vをJSONにして、ステータスコード statusで返す。
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(status)
	if err := json.MarshalWrite(w, v); err != nil {
		log.Printf("JSONを書けません：%v", err)
	}
}

// writeError は、 エラーの種類にあったステータスコードで、エラーのJSONを返す。
func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, store.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, memo.ErrEmptyTitle), errors.Is(err, memo.ErrTitleTooLong):
		status = http.StatusBadRequest
	}
	writeJSON(w, status, errorBody{Error: err.Error()})
}

// logRequest は、 リクエストごとに、メソッド・パス・かかった時間を記録するミドルウェア
func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start).Round(time.Millisecond))
	})
}
