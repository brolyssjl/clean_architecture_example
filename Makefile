.Phony: begin
begin:
	@~/.air -d -c .air.conf

.phony: build
build:
	go build -o bin/main api/main.go