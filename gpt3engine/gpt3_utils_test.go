package gpt3engine

import "testing"

func TestDesperateParse(t *testing.T) {
	testDesperateParse(t, "  unknown..", "[Jess][2022 Jan 14 22:55:37][unknown]: [[unknown..[END]")
	testDesperateParse(t, "  unknown  unknown  friends & partners...", "[Jess][2022 Jan 14 22:56:26][unknown]: [[unknown[[unknown[[friends & partners...[END]")
	testDesperateParse(t, " Wow, you are fast...", "[Jess][2022 Jan 14 22:54:37][friendly]:  Wow, you are fast...[END]")
	testDesperateParse(t, " |   ....", "[Jess][2022 Jan 14 22:39:56][friendly]:  | [[....[END]")
}

func testDesperateParse(t *testing.T, expected string, rawMessage string) {
	actual, err := desperateParse(rawMessage)
	if err != nil {
		t.Errorf("error parsing message: %v", err)
		return
	}
	if actual != expected {
		t.Errorf("expected: %s, actual: %s", expected, actual)
	}
}
