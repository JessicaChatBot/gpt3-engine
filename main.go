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
	if err != nil {
		fmt.Printf("get context failed: %v\n", err)
		return
	}
	messageFromUser := ""
	reader := bufio.NewReader(os.Stdin)
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
		// fmt.Printf(context)
		answer, err := gpt3engine.MessageToJess(context)
		if err != nil {
			fmt.Printf("error getting message from Jess: %v\n", err)
			return
		}
		fmt.Println(fmt.Sprintf("Jess: %s\n", answer.Message.Text))
		fmt.Println(fmt.Sprintf("Jess Mood: %s\n", answer.Mood))
		context = context + fmt.Sprintf("\n%s\n", answer.Raw)
	}
	// mssg, err := gpt3engine.Request()
	// if err != nil {
	// 	fmt.Printf("failed to get request: %v\n", err)
	// 	return
	// }
	// fmt.Printf("Jess: %s\n", mssg.Message.Text)
	// err := gpt3engine.CreateNewDialog("test-1")
	// if err != nil {
	// 	fmt.Printf("error: %v\n", err)
	// }
}
