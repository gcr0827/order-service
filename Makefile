.PHONY: run build test race tidy fmt vet clean migrate

# 导入建表 SQL。★ --default-character-set=utf8mb4 不能省：
# 容器内 mysql 客户端的 default-character-set 是 auto，而容器通常没有 LANG，
# auto 会退化成 latin1，导致文件里的 UTF-8 中文落库后变成双重编码的乱码。
MYSQL_CONTAINER ?= order-mysql
MYSQL_USER      ?= root
MYSQL_PASSWORD  ?= dev123456
MYSQL_DATABASE  ?= order_db

migrate:
	docker exec -i $(MYSQL_CONTAINER) mysql -u$(MYSQL_USER) -p$(MYSQL_PASSWORD) \
	  --default-character-set=utf8mb4 $(MYSQL_DATABASE) < migrations/001_init.sql

run:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

test:
	go test ./... -v

race:
	go test ./... -race

tidy:
	go mod tidy

fmt:
	gofmt -w .

vet:
	go vet ./...

clean:
	rm -rf bin/
