package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/JessicaChatBot/gpt3-engine/gpt3engine"
	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "jess",
		Usage: "Jessica is your friend, this CLI is all you need to understand her.",
		Commands: []*cli.Command{{
			Name:    "dialog",
			Aliases: []string{"d"},
			Usage:   "start the dialog",
			Action:  startDialog,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "dialogId",
					Aliases:  []string{"i", "id"},
					Usage:    "dialog id, if not set random is used",
					Required: false,
				},
			},
		},
			{
				Name:    "shorten",
				Aliases: []string{"s"},
				Usage:   "shorten memory of the dialog",
				Action:  showCompressedDialog,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "dialogId",
						Aliases:  []string{"i", "id"},
						Usage:    "dialog id",
						Required: true,
					},
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func showCompressedDialog(c *cli.Context) error {
	dialogId := c.String("chatId")
	ctx := context.Background()
	fireStoreClient, err := gpt3engine.GetDefaultFirestoreClinet(ctx)
	if err != nil {
		fmt.Printf("get forestore client failed: %v\n", err)
		return err
	}
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Printf("get storage client failed: %v\n", err)
		return err
	}
	dialogContext, err := gpt3engine.GetInitialContext(storageClient, ctx)
	if err != nil {
		fmt.Printf("get context failed: %v\n", err)
		return err
	}
	dialogContext, err = gpt3engine.PopulateContextWithAllMessages(dialogId, dialogContext, fireStoreClient, ctx)
	if err != nil {
		fmt.Printf("population of the context failed: %v\n", err)
		return err
	}
	gpt3Client, err := gpt3engine.GetDefaultGpt3Client()
	if err != nil {
		fmt.Printf("error creating GPT3 client: %v\n", err)
		return err
	}
	save, err := gpt3engine.DescribeWhatJessShouldSaveFromTheDialog(dialogContext, gpt3Client, ctx)
	if err != nil {
		fmt.Printf("error describing what to save: %v\n", err)
		return err
	}
	fmt.Printf("Main dialog parts: %s\n", save)
	return nil
}

func startDialog(c *cli.Context) error {
	dialogId := c.String("chatId")
	if dialogId == "" {
		dialogId = uuid.New().String()
	}
	ctx := context.Background()
	fireStoreClient, err := gpt3engine.GetDefaultFirestoreClinet(ctx)
	if err != nil {
		fmt.Printf("get forestore client failed: %v\n", err)
		return err
	}
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Printf("get storage client failed: %v\n", err)
		return err
	}
	dialogContext, err := gpt3engine.GetInitialContext(storageClient, ctx)
	if err != nil {
		fmt.Printf("get context failed: %v\n", err)
		return err
	}
	dialogContext, err = gpt3engine.PopulateContextWithAllMessages(dialogId, dialogContext, fireStoreClient, ctx)
	if err != nil {
		fmt.Printf("population of the context failed: %v\n", err)
		return err
	}
	dialogContext, err = gpt3engine.PopulateContextWithLongTermMemory(dialogId, dialogContext, fireStoreClient, ctx)
	if err != nil {
		fmt.Printf("population of the context with memory failed: %v\n", err)
		return err
	}
	err = gpt3engine.CreateNewDialogIfAbsent(dialogId, fireStoreClient, ctx)
	if err != nil {
		fmt.Printf("error creating dialog id: %v\n", err)
		return err
	}
	gpt3Client, err := gpt3engine.GetDefaultGpt3Client()
	if err != nil {
		fmt.Printf("error creating GPT3 client: %v\n", err)
		return err
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("you: ")
		messageFromUserRaw, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("error getting message from user: %v\n", err)
			return fmt.Errorf("error getting message from user: %v", err)
		}
		if strings.Contains(messageFromUserRaw, "exit") {
			return nil
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
			return fmt.Errorf("error getting message from Jess: %v", err)
		}
		fmt.Println(fmt.Sprintf("Jess: %s\n", answer.Text))
		fmt.Println(fmt.Sprintf("Jess Mood: %s\n", answer.Mood))
		messageToAddToContext = gpt3engine.ConvertToContextString(answer)
		dialogContext = fmt.Sprintf("%s\n%s", dialogContext, messageToAddToContext)

		err = gpt3engine.SaveMessage(dialogId, messageFromUser, fireStoreClient, ctx)
		if err != nil {
			fmt.Printf("error saving message: %v\n", err)
			return fmt.Errorf("error saving message: %v", err)
		}
		err = gpt3engine.SaveMessage(dialogId, answer, fireStoreClient, ctx)
		if err != nil {
			fmt.Printf("error saving message: %v\n", err)
			return fmt.Errorf("error saving message: %v", err)
		}
	}
	return nil
}
