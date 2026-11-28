# go-kiso-jissen-samples

- [go-kiso-jissen-samples](https://github.com/fairview2311/go-kiso-jissen-samples/)
- カバレッジレポート
  - go test -coverprofile=cover.out ./...
  - go tool cover -func=cover.out


## 本サンプルの起動
- サーバーを起動する
  - go run ./cmd/memos serve
  - 稼働確認
    - curl -s http://localhost:8080/api/memos
    - curl -s -i -X POST -d '@testdata/new-memo.json' http://localhost:8080/api/memos
    - curl -s http://localhost:8080/api/memos/<id>
      - curl -s http://localhost:8080/api/memos/ccb57b20-7ad9-4969-87d3-10a431634150
    - curl -s -X PUT -d '@testdata/edit-memo.json' http://localhost:8080/api/memos/<id>
      - curl -s -X PUT -d '@testdata/edit-memo.json' http://localhost:8080/api/memos/ccb57b20-7ad9-4969-87d3-10a431634150
    - curl -s -i -X DELETE http://localhost:8080/api/memos/<id>
      - curl -s -i -X DELETE http://localhost:8080/api/memos/ccb57b20-7ad9-4969-87d3-10a431634150

## SQLiteのドライバ
- go get modernc.org/sqlite