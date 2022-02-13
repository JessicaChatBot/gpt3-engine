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
const tableWithMemories = "memory"

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

func PopulateContextWithLongTermMemory(dialogId string, dialogContext string, client *firestore.Client, ctx context.Context) (string, error) {
	memory, err := GetLongTermMemory(dialogId, client, ctx)
	if err != nil {
		return "", err
	}
	if memory != "" {
		return fmt.Sprintf("%s\n# Jess Memory\n\n%s", dialogContext, memory), nil
	}
	return dialogContext, nil
}

func PopulateContextWithAllMessages(dialogId string, dialogContext string, client *firestore.Client, ctx context.Context) (string, error) {
	iter := client.Collection(tableWithDialogsHistory).
		Where(dialogIdColKey, "==", dialogId).
		OrderBy(timeColKey, firestore.Asc).
		Documents(ctx)
	currentContext := fmt.Sprintf("%s\n# Chatlog\n\n", dialogContext)
	messagesForContext := ""
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
		messagesForContext = fmt.Sprintf("%s\n%s",
			message.ConvertToString(),
			messagesForContext)
	}
	return fmt.Sprintf("%s%s", currentContext, messagesForContext), nil
}

func SaveLongTermMemory(memoryToSave string, dialogId string, client *firestore.Client, ctx context.Context) error {
	_, err := client.Collection(tableWithMemories).
		Doc(dialogId).
		Set(ctx, map[string]interface{}{
			"memoryGeneral": memoryToSave,
		})
	if err != nil {
		return err
	}
	return nil
}

func GetLongTermMemory(dialogId string, client *firestore.Client, ctx context.Context) (string, error) {
	doc, err := client.Collection(tableWithMemories).
		Doc(dialogId).
		Get(ctx)
	if !doc.Exists() {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	data := doc.Data()
	if data == nil {
		return "", nil
	}
	return data["memoryGeneral"].(string), nil
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
