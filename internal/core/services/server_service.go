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
	"context"
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
	"github.com/aryasoni98/wooak/internal/core/services/tracing"
	"go.uber.org/zap"
)

const (
	// DefaultPingTimeout is the default timeout for ping operations
	DefaultPingTimeout = 3 * time.Second
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
	ctx := context.Background()
	traceID := tracing.GetTraceIDOrNew(ctx)
	errorCtx := NewErrorContext("list servers").
		WithTraceID(string(traceID)).
		WithField("query", query)

	servers, err := s.serverRepository.ListServers(query)
	if err != nil {
		s.logger.Errorw("failed to list servers", "error", err, "trace_id", traceID, "query", query)
		return nil, WrapError(err, errorCtx)
	}

	// Sort servers using a multi-level sorting strategy:
	//
	// Level 1: Pinned vs Unpinned
	//   - Pinned servers (PinnedAt != zero) appear first
	//   - Unpinned servers appear after all pinned servers
	//
	// Level 2: Within Pinned Servers
	//   - Sort by PinnedAt timestamp descending (newest pinned first)
	//   - Most recently pinned servers appear at the top
	//
	// Level 3: Within Unpinned Servers
	//   - Sort alphabetically by Alias (A-Z)
	//   - Provides predictable, consistent ordering
	//
	// This ensures:
	//   - Frequently used (pinned) servers are easily accessible
	//   - Recent pins are more prominent than old pins
	//   - Unpinned servers maintain alphabetical order for easy scanning
	//
	// Example ordering:
	//   1. server-c (pinned 2 hours ago)
	//   2. server-a (pinned 1 day ago)
	//   3. server-b (pinned 1 week ago)
	//   4. alpha-server (unpinned, alphabetical)
	//   5. beta-server (unpinned, alphabetical)
	sort.SliceStable(servers, func(i, j int) bool {
		pi := !servers[i].PinnedAt.IsZero()
		pj := !servers[j].PinnedAt.IsZero()
		// First level: pinned vs unpinned
		if pi != pj {
			return pi // Pinned servers come first
		}
		// Second level: if both pinned, sort by pin time (newest first)
		if pi && pj {
			return servers[i].PinnedAt.After(servers[j].PinnedAt)
		}
		// Third level: if both unpinned, sort alphabetically
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
	ctx := context.Background()
	traceID := tracing.GetTraceIDOrNew(ctx)
	errorCtx := NewErrorContext("update server").
		WithTraceID(string(traceID)).
		WithFields(map[string]interface{}{
			"old_alias": server.Alias,
			"new_alias": newServer.Alias,
		})

	if err := validateServer(newServer); err != nil {
		s.logger.Warnw("validation failed on update", "error", err, "trace_id", traceID, "old_alias", server.Alias, "new_alias", newServer.Alias)
		return WrapErrorf(err, errorCtx, "validation failed for server update")
	}
	err := s.serverRepository.UpdateServer(server, newServer)
	if err != nil {
		s.logger.Errorw("failed to update server", "error", err, "trace_id", traceID, "old_alias", server.Alias, "new_alias", newServer.Alias)
		return WrapError(err, errorCtx)
	}
	return nil
}

// AddServer adds a new server to the repository.
func (s *serverService) AddServer(server domain.Server) error {
	ctx := context.Background()
	traceID := tracing.GetTraceIDOrNew(ctx)
	errorCtx := NewErrorContext("add server").
		WithTraceID(string(traceID)).
		WithFields(map[string]interface{}{
			"alias": server.Alias,
			"host":  server.Host,
		})

	if err := validateServer(server); err != nil {
		s.logger.Warnw("validation failed on add", "error", err, "trace_id", traceID, "alias", server.Alias, "host", server.Host)
		return WrapErrorf(err, errorCtx, "validation failed for server")
	}
	err := s.serverRepository.AddServer(server)
	if err != nil {
		s.logger.Errorw("failed to add server", "error", err, "trace_id", traceID, "alias", server.Alias, "host", server.Host)
		return WrapError(err, errorCtx)
	}
	return nil
}

// DeleteServer removes a server from the repository.
func (s *serverService) DeleteServer(server domain.Server) error {
	ctx := context.Background()
	traceID := tracing.GetTraceIDOrNew(ctx)
	errorCtx := NewErrorContext("delete server").
		WithTraceID(string(traceID)).
		WithField("alias", server.Alias)

	err := s.serverRepository.DeleteServer(server)
	if err != nil {
		s.logger.Errorw("failed to delete server", "error", err, "trace_id", traceID, "alias", server.Alias)
		return WrapError(err, errorCtx)
	}
	return nil
}

// SetPinned sets or clears a pin timestamp for the server alias.
func (s *serverService) SetPinned(alias string, pinned bool) error {
	err := s.serverRepository.SetPinned(alias, pinned)
	if err != nil {
		s.logger.Errorw("failed to set pin state", "error", err, "alias", alias, "pinned", pinned)
		return fmt.Errorf("failed to set pin state (alias: %q, pinned: %v): %w", alias, pinned, err)
	}
	return nil
}

// SSH starts an interactive SSH session to the given alias using the system's ssh client.
func (s *serverService) SSH(alias string) error {
	ctx := context.Background()
	traceID := tracing.GetTraceIDOrNew(ctx)

	// Validate alias format for security
	if !isValidAlias(alias) {
		errorCtx := NewErrorContext("SSH connection").
			WithTraceID(string(traceID)).
			WithField("alias", alias)
		return NewSecurityError(errorCtx, "invalid alias format: alias must contain only alphanumeric characters, dots, dashes, and underscores")
	}

	// Additional security checks
	if err := s.validateSSHAccess(alias); err != nil {
		errorCtx := NewErrorContext("SSH connection").
			WithTraceID(string(traceID)).
			WithField("alias", alias)
		return WrapSecurityError(err, errorCtx, "SSH access validation failed")
	}

	s.logger.Infow("ssh start", "trace_id", traceID, "alias", alias)
	cmd := exec.Command("ssh", alias)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		s.logger.Errorw("ssh command failed", "trace_id", traceID, "alias", alias, "error", err)
		errorCtx := NewErrorContext("SSH connection").
			WithTraceID(string(traceID)).
			WithField("alias", alias)
		return WrapError(err, errorCtx)
	}

	if err := s.serverRepository.RecordSSH(alias); err != nil {
		s.logger.Errorw("failed to record ssh metadata", "trace_id", traceID, "alias", alias, "error", err)
		// Don't fail the SSH connection if metadata recording fails
	}

	s.logger.Infow("ssh end", "trace_id", traceID, "alias", alias)
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

	dialer := net.Dialer{Timeout: DefaultPingTimeout}
	conn, err := dialer.Dial("tcp", addr)
	if err != nil {
		return false, time.Since(start), err
	}
	_ = conn.Close()
	return true, time.Since(start), nil
}

// resolveSSHDestination uses `ssh -G <alias>` to extract HostName and Port from the user's SSH config.
//
// The `ssh -G` command outputs the resolved configuration for an alias, including
// all inherited settings from Host blocks and wildcards. This function parses
// the output to extract the actual hostname and port that would be used for connection.
//
// Algorithm:
//  1. Execute `ssh -G <alias>` to get resolved config (includes all inherited settings)
//  2. Parse output line by line looking for "hostname" and "port" directives
//  3. Extract the first occurrence of each directive (SSH config allows multiple, last wins)
//  4. Return resolved values, or defaults (alias as hostname, port 22) if not found
//
// Edge Cases:
//   - Empty alias: returns ok=false immediately
//   - Command failure: returns ok=false (alias may not exist)
//   - Missing hostname: uses alias as fallback
//   - Missing port: uses default port 22
//
// Returns host, port, ok where ok=false if resolution failed (e.g., alias doesn't exist).
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

// isValidAlias validates that an alias is safe for SSH command execution.
//
// This function implements a defense-in-depth security strategy to prevent
// command injection attacks when aliases are used in shell commands.
//
// Security Layers:
//
//  1. Length validation: non-empty and max 100 characters
//     - Prevents buffer overflow attacks
//     - Limits DoS potential from extremely long strings
//
//  2. Path traversal prevention: blocks "..", "/", "\"
//     - Prevents directory traversal attacks
//     - Blocks absolute and relative path references
//
//  3. Command injection prevention: blocks shell metacharacters
//     - Blocks: ; & | ` $ ( ) < > " ' \n \r \t
//     - Prevents command chaining and execution
//     - Blocks variable expansion and redirection
//
//  4. Whitelist validation: only allows safe characters
//     - Pattern: ^[A-Za-z0-9_.-]+$
//     - Allows: letters, digits, dots, dashes, underscores
//     - Ensures only safe characters can be used
//
// This multi-layer approach ensures aliases can only contain safe characters
// that cannot be exploited to execute arbitrary commands or access unauthorized paths.
func isValidAlias(alias string) bool {
	// Check length to prevent buffer overflow and DoS
	if alias == "" || len(alias) > 100 {
		return false
	}

	// Check for path traversal attempts
	if strings.Contains(alias, "..") || strings.Contains(alias, "/") || strings.Contains(alias, "\\") {
		return false
	}

	// Check for command injection attempts - block shell metacharacters
	dangerousChars := []string{";", "&", "|", "`", "$", "(", ")", "<", ">", "\"", "'", "\n", "\r", "\t"}
	for _, char := range dangerousChars {
		if strings.Contains(alias, char) {
			return false
		}
	}

	// Whitelist: only allow safe characters (alphanumeric, dots, dashes, underscores)
	matched, _ := regexp.MatchString(`^[A-Za-z0-9_.-]+$`, alias)
	return matched
}

// validateSSHAccess performs additional security checks before SSH execution
func (s *serverService) validateSSHAccess(alias string) error {
	// Check if alias exists in repository
	servers, err := s.serverRepository.ListServers("")
	if err != nil {
		errorCtx := NewErrorContext("validate SSH access").
			WithField("alias", alias)
		return WrapError(err, errorCtx)
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
		errorCtx := NewErrorContext("validate SSH access").
			WithField("alias", alias)
		return NewSecurityError(errorCtx, "alias not found in known servers list")
	}

	return nil
}
