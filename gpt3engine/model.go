package gpt3engine

import (
	"fmt"
	"time"
)

type Message struct {
	Text   string
	Time   time.Time
	Author string
	Mood   []string
	Raw    string
}

const TimeFormatLayout string = "2006 Jan 2 15:04:05"
const UnknownMoodMarker string = "unknown"

func (msg Message) ConvertToString() string {
	return fmt.Sprintf("[%s][%s][%s]: %s\n[END]",
		msg.Author,
		msg.Time.Format(TimeFormatLayout),
		msg.convertMood(),
		msg.Text)
}

func (msg Message) convertMood() string {
	moodRaw := ""
	for _, v := range msg.Mood {
		if len(moodRaw) > 0 {
			moodRaw = moodRaw + " "
		}
		moodRaw = moodRaw + v
	}
	return moodRaw
}
