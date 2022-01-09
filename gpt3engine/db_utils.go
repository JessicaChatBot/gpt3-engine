package gpt3engine

import (
	"context"
	"math/rand"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
)

func PopulateContextWithAllMessages(dialogId string, dialogContext string) (string, error) {
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: "jessdb-337700"}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		return "", err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return "", err
	}
	iter := client.Collection("history").Where("dialogId", "==", dialogId).OrderBy("time", firestore.Desc).Documents(ctx)
	currentContext := dialogContext
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return "", err
		}

		data := doc.Data()
		dateLayout := "2006 Jan 2 15:04:05"
		date, err := time.Parse(dateLayout, data["time"].(string))
		if err != nil {
			continue
		}
		message := Message{
			string(data["text"].(string)),
			date,
		}
		if data["author"].(string) == "user" {
			currentContext = currentContext + "\n" + ConvertToString(message) + "\n"
		} else {
			jessMessage := JessMessage{
				message,
				data["mood"].([]string),
				data["raw"].(string),
			}
			currentContext = currentContext + "\n" + ConvertJessMessageToString(jessMessage) + "\n"
		}
	}
	return currentContext, nil
}

func SaveMessage(dialogId string, message Message) error {
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: "jessdb-337700"}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		return err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return err
	}
	_, err = client.Collection("history").Doc(randId(20)).Set(ctx, map[string]interface{}{
		"dialogId": dialogId,
		"text":     message.Text,
		"author":   "user",
		"time":     message.Time,
	})
	if err != nil {
		return err
	}
	return nil
}

func SaveJessMessage(dialogId string, message JessMessage) error {
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: "jessdb-337700"}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		return err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return err
	}
	_, err = client.Collection("history").Doc(randId(20)).Set(ctx, map[string]interface{}{
		"dialogId": dialogId,
		"text":     message.Message.Text,
		"author":   "Jess",
		"time":     message.Message.Time,
		"mood":     message.Mood,
		"raw":      message.Raw,
	})
	if err != nil {
		return err
	}
	return nil
}

func CreateNewDialogIfAbsent(id string) error {
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: "jessdb-337700"}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		return err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return err
	}
	_, err = client.Collection("dialogs").Doc(id).Set(ctx, map[string]interface{}{})
	if err != nil {
		return err
	}
	return nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randId(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
