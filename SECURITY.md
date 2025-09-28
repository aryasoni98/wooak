# Security Policy

## 🛡️ Supported Versions

We provide security updates for the following versions of Wooak:

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |
| < 0.1   | :x:                |

## 🚨 Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security vulnerability in Wooak, please follow these steps:

### 1. **DO NOT** create a public GitHub issue

Security vulnerabilities should be reported privately to prevent exploitation.

### 2. Report the vulnerability

Please report security vulnerabilities by:

- **Email**: [INSERT SECURITY EMAIL] (if available)
- **GitHub Security Advisory**: Use GitHub's private vulnerability reporting feature
- **Direct Contact**: Contact the maintainer directly through GitHub

### 3. Include the following information

When reporting a vulnerability, please include:

- **Description**: A clear description of the vulnerability
- **Steps to Reproduce**: Detailed steps to reproduce the issue
- **Impact**: Potential impact and severity assessment
- **Environment**: OS, version, and configuration details
- **Proof of Concept**: If applicable, include a minimal proof of concept
- **Suggested Fix**: If you have ideas for fixing the issue

### 4. Response Timeline

We will respond to security reports within:

- **Initial Response**: 24-48 hours
- **Status Update**: Within 1 week
- **Resolution**: As quickly as possible, typically within 30 days

### 5. Disclosure Process

- We will work with you to understand and reproduce the issue
- We will develop and test a fix
- We will coordinate the release of the fix
- We will credit you for the discovery (unless you prefer to remain anonymous)

## 🔒 Security Features

Wooak implements several security features:

### SSH Key Validation
- Validates key types and sizes
- Checks for weak or deprecated algorithms
- Provides security recommendations

### Audit Logging
- Tracks all security-relevant events
- Configurable retention policies
- Structured logging for analysis

### Host Security
- Allow/block list management
- Connection validation
- Security policy enforcement

### Configuration Safety
- Non-destructive configuration edits
- Automatic backups before changes
- Atomic file operations

## 🛠️ Security Best Practices

### For Users

1. **Keep Wooak Updated**: Always use the latest version
2. **Use Strong SSH Keys**: Prefer Ed25519 over RSA
3. **Regular Key Rotation**: Rotate SSH keys periodically
4. **Monitor Audit Logs**: Review security events regularly
5. **Secure Configuration**: Protect your SSH config file

### For Developers

1. **Input Validation**: Always validate user input
2. **Secure Defaults**: Use secure default configurations
3. **Error Handling**: Don't expose sensitive information in errors
4. **Dependency Management**: Keep dependencies updated
5. **Code Review**: Security-focused code reviews

## 🔍 Security Scanning

We use automated security scanning tools:

- **Dependabot**: Monitors for vulnerable dependencies
- **CodeQL**: Static analysis for security vulnerabilities
- **Manual Review**: Regular security-focused code reviews

## 📋 Security Checklist

Before submitting code, ensure:

- [ ] No hardcoded secrets or credentials
- [ ] Input validation for all user inputs
- [ ] Proper error handling without information disclosure
- [ ] Secure file operations
- [ ] No SQL injection or command injection vulnerabilities
- [ ] Proper authentication and authorization checks
- [ ] Secure communication protocols
- [ ] Regular dependency updates

## 🏆 Security Acknowledgments

We appreciate security researchers who help improve Wooak's security. Contributors will be acknowledged in:

- Security advisories
- Release notes
- Project documentation
- Hall of Fame (if desired)

## 📞 Contact Information

**Security Contact**: [Arya Soni](https://github.com/aryasoni98)

**GitHub Security**: Use GitHub's private vulnerability reporting

---

*This security policy is effective as of the date of its adoption and will be reviewed and updated as needed to ensure it continues to serve our security needs.*
