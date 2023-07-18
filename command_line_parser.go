package main

import (
	"flag"
	"log"
	"os"
	"sync"
)

type CommandLineParser struct {
	OpenAPIKey   string
	InputFile    string
	OutputFile   string
	Model        string
	GoogleAPIKey string
}

var instance *CommandLineParser
var once sync.Once

func GetInstance() *CommandLineParser {
	once.Do(func() {
		instance = createInstance()
	})
	return instance
}

func createInstance() *CommandLineParser {
	openAiApiKey := os.Getenv("OPENAI_API_KEY")
	if openAiApiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is not set.")
	}

	googleApiKey := os.Getenv("GOOGLE_API_KEY")
	if googleApiKey == "" {
		log.Fatal("GOOGLE_API_KEY environment variable is not set.")
	}

	inputFileFlag := flag.String("input", "input.txt", "Specify the input file")
	outputFileFlag := flag.String("output", "output.md", "Specify the output file")
	modelFlag := flag.String("model", "gpt3dot5", "Specify the model: gpt3dot5 or gpt4")
	flag.Parse()

	return &CommandLineParser{
		OpenAPIKey:   openAiApiKey,
		GoogleAPIKey: googleApiKey,
		InputFile:    *inputFileFlag,
		OutputFile:   *outputFileFlag,
		Model:        *modelFlag,
	}
}
