clean:
	rm -rf ./main

build: clean
	go build -o main

run: build
	./main

all: clean build run