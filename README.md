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

### Quick Install (Recommended)
```bash
go install github.com/ismailtsdln/mvctrace@v1.0.0
```

This will install MVCTrace to your `$GOPATH/bin` directory. Make sure it's in your PATH.

### From Source
```bash
git clone https://github.com/ismailtsdln/mvctrace.git
cd mvctrace
go build -o mvctrace .
sudo mv mvctrace /usr/local/bin/  # Optional: Move to system bin
```

### Verify Installation
```bash
mvctrace -h
```

## Usage

### Basic Usage
```bash
mvctrace https://example.com
```

### Command Line Flags

| Flag | Type | Description | Default |
|------|------|-------------|---------|
| `-json` | bool | Output results in JSON format (machine-readable) | false |
| `-timeout` | duration | HTTP request timeout for each probe | 10s |
| `-proxy` | string | HTTP proxy URL (e.g., `http://127.0.0.1:8080`) | "" |
| `-silent` | bool | Minimal output (only MVC detection result) | false |

### Examples
```bash
# Basic scan
mvctrace https://example.com

# Export to JSON for further processing
mvctrace -json https://example.com > results.json

# Use through proxy
mvctrace -proxy http://127.0.0.1:8080 https://example.com

# Extended timeout for slow targets
mvctrace -timeout 30s https://example.com

# Minimal output for scripting
mvctrace -silent https://example.com && echo "MVC Found" || echo "Not MVC"

# Scan multiple targets
for url in $(cat targets.txt); do
  mvctrace -json "$url" >> results.jsonl
done
```

## Example Output

### MVC Detected (High Confidence)
```
Target: https://example.com
Framework: ASP.NET MVC
MVC Version: 5.2 (High Confidence)
Version Source: HTTP Header: X-AspNetMvc-Version
Evidence:
  • MVC Version 5.2 detected
    Source: HTTP Header: X-AspNetMvc-Version
    Value: 5.2
  • MVC validation attributes detected
    Source: HTML Body: data-val attribute
    Value: data-val="true"
  • MVC default route is accessible
    Source: HTTP Route: /Home/Index
    Value: HTTP 200
  • MVC static file structure detected
    Source: Static File: /Content/Site.css
    Value: HTTP 200
```

### Not MVC
```
Target: http://httpbin.org
Framework: Not ASP.NET MVC
Confidence: Low
```

### JSON Output with Source Tracking
```json
{
  "target": "https://example.com",
  "is_mvc": true,
  "version": "5.2",
  "version_source": "HTTP Header: X-AspNetMvc-Version",
  "confidence": 85,
  "evidence": [
    {
      "description": "MVC Version 5.2 detected",
      "source": "HTTP Header: X-AspNetMvc-Version",
      "value": "5.2",
      "confidence": 90
    },
    {
      "description": "MVC validation attributes detected",
      "source": "HTML Body: data-val attribute",
      "value": "data-val=\"true\"",
      "confidence": 60
    }
  ]
}
```

## Detection Methods

MVCTrace uses multiple detection techniques to ensure accuracy and provide detailed source information:

### 1. **HTTP Headers** (Highest Confidence)
- `X-AspNetMvc-Version`: Explicitly declares MVC version (90% confidence)
- `X-AspNet-Version`: Indicates .NET framework version (20% confidence)
- `X-Powered-By`: Generic ASP.NET indicator (10% confidence)

### 2. **HTML Body Fingerprinting**
- `data-val="true"`: MVC unobtrusive validation attribute (60% confidence)
- `__MVCFormValidation`: ASP.NET MVC form validation script (70% confidence)
- `jquery.validate.unobtrusive`: MVC validation library (50% confidence)
- `System.Web.Mvc`: Namespace reference in scripts (80% confidence)

### 3. **Default Route Probing**
- Tests standard MVC routes: `/Home/Index`, `/Account/Login`, `/Home/About`
- Success indicates MVC application structure (40% confidence)

### 4. **Error Page Analysis**
- Examines 404 error pages at `/nonexistent-path-12345`
- Looks for MVC-specific error messages and stack traces (30-50% confidence)

### 5. **Static File Detection**
- Checks for MVC bundle structure: `/Content/Site.css`, `/Scripts/jquery-*.js`, `/bundles/jquery`
- Accessible static files suggest MVC project layout (20% confidence)

## How It Works

MVCTrace follows a multi-stage detection pipeline to identify ASP.NET MVC applications:

