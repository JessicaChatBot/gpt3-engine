clean:
	rm -rf ./main

build: clean
	go build -o main

run: build
	./main

deploy-context:
	gsutil cp ./contexts/preset.md gs://contexts/preset.md

deploy-functions:
	gcloud functions deploy "jess-chat" \
		--runtime go116 \
		--project "jessdb-337700" \
		--allow-unauthenticated \
		--trigger-http \
		--entry-point="functions.Say" \
		--service-account="jessapi@jessdb-337700.iam.gserviceaccount.com"

all: clean build run