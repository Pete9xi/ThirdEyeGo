// SPDX-License-Identifier: GPL-3.0-only
package bedrock

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"thirdeyego/bedrock/packets"
	"thirdeyego/config"
	"thirdeyego/discord"

	"github.com/df-mc/go-nethernet/endpoint"
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// StartBot connects to a direct BDS NetherNet endpoint using native transport
// capabilities. Choose the dependency revision matching the server protocol;
// this file does not translate between Minecraft protocol versions.
func StartBot(cfg config.Config) error {
	serverAddress, err := netherNetServerAddress(cfg.IP, cfg.Port)
	if err != nil {
		return err
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	signaling := endpoint.ClientConfig{Logger: logger}.New()

	fmt.Printf("[Version] Client supports protocol %d (%s)\n", protocol.CurrentProtocol, protocol.CurrentVersion)
	fmt.Printf("[NetherNet] Checking %s/v1/join\n", serverAddress)
	probeCtx, probeCancel := context.WithTimeout(context.Background(), 10*time.Second)
	status, probeErr := probeNetherNetServerStatus(probeCtx, serverAddress)
	probeCancel()
	if probeErr != nil {
		// Hidden discovery metadata may produce an empty response. A failed
		// status probe alone must not prevent a permitted login attempt.
		logger.Warn("Server status unavailable; protocol cannot be checked before login", "error", probeErr)
	} else if status.Protocol <= 0 {
		logger.Warn("Server did not advertise a usable protocol; attempting login with the compiled client protocol")
	} else {
		fmt.Printf("[Version] Server advertises protocol %d (%s)\n", status.Protocol, status.Version)
		if status.Protocol != protocol.CurrentProtocol {
			return fmt.Errorf("[Version] server protocol %d (%s) does not match this client's protocol %d (%s); select a matching dependency revision and rebuild; changing ports or only a version string will not implement that protocol", status.Protocol, status.Version, protocol.CurrentProtocol, protocol.CurrentVersion)
		}
	}

	fmt.Println("[Auth] Loading Minecraft authentication...")
	tokenSource, err := getTokenSource()
	if err != nil {
		return fmt.Errorf("[Auth] failed: %w", err)
	}

	trace := &netherNetHandshakeTrace{enabled: true, logger: logger}
	dialer := minecraft.Dialer{
		TokenSource: tokenSource,
		ErrorLog:    logger,
		PacketFunc:  trace.observe,
		DownloadResourcePack: func(id uuid.UUID, version string, current, total int) bool {
			fmt.Printf("[Packs] Skipping download %d/%d: %s (%s)\n", current+1, total, id, version)
			return false
		},
	}
	// endpoint.Client already requests non-trickle ICE. Preserve the native
	// minecraft.NetherNet path so identity and packet capabilities are retained.
	network := minecraft.NetherNet{Signaling: signaling, Log: logger}

	fmt.Printf("[NetherNet] Connecting to %s\n", serverAddress)
	started := time.Now()
	dialCtx, dialCancel := context.WithTimeout(context.Background(), 60*time.Second)
	conn, err := dialer.DialContextNetwork(dialCtx, network, serverAddress)
	dialCancel()
	if err != nil {
		return fmt.Errorf("[NetherNet] connection/login failed after %s; %s: %w", time.Since(started).Round(time.Millisecond), trace.summary(), err)
	}
	defer conn.Close()
	fmt.Println("[Minecraft] Connection/login completed")
	fmt.Printf("[NetherNet] Local address: %s\n", conn.LocalAddr())
	fmt.Printf("[NetherNet] Remote address: %s\n", conn.RemoteAddr())

	fmt.Println("[Minecraft] Completing spawn...")
	spawnCtx, spawnCancel := context.WithTimeout(context.Background(), 60*time.Second)
	err = conn.DoSpawnContext(spawnCtx)
	spawnCancel()
	if err != nil {
		return fmt.Errorf("[Minecraft] spawn failed; %s: %w", trace.summary(), err)
	}
	trace.stop()
	fmt.Println("[Minecraft] Bot spawned successfully!")

	fmt.Println("[Discord] Starting integration...")
	discordClient, err := discord.Start(conn, cfg)
	if err != nil {
		return fmt.Errorf("[Discord] startup failed: %w", err)
	}
	fmt.Println("[Discord] Integration started!")
	fmt.Println("[Packets] Starting packet loop...")
	packets.StartPacketLoop(conn, cfg, discordClient)
	return nil
}

func netherNetServerAddress(configuredHost string, port int) (string, error) {
	host := strings.TrimSpace(configuredHost)
	if host == "" || strings.ContainsAny(host, "/?#@\\ \t\r\n") {
		return "", fmt.Errorf("[Config] ip must contain a bare hostname or IP address, without http://, credentials, paths, or a port")
	}
	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		host = strings.TrimSuffix(strings.TrimPrefix(host, "["), "]")
	}
	if host == "" {
		return "", fmt.Errorf("[Config] ip must not be empty")
	}
	if strings.Contains(host, ":") && net.ParseIP(host) == nil {
		return "", fmt.Errorf("[Config] put the port in the separate port field; ip must be a hostname or IP literal")
	}
	if port < 1 || port > 65535 {
		return "", fmt.Errorf("[Config] port must be between 1 and 65535")
	}
	return "http://" + net.JoinHostPort(host, strconv.Itoa(port)), nil
}

