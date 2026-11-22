# apihop-go

[![Go Reference](https://pkg.go.dev/badge/github.com/APIHOP/apihop-go.svg)](https://pkg.go.dev/github.com/APIHOP/apihop-go)

Lightweight, zero-dependency Go SDK for **[API Hop](https://apihop.ai/en)** — a unified API gateway giving developers one API key and one balance across leading LLMs (Claude, OpenAI, Gemini, DeepSeek, Kimi, MiniMax), image/video models, and **32+ real-time data APIs** (Google Search, Twitter/X, Similarweb, Airbnb, Amazon, YouTube, LinkedIn, Google Maps, and more).

## Useful Links

- **Platform Home**: [API Hop Unified AI & Data APIs](https://apihop.ai/en)
- **32+ Real-Time Data APIs**: [API Hop Data APIs Directory](https://apihop.ai/en/data-apis)
- **AI Models & Pricing**: [Models Directory](https://apihop.ai/en/models) · [Pricing](https://apihop.ai/en/pricing)
- **API Installation Guide**: [Getting Started Docs](https://apihop.ai/en/docs/api-reference/ai-installation-guide)
- **Machine-Readable Catalog**: [public-api-catalog.json](https://apihop.ai/docs/public-api-catalog.json)

## Installation

```bash
go get github.com/APIHOP/apihop-go
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	apihop "github.com/APIHOP/apihop-go"
)

func main() {
	// Reads APIHOP_API_KEY from the environment when empty string is passed.
	client := apihop.NewClient("")
	ctx := context.Background()

	// 1. List available models
	models, err := client.ListModels(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Models:", models)

	// 2. Call Chat Completions
	chat, err := client.ChatCompletions(ctx, apihop.ChatRequest{
		Model: "deepseek-chat",
		Messages: []apihop.ChatMessage{
			{Role: "user", Content: "Summarize what a unified AI and data API gateway does in one sentence."},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Chat:", chat)

	// 3. Call Twitter User Data API with the same API key
	user, err := client.TwitterUserByUsername(ctx, "openai")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Twitter User:", user)
}
```

## Calling Any Data API Endpoint

Use `client.DoJSON` to call any of the 285+ endpoints listed in the [API Hop API Reference](https://apihop.ai/en/docs/api-reference/models) (including Google Search, Airbnb, Amazon, Similarweb, Google Maps, and YouTube Data APIs).

## License

MIT
