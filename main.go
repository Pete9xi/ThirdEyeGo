package main

import (
	"fmt"
	"thirdeyego/bedrock"
	"thirdeyego/config"
)

func main() {
	fmt.Println("=== ThirdEyeGo starting ===")
	cfg, err := config.Load("config.json")
	if err != nil {
		fmt.Println("Config error:", err)
		return
	}

	fmt.Println("Config loaded!")

	token := cfg.Token

	if token == "" {
		fmt.Println("DISCORD_TOKEN is not set")
		return
	}

	if err := bedrock.StartBot(cfg); err != nil {
		fmt.Printf("Bot stopped: %v\n", err)
	}
}
