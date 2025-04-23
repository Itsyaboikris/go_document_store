package config

import (
	"os"
	"strings"
)

func GetPeers() []string {
	peers := os.Getenv("PEERS")
	if peers == "" {
		return []string{}
	}
	return strings.Split(peers, ",")
}
