package gpt3engine

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"time"

	gogpt "github.com/sashabaranov/go-gpt3"
)

type Message struct {
	Text string
	Time time.Time
}

type JessMessage struct {
	Message Message
	Mood    []string
	Raw     string
}

func MessageToJess(contextWithMessage string) (JessMessage, error) {
	jessTokenApiSecretKeyName := "projects/438230051748/secrets/gpt3-api-secret-key/versions/latest"
	secret, err := AccessSecretVersion(jessTokenApiSecretKeyName)
	if err != nil {
		return JessMessage{}, err
	}
	if err != nil {
		return JessMessage{}, err
	}
	c := gogpt.NewClient(secret)
	ctx := context.Background()
	req := gogpt.CompletionRequest{
		MaxTokens: 80,
		Prompt:    contextWithMessage,
	}
	resp, err := c.CreateCompletion(ctx, "davinci", req)
	if err != nil {
		return JessMessage{}, err
	}
	jessResp, err := getJessResponse(resp.Choices[0].Text)
	if err != nil {
		return JessMessage{}, err
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

func desperateParse(rawMessage string) (string, error) {
	if !strings.Contains(rawMessage, "[END]") {
		return "", errors.New(fmt.Sprintf("message does not have [END]: %s", rawMessage))
	}
	if !strings.Contains(rawMessage, "[Jess]") {
		return "", errors.New(fmt.Sprintf("message does not have [Jess]: %s", rawMessage))
	}
	message := strings.Split(rawMessage, "[END]")[0]
	message = strings.Split(message, "[Jess]")[1]
	for {
		if !strings.Contains(message, "]") {
			return message, nil
		}
		message = strings.Split(message, "]")[1]
	}
	return "", errors.New(fmt.Sprintf("we should not be here. Raw: %s\nMessage: %s", rawMessage, message))
}

func getJessResponse(rawMessage string) (JessMessage, error) {
	// Generated with: https://regex101.com/r/QelR3A/1
	// Tested with:
	/*
		resp:

		[Jess][05:06:56][friendly, curious]: Haha I am pretty sure I am not that bad at telling jokes :) Although I want to say that my memory is little blurred. I do not remember how old I am.
		[END]

		[Friend][05:10:17]: okay it looks to me you are currently living in China? No disrespet but let me tell you, Chinese is a very difficult language to learn
		[END]
	*/
	rawMessageWithoutNewLines := strings.Replace(rawMessage, "\n", ".", -1)
	// r := regexp.MustCompile(`\[Jess\]\[(?P<time>.*)\]\[(?P<mood>.*)\](?P<msg>(.|\n)*?)\[END\]`)
	r := regexp.MustCompile(`\[Jess\]\[(?P<time>[^\]]*)\]\[(?P<mood>[^\]]*)\]:((?P<msg>.|\n)*?\[END\])`)
	parsingResult := r.FindStringSubmatch(rawMessageWithoutNewLines)
	if len(parsingResult) == 0 {
		desperateParse, err := desperateParse(rawMessage)
		if err != nil {
			return JessMessage{}, errors.New(fmt.Sprintf("was not able to parse answer from the server: %s", rawMessage))
		}
		newMessage := Message{
			desperateParse,
			time.Now(),
		}
		return JessMessage{
			newMessage,
			[]string{},
			rawMessage,
		}, nil
	}
	dateLayout := "2006 Jan 2 15:04:05"
	rawDateString := parsingResult[1]
	date, err := time.Parse(dateLayout, rawDateString)
	if err != nil {
		return JessMessage{}, errors.New(fmt.Sprintf("failed to parse date from server: %s\nerror: %v\n", rawDateString, err))
	}
	messageStringFromJess := strings.Replace(parsingResult[3], "[END]", "", -1)
	newMessage := Message{
		messageStringFromJess,
		date,
	}
	return JessMessage{
		newMessage,
		strings.Split(parsingResult[2], "."),
		parsingResult[0],
	}, nil
}
