# Task 010: Gestión de Base de Datos

## Descripción
Implementar un sistema completo de gestión de base de datos con soporte para múltiples drivers, migraciones automáticas, repositorios, y sincronización entre servicios en arquitectura distribuida.

## Prioridad
Alta

## Estimación
6-7 días

## Dependencias
- Task 012: Sistema de configuración (parcial)
- Task 008: Sistema de componentes

## Subtasks

### 10.1 Abstracción de base de datos
- [ ] Interface común para múltiples drivers
- [ ] Soporte para PostgreSQL, MySQL, SQLite
- [ ] Connection pooling
- [ ] Transacciones distribuidas
- [ ] Query builder básico

### 10.2 Sistema de migraciones
- [ ] CLI para crear migraciones
- [ ] Ejecutor de migraciones up/down
- [ ] Tracking de migraciones aplicadas
- [ ] Rollback de migraciones
- [ ] Migraciones automáticas desde modelos

### 10.3 ORM/Repository pattern
- [ ] Definir interface Repository base
- [ ] Implementar CRUD genérico
- [ ] Soporte para relaciones
- [ ] Lazy loading
- [ ] Query optimization

### 10.4 Estrategias multi-database
- [ ] Database por servicio
- [ ] Database compartida
- [ ] Configuración híbrida
- [ ] Connection routing
- [ ] Database sharding

### 10.5 Sincronización entre servicios
- [ ] Event sourcing
- [ ] Change Data Capture (CDC)
- [ ] Saga pattern
- [ ] Eventual consistency
- [ ] Conflict resolution

### 10.6 Herramientas de desarrollo
- [ ] Database seeding
- [ ] Fixtures para testing
- [ ] Database inspector
- [ ] Performance monitoring
- [ ] Backup/restore helpers

## Criterios de Aceptación

1. **Configuración de base de datos**
   ```yaml
   # gox.config.yaml
   database:
     strategy: "per-service"
     default:
       driver: "postgres"
       host: "localhost"
       port: 5432
       
     services:
       auth:
         name: "auth_db"
         migrations: "./services/auth/migrations"
       users:
         name: "users_db"
         migrations: "./services/users/migrations"
         
     sync:
       method: "event-sourcing"
       bus: "nats"
   ```

2. **Modelos y repositorios**
   ```go
   // Modelo
   type User struct {
       gox.Model
       Email    string `gorm:"uniqueIndex" validate:"email"`
       Name     string `validate:"required"`
       Profile  Profile `gorm:"foreignKey:UserID"`
       Posts    []Post  `gorm:"many2many:user_posts"`
   }
   
   // Repository
   type UserRepository interface {
       gox.Repository[User]
       FindByEmail(email string) (*User, error)
       FindActive() ([]User, error)
   }
   
   // Implementación
   func (r *userRepository) FindByEmail(email string) (*User, error) {
       var user User
       err := r.db.Where("email = ?", email).
           Preload("Profile").
           First(&user).Error
       return &user, err
   }
   ```

3. **Migraciones**
   ```go
   // migrations/001_create_users.go
   func init() {
       gox.RegisterMigration(&gox.Migration{
           ID: "001_create_users",
           Up: func(db *gox.DB) error {
               return db.CreateTable(&User{})
           },
           Down: func(db *gox.DB) error {
               return db.DropTable(&User{})
           },
       })
   }
   ```

4. **Sincronización con eventos**
   ```go
   // Publicar evento después de crear
   func (s *UserService) CreateUser(data CreateUserDTO) (*User, error) {
       user := &User{...}
       
       err := s.db.Transaction(func(tx *gox.DB) error {
           if err := tx.Create(user).Error; err != nil {
               return err
           }
           
           // Publicar evento
           return s.eventBus.Publish("user.created", UserCreatedEvent{
               UserID: user.ID,
               Email:  user.Email,
           })
       })
       
       return user, err
   }
   ```

## Tests Necesarios

### Tests Unitarios

1. **Test repository base**
```go
func TestBaseRepository(t *testing.T) {
    db := setupTestDB()
    repo := gox.NewRepository[User](db)
    
    // Create
    user := &User{Name: "John", Email: "john@example.com"}
    err := repo.Create(user)
    assert.NoError(t, err)
    assert.NotZero(t, user.ID)
    
    // Find
    found, err := repo.FindByID(user.ID)
    assert.NoError(t, err)
    assert.Equal(t, user.Email, found.Email)
    
    // Update
    found.Name = "John Doe"
    err = repo.Update(found)
    assert.NoError(t, err)
    
    // Delete
    err = repo.Delete(user.ID)
    assert.NoError(t, err)
}
```

