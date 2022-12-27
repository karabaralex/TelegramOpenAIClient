package operations

import (
	"testing"
)

func TestResponseParsed(t *testing.T) {
	response := `{"id":"cmpl-6RtTsQp3lMdB76r3JJxksnFlVEkEg","object":"text_completion","created":1672106144,"model":"text-davinci-003","choices":
	[{"text":"\n\nIngredients:\n\n-1 1/2 cups all-purpose flour\n-1/4 teaspoon salt\n-2/3 cup cold water\n-2 tablespoons vegetable oil\n-1/2 pound ground pork or beef\n-1/4 cup minced onion\n-1 teaspoon sesame oil\n-2 cloves garlic, minced\n-1 tablespoon freshly grated ginger\n-1 tablespoon soy sauce\n-1/3 cup finely chopped cabbage\n-1/4 cup shredded carrots\n-1/4 cup chopped green onions\n-1/4 cup frozen corn kernels\n-Salt and pepper to taste\n\nInstructions:\n\n1. In a large bowl, whisk together the flour and salt. Stir in the cold","index":0,"logprobs":null,"finish_reason":"length"}],"usage":{"prompt_tokens":7,"completion_tokens":150,"total_tokens":157}}`

	actual, err := ParseOpenAIResponse(response)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if actual.Choices[0].Text == "" {
		t.Fatalf("expected not empty, got %v", actual.Choices[0].Text)
	}

	t.Log(actual.Choices[0].Text)
}
