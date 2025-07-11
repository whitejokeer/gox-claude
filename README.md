# GOX Framework

[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/gox-framework/gox)](https://goreportcard.com/report/github.com/gox-framework/gox)

GOX is a modern web development framework that combines Go's power with HTMX's simplicity for ultra-fast Server-Side Rendered applications. Inspired by Vue.js SFC, React and Svelte, but designed to be "easy to learn, difficult to master".

## 🚀 Features

- **🔥 Ultra Performance**: SSR native + Go + HTMX = extreme speed
- **🎨 Developer Experience**: Integrated tools that accelerate development  
- **🏗️ Distributed Architecture**: Scalable from day one
- **📦 Zero Bundle**: No heavy client-side JavaScript
- **🔧 Easy to Learn**: Smooth learning curve
- **💎 Hard to Master**: Powerful advanced options

## 🏗️ Architecture

```
my-app/
├── pages/              # Frontend (.gox pages) - UI Layer
├── services/           # Backend services/modules - Business Logic  
├── routing/            # Route configuration - Presentation Layer
├── common/             # Shared code - Shared Layer
├── infra/              # Infrastructure - Infrastructure Layer
└── gox.config.yaml     # Main configuration
```

## 📋 .gox File Format

Single File Components inspired by Vue.js:

```gox
<template>
  <div class="p-4" hx-get="/api/stats" hx-target="#stats">
    <h1>{{.Title}}</h1>
    <div id="stats">{{.Stats}}</div>
  </div>
</template>

<go>
type DashboardData struct {
  Title string
  Stats interface{}
}

func (d *Dashboard) Mount(ctx *gox.Context) error {
  d.Title = "My Dashboard"
  return nil
}
</go>

<style scoped>
.dashboard {
  padding: 2rem;
  background: var(--bg-primary);
}
</style>
```

## 🛠️ Quick Start

### Installation

```bash
go install github.com/gox-framework/gox@latest
```

### Create Project

```bash
gox new project my-app
cd my-app
gox dev
```

## 📚 Documentation

- [Getting Started](docs/README.md)
- [Architecture Guide](docs/architecture.md)
- [API Reference](docs/api.md)
- [Examples](examples/)

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

Inspired by the excellent work of:
- [Vue.js](https://vuejs.org/) - Single File Components
- [HTMX](https://htmx.org/) - HTML-driven interactivity
- [Go](https://golang.org/) - Simplicity and performance