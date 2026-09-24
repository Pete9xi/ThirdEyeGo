package bedrock

import (
	"context"
	"fmt"
	"thirdeyego/bedrock/packets"
	"thirdeyego/config"
	"time"

	"github.com/sandertv/gophertunnel/minecraft"
)

func StartBot(cfg config.Config) error {
	serverAddress := fmt.Sprintf("%s:%d", cfg.IP, cfg.Port)
	tokenSource, err := getTokenSource()
	if err != nil {
		return fmt.Errorf("authentication: %w", err)
	}

	dialer := minecraft.Dialer{
		TokenSource: tokenSource,
	}

	fmt.Printf("Connecting to %s...\n", serverAddress)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		60*time.Second,
	)
	defer cancel()

	conn, err := dialer.DialContext(
		ctx,
		"raknet",
		serverAddress,
	)
	if err != nil {
		return fmt.Errorf("connect to server: %w", err)
	}

	defer conn.Close()

	fmt.Println("RakNet connection established!")

	if err := conn.DoSpawn(); err != nil {
		return fmt.Errorf("spawn: %w", err)
	}

	fmt.Println("Connected and spawned into the world!")

	packets.StartPacketLoop(conn, cfg)

	return nil
}
