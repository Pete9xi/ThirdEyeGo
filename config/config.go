package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Token                     string   `json:"token"`
	Username                  string   `json:"username"`
	IsRealm                   bool     `json:"isRealm"`
	RealmInviteCode           string   `json:"realmInviteCode"`
	IP                        string   `json:"ip"`
	Port                      int      `json:"port"`
	Guild                     string   `json:"guild"`
	Channel                   string   `json:"channel"`
	AntiCheatEnabled          bool     `json:"antiCheatEnabled"`
	AntiCheatLogsChannel      string   `json:"antiCheatLogsChannel"`
	LogSystemCommands         bool     `json:"logSystemCommands"`
	SystemCommandsChannel     string   `json:"systemCommandsChannel"`
	SendWhisperMessages       bool     `json:"sendWhisperMessages"`
	UseEmbed                  bool     `json:"useEmbed"`
	SetColor                  []int    `json:"setColor"`
	SetTitle                  string   `json:"setTitle"`
	AuthType                  bool     `json:"authType"`
	Admins                    []string `json:"admins"`
	BlacklistDeviceTypes      []string `json:"blacklistDeviceTypes"`
	VoiceChannelCommandPrefix string   `json:"voiceChannelCommandPrefix"`
	VoiceChannelsCategory     string   `json:"voiceChannelsCategory"`
	VoiceAdminRoleID          string   `json:"voiceAdminRoleID"`
	LogBadActors              bool     `json:"logBadActors"`
	LogoURL                   string   `json:"logoURL"`
	OperatorsRoleID           string   `json:"operatorsRoleID"`
	ProfanityLogsChannel      string   `json:"profanityLogsChannel"`
	ProfanityFilter           bool     `json:"profanityFilter"`
	Debug                     bool     `json:"Debug"`
	IsDev                     bool     `json:"isDev"`
}

func Default() Config {
	return Config{
		Token:                     "",
		Username:                  "",
		IsRealm:                   false,
		RealmInviteCode:           "",
		IP:                        "",
		Port:                      19132,
		Guild:                     "",
		Channel:                   "",
		AntiCheatEnabled:          true,
		AntiCheatLogsChannel:      "",
		LogSystemCommands:         false,
		SystemCommandsChannel:     "",
		SendWhisperMessages:       false,
		UseEmbed:                  true,
		SetColor:                  []int{0, 153, 255},
		SetTitle:                  "My Servers Name!",
		AuthType:                  false,
		Admins:                    []string{},
		BlacklistDeviceTypes:      []string{},
		VoiceChannelCommandPrefix: "$",
		VoiceChannelsCategory:     "Voice Channels",
		VoiceAdminRoleID:          "",
		LogBadActors:              true,
		LogoURL:                   "https://i.imgur.com/XfoZ8XS.jpg",
		OperatorsRoleID:           "",
		ProfanityLogsChannel:      "",
		ProfanityFilter:           true,
		Debug:                     false,
		IsDev:                     false,
	}
}

func Load(path string) (Config, error) {
	cfg := Default()

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("read config file: %w", err)
	}

	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("parse config file: %w", err)
	}

	return cfg, nil
}
