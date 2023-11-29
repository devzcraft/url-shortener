run:
	go run ./cmd/url-shortener --config=./config/local.yaml
migrate:
	go run ./cmd/migrator --storage-path=./storage/storage.db --migrations-path=./migrations
gen:
	go generate ./...
test:
	go test -v -race ./...