2. **Test migraciones**
```go
func TestMigrations(t *testing.T) {
    db := setupTestDB()
    migrator := gox.NewMigrator(db)
    
    // Registrar migraciones
    migrator.Register(&gox.Migration{
        ID: "001_test",
        Up: func(db *gox.DB) error {
            return db.Exec(`CREATE TABLE test (id INT PRIMARY KEY)`).Error
        },
        Down: func(db *gox.DB) error {
            return db.Exec(`DROP TABLE test`).Error
        },
    })
    
    // Ejecutar
    err := migrator.Up()
    assert.NoError(t, err)
    
    // Verificar
    var exists bool
    db.Raw("SELECT EXISTS (SELECT FROM pg_tables WHERE tablename = 'test')").Scan(&exists)
    assert.True(t, exists)
    
    // Rollback
    err = migrator.Down(1)
    assert.NoError(t, err)
}
```

3. **Test transacciones**
```go
func TestDistributedTransaction(t *testing.T) {
    db1 := setupTestDB("db1")
    db2 := setupTestDB("db2")
    
    tx := gox.NewDistributedTx(db1, db2)
    
    err := tx.Execute(func() error {
        // Operación en db1
        if err := db1.Create(&User{Name: "User1"}).Error; err != nil {
            return err
        }
        
        // Operación en db2
        if err := db2.Create(&Order{Total: 100}).Error; err != nil {
            return err
        }
        
        return nil
    })
    
    assert.NoError(t, err)
}
```

### Tests de Integración

1. **Test sincronización entre servicios**
```go
func TestEventSourcingSync(t *testing.T) {
    // Setup servicios
    userService := setupUserService()
    orderService := setupOrderService()
    
    // Suscribir a eventos
    var receivedEvent UserCreatedEvent
    orderService.Subscribe("user.created", func(event UserCreatedEvent) {
        receivedEvent = event
    })
    
    // Crear usuario
    user, err := userService.CreateUser(CreateUserDTO{
        Name: "Test User",
        Email: "test@example.com",
    })
    assert.NoError(t, err)
    
    // Verificar evento recibido
    time.Sleep(100 * time.Millisecond)
    assert.Equal(t, user.ID, receivedEvent.UserID)
}
```

2. **Test estrategia multi-database**
```go
func TestMultiDatabaseStrategy(t *testing.T) {
    config := &gox.DatabaseConfig{
        Strategy: "per-service",
        Services: map[string]gox.DBConfig{
            "users": {Driver: "postgres", Name: "users_db"},
            "orders": {Driver: "mysql", Name: "orders_db"},
        },
    }
    
    dbManager := gox.NewDatabaseManager(config)
    
    // Obtener conexiones
    usersDB := dbManager.GetDB("users")
    ordersDB := dbManager.GetDB("orders")
    
    assert.NotNil(t, usersDB)
    assert.NotNil(t, ordersDB)
    
    // Verificar aislamiento
    var userCount, orderCount int64
    usersDB.Model(&User{}).Count(&userCount)
    ordersDB.Model(&Order{}).Count(&orderCount)
}
```

### Tests de Performance

```go
func BenchmarkRepository(b *testing.B) {
    db := setupBenchDB()
    repo := gox.NewRepository[User](db)
    
    // Preparar datos
    users := make([]User, 1000)
    for i := range users {
        users[i] = User{
            Name: fmt.Sprintf("User%d", i),
            Email: fmt.Sprintf("user%d@example.com", i),
        }
    }
    db.Create(&users)
    
    b.Run("FindByID", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            repo.FindByID(uint(i%1000 + 1))
        }
    })
    
    b.Run("FindWithPreload", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            var user User
            db.Preload("Profile").Preload("Posts").First(&user, uint(i%1000+1))
        }
    })
}

func BenchmarkMigrations(b *testing.B) {
    for i := 0; i < b.N; i++ {
        db := setupTestDB()
        migrator := gox.NewMigrator(db)
        migrator.Up()
        dropTestDB(db)
    }
}
```

## Definición de Done

- [ ] Abstracción de DB con múltiples drivers
- [ ] Sistema de migraciones completo
- [ ] Repository pattern implementado
- [ ] Estrategias multi-database funcionando
- [ ] Sincronización entre servicios
- [ ] Herramientas de desarrollo
- [ ] Tests con cobertura > 80%
- [ ] Documentación y ejemplos

## Notas Adicionales

- Usar GORM como base pero con abstracción
- Las migraciones deben ser idempotentes
- Considerar performance desde el diseño
- Implementar connection pooling eficiente
- Documentar patrones de acceso a datos
- Pensar en monitoring y debugging de queries