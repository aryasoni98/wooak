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
	"sync"
	"time"

	"github.com/aryasoni98/wooak/internal/core/domain"
)

// ServerPool provides object pooling for Server instances to reduce memory allocations
type ServerPool struct {
	pool sync.Pool
}

// NewServerPool creates a new server object pool
func NewServerPool() *ServerPool {
	return &ServerPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &domain.Server{}
			},
		},
	}
}

// Get retrieves a Server instance from the pool
func (p *ServerPool) Get() *domain.Server {
	server := p.pool.Get().(*domain.Server)
	// Reset the server to its zero state
	p.resetServer(server)
	return server
}

// Put returns a Server instance to the pool
func (p *ServerPool) Put(server *domain.Server) {
	if server != nil {
		p.pool.Put(server)
	}
}

// resetServer resets a server instance to its zero state
func (p *ServerPool) resetServer(server *domain.Server) {
	// Reset all fields to their zero values
	server.Aliases = nil
	server.Host = ""
	server.User = ""
	server.Port = 0
	server.IdentityFiles = nil
	server.Tags = nil
	server.PinnedAt = time.Time{}
	server.LastSeen = time.Time{}
	server.SSHCount = 0

	// Reset all advanced fields
	server.ProxyJump = ""
	server.ProxyCommand = ""
	server.RemoteCommand = ""
	server.RequestTTY = ""
	server.SessionType = ""
	server.ConnectTimeout = ""
	server.ConnectionAttempts = ""
	server.BindAddress = ""
	server.BindInterface = ""
	server.AddressFamily = ""
	server.ExitOnForwardFailure = ""
	server.IPQoS = ""
	server.CanonicalizeHostname = ""
	server.CanonicalDomains = ""
	server.CanonicalizeFallbackLocal = ""
	server.CanonicalizeMaxDots = ""
	server.CanonicalizePermittedCNAMEs = ""
	server.ServerAliveInterval = ""
	server.ServerAliveCountMax = ""
	server.Compression = ""
	server.TCPKeepAlive = ""
	server.BatchMode = ""
	server.ControlMaster = ""
	server.ControlPath = ""
	server.ControlPersist = ""
	server.PubkeyAuthentication = ""
	server.PubkeyAcceptedAlgorithms = ""
	server.HostbasedAcceptedAlgorithms = ""
	server.PasswordAuthentication = ""
	server.PreferredAuthentications = ""
	server.IdentitiesOnly = ""
	server.AddKeysToAgent = ""
	server.IdentityAgent = ""
	server.KbdInteractiveAuthentication = ""
	server.NumberOfPasswordPrompts = ""
	server.ForwardAgent = ""
	server.ForwardX11 = ""
	server.ForwardX11Trusted = ""
	server.LocalForward = nil
	server.RemoteForward = nil
	server.DynamicForward = nil
	server.ClearAllForwardings = ""
	server.GatewayPorts = ""
	server.StrictHostKeyChecking = ""
	server.CheckHostIP = ""
	server.FingerprintHash = ""
	server.UserKnownHostsFile = ""
	server.HostKeyAlgorithms = ""
	server.Ciphers = ""
	server.MACs = ""
	server.KexAlgorithms = ""
	server.VerifyHostKeyDNS = ""
	server.UpdateHostKeys = ""
	server.HashKnownHosts = ""
	server.VisualHostKey = ""
	server.LocalCommand = ""
	server.PermitLocalCommand = ""
	server.EscapeChar = ""
	server.SendEnv = nil
	server.SetEnv = nil
	server.LogLevel = ""
}

// Global server pool instance
var globalServerPool = NewServerPool()

// GetServerFromPool retrieves a Server instance from the global pool
func GetServerFromPool() *domain.Server {
	return globalServerPool.Get()
}

// PutServerToPool returns a Server instance to the global pool
func PutServerToPool(server *domain.Server) {
	globalServerPool.Put(server)
}
