SERVICE_NAME = test_rss_service
PID          = /tmp/$(SERVICE_NAME).pid
BINARY_PATH  = bin/$(SERVICE_NAME)
HTTP_PORT 	:= $(if $(HTTP_PORT),$(HTTP_PORT),8081)

run: restart
	fswatch -o cmd pkg config | xargs -n1 -I{} make restart || make kill

kill:
	kill `cat $(PID)` || true

build:
	GO111MODULE=on go build -v -o $(BINARY_PATH) ./cmd/rssservice/main.go

restart: kill build
	HTTP_PORT=$(HTTP_PORT) $(BINARY_PATH) & echo $$! > $(PID)

fmt:
	go fmt ./...

test:
	APP_ENV=test go test --race -p=1 -count=1 ./...
