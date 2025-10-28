// Copyright 2025.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ssh_config_file

import (
	"strconv"
	"strings"
	"time"

	"github.com/aryasoni98/wooak/internal/core/domain"
	"github.com/kevinburke/ssh_config"
)

// toDomainServer converts ssh_config.Config to a slice of domain.Server.
func (r *Repository) toDomainServer(cfg *ssh_config.Config) []domain.Server {
	servers := make([]domain.Server, 0, len(cfg.Hosts))
	for _, host := range cfg.Hosts {

		aliases := make([]string, 0, len(host.Patterns))

		for _, pattern := range host.Patterns {
			alias := pattern.String()
			// Skip if alias contains wildcards (not a concrete Host)
			if strings.ContainsAny(alias, "!*?[]") {
				continue
			}
			aliases = append(aliases, alias)
		}
		if len(aliases) == 0 {
			continue
		}
		server := &domain.Server{
			Alias:         aliases[0],
			Aliases:       aliases,
			Port:          22,
			IdentityFiles: []string{},
		}

		for _, node := range host.Nodes {
			kvNode, ok := node.(*ssh_config.KV)
			if !ok {
				continue
			}

			r.mapKVToServer(server, kvNode)
		}

		servers = append(servers, *server)
	}

	return servers
}

// mapKVToServer maps an ssh_config.KV node to the corresponding fields in domain.Server.
func (r *Repository) mapKVToServer(server *domain.Server, kvNode *ssh_config.KV) {
	key := strings.ToLower(kvNode.Key)
	value := kvNode.Value

	// Try mapping in order of categories
	if r.mapBasicConfig(server, key, value) {
		return
	}
	if r.mapConnectionConfig(server, key, value) {
		return
	}
	if r.mapForwardingConfig(server, key, value) {
		return
	}
	if r.mapAuthenticationConfig(server, key, value) {
		return
	}
	if r.mapSecurityConfig(server, key, value) {
		return
	}
	if r.mapEnvironmentConfig(server, key, value) {
		return
	}
	r.mapDebugConfig(server, key, value)
}

// mapBasicConfig maps basic SSH configuration fields
func (r *Repository) mapBasicConfig(server *domain.Server, key, value string) bool {
	switch key {
	case "hostname":
		server.Host = value
	case "user":
		server.User = value
	case "port":
		port, err := strconv.Atoi(value)
		if err == nil {
			server.Port = port
		}
	case "identityfile":
		server.IdentityFiles = append(server.IdentityFiles, value)
	default:
		return false
	}
	return true
}

// mapConnectionConfig maps connection and proxy configuration fields
func (r *Repository) mapConnectionConfig(server *domain.Server, key, value string) bool {
	switch key {
	case "proxycommand":
		server.ProxyCommand = value
	case "proxyjump":
		server.ProxyJump = value
	case "remotecommand":
		server.RemoteCommand = value
	case "requesttty":
		server.RequestTTY = value
	case "sessiontype":
		server.SessionType = value
	case "connecttimeout":
		server.ConnectTimeout = value
	case "connectionattempts":
		server.ConnectionAttempts = value
	case "bindaddress":
		server.BindAddress = value
	case "bindinterface":
		server.BindInterface = value
	case "addressfamily":
		server.AddressFamily = value
	case "exitonforwardfailure":
		server.ExitOnForwardFailure = value
	case "ipqos":
		server.IPQoS = value
	case "canonicalizehostname":
		server.CanonicalizeHostname = value
	case "canonicaldomains":
		server.CanonicalDomains = value
	case "canonicalizefallbacklocal":
		server.CanonicalizeFallbackLocal = value
	case "canonicalizemaxdots":
		server.CanonicalizeMaxDots = value
	case "canonicalizepermittedcnames":
		server.CanonicalizePermittedCNAMEs = value
	case "serveraliveinterval":
		server.ServerAliveInterval = value
	case "serveralivecountmax":
		server.ServerAliveCountMax = value
	case "compression":
		server.Compression = value
	case "tcpkeepalive":
		server.TCPKeepAlive = value
	case "batchmode":
		server.BatchMode = value
	case "controlmaster":
		server.ControlMaster = value
	case "controlpath":
		server.ControlPath = value
	case "controlpersist":
		server.ControlPersist = value
	default:
		return false
	}
	return true
}

// mapForwardingConfig maps port forwarding and agent forwarding fields
func (r *Repository) mapForwardingConfig(server *domain.Server, key, value string) bool {
	switch key {
	case "localforward":
		cliFormat := r.convertConfigForwardToCLIFormat(value)
		server.LocalForward = append(server.LocalForward, cliFormat)
	case "remoteforward":
		cliFormat := r.convertConfigForwardToCLIFormat(value)
		server.RemoteForward = append(server.RemoteForward, cliFormat)
	case "dynamicforward":
		server.DynamicForward = append(server.DynamicForward, value)
	case "clearallforwardings":
		server.ClearAllForwardings = value
	case "gatewayports":
		server.GatewayPorts = value
	case "forwardagent":
		server.ForwardAgent = value
	case "forwardx11":
		server.ForwardX11 = value
	case "forwardx11trusted":
		server.ForwardX11Trusted = value
	default:
		return false
	}
	return true
}

