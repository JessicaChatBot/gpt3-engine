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

const googleProjectId = "jessdb-337700"

const tableWithDialogsIds = "dialogs"
const tableWithDialogsHistory = "history"

const dialogIdColKey = "dialogId"
const msgTextColKey = "text"
const moodColKey = "mood"
const timeColKey = "time"
const rawMessageColKey = "raw"
const authorColKey = "author"

func GetDefaultFirestoreClinet(ctx context.Context) (*firestore.Client, error) {
	conf := &firebase.Config{ProjectID: googleProjectId}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		return nil, err
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func PopulateContextWithAllMessages(dialogId string, dialogContext string, client *firestore.Client, ctx context.Context) (string, error) {
	iter := client.Collection(tableWithDialogsHistory).
		Where(dialogIdColKey, "==", dialogId).
		OrderBy(timeColKey, firestore.Desc).
		Documents(ctx)
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
		date := data[timeColKey].(time.Time)
		mood := []string{UnknownMoodMarker}
		if data[moodColKey] != nil {
			mood = make(
				[]string,
				len(data[moodColKey].([]interface{})))
			for i, v := range data[moodColKey].([]interface{}) {
				mood[i] = fmt.Sprint(v)
			}
		}
		raw := ""
		if data[rawMessageColKey] != nil {
			raw = data[rawMessageColKey].(string)
		}
		message := Message{
			Text:   string(data[msgTextColKey].(string)),
			Time:   date,
			Author: data[authorColKey].(string),
			Mood:   mood,
			Raw:    raw,
		}
		currentContext = fmt.Sprintf("%s\n%s",
			currentContext,
			message.ConvertToString())
	}
	return currentContext, nil
}

func SaveMessage(dialogId string, message Message, client *firestore.Client, ctx context.Context) error {
	_, err := client.Collection(tableWithDialogsHistory).
		Doc(randId(20)).
		Set(ctx, map[string]interface{}{
			dialogIdColKey:   dialogId,
			msgTextColKey:    message.Text,
			authorColKey:     message.Author,
			timeColKey:       message.Time,
			moodColKey:       message.Mood,
			rawMessageColKey: message.Raw,
		})
	if err != nil {
		return err
	}
	return nil
}

func CreateNewDialogIfAbsent(id string, client *firestore.Client, ctx context.Context) error {
	_, err := client.Collection(tableWithDialogsIds).
		Doc(id).
		Set(ctx, map[string]interface{}{})
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
