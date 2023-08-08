package main

import (
	"context"
	openai "github.com/sashabaranov/go-openai"
	"log"
	"regexp"
	"strings"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9#]+`)

func sanitizeString(str *string) string {
	return nonAlphanumericRegex.ReplaceAllString(*str, "")
}

func aiPaletteCall(apiKey, prompt string) []string {
	// Sending the prompt to OpenAI
	openaiClient := openai.NewClient(apiKey)
	bgCxt := context.Background()

	prompt = "Give me a palette, limited to 6 colors max, in the form of a comma deliminated list, meaning that there should be a comma and space, that also includes the pound symbol, and the vibe or essence of the palette should be " + prompt + ". Under no circumstances are you to return anything that isn't hex codes colors, pound symbols, or other formating characters that were explicitly asked for."

	compReq := openai.CompletionRequest{
		Model:     openai.GPT3TextDavinci003,
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
