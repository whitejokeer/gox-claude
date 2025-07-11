# Task 015: Testing Framework

## Descripción
Implementar un framework de testing completo para GOX que incluya unit tests, integration tests, E2E tests, y herramientas de mocking/stubbing específicas para componentes .gox.

## Prioridad
Alta

## Estimación
5-6 días

## Dependencias
- Task 008: Sistema de componentes
- Task 005: Sistema de routing

## Subtasks

### 15.1 Testing framework base
- [ ] Integración con testing nativo de Go
- [ ] Test runner personalizado
- [ ] Assertions helpers
- [ ] Test discovery automático
- [ ] Parallel test execution

### 15.2 Component testing
- [ ] Testing de componentes .gox
- [ ] Rendering tests
- [ ] Props validation testing
- [ ] Event handling tests
- [ ] Lifecycle testing

### 15.3 HTTP/API testing
- [ ] HTTP test helpers
- [ ] Request/response mocking
- [ ] Middleware testing
- [ ] Authentication testing
- [ ] HTMX interaction testing

### 15.4 Database testing
- [ ] Test database setup/teardown
- [ ] Transaction rollback
- [ ] Fixtures y seeding
- [ ] Repository testing
- [ ] Migration testing

### 15.5 E2E testing
- [ ] Browser automation (Playwright/Selenium)
- [ ] Page object patterns
- [ ] Visual regression testing
- [ ] Performance testing
- [ ] Accessibility testing

### 15.6 Testing utilities
- [ ] Test project scaffolding
- [ ] Mock/stub generators
- [ ] Test data factories
- [ ] Coverage reporting
- [ ] Test parallelization

## Criterios de Aceptación

1. **Component testing**
   ```go
   // component_test.gox
   func TestUserCard(t *testing.T) {
       comp := &UserCard{
           User: User{
               Name:  "John Doe",
               Email: "john@example.com",
           },
           Featured: true,
       }
       
       // Test rendering
       html, err := gox.RenderComponent(comp)
       assert.NoError(t, err)
       assert.Contains(t, html, "John Doe")
       assert.Contains(t, html, "featured")
       
       // Test props validation
       comp.User = User{} // Invalid
       err = comp.ValidateProps()
       assert.Error(t, err)
   }
   ```

2. **HTTP testing**
   ```go
   func TestUserAPI(t *testing.T) {
       app := gox.NewTestApp()
       
       // Mock database
       mockDB := gox.NewMockDB()
       mockDB.On("FindUser", 1).Return(&User{ID: 1, Name: "John"}, nil)
       app.SetDB(mockDB)
       
       // Test endpoint
       req := gox.NewTestRequest("GET", "/api/users/1", nil)
       resp := app.Request(req)
       
       assert.Equal(t, 200, resp.StatusCode)
       assert.JSONContains(t, resp.Body, `{"name": "John"}`)
   }
   ```

3. **E2E testing**
   ```go
   func TestUserFlow(t *testing.T) {
       browser := gox.NewBrowser(t)
       defer browser.Close()
       
       page := browser.NewPage()
       
       // Navigate to login
       page.Navigate("/login")
       page.Fill("#email", "user@example.com")
       page.Fill("#password", "password")
       page.Click("#login-btn")
       
       // Verify redirect
       assert.Equal(t, "/dashboard", page.URL().Path)
       assert.Contains(t, page.TextContent("h1"), "Welcome")
   }
   ```

4. **Test CLI commands**
   ```bash
   # Run all tests
   gox test
   
   # Run specific test types
   gox test --unit
   gox test --integration
   gox test --e2e
   
   # With coverage
   gox test --coverage
   
   # Watch mode
   gox test --watch
   ```

## Tests Necesarios

### Tests Unitarios

1. **Test component rendering**
```go
func TestComponentRendering(t *testing.T) {
    tests := []struct {
        name      string
        component gox.Component
        expected  []string
        notExpected []string
    }{
        {
            name: "button with label",
            component: &Button{
                Label: "Click me",
                Variant: "primary",
            },
            expected: []string{"Click me", "btn-primary"},
        },
        {
            name: "disabled button",
            component: &Button{
                Label: "Disabled",
                Disabled: true,
            },
            expected: []string{"disabled"},
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            html, err := gox.RenderComponent(tt.component)
            assert.NoError(t, err)
            
            for _, exp := range tt.expected {
                assert.Contains(t, html, exp)
            }
            for _, notExp := range tt.notExpected {
                assert.NotContains(t, html, notExp)
            }
        })
    }
}
```

