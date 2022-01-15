package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/JessicaChatBot/gpt3-engine/gpt3engine"
)

const dialogId = "vsk-9"

func main() {
	ctx := context.Background()
	fireStoreClient, err := gpt3engine.GetDefaultFirestoreClinet(ctx)
	if err != nil {
		fmt.Printf("get forestore client failed: %v\n", err)
		return
	}
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Printf("get storage client failed: %v\n", err)
		return
	}
	dialogContext, err := gpt3engine.GetInitialContext(storageClient, ctx)
	if err != nil {
		fmt.Printf("get context failed: %v\n", err)
		return
	}
	dialogContext, err = gpt3engine.PopulateContextWithAllMessages(dialogId, dialogContext, fireStoreClient, ctx)
	if err != nil {
		fmt.Printf("population of the context failed: %v\n", err)
		return
	}
	err = gpt3engine.CreateNewDialogIfAbsent(dialogId, fireStoreClient, ctx)
	if err != nil {
		fmt.Printf("error creating dialog id: %v\n", err)
		return
	}
	gpt3Client, err := gpt3engine.GetDefaultGpt3Client()
	if err != nil {
		fmt.Printf("error creating GPT3 client: %v\n", err)
		return
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("you: ")
		messageFromUserRaw, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("error getting message from user: %v\n", err)
			return
		}
		if strings.Contains(messageFromUserRaw, "exit") {
			return
		}
		if strings.Contains(messageFromUserRaw, "context") {
			fmt.Printf("context: %s\n", dialogContext)
			continue
		}
		messageFromUser := gpt3engine.Message{
			Text:   messageFromUserRaw,
			Time:   time.Now(),
			Author: "Terminal",
			Mood:   []string{"unknown"},
			Raw:    messageFromUserRaw,
		}
		messageToAddToContext := gpt3engine.ConvertToContextString(messageFromUser)
		dialogContext = fmt.Sprintf("%s\n%s", dialogContext, messageToAddToContext)
		answer, err := gpt3engine.MessageToJess(dialogContext, gpt3Client, ctx)
		if err != nil {
			fmt.Printf("error getting message from Jess: %v\n", err)
			return
		}
		fmt.Println(fmt.Sprintf("Jess: %s\n", answer.Text))
		fmt.Println(fmt.Sprintf("Jess Mood: %s\n", answer.Mood))
		messageToAddToContext = gpt3engine.ConvertToContextString(answer)
		dialogContext = fmt.Sprintf("%s\n%s", dialogContext, messageToAddToContext)

		err = gpt3engine.SaveMessage(dialogId, messageFromUser, fireStoreClient, ctx)
		if err != nil {
			fmt.Printf("error saving message: %v\n", err)
			return
		}
		err = gpt3engine.SaveMessage(dialogId, answer, fireStoreClient, ctx)
		if err != nil {
			fmt.Printf("error saving message: %v\n", err)
			return
		}
	}
}
