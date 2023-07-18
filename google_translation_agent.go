package main

import (
	"cloud.google.com/go/translate"
	"context"
	"fmt"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

type GoogleTranslationAgent struct {
	GoogleApiKey string
	Text         string
}

func NewGoogleTranslationAgent(spec TaskSpec) *GoogleTranslationAgent {
	agent := GoogleTranslationAgent{}
	parser := GetInstance()
	agent.Text = spec.Args["text"].(string)
	agent.GoogleApiKey = parser.GoogleAPIKey
	return &agent
}

func (e *GoogleTranslationAgent) Execute(ctx context.Context, task *Task) error {
	client, err := translate.NewClient(ctx, option.WithAPIKey(e.GoogleApiKey))
	if err != nil {
		return err
	}
	defer client.Close()

	source, target := language.English, language.French
	//if val, ok := task.Args["source"]; ok {
	//	source = val.(language.Tag)
	//}
	//if val, ok := task.Args["target"]; ok {
	//	target = val.(language.Tag)
	//}

	resp, err := client.Translate(ctx, []string{e.Text}, target,
		&translate.Options{Format: translate.Text, Source: source})
	if err != nil {
		if err == context.Canceled {
			fmt.Println("Google Translation request was cancelled")
		} else {
			fmt.Println("Google Translation request failed with error:", err)
		}
		return err
	}

	fmt.Println("Translation from Google: ", resp[0].Text)

	return nil
}

// ... rest of the code (Task struct, etc) ...
