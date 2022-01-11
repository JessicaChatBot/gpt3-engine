clean:
	rm -rf ./main

build: clean
	go build -o main

run: build
	./main

deploy-context:
	gsutil cp ./contexts/preset.md gs://contexts/preset.md

all: clean build run