package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/JessicaChatBot/gpt3-engine/gpt3engine"
)

func main() {
	context, err := gpt3engine.GetInitialContext()
	dialogId := "vsk-1"
	context, err = gpt3engine.PopulateContextWithAllMessages(dialogId, context)
	fmt.Println(context)
	if err != nil {
		fmt.Printf("get context failed: %v\n", err)
		return
	}
	messageFromUser := ""
	reader := bufio.NewReader(os.Stdin)
	err = gpt3engine.CreateNewDialogIfAbsent(dialogId)
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
		context = context + messageToAddToContext
		answer, err := gpt3engine.MessageToJess(context)
		if err != nil {
			fmt.Printf("error getting message from Jess: %v\n", err)
			return
		}
		fmt.Println(fmt.Sprintf("Jess: %s\n", answer.Message.Text))
		fmt.Println(fmt.Sprintf("Jess Mood: %s\n", answer.Mood))
		context = context + fmt.Sprintf("\n%s\n", answer.Raw)
		gpt3engine.SaveMessage(dialogId, gpt3engine.Message{
			messageFromUser,
			time.Now(),
		})
		gpt3engine.SaveJessMessage(dialogId, answer)
	}
}
