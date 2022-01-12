package gpt3engine

import (
	"context"
	"fmt"
	"io/ioutil"

	"cloud.google.com/go/storage"
)

const bucketName = "contexts"
const defaultContextObjectName = "preset.md"

func GetInitialContext(client *storage.Client, ctx context.Context) (string, error) {
	rc, err := client.
		Bucket(bucketName).
		Object(defaultContextObjectName).
		NewReader(ctx)
	if err != nil {
		return "", fmt.Errorf("Object(%q).NewReader: %v", defaultContextObjectName, err)
	}
	defer rc.Close()
	content, err := ioutil.ReadAll(rc)

	if err != nil {
		return "", err
	}

	return string(content), nil
}

func ConvertToContextString(message Message) string {
	return fmt.Sprintf("\n[%s][%s][%s]: %s\n[END]\n", message.Author, message.Time.Format(TimeFormatLayout), message.Mood, message.Text)
}
