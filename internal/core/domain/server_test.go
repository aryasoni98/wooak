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

package domain

import (
	"testing"
	"time"
)

func TestServerCreation(t *testing.T) {
	tests := []struct {
		name     string
		server   Server
		expected Server
	}{
		{
			name: "basic server",
			server: Server{
				Alias: "test-server",
				Host:  "example.com",
				User:  "testuser",
				Port:  22,
			},
			expected: Server{
				Alias: "test-server",
				Host:  "example.com",
				User:  "testuser",
				Port:  22,
			},
		},
		{
			name: "server with identity files",
			server: Server{
				Alias:         "secure-server",
				Host:          "secure.example.com",
				User:          "admin",
				Port:          2222,
				IdentityFiles: []string{"/home/user/.ssh/id_rsa", "/home/user/.ssh/id_ed25519"},
			},
			expected: Server{
				Alias:         "secure-server",
				Host:          "secure.example.com",
				User:          "admin",
				Port:          2222,
				IdentityFiles: []string{"/home/user/.ssh/id_rsa", "/home/user/.ssh/id_ed25519"},
			},
		},
		{
			name: "server with proxy jump",
			server: Server{
				Alias:     "proxy-server",
				Host:      "internal.example.com",
				User:      "user",
				Port:      22,
				ProxyJump: "bastion.example.com",
			},
			expected: Server{
				Alias:     "proxy-server",
				Host:      "internal.example.com",
				User:      "user",
				Port:      22,
				ProxyJump: "bastion.example.com",
			},
		},
		{
			name: "server with tags",
			server: Server{
				Alias: "tagged-server",
				Host:  "tagged.example.com",
				User:  "user",
				Port:  22,
				Tags:  []string{"production", "web", "critical"},
			},
			expected: Server{
				Alias: "tagged-server",
				Host:  "tagged.example.com",
				User:  "user",
				Port:  22,
				Tags:  []string{"production", "web", "critical"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.server.Alias != tt.expected.Alias {
				t.Errorf("Alias = %v, want %v", tt.server.Alias, tt.expected.Alias)
			}
			if tt.server.Host != tt.expected.Host {
				t.Errorf("Host = %v, want %v", tt.server.Host, tt.expected.Host)
			}
			if tt.server.User != tt.expected.User {
				t.Errorf("User = %v, want %v", tt.server.User, tt.expected.User)
			}
			if tt.server.Port != tt.expected.Port {
				t.Errorf("Port = %v, want %v", tt.server.Port, tt.expected.Port)
			}
		})
	}
}

func TestServerWithTimestamps(t *testing.T) {
	now := time.Now()
	server := Server{
		Alias:    "timestamped-server",
		Host:     "example.com",
		User:     "user",
		Port:     22,
		LastSeen: now,
		PinnedAt: now,
		SSHCount: 5,
	}

	if server.LastSeen.IsZero() {
		t.Error("LastSeen should not be zero")
	}
	if server.PinnedAt.IsZero() {
		t.Error("PinnedAt should not be zero")
	}
	if server.SSHCount != 5 {
		t.Errorf("SSHCount = %v, want 5", server.SSHCount)
	}
}

func TestServerWithSecuritySettings(t *testing.T) {
	server := Server{
		Alias:                 "secure-server",
		Host:                  "secure.example.com",
		User:                  "admin",
		Port:                  22,
		StrictHostKeyChecking: "yes",
		CheckHostIP:           "yes",
		FingerprintHash:       "sha256",
		VerifyHostKeyDNS:      "yes",
	}

	if server.StrictHostKeyChecking != "yes" {
		t.Errorf("StrictHostKeyChecking = %v, want yes", server.StrictHostKeyChecking)
	}
	if server.CheckHostIP != "yes" {
		t.Errorf("CheckHostIP = %v, want yes", server.CheckHostIP)
	}
	if server.FingerprintHash != "sha256" {
		t.Errorf("FingerprintHash = %v, want sha256", server.FingerprintHash)
	}
}

func TestServerWithPortForwarding(t *testing.T) {
	server := Server{
		Alias:          "forwarding-server",
		Host:           "example.com",
		User:           "user",
		Port:           22,
		LocalForward:   []string{"8080:localhost:80", "3306:localhost:3306"},
		RemoteForward:  []string{"9000:localhost:9000"},
		DynamicForward: []string{"1080"},
	}

	if len(server.LocalForward) != 2 {
		t.Errorf("LocalForward length = %v, want 2", len(server.LocalForward))
	}
	if len(server.RemoteForward) != 1 {
		t.Errorf("RemoteForward length = %v, want 1", len(server.RemoteForward))
	}
	if len(server.DynamicForward) != 1 {
		t.Errorf("DynamicForward length = %v, want 1", len(server.DynamicForward))
	}
}

func TestServerWithConnectionSettings(t *testing.T) {
	server := Server{
		Alias:               "connection-server",
		Host:                "example.com",
		User:                "user",
		Port:                22,
		ConnectTimeout:      "30",
		ConnectionAttempts:  "3",
		ServerAliveInterval: "60",
		ServerAliveCountMax: "3",
		TCPKeepAlive:        "yes",
		Compression:         "yes",
	}

	if server.ConnectTimeout != "30" {
		t.Errorf("ConnectTimeout = %v, want 30", server.ConnectTimeout)
	}
	if server.ConnectionAttempts != "3" {
		t.Errorf("ConnectionAttempts = %v, want 3", server.ConnectionAttempts)
	}
	if server.ServerAliveInterval != "60" {
		t.Errorf("ServerAliveInterval = %v, want 60", server.ServerAliveInterval)
	}
}

func TestServerWithAuthSettings(t *testing.T) {
	server := Server{
		Alias:                        "auth-server",
		Host:                         "example.com",
		User:                         "user",
		Port:                         22,
		PubkeyAuthentication:         "yes",
		PasswordAuthentication:       "no",
		KbdInteractiveAuthentication: "no",
		IdentitiesOnly:               "yes",
	}

	if server.PubkeyAuthentication != "yes" {
		t.Errorf("PubkeyAuthentication = %v, want yes", server.PubkeyAuthentication)
	}
	if server.PasswordAuthentication != "no" {
		t.Errorf("PasswordAuthentication = %v, want no", server.PasswordAuthentication)
	}
	if server.IdentitiesOnly != "yes" {
		t.Errorf("IdentitiesOnly = %v, want yes", server.IdentitiesOnly)
	}
}

func TestServerWithEnvironment(t *testing.T) {
	server := Server{
		Alias:   "env-server",
		Host:    "example.com",
		User:    "user",
		Port:    22,
		SendEnv: []string{"LANG", "LC_*"},
		SetEnv:  []string{"FOO=bar", "BAZ=qux"},
	}

	if len(server.SendEnv) != 2 {
		t.Errorf("SendEnv length = %v, want 2", len(server.SendEnv))
	}
	if len(server.SetEnv) != 2 {
		t.Errorf("SetEnv length = %v, want 2", len(server.SetEnv))
	}
}

func TestServerWithMultiplexing(t *testing.T) {
	server := Server{
		Alias:          "mux-server",
		Host:           "example.com",
		User:           "user",
		Port:           22,
		ControlMaster:  "auto",
		ControlPath:    "~/.ssh/controlmasters/%r@%h:%p",
		ControlPersist: "10m",
	}

	if server.ControlMaster != "auto" {
		t.Errorf("ControlMaster = %v, want auto", server.ControlMaster)
	}
	if server.ControlPath != "~/.ssh/controlmasters/%r@%h:%p" {
		t.Errorf("ControlPath = %v, want ~/.ssh/controlmasters/%%r@%%h:%%p", server.ControlPath)
	}
	if server.ControlPersist != "10m" {
		t.Errorf("ControlPersist = %v, want 10m", server.ControlPersist)
	}
}