1. **Initial Request**: Sends HTTP GET request to target URL with custom User-Agent
2. **Parallel Detection**: Runs all 5 detection methods simultaneously for speed
3. **Evidence Collection**: Each detection method returns source-tracked evidence
4. **Confidence Aggregation**: Combines individual confidence scores (capped at 100)
5. **Version Extraction**: Identifies MVC version and its source (if available)
6. **Result Assembly**: Aggregates findings with detailed source information

### Key Improvements

- ✅ **Source Tracking**: Every piece of evidence shows exactly where it was found
- ✅ **Version Attribution**: MVC version includes source header/location
- ✅ **No False Positives**: Requires multiple corroborating pieces of evidence
- ✅ **Detailed Reporting**: JSON and human-readable formats with full context
- ✅ **Security-Focused**: Designed for penetration testing and reconnaissance

## Development

### Prerequisites
- Go 1.21 or later

### Quick Start
```bash
git clone https://github.com/ismailtsdln/mvctrace.git
cd mvctrace

# Build locally
go build -o mvctrace .

# Run tests
go test ./...

# Run with local binary
./mvctrace https://example.com
```

### Project Structure
```
mvctrace/
├── main.go                      # CLI interface with colored output
├── internal/
│   ├── detector/
│   │   ├── detector.go          # Core detection logic (5 methods)
│   │   └── detector_test.go     # Unit tests
│   ├── httpclient/
│   │   └── client.go            # HTTP client with proxy support
│   └── result/
│       └── result.go            # Result structures with source tracking
├── go.mod                       # Go module definition
├── VERSION                      # Version file (1.0.0)
├── LICENSE                      # MIT License
├── README.md                    # This file
└── .gitignore                   # Git ignore rules
```

### Code Quality
- Passes `go vet` and `go fmt`
- Comprehensive unit tests
- Modular design for extensibility
- Clean error handling

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Version Information

**Current Version:** 1.0.0

Version information is maintained in the `VERSION` file and follows [semantic versioning](https://semver.org/).

### Release History
- **v1.0.0** (Jan 5, 2026): Initial release with full detection pipeline, source tracking, and go install support

## FAQ

**Q: Why does MVCTrace detect my .NET application as "Not MVC"?**  
A: MVCTrace specifically targets ASP.NET MVC (System.Web.MVC). Other frameworks like ASP.NET Core, Razor Pages, or WebForms won't be detected. Check your detection evidence to see what was found.

**Q: Can I use MVCTrace with IPv6 addresses?**  
A: Yes, MVCTrace works with IPv6. Use format: `mvctrace [::1]:8080` or provide the full URL.

**Q: Is MVCTrace safe to use on production systems?**  
A: MVCTrace only sends GET requests to publicly accessible endpoints. It doesn't exploit or modify anything. However, always obtain proper authorization before scanning.

**Q: How accurate is MVCTrace?**  
A: Accuracy depends on the target. Modern ASP.NET MVC applications with version headers are detected with >90% confidence. Older versions or heavily customized applications may require multiple evidence points.

**Q: Can I contribute detection methods?**  
A: Yes! We welcome contributions. See the Contributing section below.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

```
MIT License

Copyright (c) 2026 İsmail Taşdelen

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software")...
```

## Disclaimer

**Important:** MVCTrace is designed for:
- ✅ Authorized security research and penetration testing
- ✅ Web application analysis with proper permissions
- ✅ Learning about ASP.NET MVC detection techniques

**Not for:**
- ❌ Unauthorized scanning of systems you don't own
- ❌ Circumventing security systems
- ❌ Any illegal activities

Users are **fully responsible** for complying with applicable laws and regulations when using this tool. The authors accept no liability for misuse or damage caused by this software.

## Support & Contributing

### Getting Help
- 🐛 Found a bug? [Open an issue](https://github.com/ismailtsdln/mvctrace/issues)
- 💡 Have an idea? [Start a discussion](https://github.com/ismailtsdln/mvctrace/discussions)
- 📖 Need documentation? Check the [wiki](https://github.com/ismailtsdln/mvctrace/wiki)

### Contributing
We welcome contributions! Here's how:

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/your-feature`)
3. **Make** your changes with clear commit messages
4. **Test** your changes (`go test ./...`)
5. **Push** to your fork (`git push origin feature/your-feature`)
6. **Open** a Pull Request with a clear description

### Development Guidelines
- Follow Go conventions and best practices
- Add tests for new features
- Update documentation as needed
- Keep commits atomic and well-documented

## Author

**İsmail Taşdelen**
- GitHub: [@ismailtsdln](https://github.com/ismailtsdln)
- Email: contact via GitHub
- Twitter: [@ismailtsdln](https://twitter.com/ismailtsdln)

## Acknowledgments

- Inspired by security research community best practices
- Built with Go's excellent standard library
- Special thanks to all contributors and testers