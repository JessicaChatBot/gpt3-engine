package functions

import (
	"encoding/json"
	"fmt"
	"log"

	"context"
	"time"

	"net/http"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"github.com/JessicaChatBot/gpt3-engine/gpt3engine"
	gogpt "github.com/sashabaranov/go-gpt3"
)

var fireStoreClient *firestore.Client
var storageClient *storage.Client
var gpt3Client *gogpt.Client

func init() {
	ctx := context.Background()
	var err error
	fireStoreClient, err = gpt3engine.GetDefaultFirestoreClinet(ctx)
	if err != nil {
		log.Fatalf("get forestore client failed: %v\n", err)
	}
	storageClient, err = storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("get storage client failed: %v\n", err)
	}
	gpt3Client, err = gpt3engine.GetDefaultGpt3Client()
	if err != nil {
		log.Fatalf("error creating GPT3 client: %v\n", err)
	}
}

type message struct {
	DialogId string `json:"dialogId"`
	Text     string `json:"text"`
}

func Say(w http.ResponseWriter, r *http.Request) {
	m := message{}
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		log.Printf("json.NewDecoder: %v", err)
		http.Error(w, "Error parsing request", http.StatusBadRequest)
		return
	}
	dialogContext, err := gpt3engine.GetInitialContext(storageClient, ctx)
	if err != nil {
		log.Printf("get context failed: %v\n", err)
		http.Error(w, "Error getting the context", http.StatusBadRequest)
		return
	}
	dialogContext, err = gpt3engine.PopulateContextWithAllMessages(m.DialogId, dialogContext, fireStoreClient, ctx)
	if err != nil {
		log.Printf("population of the context failed: %v\n", err)
		http.Error(w, "Error populating the context", http.StatusBadRequest)
		return
	}
	err = gpt3engine.CreateNewDialogIfAbsent(m.DialogId, fireStoreClient, ctx)
	if err != nil {
		log.Printf("error creating dialog id: %v\n", err)
		http.Error(w, "Error createing new dialogId", http.StatusBadRequest)
		return
	}
	messageFromUser := gpt3engine.Message{
		Text:   m.Text,
		Time:   time.Now(),
		Author: "Friend",
		Mood:   []string{"unknown"},
		Raw:    m.Text,
	}
	messageToAddToContext := gpt3engine.ConvertToContextString(messageFromUser)
	dialogContext = fmt.Sprintf("%s\n%s", dialogContext, messageToAddToContext)
	answer, err := gpt3engine.MessageToJess(dialogContext, gpt3Client, ctx)
	if err != nil {
		log.Printf("error getting message from Jess: %v\n", err)
		http.Error(w, "Error getting message from Jess", http.StatusBadRequest)
		return
	}
	err = gpt3engine.SaveMessage(m.DialogId, messageFromUser, fireStoreClient, ctx)
	if err != nil {
		log.Printf("error saving message: %v\n", err)
		http.Error(w, "Error saving user message", http.StatusBadRequest)
		return
	}
	err = gpt3engine.SaveMessage(m.DialogId, answer, fireStoreClient, ctx)
	if err != nil {
		log.Printf("error saving message: %v\n", err)
		http.Error(w, "Error saving message from Jess", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, answer.Text)
}
