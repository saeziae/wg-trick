package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Define the structure for the configuration
type Interface struct {
	PublicKey  string
	PrivateKey string
	Address    string
	Endpoint   string
	ListenPort string
	Mask       string
}

type Peer struct {
	PublicKey           string
	AllowedIPs          string
	Endpoint            string
	IsGateway           bool
	PersistentKeepalive string
}

type WireGuardConfig struct {
	Interface
	Peers []Peer
}

func ReadWireGuardConfig(filename string) (*WireGuardConfig, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config WireGuardConfig

	scanner := bufio.NewScanner(file)

	var section string
	var currentPeer *Peer

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines or comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.Split(line, "#")[0] // Remove comments

		// Check for section headers
		if strings.HasPrefix(line, "[") {
			section = strings.Split(line, "]")[0]
			section = line[1:] // Extract section name
			// new section, add the last peer
			if currentPeer != nil {
				config.Peers = append(config.Peers, *currentPeer)
			}
			if section == "Interface" {
				currentPeer = nil
			} else if section == "Peer" {
				// Create a new PeerConfig when encountering a new [Peer] section
				currentPeer = &Peer{}
			}
			continue
		}

		// Parse key-value pairs within sections
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])

		switch section {
		case "Interface":
			switch key {
			case "PublicKey":
				config.Interface.PublicKey = value
			case "PrivateKey":
				config.Interface.PrivateKey = value
			case "Address":
				config.Interface.Address = value
			case "Endpoint":
				config.Interface.Endpoint = value
			case "ListenPort":
				config.Interface.ListenPort = value
			case "Mask":
				config.Interface.Mask = value
			}
		case "Peer":
			switch key {
			case "PublicKey":
				currentPeer.PublicKey = value
			case "AllowedIPs":
				currentPeer.AllowedIPs = value
			case "Endpoint":
				currentPeer.Endpoint = value
			case "PersistentKeepalive":
				currentPeer.PersistentKeepalive = value
			case "IsGateway":
				currentPeer.IsGateway = (value == "true" || value == "True")
			}
		}
	}
	// don't forget to add the last peer
	if currentPeer != nil {
		config.Peers = append(config.Peers, *currentPeer)
	}

	fmt.Print(config)
	return &config, nil
}

// GeneratePeerIni generates an INI-formatted string for a single Peer configuration
func GeneratePeerIni(peer Peer, ip string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Address = %s\n", ip))
	sb.WriteString("[Peer]\n")
	sb.WriteString(fmt.Sprintf("PublicKey = %s\n", peer.PublicKey))
	sb.WriteString(fmt.Sprintf("AllowedIPs = %s\n", peer.AllowedIPs))
	sb.WriteString(fmt.Sprintf("Endpoint = %s\n", peer.Endpoint))
	sb.WriteString(fmt.Sprintf("PersistentKeepalive = %s\n", peer.PersistentKeepalive))
	return sb.String()
}