// mapAuthenticationConfig maps authentication-related fields
func (r *Repository) mapAuthenticationConfig(server *domain.Server, key, value string) bool {
	switch key {
	case "pubkeyauthentication":
		server.PubkeyAuthentication = value
	case "pubkeyacceptedalgorithms", "pubkeyacceptedkeytypes":
		// PubkeyAcceptedKeyTypes is deprecated alias for PubkeyAcceptedAlgorithms (since OpenSSH 8.5)
		server.PubkeyAcceptedAlgorithms = value
	case "hostbasedacceptedalgorithms", "hostbasedkeytypes", "hostbasedacceptedkeytypes":
		// HostbasedKeyTypes and HostbasedAcceptedKeyTypes are deprecated aliases (since OpenSSH 8.5)
		server.HostbasedAcceptedAlgorithms = value
	case "passwordauthentication":
		server.PasswordAuthentication = value
	case "preferredauthentications":
		server.PreferredAuthentications = value
	case "identitiesonly":
		server.IdentitiesOnly = value
	case "addkeystoagent":
		server.AddKeysToAgent = value
	case "identityagent":
		server.IdentityAgent = value
	case "kbdinteractiveauthentication", "challengeresponseauthentication":
		// ChallengeResponseAuthentication is deprecated alias for KbdInteractiveAuthentication
		server.KbdInteractiveAuthentication = value
	case "numberofpasswordprompts":
		server.NumberOfPasswordPrompts = value
	default:
		return false
	}
	return true
}

// mapSecurityConfig maps security-related fields
func (r *Repository) mapSecurityConfig(server *domain.Server, key, value string) bool {
	switch key {
	case "stricthostkeychecking":
		server.StrictHostKeyChecking = value
	case "checkhostip":
		server.CheckHostIP = value
	case "fingerprinthash":
		server.FingerprintHash = value
	case "userknownhostsfile":
		server.UserKnownHostsFile = value
	case "hostkeyalgorithms":
		server.HostKeyAlgorithms = value
	case "macs":
		server.MACs = value
	case "ciphers":
		server.Ciphers = value
	case "kexalgorithms":
		server.KexAlgorithms = value
	case "verifyhostkeydns":
		server.VerifyHostKeyDNS = value
	case "updatehostkeys":
		server.UpdateHostKeys = value
	case "hashknownhosts":
		server.HashKnownHosts = value
	case "visualhostkey":
		server.VisualHostKey = value
	default:
		return false
	}
	return true
}

// mapEnvironmentConfig maps environment and command execution fields
func (r *Repository) mapEnvironmentConfig(server *domain.Server, key, value string) bool {
	switch key {
	case "localcommand":
		server.LocalCommand = value
	case "permitlocalcommand":
		server.PermitLocalCommand = value
	case "escapechar":
		server.EscapeChar = value
	case "sendenv":
		server.SendEnv = append(server.SendEnv, value)
	case "setenv":
		server.SetEnv = append(server.SetEnv, value)
	default:
		return false
	}
	return true
}

// mapDebugConfig maps debugging-related fields
func (r *Repository) mapDebugConfig(server *domain.Server, key, value string) bool {
	switch key {
	case "loglevel":
		server.LogLevel = value
	default:
		return false
	}
	return true
}

// mergeMetadata merges additional metadata into the servers.
func (r *Repository) mergeMetadata(servers []domain.Server, metadata map[string]ServerMetadata) []domain.Server {
	for i, server := range servers {
		servers[i].LastSeen = time.Time{}

		if meta, exists := metadata[server.Alias]; exists {
			servers[i].Tags = meta.Tags
			servers[i].SSHCount = meta.SSHCount

			if meta.LastSeen != "" {
				if lastSeen, err := time.Parse(time.RFC3339, meta.LastSeen); err == nil {
					servers[i].LastSeen = lastSeen
				}
			}

			if meta.PinnedAt != "" {
				if pinnedAt, err := time.Parse(time.RFC3339, meta.PinnedAt); err == nil {
					servers[i].PinnedAt = pinnedAt
				}
			}
		}
	}
	return servers
}

// convertConfigForwardToCLIFormat converts SSH config format forwarding spec to CLI format.
// Config format: [bind_address:]port host:hostport
// CLI format: [bind_address:]port:host:hostport
func (r *Repository) convertConfigForwardToCLIFormat(forward string) string {
	// Find the last space which separates the local part from the remote part
	lastSpace := strings.LastIndex(forward, " ")
	if lastSpace != -1 {
		localPart := forward[:lastSpace]
		remotePart := forward[lastSpace+1:]
		// Join them with a colon for CLI format
		return localPart + ":" + remotePart
	}
	// If no space found, return as-is (might already be in CLI format)
	return forward
}
