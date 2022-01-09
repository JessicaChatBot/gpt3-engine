package gpt3engine

import (
	"context"

	firebase "firebase.google.com/go"
)

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
