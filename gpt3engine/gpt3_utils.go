package gpt3engine

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

const jessStartToken = "[Jess]"
const endToken = "[END]"

func deriveToken(rawMessage string, tokens []string, currentIndex int) (string, error) {
	if currentIndex >= len(tokens) {
		return "", errors.New(fmt.Sprintf("no token found: %s", rawMessage))
	}
	if strings.Contains(rawMessage, tokens[currentIndex]) {
		return tokens[currentIndex], nil
	}
	return deriveToken(rawMessage, tokens, currentIndex+1)
}

func desperateParse(rawMessage string) (string, error) {
	derivedEndToken, err := deriveToken(rawMessage, []string{endToken, "END"}, 0)
	if err != nil {
		return "", err
	}
	derivedStartToken, err := deriveToken(rawMessage, []string{jessStartToken, "Jess"}, 0)
	if err != nil {
		return "", err
	}
	message := strings.Split(rawMessage, derivedEndToken)[0]
	message = strings.Split(message, derivedStartToken)[1]
	for {
		if !strings.Contains(message, "]") {
			return message, nil
		}
		message = strings.Split(message, "]")[1]
	}
}

func ParseJessResponse(rawMessage string) (Message, error) {
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
	r := regexp.MustCompile(`\[Jess\]\[(?P<time>[^\]]*)\]\[(?P<mood>[^\]]*)\]:((?P<msg>.|\n)*?\[END\])`)
	parsingResult := r.FindStringSubmatch(rawMessageWithoutNewLines)
	if len(parsingResult) == 0 {
		desperateParse, err := desperateParse(rawMessage)
		if err != nil {
			return Message{}, errors.New(fmt.Sprintf("was not able to parse answer from the server: %s", rawMessage))
		}
		return Message{
			Text:   desperateParse,
			Time:   time.Now(),
			Author: "Jess",
			Mood:   []string{"unknown"},
			Raw:    rawMessage,
		}, nil
	}
	rawDateString := parsingResult[1]
	date, err := time.Parse(TimeFormatLayout, rawDateString)
	if err != nil {
		return Message{}, errors.New(fmt.Sprintf("failed to parse date from server: %s\nerror: %v\n", rawDateString, err))
	}
	messageStringFromJess := strings.Replace(parsingResult[3], "[END]", "", -1)
	return Message{
		Text:   messageStringFromJess,
		Time:   date,
		Author: "Jess",
		Mood:   strings.Split(parsingResult[2], "."),
		Raw:    parsingResult[0],
	}, nil
}
