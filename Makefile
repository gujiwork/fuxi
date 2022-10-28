.PHONY: run build image push clean

tag = v0.1
releaseName = fuxi
dockerhubUser = dnsjia

ALL: run

run: build
	./fuxi

build:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o $(releaseName) ./cmd/main.go

image:
	docker build -t $(dockerhubUser)/$(releaseName):$(tag) .

push: image
	docker push $(dockerhubUser)/$(releaseName):$(tag)

clean:
	-rm -f ./$(releaseName)