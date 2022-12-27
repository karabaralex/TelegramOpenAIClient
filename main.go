package main

import (
	"fmt"
	"runtime/debug"

	"github.com/telegram-command-reader/bot"
	"github.com/telegram-command-reader/config"
	"github.com/telegram-command-reader/operations"
)

// safely call function without panic
func safeCall(f func(), onError func(string)) {
	defer func() {
		if r := recover(); r != nil {
			// print recovered stack trace
			fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))
			fmt.Println("Recovered in f", r)
			stackTraceUrl, sterr := operations.SendStringToPastebin(string(debug.Stack()))
			if sterr != nil {
				fmt.Println("Send stacktrace error ", sterr)
			} else {
				result := fmt.Sprintf("Error: %s, Stacktrace url: %s", r, stackTraceUrl)
				onError(result)
			}
		}
	}()
	f()
}

func main() {
	version := "Telegram OpenAI client"
	fmt.Println(version)
	envConfig, envError := config.Read()
	if envError != nil {
		fmt.Println("Load config error ", envError)
		return
	}

	bot.API_TOKEN = envConfig.TelegramBotToken

	outputChannel := make(chan bot.OutMessage)

	bot.AddHandler(bot.NewCommandMatcher("/version"), func(message *bot.Info) {
		outputChannel <- bot.OutMessage{OriginalMessage: message, Text: version}
	})

	bot.AddHandler(bot.NewTextMatcher(".*"), func(message *bot.Info) {
		fmt.Println("Command .*", message.Text)
		fmt.Println("Sender .*", message.FromID)
		bot.SendTypingStatus(message)
		if findIntInArray(envConfig.WhitelistedUserId, message.FromID) == false {
			reply := bot.OutMessage{OriginalMessage: message, Text: "Not whitelisted", Html: false}
			outputChannel <- reply
			return
		}

		go safeCall(func() {
			res, err := operations.AskOpenAI(message.Text, envConfig.OpenAIToken)
			if err != nil {
				fmt.Println(res)
				reply := bot.OutMessage{OriginalMessage: message, Text: err.Error()}
				outputChannel <- reply
			} else {
				reply := bot.OutMessage{OriginalMessage: message, Text: res, Html: true}
				outputChannel <- reply
			}
		}, func(s string) {
			reply := bot.OutMessage{OriginalMessage: message, Text: s}
			outputChannel <- reply
		})
	})

	go bot.Sender(outputChannel)
	bot.RequestUpdates()
}

// This function takes an array of integers and a value to find as arguments,
// and returns a boolean indicating whether the value was found in the array.
func findIntInArray(array []int64, valueToFind int64) bool {
	// Set a flag to track whether we have found the value
	found := false

	// Iterate over the array
	for _, v := range array {
		// Check if the current element is the value we are looking for
		if v == valueToFind {
			found = true
			break
		}
	}

	return found
}
