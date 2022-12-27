package config

import (
	"os"
	"testing"
)

func TestNotSet(t *testing.T) {
	os.Setenv("TELEGRAM_TOKEN", "")
	os.Setenv("OPENAI_TOKEN", "123")
	_, error := Read()
	if error == nil {
		t.Fatalf("expected error")
	}
}

func TestNotSet2(t *testing.T) {
	os.Setenv("TELEGRAM_TOKEN", "123")
	os.Setenv("OPENAI_TOKEN", "")
	_, error := Read()
	if error == nil {
		t.Fatalf("expected error")
	}
}

func TestId1(t *testing.T) {
	os.Setenv("TELEGRAM_TOKEN", "123")
	os.Setenv("OPENAI_TOKEN", "123")
	os.Setenv("WHITELISTED_IDS", "123")
	res, error := Read()
	if error != nil {
		t.Fatalf("expected error")
	}

	if len(res.WhitelistedUserId) == 0 || res.WhitelistedUserId[0] != 123 {
		t.Fatalf("not parsed")
	}
}

func TestId2(t *testing.T) {
	os.Setenv("TELEGRAM_TOKEN", "123")
	os.Setenv("OPENAI_TOKEN", "123")
	os.Setenv("WHITELISTED_IDS", "123,565,838")
	res, error := Read()
	if error != nil {
		t.Fatalf("expected error")
	}

	if len(res.WhitelistedUserId) == 0 ||
		res.WhitelistedUserId[0] != 123 ||
		res.WhitelistedUserId[1] != 565 ||
		res.WhitelistedUserId[2] != 838 {
		t.Fatalf("not parsed")
	}
}
