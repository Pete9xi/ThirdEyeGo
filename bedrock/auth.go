package bedrock

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sandertv/gophertunnel/minecraft/auth"
	"golang.org/x/oauth2"
)

const tokenFile = "token.json"

func getTokenSource() (oauth2.TokenSource, error) {
	var token oauth2.Token

	data, err := os.ReadFile(tokenFile)

	if err == nil {
		if err := json.Unmarshal(data, &token); err != nil {
			return nil, fmt.Errorf("decode token: %w", err)
		}

		fmt.Println("Loaded saved authentication token.")

		return auth.RefreshTokenSource(&token), nil
	}

	if !os.IsNotExist(err) {
		return nil, fmt.Errorf("read token file: %w", err)
	}

	fmt.Println("No saved authentication found.")
	fmt.Println("Starting Microsoft authentication...")

	tokenPtr, err := auth.RequestLiveToken()
	if err != nil {
		return nil, fmt.Errorf("authenticate: %w", err)
	}

	token = *tokenPtr

	data, err = json.MarshalIndent(&token, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode token: %w", err)
	}

	if err := os.WriteFile(tokenFile, data, 0600); err != nil {
		return nil, fmt.Errorf("save token: %w", err)
	}

	fmt.Println("Authentication token saved to", tokenFile)

	return auth.RefreshTokenSource(&token), nil
}
