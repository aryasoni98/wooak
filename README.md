<div align="center">

# Wooak - Intelligent SSH Management

**A modern, AI-powered terminal-based SSH manager with enterprise-grade security**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue?style=for-the-badge)](LICENSE)
[![Release](https://img.shields.io/badge/Release-v0.1.0-green?style=for-the-badge)](https://github.com/aryasoni98/wooak/releases)

</div>

---

## 🚀 Overview

Wooak is a next-generation terminal-based SSH manager that combines the power of modern AI with enterprise-grade security features. Built for developers and system administrators who manage multiple servers, Wooak provides an intuitive interface for SSH server management with intelligent recommendations and comprehensive security analysis.

### 🎯 Key Capabilities

- **🤖 AI-Powered Assistant**: Get intelligent recommendations for SSH configurations
- **🔐 Advanced Security**: Comprehensive security analysis and audit logging
- **⚡ Lightning Fast**: Optimized for speed with intelligent caching
- **🎨 Beautiful UI**: Clean, keyboard-driven interface inspired by k9s and lazydocker
- **🔧 Highly Configurable**: Extensive customization options for every use case

---

## 📋 Table of Contents

- [Quick Start](#-quick-start)
- [Features](#-features)
- [Architecture](#-architecture)
- [Installation](#-installation)
- [Usage Guide](#-usage-guide)
- [Development](#-development)
- [Security](#-security)
- [Contributing](#-contributing)
- [Support](#-support)

---

## ⚡ Quick Start

### 1. Install Wooak

```bash
# Option 1: Homebrew (macOS)
brew install aryasoni98/homebrew-tap/wooak

# Option 2: Download Binary
curl -fsSL https://raw.githubusercontent.com/aryasoni98/wooak/main/install.sh | bash

# Option 3: Build from Source
git clone https://github.com/aryasoni98/wooak.git
cd wooak && make build
```

### 2. Setup AI Assistant (Optional)

```bash
# Install Ollama and AI models
make ai-setup

# Check AI status
make ai-status
```

### 3. Launch Wooak

```bash
wooak
```

### 4. Essential Commands

| Key | Action |
|-----|--------|
| `i` | Open AI Assistant |
| `z` | Open Security Panel |
| `a` | Add new server |
| `Enter` | Connect to server |
| `q` | Quit |

---

## ✨ Features

### 🤖 AI-Powered Intelligence

- **Natural Language Search**: "Find all production web servers"
- **Security Analysis**: AI-powered vulnerability detection
- **Configuration Optimization**: Intelligent performance recommendations
- **Smart Suggestions**: Personalized recommendations based on usage patterns

### 🔐 Enterprise Security

- **SSH Key Validation**: Comprehensive key security analysis
- **Audit Logging**: Complete security event tracking
- **Host Security**: Allow/block list management
- **Policy Enforcement**: Configurable security policies

### 🖥️ Server Management

- **Visual Server List**: Clean, organized server display
- **Fuzzy Search**: Quick server discovery
- **Tagging System**: Organize servers by environment/role
- **Connection Multiplexing**: Faster subsequent connections
- **Port Forwarding**: Local, remote, and dynamic forwarding

### ⚙️ Advanced Configuration

- **Tabbed Interface**: Organized configuration options
- **Auto-completion**: Smart SSH key detection
- **Backup System**: Automatic configuration backups
- **Non-destructive Edits**: Preserves existing formatting

---

## 🏗️ Architecture

### System Architecture

```mermaid

graph TB
    subgraph "User Interface Layer"
        UI[TUI Interface]
        AI_UI[AI Assistant Panel]
        SEC_UI[Security Panel]
    end
    
    subgraph "Application Layer"
        HANDLERS[Event Handlers]
        SERVICES[Core Services]
    end
    
    subgraph "Business Logic Layer"
        SERVER_SVC[Server Service]
        AI_SVC[AI Service]
        SEC_SVC[Security Service]
    end
    
    subgraph "Data Layer"
        REPO[SSH Config Repository]
        CACHE[AI Cache]
        AUDIT[Audit Logger]
    end
    
    subgraph "External Services"
        OLLAMA[Ollama AI]
        SSH[OpenSSH Binary]
        FILES[File System]
    end
    
    UI --> HANDLERS
    AI_UI --> HANDLERS
    SEC_UI --> HANDLERS
    
    HANDLERS --> SERVICES
    SERVICES --> SERVER_SVC
    SERVICES --> AI_SVC
    SERVICES --> SEC_SVC
    
    SERVER_SVC --> REPO
    AI_SVC --> CACHE
    AI_SVC --> OLLAMA
    SEC_SVC --> AUDIT
    
    REPO --> FILES
    AUDIT --> FILES
    
    classDef uiLayer fill:#4A90E2,stroke:#2E5BBA,stroke-width:3px,color:#fff
    classDef appLayer fill:#7ED321,stroke:#5BA517,stroke-width:3px,color:#fff
    classDef businessLayer fill:#F5A623,stroke:#D68910,stroke-width:3px,color:#fff
    classDef dataLayer fill:#BD10E0,stroke:#9013FE,stroke-width:3px,color:#fff
    classDef externalLayer fill:#D0021B,stroke:#A00000,stroke-width:3px,color:#fff
    
    class UI,AI_UI,SEC_UI uiLayer
    class HANDLERS,SERVICES appLayer
    class SERVER_SVC,AI_SVC,SEC_SVC businessLayer
    class REPO,CACHE,AUDIT dataLayer
    class OLLAMA,SSH,FILES externalLayer
```

### Data Flow

```mermaid

sequenceDiagram
    participant U as User
    participant UI as TUI Interface
    participant H as Handlers
    participant S as Services
    participant R as Repository
    participant AI as AI Service
    participant O as Ollama
    
    U->>UI: Press 'i' (AI Assistant)
    UI->>H: handleAIPanel()
    H->>AI: Initialize AI Service
    AI->>O: Check Connection
    O-->>AI: Connection Status
    AI-->>H: AI Ready
    H-->>UI: Show AI Panel
    
    U->>UI: Ask Question
    UI->>H: processAIQuery()
    H->>AI: Generate Response
    AI->>O: Send Prompt
    O-->>AI: AI Response
    AI-->>H: Processed Response
    H-->>UI: Display Result
    
```

---

## 📦 Installation

### Prerequisites

- **Go 1.21+** (for building from source)
- **OpenSSH** (for SSH connections)
- **Ollama** (optional, for AI features)

### Installation Methods

#### Option 1: Homebrew (macOS)

```bash
brew install aryasoni98/homebrew-tap/wooak
```

#### Option 2: Download Binary

```bash
# Auto-install script
curl -fsSL https://raw.githubusercontent.com/aryasoni98/wooak/main/install.sh | bash

# Manual download
LATEST_TAG=$(curl -fsSL https://api.github.com/repos/aryasoni98/wooak/releases/latest | jq -r .tag_name)
curl -LJO "https://github.com/aryasoni98/wooak/releases/download/${LATEST_TAG}/wooak_$(uname)_$(uname -m).tar.gz"
tar -xzf wooak_$(uname)_$(uname -m).tar.gz
sudo mv wooak /usr/local/bin/
```

#### Option 3: Build from Source

```bash
git clone https://github.com/aryasoni98/wooak.git
cd wooak
make dev-setup  # Setup development environment
make build      # Build the binary
```

---

## 📖 Usage Guide

### Basic Workflow

```mermaid

flowchart TD
    A[Launch Wooak] --> B[View Server List]
    B --> C{Action Needed?}
    C -->|Search| D[Press '/' - Fuzzy Search]
    C -->|Add Server| E[Press 'a' - Add Server]
    C -->|Connect| F[Press Enter - SSH Connect]
    C -->|AI Help| G[Press 'i' - AI Assistant]
    C -->|Security| H[Press 'z' - Security Panel]
    
    D --> B
    E --> I[Configure Server]
    I --> B
    F --> J[SSH Session]
    J --> B
    G --> K[AI Recommendations]
    K --> B
    H --> L[Security Analysis]
    L --> B
    
    classDef startNode fill:#4CAF50,stroke:#2E7D32,stroke-width:3px,color:#fff
    classDef actionNode fill:#2196F3,stroke:#1565C0,stroke-width:3px,color:#fff
    classDef decisionNode fill:#FF9800,stroke:#E65100,stroke-width:3px,color:#fff
    classDef processNode fill:#9C27B0,stroke:#6A1B9A,stroke-width:3px,color:#fff
    classDef aiNode fill:#00BCD4,stroke:#006064,stroke-width:3px,color:#fff
    classDef securityNode fill:#F44336,stroke:#C62828,stroke-width:3px,color:#fff
    
    class A startNode
    class B,C actionNode
    class D,E,F decisionNode
    class I,J processNode
    class G,K aiNode
    class H,L securityNode

```

### Key Bindings

#### Main Interface

| Key | Action | Description |
|-----|--------|-------------|
| `/` | Search | Toggle fuzzy search bar |
| `↑↓` / `jk` | Navigate | Move through server list |
| `Enter` | Connect | SSH into selected server |
| `a` | Add | Add new server |
| `e` | Edit | Edit selected server |
| `d` | Delete | Delete selected server |
| `p` | Pin | Pin/unpin server |
| `t` | Tags | Edit server tags |
| `s` | Sort | Toggle sort field |
| `S` | Reverse | Reverse sort order |
| `c` | Copy | Copy SSH command |
| `g` | Ping | Ping selected server |
| `r` | Refresh | Refresh server data |
| `i` | AI | Open AI Assistant |
| `z` | Security | Open Security Panel |
| `q` | Quit | Exit application |

#### AI Assistant Panel

| Key | Action | Description |
|-----|--------|-------------|
| `Enter` | Send | Send message to AI |
| `Esc` | Close | Close AI panel |
| `Tab` | Switch | Switch between panels |

#### Security Panel

| Key | Action | Description |
|-----|--------|-------------|
| `Tab` | Navigate | Move between fields |
| `Enter` | Save | Save configuration |
| `Esc` | Close | Close security panel |

### Configuration

Wooak automatically reads from your `~/.ssh/config` file. No additional configuration is required, but you can customize:

- **AI Settings**: Configure AI models and providers
- **Security Policies**: Set security validation rules
- **UI Preferences**: Customize display options

---

## 🛠️ Development

### Development Workflow

```mermaid

graph LR
    A[Clone Repo] --> B[Setup Environment]
    B --> C[Make Changes]
    C --> D[Run Tests]
    D --> E{Quality Checks}
    E -->|Pass| F[Build]
    E -->|Fail| C
    F --> G[Test Features]
    G --> H[Submit PR]
    
    classDef startNode fill:#4CAF50,stroke:#2E7D32,stroke-width:3px,color:#fff
    classDef processNode fill:#2196F3,stroke:#1565C0,stroke-width:3px,color:#fff
    classDef decisionNode fill:#FF9800,stroke:#E65100,stroke-width:3px,color:#fff
    classDef successNode fill:#8BC34A,stroke:#558B2F,stroke-width:3px,color:#fff
    classDef failNode fill:#F44336,stroke:#C62828,stroke-width:3px,color:#fff
    
    class A startNode
    class B,C,D,G,H processNode
    class E decisionNode
    class F successNode

```

### Available Make Targets

```bash
# Development Setup
make dev-setup      # Setup complete development environment
make tools          # Install development tools
make deps           # Download dependencies

# Building
make build          # Build binary with quality checks
make build-all      # Build for all platforms
make run            # Run from source
make demo           # Run demo with sample data

# Quality Assurance
make quality        # Run all quality checks
make test           # Run unit tests
make coverage       # Generate coverage report
make lint           # Run linter
make security-scan  # Run security checks

# AI Features
make ai-setup       # Setup AI dependencies
make ai-status      # Check AI service status
make ai-models      # List available AI models
make ai-test        # Test AI functionality

# Security Features
make security-test  # Test security features
make security-scan  # Run security analysis

# Maintenance
make clean          # Clean build artifacts
make update-deps    # Update dependencies
make help           # Show all available targets
```

### Project Structure

```
wooak/
├── cmd/                    # Application entry point
│   └── main.go
├── internal/               # Private application code
│   ├── adapters/          # External interface adapters
│   │   ├── data/          # Data layer adapters
│   │   └── ui/            # User interface adapters
│   │       ├── ai/        # AI UI components
│   │       └── security/  # Security UI components
│   ├── core/              # Business logic
│   │   ├── domain/        # Domain models
│   │   │   ├── ai/        # AI domain models
│   │   │   └── security/  # Security domain models
│   │   ├── ports/         # Interface definitions
│   │   └── services/      # Business services
│   │       ├── ai/        # AI services
│   │       └── security/  # Security services
│   └── logger/            # Logging utilities
├── docs/                  # Documentation and screenshots
├── makefile              # Build automation
├── .goreleaser.yaml      # Release configuration
└── README.md             # This file
```

### Adding New Features

1. **Create Domain Models** (if needed)
2. **Implement Services** in `internal/core/services/`
3. **Add UI Components** in `internal/adapters/ui/`
4. **Update Handlers** in `internal/adapters/ui/handlers.go`
5. **Add Tests** and ensure quality checks pass
6. **Update Documentation**

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run specific test packages
go test ./internal/core/services/...

# Run benchmarks
make benchmark
```

---

## 🔐 Security

### Security Features

Wooak implements multiple layers of security:

#### 1. SSH Key Validation
- Validates key types and sizes
- Checks for weak or deprecated algorithms
- Provides security recommendations

#### 2. Audit Logging
- Tracks all security-relevant events
- Configurable retention policies
- Structured logging for analysis

#### 3. Host Security
- Allow/block list management
- Connection validation
- Security policy enforcement

#### 4. Configuration Safety
- Non-destructive configuration edits
- Automatic backups before changes
- Atomic file operations

### Security Workflow

```mermaid

graph TD
    A[SSH Connection Request] --> B[Security Validation]
    B --> C{Key Valid?}
    C -->|No| D[Block Connection]
    C -->|Yes| E{Host Allowed?}
    E -->|No| D
    E -->|Yes| F[Log Event]
    F --> G[Allow Connection]
    
    D --> H[Log Security Event]
    H --> I[Update Audit Log]
    G --> I
    
    classDef requestNode fill:#2196F3,stroke:#1565C0,stroke-width:3px,color:#fff
    classDef validationNode fill:#FF9800,stroke:#E65100,stroke-width:3px,color:#fff
    classDef decisionNode fill:#9C27B0,stroke:#6A1B9A,stroke-width:3px,color:#fff
    classDef blockNode fill:#F44336,stroke:#C62828,stroke-width:3px,color:#fff
    classDef allowNode fill:#4CAF50,stroke:#2E7D32,stroke-width:3px,color:#fff
    classDef logNode fill:#00BCD4,stroke:#006064,stroke-width:3px,color:#fff
    
    class A requestNode
    class B validationNode
    class C,E decisionNode
    class D blockNode
    class G allowNode
    class F,H,I logNode

```

### Security Best Practices

1. **Regular Key Rotation**: Use AI recommendations for key management
2. **Monitor Audit Logs**: Review security events regularly
3. **Update Security Policies**: Keep policies current with best practices
4. **Use Strong Keys**: Prefer Ed25519 over RSA when possible
5. **Enable Host Verification**: Always verify host keys

---

## 🤝 Contributing

We welcome contributions! Please follow these guidelines:

### Development Process


```mermaid

graph LR
    A[Fork Repository] --> B[Create Feature Branch]
    B --> C[Make Changes]
    C --> D[Run Tests]
    D --> E[Submit PR]
    E --> F[Code Review]
    F --> G[Merge]
    
    classDef startNode fill:#4CAF50,stroke:#2E7D32,stroke-width:3px,color:#fff
    classDef processNode fill:#2196F3,stroke:#1565C0,stroke-width:3px,color:#fff
    classDef testNode fill:#FF9800,stroke:#E65100,stroke-width:3px,color:#fff
    classDef reviewNode fill:#9C27B0,stroke:#6A1B9A,stroke-width:3px,color:#fff
    classDef successNode fill:#8BC34A,stroke:#558B2F,stroke-width:3px,color:#fff
    
    class A startNode
    class B,C processNode
    class D testNode
    class E,F reviewNode
    class G successNode

```

### Pull Request Guidelines

1. **Use Semantic PR Titles**:
   - `feat(scope): description` - New features
   - `fix(scope): description` - Bug fixes
   - `improve(scope): description` - Improvements
   - `docs: description` - Documentation

2. **Ensure Quality**:
   ```bash
   make quality  # Run all quality checks
   make test     # Run tests
   ```

3. **Update Documentation**:
   - Update README if needed
   - Add/update code comments
   - Update help text

### Available Scopes

- `ui` - User interface changes
- `ai` - AI-related features
- `security` - Security features
- `config` - Configuration handling
- `parser` - SSH config parsing

### Examples

```bash
feat(ai): add natural language search
fix(security): resolve key validation edge case
improve(ui): enhance server list performance
docs: update installation instructions
```

---

## ⭐ Support

If you find Wooak useful, please consider:

- ⭐ **Starring** the repository
- 🐛 **Reporting** bugs via issues
- 💡 **Suggesting** new features
- 🤝 **Contributing** code improvements

### Community

- 📧 **Issues**: [GitHub Issues](https://github.com/aryasoni98/wooak/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/aryasoni98/wooak/discussions)

---

## 🙏 Acknowledgments

- Built with [tview](https://github.com/rivo/tview) and [tcell](https://github.com/gdamore/tcell)
- Inspired by [k9s](https://github.com/derailed/k9s) and [lazydocker](https://github.com/jesseduffield/lazydocker)
- AI powered by [Ollama](https://ollama.ai/)

---

<div align="center">

**Made with ❤️ for the developer community**

[⬆ Back to Top](#-overview)

</div>