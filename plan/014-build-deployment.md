# Task 014: Build y Deployment

## Descripción
Implementar sistema de build para producción con optimizaciones, y herramientas de deployment para múltiples plataformas (Docker, Kubernetes, Cloud providers).

## Prioridad
Alta

## Estimación
5-6 días

## Dependencias
- Task 004: Compilador .gox a Go
- Task 012: Sistema de configuración

## Subtasks

### 14.1 Sistema de build
- [ ] Comando `gox build`
- [ ] Compilación optimizada para producción
- [ ] Minificación de assets (CSS, JS)
- [ ] Tree-shaking de código no usado
- [ ] Bundling de static assets
- [ ] Build caching

### 14.2 Optimizaciones de producción
- [ ] Compilación de templates
- [ ] Compresión de assets
- [ ] Asset fingerprinting/hashing
- [ ] Code splitting automático
- [ ] Dead code elimination
- [ ] Binary size optimization

### 14.3 Multi-platform builds
- [ ] Cross-compilation para diferentes OS/arch
- [ ] Docker builds
- [ ] ARM64 support
- [ ] Static binary generation
- [ ] CGO handling

### 14.4 Docker integration
- [ ] Dockerfile generation automática
- [ ] Multi-stage builds
- [ ] Distroless images
- [ ] Health checks
- [ ] Security scanning

### 14.5 Deployment tools
- [ ] Comando `gox deploy`
- [ ] Deploy a Docker Registry
- [ ] Kubernetes manifests generation
- [ ] Cloud Run deployment
- [ ] AWS Lambda deployment
- [ ] Vercel/Netlify support

### 14.6 CI/CD integration
- [ ] GitHub Actions workflows
- [ ] GitLab CI templates
- [ ] Build artifacts
- [ ] Release automation
- [ ] Deployment pipelines

## Criterios de Aceptación

1. **Build básico funcionando**
   ```bash
   # Build para producción
   gox build
   # Output: dist/my-app (binary) + dist/static/ (assets)
   
   # Build con opciones
   gox build --target=linux/amd64 --compress --static
   
   # Build para Docker
   gox build --docker --tag=my-app:latest
   ```

2. **Assets optimizados**
   ```
   dist/
   ├── my-app                    # Binary optimizado
   ├── static/
   │   ├── css/
   │   │   └── app.a1b2c3.css   # CSS minificado con hash
   │   ├── js/
   │   │   └── app.d4e5f6.js    # JS minificado con hash
   │   └── images/
   │       └── logo.789abc.png   # Assets con hash
   └── Dockerfile                # Dockerfile generado
   ```

3. **Docker integration**
   ```dockerfile
   # Dockerfile generado automáticamente
   FROM scratch
   COPY my-app /app
   COPY static /static
   EXPOSE 8080
   HEALTHCHECK --interval=30s --timeout=3s \
     CMD ["/app", "health"]
   CMD ["/app", "serve"]
   ```

4. **Deployment commands**
   ```bash
   # Deploy a diferentes targets
   gox deploy --target=docker --registry=ghcr.io/user/app
   gox deploy --target=k8s --namespace=production
   gox deploy --target=cloud-run --region=us-central1
   gox deploy --target=lambda --region=us-east-1
   ```

## Tests Necesarios

### Tests Unitarios

1. **Test build básico**
```go
func TestBasicBuild(t *testing.T) {
    project := createTestProject(t)
    defer cleanupProject(project)
    
    builder := NewBuilder(BuildConfig{
        Source: project.Root,
        Output: filepath.Join(project.Root, "dist"),
        Target: "linux/amd64",
    })
    
    err := builder.Build()
    assert.NoError(t, err)
    
    // Verificar binary
    binaryPath := filepath.Join(project.Root, "dist", "app")
    assert.FileExists(t, binaryPath)
    
    // Verificar que es ejecutable
    info, err := os.Stat(binaryPath)
    assert.NoError(t, err)
    assert.True(t, info.Mode()&0111 != 0)
}
```

