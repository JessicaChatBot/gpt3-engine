package gpt3engine

// func Say(message string, dialogId string) (string, error) {
// 	ctx := context.Background()
// 	fireStoreClient, err := GetDefaultFirestoreClinet(ctx)
// 	if err != nil {
// 		return "", err
// 	}
// 	storageClient, err := storage.NewClient(ctx)
// 	if err != nil {
// 		return "", err
// 	}
// 	dialogContext, err := getDialogContext(dialogId, fireStoreClient, ctx)
// 	if err != nil {
// 		return "", err
// 	}
// 	currentTime := time.Now()
// 	messageToAddToContext := fmt.Sprintf("\n[Friend][%s]: %s\n[END]\n", currentTime.Format("2006 Jan 2 15:04:05"), messageFromUser)
// 	gpt3Client, err := gpt3engine.GetDefaultGpt3Client()
// 	if err != nil {
// 		return "", err
// 	}
// 	gpt3engine.MessageToJess(dialogContext, gpt3Client, ctx)
// }

// func getDialogContext(dialogId string, fireStoreClient *firestore.Client, ctx context.Context) (string, error) {
// 	storageClient, err := storage.NewClient(ctx)
// 	if err != nil {
// 		return "", err
// 	}
// 	dialogContext, err := gpt3engine.GetInitialContext(storageClient, ctx)
// 	if err != nil {
// 		return "", err
// 	}
// 	fullContext, err := gpt3engine.PopulateContextWithAllMessages(dialogId, dialogContext, fireStoreClient, ctx)
// 	if err != nil {
// 		return "", err
// 	}
// 	return fullContext, nil
// }
