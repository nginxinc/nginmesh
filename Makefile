DOCKER_RUN = docker run --rm -v $(shell pwd):/go/src/github.com/nginxinc/nginmesh -w /go/src/github.com/nginxinc/nginmesh/
GOLANG_CONTAINER = golang:1.8

darwin:
	mkdir -p build
	$(DOCKER_RUN) -e GOOS=darwin $(GOLANG_CONTAINER) go build -o build/nginx-inject cmd/inject/*.go

linux:
	mkdir -p build
	$(DOCKER_RUN) -e GOOS=linux $(GOLANG_CONTAINER) go build -o build/nginx-inject-linux cmd/inject/*.go

clean:
	rm -rf build