package main

import (
	"context"
	"log"
	"regexp"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9#]+`)

func sanitizeString(str *string) string {
	return nonAlphanumericRegex.ReplaceAllString(*str, "")
}

func aiPaletteCall(apiKey, prompt string) []string {
	// Sending the prompt to OpenAI
	openaiClient := openai.NewClient(apiKey)
	bgCxt := context.Background()

	prompt = "Give me a palette in the form of a comma deliminated list of hex colors, limited to a maximum of 6 colors, that also includes the pound symbol for each hex code, and the vibe or essence of the palette should be " + prompt + ". Do not return anything that isn't the hex codes."

	compReq := openai.CompletionRequest{
		Model:     openai.GPT3Dot5TurboInstruct,
		MaxTokens: 64,
		Prompt:    prompt,
	}

	responce, err := openaiClient.CreateCompletion(bgCxt, compReq)
	if err != nil {
		log.Println("Completion error: ", err)
	}

	log.Println(strings.Split(responce.Choices[0].Text, ", "))

	// Saving the responce in preferred format
	responcePaletteSlice := strings.Split(responce.Choices[0].Text, ", ")

	// Sanitizing responce
	responcePaletteSlice[0] = responcePaletteSlice[0][3:]
	responcePaletteSlice[0] = "#" + responcePaletteSlice[0]

	for _, valStr := range responcePaletteSlice {
		sanitizeString(&valStr)
		strings.ReplaceAll(valStr, "\n", "")
		log.Printf("Color: -%v-", valStr)
	}

	return responcePaletteSlice
}