2. **Test HTTP handlers**
```go
func TestHTTPHandlers(t *testing.T) {
    app := gox.NewTestApp()
    
    tests := []struct {
        name           string
        method         string
        path           string
        body           interface{}
        expectedStatus int
        expectedBody   string
    }{
        {
            name:           "get users",
            method:         "GET",
            path:           "/api/users",
            expectedStatus: 200,
            expectedBody:   `"users"`,
        },
        {
            name:           "create user",
            method:         "POST",
            path:           "/api/users",
            body:           map[string]string{"name": "John"},
            expectedStatus: 201,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := gox.NewTestRequest(tt.method, tt.path, tt.body)
            resp := app.Request(req)
            
            assert.Equal(t, tt.expectedStatus, resp.StatusCode)
            if tt.expectedBody != "" {
                assert.Contains(t, resp.Body, tt.expectedBody)
            }
        })
    }
}
```

3. **Test database operations**
```go
func TestDatabaseOperations(t *testing.T) {
    db := gox.NewTestDB(t)
    defer db.Close()
    
    repo := NewUserRepository(db)
    
    // Test create
    user := &User{Name: "John", Email: "john@example.com"}
    err := repo.Create(user)
    assert.NoError(t, err)
    assert.NotZero(t, user.ID)
    
    // Test find
    found, err := repo.FindByEmail("john@example.com")
    assert.NoError(t, err)
    assert.Equal(t, user.ID, found.ID)
    
    // Test update
    found.Name = "John Doe"
    err = repo.Update(found)
    assert.NoError(t, err)
    
    // Verify update
    updated, err := repo.FindByID(found.ID)
    assert.NoError(t, err)
    assert.Equal(t, "John Doe", updated.Name)
}
```

### Tests de Integración

1. **Test full request flow**
```go
func TestFullRequestFlow(t *testing.T) {
    app := gox.NewTestApp()
    db := gox.NewTestDB(t)
    app.SetDB(db)
    
    // Seed data
    user := &User{Email: "test@example.com", Password: hashPassword("password")}
    db.Create(user)
    
    // Test login
    loginReq := gox.NewTestRequest("POST", "/auth/login", map[string]string{
        "email": "test@example.com",
        "password": "password",
    })
    loginResp := app.Request(loginReq)
    assert.Equal(t, 200, loginResp.StatusCode)
    
    // Extract token
    var loginData map[string]interface{}
    json.Unmarshal([]byte(loginResp.Body), &loginData)
    token := loginData["token"].(string)
    
    // Test protected endpoint
    protectedReq := gox.NewTestRequest("GET", "/api/profile", nil)
    protectedReq.SetHeader("Authorization", "Bearer "+token)
    protectedResp := app.Request(protectedReq)
    assert.Equal(t, 200, protectedResp.StatusCode)
}
```

2. **Test component integration**
```go
func TestComponentIntegration(t *testing.T) {
    app := gox.NewTestApp()
    
    // Register components
    app.RegisterComponent("user-card", &UserCard{})
    app.RegisterComponent("user-list", &UserList{})
    
    // Test page with nested components
    page := `
    <gox-component src="user-list" :users="{{.Users}}">
      <template slot="item" slot-scope="user">
        <gox-component src="user-card" :user="user" />
      </template>
    </gox-component>
    `
    
    data := map[string]interface{}{
        "Users": []User{{Name: "John"}, {Name: "Jane"}},
    }
    
    html, err := app.RenderTemplate(page, data)
    assert.NoError(t, err)
    assert.Contains(t, html, "John")
    assert.Contains(t, html, "Jane")
}
```

### Tests E2E

```go
func TestE2EUserRegistration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping E2E test in short mode")
    }
    
    browser := gox.NewBrowser(t)
    defer browser.Close()
    
    page := browser.NewPage()
    
    // Navigate to registration
    page.Navigate("/register")
    
    // Fill form
    page.Fill("#name", "John Doe")
    page.Fill("#email", "john@example.com")
    page.Fill("#password", "password123")
    page.Click("#register-btn")
    
    // Wait for redirect
    page.WaitForURL("/dashboard")
    
    // Verify user is logged in
    assert.Contains(t, page.TextContent(".user-name"), "John Doe")
    
    // Test navigation
    page.Click("a[href='/profile']")
    page.WaitForURL("/profile")
    
    // Verify profile page
    assert.Equal(t, "john@example.com", page.InputValue("#email"))
}
```

## Definición de Done

- [ ] Framework de testing completo
- [ ] Component testing funcionando
- [ ] HTTP/API testing utilities
- [ ] Database testing helpers
- [ ] E2E testing con browser automation
- [ ] CLI commands para testing
- [ ] Coverage reporting
- [ ] Tests con cobertura > 90%

## Notas Adicionales

- Integrar con testing tools estándar de Go
- Los tests deben ser rápidos y confiables
- Proveer helpers para casos comunes
- Documentar patrones de testing
- El framework debe ser extensible
- Considerar property-based testing
- Implementar test parallelization eficiente