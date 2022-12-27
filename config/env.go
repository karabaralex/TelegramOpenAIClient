package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	TelegramBotToken  string
	OpenAIToken       string
	WhitelistedUserId []int64
}

func Read() (Config, error) {
	result := Config{}
	result.TelegramBotToken = os.Getenv("TELEGRAM_TOKEN")
	if result.TelegramBotToken == "" {
		return result, errors.New("no telegram token")
	}

	result.OpenAIToken = os.Getenv("OPENAI_TOKEN")
	if result.OpenAIToken == "" {
		return result, errors.New("no OpenAI token")
	}

	idsString := os.Getenv("WHITELISTED_IDS")
	if idsString != "" {
		idsStrArray := strings.Split(idsString, ",")
		idsIntArray, err := convertStringsToInts(idsStrArray)
		if err != nil {
			return result, err
		}

		result.WhitelistedUserId = idsIntArray
	}

	return result, nil
}

func convertStringsToInts(input []string) ([]int64, error) {
	var t2 = []int64{}

	for _, i := range input {
		j, err := strconv.ParseInt(i, 10, 0)
		if err != nil {
			panic(err)
		}
		t2 = append(t2, j)
	}
	return t2, nil
}
