static: cmd/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-extldflags "-static"' -o main cmd/main.go

image: static
	docker buildx build --platform linux/amd64 -f Dockerfile -t xxx.xxx.xxx/spider:v1 .

clean:
	rm main

push: image
	docker image push xxx.xxx.xxx/spider:v1
