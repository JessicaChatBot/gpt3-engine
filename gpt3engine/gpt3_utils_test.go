package gpt3engine

import "testing"

func TestParse(t *testing.T) {
	testDesperateParse(t, "Haha\n", "[Jess][05:06:56][friendly, curious]: Haha\n[END]", parseTestingWrapper)
	testDesperateParse(t, "Hi, are you there?", "[Jess][05:06:56][friendly, curious]: Hi, are you there?[END]", parseTestingWrapper)
}

func parseTestingWrapper(rawMessage string) (string, error) {
	resp, err := parseJessResponse(rawMessage)
	return resp.Text, err
}

func TestDesperateParse(t *testing.T) {
	testDesperateParse(t, "  unknown..", "[Jess][2022 Jan 14 22:55:37][unknown]: [[unknown..[END]", desperateParse)
	testDesperateParse(t, "  unknown  unknown  friends & partners...", "[Jess][2022 Jan 14 22:56:26][unknown]: [[unknown[[unknown[[friends & partners...[END]", desperateParse)
	testDesperateParse(t, " Wow, you are fast...", "[Jess][2022 Jan 14 22:54:37][friendly]:  Wow, you are fast...[END]", desperateParse)
	testDesperateParse(t, " |   ....", "[Jess][2022 Jan 14 22:39:56][friendly]:  | [[....[END]", desperateParse)
}

func testDesperateParse(t *testing.T, expected string, rawMessage string, funcToTest func(string) (string, error)) {
	actual, err := desperateParse(rawMessage)
	if err != nil {
		t.Errorf("error parsing message: %v", err)
		return
	}
	if actual != expected {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}
}
