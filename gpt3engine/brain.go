package gpt3engine

import (
	"context"
	"io/ioutil"

	gogpt "github.com/sashabaranov/go-gpt3"
)

func MessageToJess(contextWithMessage string) (Message, error) {
	jessTokenApiSecretKeyName := "projects/438230051748/secrets/gpt3-api-secret-key/versions/latest"
	secret, err := AccessSecretVersion(jessTokenApiSecretKeyName)
	if err != nil {
		return Message{}, err
	}
	if err != nil {
		return Message{}, err
	}
	c := gogpt.NewClient(secret)
	ctx := context.Background()
	req := gogpt.CompletionRequest{
		MaxTokens: 80,
		Prompt:    contextWithMessage,
		Stop:      []string{"[END]"},
	}
	resp, err := c.CreateCompletion(ctx, "davinci", req)
	if err != nil {
		return Message{}, err
	}
	jessResp, err := ParseJessResponse(resp.Choices[0].Text + "[END]")
	if err != nil {
		return Message{}, err
	}
	return jessResp, nil
}

func GetInitialContext() (string, error) {
	content, err := ioutil.ReadFile("gpt3engine/preset.md")

	if err != nil {
		return "", err
	}

	return string(content), nil
}
