package gpt3engine

import (
	"context"
	"fmt"
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
		date := data["time"].(time.Time)
		if err != nil {
			continue
		}
		mood := []string{"unknown"}
		if data["mood"] != nil {
			mood = make([]string, len(data["mood"].([]interface{})))
			for i, v := range data["mood"].([]interface{}) {
				mood[i] = fmt.Sprint(v)
			}
		}
		raw := ""
		if data["raw"] != nil {
			raw = data["raw"].(string)
		}
		message := Message{
			Text:   string(data["text"].(string)),
			Time:   date,
			Author: data["author"].(string),
			Mood:   mood,
			Raw:    raw,
		}
		currentContext = currentContext + "\n" + message.ConvertToString()
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
		"author":   message.Author,
		"time":     message.Time,
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
