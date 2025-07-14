# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

GOX Framework is a modern web development framework that combines Go's power with HTMX's simplicity for ultra-fast Server-Side Rendered applications. It follows a "distributed by default" architecture and uses Single File Components (.gox files) inspired by Vue.js.

## Core Architecture

### .gox File Format
The framework's core abstraction is the .gox file (similar to Vue SFC):
- `<template>`: HTML with Go template syntax + HTMX attributes
- `<go>`: Go code with component logic and HTTP handlers
- `<style>`: CSS with optional scoped styling and Tailwind support
- All sections are optional

### Project Structure
```
my-app/
├── app/                # Frontend application (UI layer)
│   ├── pages/         # Frontend pages (.gox files) - auto-routed
│   ├── components/    # Reusable UI components
│   ├── shared/        # Shared UI resources
│   │   ├── ui/       # Shared components (--shared flag)
│   │   └── layouts/  # Layout components
│   └── routing/      # Page routing configuration
├── services/          # Backend microservices
├── common/            # Shared code between app and services
│   ├── middleware/   # HTTP middleware
│   └── discovery/    # Service discovery
├── infra/             # Infrastructure (database, cache, messaging)
├── docker-compose.yml # Development orchestration
├── go.work           # Go workspace for multi-module development
└── gox.config.yaml   # Main configuration
```

### Distributed Architecture
- **Services**: Independent Go services that can include .gox components
- **Modules**: Services with UI components that can be reused
- **Communication**: gRPC/REST between services with auto-generated clients
- **Database**: Per-service or shared strategies with event sourcing sync

## Development Commands

```bash
# Project creation
gox new project my-app                    # Full application
gox new module user-management             # Reusable module
gox new service auth-service               # Pure backend service

# Code generation
gox generate page dashboard --auth         # Generate authenticated page
gox generate component user-card --props="user:User,featured:bool"
gox generate service users --crud --model=User

# Development
gox dev                                    # Start dev server with hot reload
gox dev --storybook                        # Include integrated Storybook
gox build                                  # Production build
gox test                                   # Run test suite
gox migrate up                             # Run database migrations

# Deployment
gox deploy --target=docker --registry=ghcr.io/user/app
gox deploy --target=k8s --namespace=production
```

## Key Design Principles

### Performance First
- Zero client-side JavaScript bundle
- Server-Side Rendering with HTMX for interactivity
- Sub-50ms TTFB target, <100ms FCP target

### Developer Experience
- Integrated Storybook with auto-discovery of components
- Hot reload <500ms for .gox files
- Type-safe Go throughout with automatic HTMX endpoint generation

### Routing System
File-based routing (pages/about.gox → /about) with support for:
- Dynamic routes: [id].gox → /:id
- Catch-all: [...].gox → /*
- Manual routing via routing/routes.go when needed

### Component System
- Props with Go struct tags for validation and Storybook generation
- Lifecycle hooks: BeforeMount, Mount, AfterMount
- HTMX handlers as Go methods
- Scoped CSS and Tailwind integration

## Configuration

The gox.config.yaml controls:
- Architecture pattern (distributed/monolithic)
- Database strategy (per-service/shared)
- Authentication provider (JWT/OAuth2/session)
- Communication protocol (gRPC/REST/GraphQL)
- Development settings (hot_reload, storybook, port)

## Implementation Status

Task 2 (CLI Básico) has been completed. The project now uses an `app` directory instead of `gateway` for better frontend developer familiarity.

Current status:
- ✅ Task 2: CLI básico with context detection, new/generate commands
- Project structure uses `app/` for UI layer (pages, components, routing)
- Remaining tasks in /plan/ directory for implementation
- Each task includes subtasks, acceptance criteria, and comprehensive tests

## Project Goals

Create a framework that is "easy to learn, difficult to master" following Go's philosophy, providing an alternative to JavaScript-heavy frameworks while maintaining modern developer experience and performance.