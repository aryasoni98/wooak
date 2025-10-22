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

package services

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aryasoni98/wooak/internal/core/domain"
	"github.com/aryasoni98/wooak/internal/core/ports"
	"go.uber.org/zap"
)

type serverService struct {
	serverRepository ports.ServerRepository
	logger           *zap.SugaredLogger
}

// NewServerService creates a new instance of serverService.
func NewServerService(logger *zap.SugaredLogger, sr ports.ServerRepository) ports.ServerService {
	return &serverService{
		logger:           logger,
		serverRepository: sr,
	}
}

// ListServers returns a list of servers sorted with pinned on top.
func (s *serverService) ListServers(query string) ([]domain.Server, error) {
	servers, err := s.serverRepository.ListServers(query)
	if err != nil {
		s.logger.Errorw("failed to list servers", "error", err)
		return nil, err
	}

	// Sort: pinned first (PinnedAt non-zero), then by PinnedAt desc, then by Alias asc.
	sort.SliceStable(servers, func(i, j int) bool {
		pi := !servers[i].PinnedAt.IsZero()
		pj := !servers[j].PinnedAt.IsZero()
		if pi != pj {
			return pi
		}
		if pi && pj {
			return servers[i].PinnedAt.After(servers[j].PinnedAt)
		}
		return servers[i].Alias < servers[j].Alias
	})

	return servers, nil
}

// validateServer performs core validation of server fields.
func validateServer(srv domain.Server) error {
	if strings.TrimSpace(srv.Alias) == "" {
		return fmt.Errorf("alias is required")
	}
	if ok, _ := regexp.MatchString(`^[A-Za-z0-9_.-]+$`, srv.Alias); !ok {
		return fmt.Errorf("alias may contain letters, digits, dot, dash, underscore")
	}
	if strings.TrimSpace(srv.Host) == "" {
		return fmt.Errorf("Host/IP is required")
	}
	if ip := net.ParseIP(srv.Host); ip == nil {
		if strings.Contains(srv.Host, " ") {
			return fmt.Errorf("host must not contain spaces")
		}
		if ok, _ := regexp.MatchString(`^[A-Za-z0-9.-]+$`, srv.Host); !ok {
			return fmt.Errorf("host contains invalid characters")
		}
		if strings.HasPrefix(srv.Host, ".") || strings.HasSuffix(srv.Host, ".") {
			return fmt.Errorf("host must not start or end with a dot")
		}
		for _, lbl := range strings.Split(srv.Host, ".") {
			if lbl == "" {
				return fmt.Errorf("host must not contain empty labels")
			}
			if strings.HasPrefix(lbl, "-") || strings.HasSuffix(lbl, "-") {
				return fmt.Errorf("hostname labels must not start or end with a hyphen")
			}
		}
	}
	if srv.Port != 0 && (srv.Port < 1 || srv.Port > 65535) {
		return fmt.Errorf("port must be a number between 1 and 65535")
	}
	return nil
}

// UpdateServer updates an existing server with new details.
func (s *serverService) UpdateServer(server domain.Server, newServer domain.Server) error {
	if err := validateServer(newServer); err != nil {
		s.logger.Warnw("validation failed on update", "error", err, "server", newServer)
		return err
	}
	err := s.serverRepository.UpdateServer(server, newServer)
	if err != nil {
		s.logger.Errorw("failed to update server", "error", err, "server", server)
	}
	return err
}

// AddServer adds a new server to the repository.
func (s *serverService) AddServer(server domain.Server) error {
	if err := validateServer(server); err != nil {
		s.logger.Warnw("validation failed on add", "error", err, "server", server)
		return err
	}
	err := s.serverRepository.AddServer(server)
	if err != nil {
		s.logger.Errorw("failed to add server", "error", err, "server", server)
	}
	return err
}

// DeleteServer removes a server from the repository.
func (s *serverService) DeleteServer(server domain.Server) error {
	err := s.serverRepository.DeleteServer(server)
	if err != nil {
		s.logger.Errorw("failed to delete server", "error", err, "server", server)
	}
	return err
}

