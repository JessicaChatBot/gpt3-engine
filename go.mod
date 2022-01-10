module github.com/JessicaChatBot/gpt3-engine

go 1.16

require (
	cloud.google.com/go v0.97.0
	cloud.google.com/go/firestore v1.6.1 // indirect
	cloud.google.com/go/secretmanager v1.0.0 // indirect
	firebase.google.com/go v3.13.0+incompatible // indirect
	github.com/sashabaranov/go-gpt3 v0.0.0-20211215192434-7ff9fedf93e5 // indirect
	google.golang.org/api v0.59.0 // indirect
	google.golang.org/genproto v0.0.0-20211028162531-8db9c33dc351
)

replace github.com/sashabaranov/go-gpt3 => github.com/b0noi/go-gpt3 v0.0.0-20220110024631-0a3357bc7c78
