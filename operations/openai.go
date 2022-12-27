package operations

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

/**
{
  "model": "text-davinci-003",
  "prompt": "The following is a conversation with an AI assistant. The assistant is helpful, creative, clever, and very friendly.\n\nHuman: Hello, who are you?\nAI: I am an AI created by OpenAI. How can I help you today?\nHuman: I'd like to cancel my subscription.\nAI:",
  "temperature": 0.9,
  "max_tokens": 150,
  "top_p": 1,
  "frequency_penalty": 0.0,
  "presence_penalty": 0.6,
  "stop": [" Human:", " AI:"]
**/

type openAiRequest struct {
	Model            string   `json:"model"`
	Prompt           string   `json:"prompt"`
	Temperature      float32  `json:"temperature"`
	MaxTokens        int      `json:"max_tokens"`
	TopP             float32  `json:"top_p"`
	FrequencyPenalty float32  `json:"frequency_penalty"`
	PresencePenalty  float32  `json:"presence_penalty"`
	Stop             []string `json:"stop"`
}

type OpenAiResponse struct {
	ID      string                 `json:"id"`
	Prompt  string                 `json:"prompt"`
	Choices []OpenAiResponseChoice `json:"choices"`
}

type OpenAiResponseChoice struct {
	Text         string `json:"text"`
	FinishReason string `json:"finish_reason"`
}

/**
send request to open ai and get response

response example

{"id":"cmpl-6RtTsQp3lMdB76r3JJxksnFlVEkEg","object":"text_completion","created":1672106144,"model":"text-davinci-003","choices":
[{"text":"\n\nIngredients:\n\n-1 1/2 cups all-purpose flour\n-1/4 teaspoon salt\n-2/3 cup cold water\n-2 tablespoons vegetable oil\n-1/2 pound ground pork or beef\n-1/4 cup minced onion\n-1 teaspoon sesame oil\n-2 cloves garlic, minced\n-1 tablespoon freshly grated ginger\n-1 tablespoon soy sauce\n-1/3 cup finely chopped cabbage\n-1/4 cup shredded carrots\n-1/4 cup chopped green onions\n-1/4 cup frozen corn kernels\n-Salt and pepper to taste\n\nInstructions:\n\n1. In a large bowl, whisk together the flour and salt. Stir in the cold","index":0,"logprobs":null,"finish_reason":"length"}],"usage":{"prompt_tokens":7,"completion_tokens":150,"total_tokens":157}}

*/
func AskOpenAI(body string, token string) (string, error) {
	url := "https://api.openai.com/v1/completions"

	stop := [2]string{" Human:", " AI:"}
	foo_marshalled, err := json.Marshal(openAiRequest{Model: "text-davinci-003",
		Prompt: body, Temperature: 0.9, MaxTokens: 2000, TopP: 1, FrequencyPenalty: 0.0, PresencePenalty: 0.6, Stop: stop[:]})

	// Create the request
	req, err := http.NewRequest("POST", url,
		strings.NewReader(string(foo_marshalled)))
	if err != nil {
		return "", err
	}

	// Set the headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Send the request via a client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	jsonResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Close the response body
	defer resp.Body.Close()

	openAiResponse, err := ParseOpenAIResponse(string(jsonResponse))
	if err != nil {
		return "", err
	}

	if len(openAiResponse.Choices) == 0 {
		return "", fmt.Errorf("empty choices in response")
	}

	// convertedString := string(pastedUrl)
	// // trim and remove new lines
	// convertedString = strings.Trim(convertedString, "\n")
	return openAiResponse.Choices[0].Text, nil
}

func ParseOpenAIResponse(response string) (*OpenAiResponse, error) {
	// Create a struct to hold the parsed data
	var p OpenAiResponse

	// Use json.Unmarshal to parse the JSON data
	err := json.Unmarshal([]byte(response), &p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}
