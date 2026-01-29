# AGENTS.md

## Build/Lint/Test Commands

### Default Commands
- `go mod tidy` - Download dependencies
- `go build` - Build project
- `go vet` - Lint code (Go vet)
- `go test` - Run all tests
- `go test -v` - Run tests with verbose output
- `go test -cover` - Run tests with coverage reporting

### Running a Single Test
To run a specific test, use:
```bash
go test -run TestFunctionName
```

## Code Style Guidelines

### General
- Use 2 spaces for indentation
- Keep lines under 120 characters (Go standard)
- Use semantic line breaks
- Use camelCase for variables and functions (Go convention)
- Use PascalCase for exported functions and types (Go convention)
- Use snake_case for package names (Go convention)
- Use kebab-case for filenames (Go convention)

### Imports
- Use standard Go import syntax (relative imports for local packages)
- Group imports by standard library, external libraries, and local packages
- Use import aliases when needed
- Import all required packages (no unused imports)

### Types
- Use Go's built-in types for standard operations
- Define custom types with clear naming conventions
- Use interfaces for abstraction and decoupling
- Use structs for data containers
- Use pointers where necessary for performance

### Formatting
- Use gofmt for formatting (Go standard)
- Use golint for linting (optional) 
- Use go vet for static analysis
- Use consistent naming conventions across the codebase
- Use consistent spacing around operators
- Use consistent braces for control structures

### Error Handling
- Use Go's error handling patterns with explicit error checks
- Don't ignore errors
- Use error wrapping with fmt.Errorf for better error messages
- Use custom error types when necessary
- Use panic only for unrecoverable errors

## Folder Structure

venue/
+-- cmd/                          # Application entry points
│   +-- server/
│   │   +-- main.go
|   +-- rand_inputs/  # Tool to generate random inputs to VENUE
|   +-- venue_cli/
|       +-- main.go
+-- internal/  # Core application logic
|   +-- math/
+-- codes/  # Error code definitions
+-- docs/  # Documentation
+-- ping/
+-- router/
+-- touchosc/  # Server code for hosting a TouchOSC interface
+-- venuelib/  # Core library functions and utilities
+-- venue/  # Client code for interacting with a VENUE system
+-- vnc/  # Code for connecting to VENUE VNC server
+-- go.mod
+-- go.sum
+-- README.md

### Test Structure
- All tests use standard Go `*_test.go` naming
- Tests should be in the same package as the code being tested
- Use table-driven tests for multiple test cases
- Use go test flags for coverage and verbose output

## Cursor Rules

### .cursor/rules/
No Cursor rules file found. Create a `.cursor/rules/` directory and add your custom rules there. Example rules might include:
- "no-unused-vars": true
- "no-console": true
- "prefer-const": true

### .cursorrules
No Cursor rules file found. Create a `.cursorrules` file with your custom rules. Example rules might include:
- "no-unused-vars": true
- "no-console": true
- "prefer-const": true

## Copilot Instructions

### .github/copilot-instructions.md
No Copilot instructions file found. Create a `.github/copilot-instructions.md` file with your Copilot configuration. Example instructions might include:
- "ignore": ["/node_modules/"]
- "exclude": ["/test/"]
- "include": ["/src/"]

## Notes
- If you have existing rules in .cursor/rules/ or .github/copilot-instructions.md, please add them to this file
- This file will be used by agentic coding agents to understand the codebase structure and conventions
- Make sure to update this file when adding new rules or changing existing ones