2. **Test optimización de assets**
```go
func TestAssetOptimization(t *testing.T) {
    builder := NewBuilder(BuildConfig{
        Optimize: true,
        Minify:   true,
    })
    
    // CSS sin optimizar
    originalCSS := `
    .component {
        background-color: #ffffff;
        padding: 1rem;
        margin: 0;
    }
    `
    
    optimized, err := builder.OptimizeCSS(originalCSS)
    assert.NoError(t, err)
    assert.Less(t, len(optimized), len(originalCSS))
    assert.Contains(t, optimized, ".component{background-color:#fff")
}
```

3. **Test cross-compilation**
```go
func TestCrossCompilation(t *testing.T) {
    targets := []string{
        "linux/amd64",
        "linux/arm64", 
        "darwin/amd64",
        "darwin/arm64",
        "windows/amd64",
    }
    
    for _, target := range targets {
        t.Run(target, func(t *testing.T) {
            builder := NewBuilder(BuildConfig{
                Target: target,
            })
            
            err := builder.Build()
            assert.NoError(t, err)
            
            binary := getBinaryName(target)
            assert.FileExists(t, binary)
        })
    }
}
```

### Tests de Integración

1. **Test Docker build**
```go
func TestDockerBuild(t *testing.T) {
    if !dockerAvailable() {
        t.Skip("Docker not available")
    }
    
    project := createTestProject(t)
    defer cleanupProject(project)
    
    builder := NewBuilder(BuildConfig{
        Docker: true,
        Tag:    "gox-test:latest",
    })
    
    err := builder.Build()
    assert.NoError(t, err)
    
    // Verificar imagen creada
    images, err := dockerClient.ImageList(context.Background(), types.ImageListOptions{})
    assert.NoError(t, err)
    
    found := false
    for _, img := range images {
        for _, tag := range img.RepoTags {
            if tag == "gox-test:latest" {
                found = true
                break
            }
        }
    }
    assert.True(t, found)
}
```

2. **Test deployment a Kubernetes**
```go
func TestKubernetesDeployment(t *testing.T) {
    deployer := NewKubernetesDeployer(K8sConfig{
        Namespace: "test",
        Image:     "my-app:latest",
        Replicas:  2,
    })
    
    manifests, err := deployer.GenerateManifests()
    assert.NoError(t, err)
    
    // Verificar deployment manifest
    deployment := manifests["deployment.yaml"]
    assert.Contains(t, deployment, "kind: Deployment")
    assert.Contains(t, deployment, "replicas: 2")
    assert.Contains(t, deployment, "image: my-app:latest")
    
    // Verificar service manifest
    service := manifests["service.yaml"]
    assert.Contains(t, service, "kind: Service")
    assert.Contains(t, service, "port: 80")
}
```

3. **Test CI/CD workflow generation**
```go
func TestGitHubActionsGeneration(t *testing.T) {
    generator := NewCIGenerator(CIConfig{
        Provider: "github",
        Registry: "ghcr.io",
        Stages:   []string{"test", "build", "deploy"},
    })
    
    workflow, err := generator.GenerateWorkflow()
    assert.NoError(t, err)
    
    assert.Contains(t, workflow, "name: CI/CD Pipeline")
    assert.Contains(t, workflow, "runs-on: ubuntu-latest")
    assert.Contains(t, workflow, "gox test")
    assert.Contains(t, workflow, "gox build")
    assert.Contains(t, workflow, "docker push")
}
```

### Tests de Performance

```go
func BenchmarkBuildTime(b *testing.B) {
    project := createLargeProject()
    builder := NewBuilder(BuildConfig{
        Source: project.Root,
        Cache:  true,
    })
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        builder.Build()
    }
}

func BenchmarkAssetOptimization(b *testing.B) {
    css := generateLargeCSS()
    builder := NewBuilder(BuildConfig{Optimize: true})
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        builder.OptimizeCSS(css)
    }
}
```

## Definición de Done

- [ ] Build system completo y optimizado
- [ ] Cross-compilation funcionando
- [ ] Docker integration completa
- [ ] Deployment a múltiples platforms
- [ ] CI/CD templates generados
- [ ] Assets optimizados automáticamente
- [ ] Tests con cobertura > 80%
- [ ] Documentación de deployment

## Notas Adicionales

- Los builds deben ser reproducibles
- Optimizar para tamaño de binary en producción
- Considerar security scanning automático
- Documentar best practices de deployment
- El sistema debe ser extensible para nuevas platforms
- Pensar en rollback strategies
- Implementar health checks por defecto