type netherNetServerStatus struct {
	Protocol int    `json:"protocol"`
	Version  string `json:"version"`
}

// Read the documented JSON status directly: go-nethernet changed the return
// type of Client.PingContext between the old and current dependency versions.
func probeNetherNetServerStatus(ctx context.Context, serverAddress string) (netherNetServerStatus, error) {
	var status netherNetServerStatus
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverAddress+"/v1/join", nil)
	if err != nil {
		return status, err
	}
	req.Header.Set("User-Agent", "libhttpclient/1.0.0.0")
	client := &http.Client{CheckRedirect: func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse }}
	response, err := client.Do(req)
	if err != nil {
		return status, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return status, fmt.Errorf("status endpoint returned HTTP %d", response.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(response.Body, 16385))
	if err != nil {
		return status, err
	}
	if len(body) > 16384 {
		return status, fmt.Errorf("status response exceeded 16 KiB")
	}
	if len(body) == 0 {
		return status, nil
	}
	err = json.Unmarshal(body, &status)
	return status, err
}

// Trace packet names only. Outgoing callbacks run while packets are being
// queued; incoming callbacks run before decoding. Neither proves acceptance.
// Login payloads, tokens and raw packet contents are never printed here.
type netherNetHandshakeTrace struct {
	mu           sync.Mutex
	enabled      bool
	logger       *slog.Logger
	lastOutgoing string
	lastIncoming string
}

func (t *netherNetHandshakeTrace) observe(header packet.Header, _ []byte, _, _ net.Addr) {
	var name, direction string
	switch header.PacketID {
	case packet.IDRequestNetworkSettings:
		name, direction = "RequestNetworkSettings", "outgoing queued"
	case packet.IDLogin:
		name, direction = "Login", "outgoing queued"
	case packet.IDClientToServerHandshake:
		name, direction = "ClientToServerHandshake", "outgoing queued"
	case packet.IDResourcePackClientResponse:
		name, direction = "ResourcePackClientResponse", "outgoing queued"
	case packet.IDRequestChunkRadius:
		name, direction = "RequestChunkRadius", "outgoing queued"
	case packet.IDSetLocalPlayerAsInitialised:
		name, direction = "SetLocalPlayerAsInitialised", "outgoing queued"
	case packet.IDNetworkSettings:
		name, direction = "NetworkSettings", "incoming observed"
	case packet.IDServerToClientHandshake:
		name, direction = "ServerToClientHandshake", "incoming observed"
	case packet.IDPlayStatus:
		name, direction = "PlayStatus", "incoming observed"
	case packet.IDDisconnect:
		name, direction = "Disconnect", "incoming observed"
	case packet.IDResourcePacksInfo:
		name, direction = "ResourcePacksInfo", "incoming observed"
	case packet.IDResourcePackStack:
		name, direction = "ResourcePackStack", "incoming observed"
	case packet.IDStartGame:
		name, direction = "StartGame", "incoming observed"
	case packet.IDChunkRadiusUpdated:
		name, direction = "ChunkRadiusUpdated", "incoming observed"
	default:
		return
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if !t.enabled {
		return
	}
	if direction == "outgoing queued" {
		t.lastOutgoing = name
	} else {
		t.lastIncoming = name
	}
	t.logger.Debug("Minecraft handshake packet", "direction", direction, "packet", name, "id", header.PacketID)
}

func (t *netherNetHandshakeTrace) summary() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	outgoing, incoming := t.lastOutgoing, t.lastIncoming
	if outgoing == "" {
		outgoing = "none"
	}
	if incoming == "" {
		incoming = "none"
	}
	return fmt.Sprintf("last outgoing packet queued=%s; last incoming packet observed=%s", outgoing, incoming)
}

func (t *netherNetHandshakeTrace) stop() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.enabled = false
}
