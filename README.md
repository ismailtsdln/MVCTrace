# MVCTrace

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://img.shields.io/badge/build-passing-green.svg)]()

A production-grade CLI tool for detecting ASP.NET MVC applications and inferring their versions. Designed for security reconnaissance, penetration testing, and web application analysis. MVCTrace helps identify MVC frameworks with high accuracy while minimizing false positives.

## Features

- 🔍 **Advanced Detection**: Multi-layered fingerprinting using HTTP headers, HTML analysis, route probing, error pages, and static files
- 📊 **Confidence Scoring**: Intelligent scoring system (0-100) with High/Medium/Low confidence levels
- 🎨 **Colored Output**: Beautiful, color-coded CLI output for better readability
- 📋 **Multiple Formats**: Human-readable output and machine-readable JSON
- ⚡ **Fast & Lightweight**: Optimized HTTP requests with configurable timeouts
- 🛡️ **Security-Focused**: Supports proxy configuration and TLS verification bypass for testing
- 🔧 **Extensible**: Modular architecture for easy addition of new detection methods

## Installation

### From Source
```bash
git clone https://github.com/ismailtsdln/mvctrace.git
cd mvctrace
go build -o mvctrace .
```

### Using Go Install
```bash
go install github.com/ismailtsdln/mvctrace@main
```

Make sure your Go bin directory is in your PATH.

## Usage

### Basic Usage
```bash
mvctrace https://example.com
```

### Advanced Options
```bash
# JSON output for scripting
mvctrace -json https://example.com

# Use proxy for testing
mvctrace -proxy http://127.0.0.1:8080 https://example.com

# Minimal output
mvctrace -silent https://example.com

# Custom timeout
mvctrace -timeout 30s https://example.com
```

### Command Line Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-json` | Output results in JSON format | false |
| `-timeout` | HTTP request timeout | 10s |
| `-proxy` | HTTP proxy URL | "" |
| `-silent` | Minimal output mode | false |

## Example Output

### MVC Detected (High Confidence)
```
Target: https://example.com
Framework: ASP.NET MVC
MVC Version: 5.2 (High Confidence)
Evidence:
  • X-AspNetMvc-Version header detected: 5.2
  • MVC validation attributes (data-val) found in HTML
  • Default MVC route /Home/Index responded with 200
  • MVC static file /Content/Site.css found
```

### Not MVC
```
Target: https://example.com
Framework: Not ASP.NET MVC
Confidence: Low
Evidence:
  • No MVC-specific headers detected
```

### JSON Output
```json
{
  "target": "https://example.com",
  "is_mvc": true,
  "version": "5.2",
  "confidence": 85,
  "evidence": [
    {
      "description": "X-AspNetMvc-Version header detected: 5.2",
      "confidence": 90
    }
  ]
}
```

## Detection Methods

MVCTrace uses multiple detection techniques to ensure accuracy:

1. **HTTP Headers**: Analyzes `X-AspNetMvc-Version`, `X-AspNet-Version`, `X-Powered-By`
2. **HTML Fingerprinting**: Searches for MVC-specific attributes like `data-val="true"`, `__MVCFormValidation`
3. **Route Probing**: Tests common MVC routes (`/Home/Index`, `/Account/Login`)
4. **Error Analysis**: Examines 404 pages for MVC error patterns
5. **Static Files**: Checks for MVC bundle files and CSS/JS assets

## Confidence Levels

- **High (≥70)**: Confirmed ASP.NET MVC with strong evidence
- **Medium (40-69)**: Likely MVC with moderate evidence
- **Low (<40)**: Unlikely MVC or insufficient evidence

## Development

### Prerequisites
- Go 1.21 or later

### Building
```bash
go build -o mvctrace .
```

### Testing
```bash
go test ./...
```

### Project Structure
```
mvctrace/
├── main.go                 # CLI entry point
├── internal/
│   ├── detector/           # Detection logic
│   ├── httpclient/         # HTTP client wrapper
│   └── result/             # Result structures
├── go.mod
├── README.md
└── VERSION
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Version Information

Current version: 1.0.0 (see VERSION file)

Version information follows semantic versioning.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Disclaimer

MVCTrace is designed for security research and penetration testing purposes. Users are responsible for complying with applicable laws and regulations when using this tool. The authors are not responsible for any misuse or damage caused by this software.

## Author

**İsmail Taşdelen** - [GitHub](https://github.com/ismailtsdln)