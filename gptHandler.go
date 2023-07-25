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

	prompt = "Give me a pallete, limited to 6 colors max, in the form of a comma deliminated list, meaning that there should be a comma and space, that also includes the pound symbol, and the vibe or essence of the pallete should be " + prompt

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
	responcePalleteSlice := strings.Split(responce.Choices[0].Text, ", ")

	// Sanitizing responce
	responcePalleteSlice[0] = responcePalleteSlice[0][3:]
	responcePalleteSlice[0] = "#" + responcePalleteSlice[0]

	for _, valStr := range responcePalleteSlice {
		sanitizeString(&valStr)
		strings.ReplaceAll(valStr, "\n", "")
		log.Printf("Color: -%v-", valStr)
	}

	return responcePalleteSlice
}
