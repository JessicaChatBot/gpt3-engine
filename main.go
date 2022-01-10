package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/JessicaChatBot/gpt3-engine/gpt3engine"
)

func main() {
	ctx := context.Background()
	client, err := gpt3engine.GetDefaultFirestoreClinet(ctx)
	if err != nil {
		fmt.Printf("get db client failed: %v\n", err)
		return
	}
	dialogContext, err := gpt3engine.GetInitialContext()
	dialogId := "vsk-1"
	dialogContext, err = gpt3engine.PopulateContextWithAllMessages(dialogId, dialogContext, client, ctx)
	fmt.Println(dialogContext)
	if err != nil {
		fmt.Printf("get context failed: %v\n", err)
		return
	}
	messageFromUser := ""
	reader := bufio.NewReader(os.Stdin)
	err = gpt3engine.CreateNewDialogIfAbsent(dialogId, client, ctx)
	if err != nil {
		fmt.Printf("error creating dialog id: %v\n", err)
		return
	}
	for {
		fmt.Print("you: ")
		messageFromUser, err = reader.ReadString('\n')
		if err != nil {
			fmt.Printf("error getting message from user: %v\n", err)
			return
		}
		if strings.Contains(messageFromUser, "exit") {
			return
		}
		currentTime := time.Now()
		messageToAddToContext := fmt.Sprintf("\n[Friend][%s]: %s\n[END]\n", currentTime.Format("2006 Jan 2 15:04:05"), messageFromUser)
		dialogContext = dialogContext + messageToAddToContext
		answer, err := gpt3engine.MessageToJess(dialogContext)
		if err != nil {
			fmt.Printf("error getting message from Jess: %v\n", err)
			return
		}
		fmt.Println(fmt.Sprintf("Jess: %s\n", answer.Text))
		fmt.Println(fmt.Sprintf("Jess Mood: %s\n", answer.Mood))
		dialogContext = dialogContext + fmt.Sprintf("\n%s\n", answer.Raw)
		gpt3engine.SaveMessage(dialogId, gpt3engine.Message{
			Text:   messageFromUser,
			Time:   time.Now(),
			Author: "Friend",
			Mood:   []string{"unknown"},
			Raw:    messageFromUser,
		}, client, ctx)
		gpt3engine.SaveMessage(dialogId, answer, client, ctx)
	}
}
