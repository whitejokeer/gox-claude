# GOX Framework Project Structure

This document explains the purpose of each file and directory in the GOX Framework project.

## Root Directory Files

### Documentation
- `README.md` - Main project documentation with installation and usage instructions
- `CHANGELOG.md` - Version history and release notes
- `LICENSE` - MIT license for the project
- `CONTRIBUTING.md` - Guidelines for contributing to the project
- `PRD.md` - Product Requirements Document outlining the framework's vision

### Configuration
- `go.mod` / `go.sum` - Go module dependencies
- `Makefile` - Build automation and common tasks
- `Dockerfile` - Container image for GOX CLI
- `CLAUDE.md` - Project-specific instructions for AI assistants

## Directory Structure

### `/cmd/gox/`
The GOX CLI implementation.
- `main.go` - CLI entry point
- `commands/` - Command implementations
  - `context/` - Project context detection (inside/outside project)
  - `new/` - Create new projects and services
  - `generate/` - Generate pages, components, services, middleware
  - `dev/` - Development server (placeholder for future)
  - `root.go` - Root command configuration

### `/internal/`
Internal packages not meant for external use.
- `parser/` - GOX file parser implementation
  - Parser for .gox single file components
  - Lexer, AST, template/style/go parsers
  - Component detection and validation
- `project/` - Project management utilities (future)

### `/pkg/`
Public packages that can be imported by other projects.
- `version/` - Version information for the framework

### `/docs/`
Additional documentation.
- `GOX_CONVENTIONS.md` - Coding conventions and best practices
- `PROJECT_STRUCTURE.md` - This file

### `/plan/`
Implementation planning documents.
- `001-015-*.md` - Detailed task breakdown for framework implementation
- Each file represents a major feature/component to be developed

### `/scripts/`
Utility scripts for development and testing.
- `install-hooks.sh` - Git hooks installation
- `test-e2e.sh` - End-to-end testing script
- `verify-task2.sh` - Task 2 verification script
- `demo-task2.sh` - Task 2 demonstration script

### `/examples/`
Example applications demonstrating GOX Framework usage.
- `blog-app/` - Complete blog application example
  - `app/` - Frontend application with pages and components
  - `services/` - Microservices (posts service)
  - `common/` - Shared middleware and utilities
  - `docker-compose.yml` - Multi-service development setup

## Key Concepts

### App-Based Architecture
The framework uses an "app-based" architecture where:
- `app/` contains the frontend (pages, components, routing)
- `services/` contains backend microservices
- `common/` contains shared code between app and services

### Single File Components (.gox)
Components use a Vue-inspired format with three sections:
- `<template>` - HTML with Go templates and HTMX
- `<go>` - Go code for component logic
- `<style>` - Scoped CSS styles

### Context-Aware CLI
Commands are context-aware:
- `gox new` must be run OUTSIDE a project
- `gox generate` must be run INSIDE a project

## Development Workflow

1. **Create a project**: `gox new project my-app`
2. **Generate components**: `gox generate component user-card`
3. **Generate pages**: `gox generate page dashboard --auth`
4. **Generate services**: `gox generate service users --api`
5. **Run development**: `docker-compose up` (future: `gox dev`)

## Testing

Run all tests:
```bash
go test ./...
```

Run specific package tests:
```bash
go test ./cmd/gox/commands/...
```

## Building

Build the CLI:
```bash
go build -o gox ./cmd/gox/
```

## Status

✅ **Completed**:
- Task 001: Initial setup
- Task 002: Basic CLI with context detection

🚧 **In Progress**:
- Task 003: GOX file parser

📋 **Planned**:
- Tasks 004-015: Compiler, routing, dev server, etc.