// SetPinned sets or clears a pin timestamp for the server alias.
func (s *serverService) SetPinned(alias string, pinned bool) error {
	err := s.serverRepository.SetPinned(alias, pinned)
	if err != nil {
		s.logger.Errorw("failed to set pin state", "error", err, "alias", alias, "pinned", pinned)
	}
	return err
}

// SSH starts an interactive SSH session to the given alias using the system's ssh client.
func (s *serverService) SSH(alias string) error {
	// Validate alias format for security
	if !isValidAlias(alias) {
		return fmt.Errorf("invalid alias format: %s", alias)
	}

	// Additional security checks
	if err := s.validateSSHAccess(alias); err != nil {
		return fmt.Errorf("ssh access validation failed: %w", err)
	}

	s.logger.Infow("ssh start", "alias", alias)
	cmd := exec.Command("ssh", alias)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		s.logger.Errorw("ssh command failed", "alias", alias, "error", err)
		return err
	}

	if err := s.serverRepository.RecordSSH(alias); err != nil {
		s.logger.Errorw("failed to record ssh metadata", "alias", alias, "error", err)
	}

	s.logger.Infow("ssh end", "alias", alias)
	return nil
}

// Ping checks if the server is reachable on its SSH port.
func (s *serverService) Ping(server domain.Server) (bool, time.Duration, error) {
	start := time.Now()

	host, port, ok := resolveSSHDestination(server.Alias)
	if !ok {

		host = strings.TrimSpace(server.Host)
		if host == "" {
			host = server.Alias
		}
		if server.Port > 0 {
			port = server.Port
		} else {
			port = 22
		}
	}
	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))

	dialer := net.Dialer{Timeout: 3 * time.Second}
	conn, err := dialer.Dial("tcp", addr)
	if err != nil {
		return false, time.Since(start), err
	}
	_ = conn.Close()
	return true, time.Since(start), nil
}

// resolveSSHDestination uses `ssh -G <alias>` to extract HostName and Port from the user's SSH config.
// Returns host, port, ok where ok=false if resolution failed.
func resolveSSHDestination(alias string) (string, int, bool) {
	alias = strings.TrimSpace(alias)
	if alias == "" {
		return "", 0, false
	}
	cmd := exec.Command("ssh", "-G", alias)
	out, err := cmd.Output()
	if err != nil {
		return "", 0, false
	}
	host := ""
	port := 0
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "hostname ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				host = parts[1]
			}
		}
		if strings.HasPrefix(line, "port ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				if p, err := strconv.Atoi(parts[1]); err == nil {
					port = p
				}
			}
		}
	}
	if host == "" {
		host = alias
	}
	if port == 0 {
		port = 22
	}
	return host, port, true
}

// isValidAlias validates that an alias is safe for SSH command execution
func isValidAlias(alias string) bool {
	// Check length
	if alias == "" || len(alias) > 100 {
		return false
	}

	// Check for dangerous characters
	if strings.Contains(alias, "..") || strings.Contains(alias, "/") || strings.Contains(alias, "\\") {
		return false
	}

	// Check for command injection attempts
	dangerousChars := []string{";", "&", "|", "`", "$", "(", ")", "<", ">", "\"", "'", "\n", "\r", "\t"}
	for _, char := range dangerousChars {
		if strings.Contains(alias, char) {
			return false
		}
	}

	// Allow only alphanumeric, dots, dashes, underscores
	matched, _ := regexp.MatchString(`^[A-Za-z0-9_.-]+$`, alias)
	return matched
}

// validateSSHAccess performs additional security checks before SSH execution
func (s *serverService) validateSSHAccess(alias string) error {
	// Check if alias exists in repository
	servers, err := s.serverRepository.ListServers("")
	if err != nil {
		return fmt.Errorf("failed to list servers for validation: %w", err)
	}

	// Verify alias exists in our known servers
	aliasExists := false
	for _, server := range servers {
		if server.Alias == alias {
			aliasExists = true
			break
		}
	}

	if !aliasExists {
		return fmt.Errorf("alias '%s' not found in known servers", alias)
	}

	return nil
}
