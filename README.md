# Go Notes API (Deep Dive)

**Developer:** Solo Lead | **Methodology:** TDD | **Focus:** Authentication & Security

This project implements a secure Notes API using Go, following a strict Test-Driven Development (TDD) approach.

## Progress Tracking

### 🛠️ Fase 0: Configuración del Entorno de Pruebas

- [x] **Tarea 0.1: Project Layout & Test Infrastructure**
  - [x] `go mod init`
  - [x] Create directories: `cmd/api`, `internal/auth`, `internal/store`, `internal/api`
  - [x] Create `Makefile`

### 🛡️ Fase 1: Módulo de Autenticación (Core Logic)

#### Tarea 1.1: Servicio de Hashing de Contraseñas (Bcrypt)
- [x] **Paso 1: Write Tests (RED)**
  - [x] `TestHashPassword_Structure`
  - [x] `TestComparePassword_Success`
  - [x] `TestComparePassword_Failure`
- [x] **Paso 2: Implement (GREEN)**
  - [x] `Hash` function
  - [x] `Compare` function

#### Tarea 1.2: Servicio JWT (Generación y Validación)
- [x] **Paso 1: Write Tests (RED)**
  - [x] `TestGenerateToken_ContainsClaims`
  - [x] `TestValidateToken_Valid`
  - [x] `TestValidateToken_Expired`
  - [x] `TestValidateToken_TamperedSignature`
- [x] **Paso 2: Implement (GREEN)**
  - [x] `Claims` struct
  - [x] `GenerateToken`
  - [x] `ValidateToken`

### 💾 Fase 2: Capa de Persistencia (Store)

#### Tarea 2.1: Modelo de Usuario y Validación
- [x] **Paso 1: Write Tests (RED)**
  - [x] `TestUser_Validate_Success`
  - [x] `TestUser_Validate_BadEmail`
  - [x] `TestUser_Validate_ShortPassword`
- [x] **Paso 2: Implement (GREEN)**
  - [x] `User` struct
  - [x] `Validate` method

#### Tarea 2.2: User Store (Repository Pattern)
- [x] **Paso 1: Write Tests (RED - Integration)**
  - [x] `TestCreateUser_HappyPath`
  - [x] `TestCreateUser_DuplicateEmail`
  - [x] `TestGetByEmail_Found`
  - [x] `TestGetByEmail_NotFound`
- [x] **Paso 2: Implement (GREEN)**
  - [x] SQL connection
  - [x] `INSERT` query
  - [x] `SELECT` query
  - [x] Handle unique constraint

### 🌐 Fase 3: Capa HTTP (Handlers & Rutas)

#### Tarea 3.1: Register Handler (POST /auth/register)
- [ ] **Paso 1: Write Tests (RED)**
  - [ ] `TestRegister_Success`
  - [ ] `TestRegister_InvalidInput`
  - [ ] `TestRegister_Duplicate`
- [ ] **Paso 2: Implement (GREEN)**
  - [ ] `AuthHandler` struct
  - [ ] `Register` method

#### Tarea 3.2: Login Handler (POST /auth/login)
- [ ] **Paso 1: Write Tests (RED)**
  - [ ] `TestLogin_Success`
  - [ ] `TestLogin_UserNotFound`
  - [ ] `TestLogin_BadPassword`
- [ ] **Paso 2: Implement (GREEN)**
  - [ ] `Login` method

#### Tarea 3.3: Middleware de Autenticación
- [ ] **Paso 1: Write Tests (RED)**
  - [ ] `TestAuthMiddleware_NoHeader`
  - [ ] `TestAuthMiddleware_BadFormat`
  - [ ] `TestAuthMiddleware_InvalidToken`
  - [ ] `TestAuthMiddleware_Success`
- [ ] **Paso 2: Implement (GREEN)**
  - [ ] `WithAuth` function

### 📝 Fase 4: Notes (CRUD Rápido)

#### Tarea 4.1: Notes Store (Integration)
- [ ] `TestCreateNote`
- [ ] `TestListNotes_Isolation`

#### Tarea 4.2: Notes Handlers (Unit with Mock)
- [ ] `TestCreateNote_Authorized`
- [ ] `TestGetNotes_Format`

---

## Instructions (Original Plan)

### 📘 TDD Master Implementation Plan: Go Notes API (Deep Dive)

**Developer:** Solo Lead | **Methodology:** TDD | **Focus:** Authentication & Security

## 🛠️ Fase 0: Configuración del Entorno de Pruebas

Antes de escribir el primer test, necesitamos la infraestructura para ejecutarlos.

### Tarea 0.1: Project Layout & Test Infrastructure

* **Objetivo:** Tener un proyecto que compile y pueda ejecutar `go test`.
* **Acción:**
1. `go mod init github.com/ivan-almanza/notes-api`
2. Crear directorios: `cmd/api`, `internal/auth`, `internal/store`, `internal/api`.
3. Crear `Makefile`:
```makefile
test:
    go test -v ./...
test-coverage:
    go test -coverprofile=coverage.out ./... ; go tool cover -html=coverage.out

```





---

## 🛡️ Fase 1: Módulo de Autenticación (Core Logic)

*En esta fase no tocamos la base de datos ni HTTP. Solo lógica pura.*

### Tarea 1.1: Servicio de Hashing de Contraseñas (Bcrypt)

**Archivo de Test:** `internal/auth/password_test.go`
**Archivo de Implementación:** `internal/auth/password.go`

#### 🔴 Paso 1: Write Tests (RED)

Crea los siguientes casos de prueba:

1. **`TestHashPassword_Structure`**:
* Input: "mysecretpassword"
* Assert: El output *no* es igual al input. El output no está vacío.


2. **`TestComparePassword_Success`**:
* Setup: Generar hash de "password123".
* Action: Llamar a `Compare("password123", hash)`.
* Assert: Debe retornar `nil` (sin error).


3. **`TestComparePassword_Failure`**:
* Setup: Generar hash de "password123".
* Action: Llamar a `Compare("wrongpassword", hash)`.
* Assert: Debe retornar error `bcrypt.ErrMismatchedHashAndPassword`.



#### 🟢 Paso 2: Implement (GREEN)

Implementa las funciones usando `golang.org/x/crypto/bcrypt`:

* `func Hash(password string) (string, error)`
* `func Compare(password, hash string) error`

---

### Tarea 1.2: Servicio JWT (Generación y Validación)

**Archivo de Test:** `internal/auth/jwt_test.go`
**Archivo de Implementación:** `internal/auth/jwt.go`

#### 🔴 Paso 1: Write Tests (RED)

Define constantes de test (ej: `Secret = "test-secret"`).

1. **`TestGenerateToken_ContainsClaims`**:
* Input: `userID = "user-123"`
* Action: Generar token. Parsear el token (sin validar firma aún).
* Assert: Los claims deben contener `sub: "user-123"` y un `exp` mayor al tiempo actual.


2. **`TestValidateToken_Valid`**:
* Setup: Generar un token válido con el secreto correcto.
* Action: Validar el token.
* Assert: Retorna el `jwt.Token` válido y sin errores.


3. **`TestValidateToken_Expired`**:
* *Truco TDD:* Necesitarás poder inyectar el tiempo o crear un token manualmente con `exp` en el pasado.
* Assert: Retorna error indicando que el token expiró.


4. **`TestValidateToken_TamperedSignature`**:
* Setup: Generar token con secreto "A". Intentar validar con secreto "B".
* Assert: Error de firma inválida.



#### 🟢 Paso 2: Implement (GREEN)

* Define `struct Claims` que embeba `jwt.RegisteredClaims`.
* Implementa `GenerateToken(userID string) (string, error)`.
* Implementa `ValidateToken(tokenString string) (*jwt.Token, error)`.

---

## 💾 Fase 2: Capa de Persistencia (Store)

*Aquí usamos Integration Testing. Necesitas Docker corriendo Postgres.*

### Tarea 2.1: Modelo de Usuario y Validación

**Archivo de Test:** `internal/store/users_model_test.go`
**Implementación:** `internal/store/users.go`

#### 🔴 Paso 1: Write Tests (RED)

Validar reglas de negocio *antes* de tocar la DB.

1. **`TestUser_Validate_Success`**: User con email y password correctos pasa.
2. **`TestUser_Validate_BadEmail`**: Email "no-arroba" retorna error.
3. **`TestUser_Validate_ShortPassword`**: Password < 6 chars retorna error.

#### 🟢 Paso 2: Implement (GREEN)

* Define `type User struct`.
* Agrega método `func (u *User) Validate() error`.

### Tarea 2.2: User Store (Repository Pattern)

**Archivo de Test:** `internal/store/users_store_test.go`
**Implementación:** `internal/store/users.go`

*Requisito:* Define una interfaz para facilitar el mocking más adelante.

```go
type UserStorer interface {
    Create(ctx context.Context, user *User) error
    GetByEmail(ctx context.Context, email string) (*User, error)
}

```

#### 🔴 Paso 1: Write Tests (RED - Integration)

*Setup:* Función `setupTestDB()` que limpia la tabla `users` antes de cada test.

1. **`TestCreateUser_HappyPath`**:
* Action: Crear usuario.
* Assert: No retorna error. El campo `ID` (UUID) ya no es nil/vacío. `CreatedAt` tiene valor.


2. **`TestCreateUser_DuplicateEmail`**:
* Action: Crear usuario "a@a.com". Intentar crear otro con "a@a.com".
* Assert: Debe retornar error específico (ej: `ErrDuplicateEmail`).


3. **`TestGetByEmail_Found`**:
* Action: Crear usuario, luego buscarlo por email.
* Assert: El ID del usuario retornado coincide con el creado.


4. **`TestGetByEmail_NotFound`**:
* Action: Buscar "ghost@user.com".
* Assert: Retorna error `ErrNotFound`.



#### 🟢 Paso 2: Implement (GREEN)

* Implementa la conexión SQL.
* Escribe las queries `INSERT` y `SELECT`.
* Maneja el error de violación de constraint `unique` de Postgres (código `23505` en `lib/pq`).

---

## 🌐 Fase 3: Capa HTTP (Handlers & Rutas)

*Aquí unimos todo. Usaremos `httptest` y Mocks del Store.*

### Tarea 3.1: Register Handler (POST /auth/register)

**Archivo de Test:** `internal/api/auth_test.go`

#### 🔴 Paso 1: Write Tests (RED)

Necesitas un **MockUserStore** para no depender de la DB real en unit tests.

1. **`TestRegister_Success`**:
* Input: JSON válido.
* Mock: `Create` retorna `nil`.
* Assert: Status 201 Created. Response Body contiene el ID.


2. **`TestRegister_InvalidInput`**:
* Input: JSON con email inválido.
* Assert: Status 400 Bad Request. Mock `Create` NUNCA se llama.


3. **`TestRegister_Duplicate`**:
* Input: JSON válido.
* Mock: `Create` retorna `ErrDuplicateEmail`.
* Assert: Status 409 Conflict.



#### 🟢 Paso 2: Implement (GREEN)

* Define `AuthHandler` que tenga acceso al `UserStore`.
* Implementa método `Register`. Decodifica JSON -> Valida -> Llama Store -> Responde.

### Tarea 3.2: Login Handler (POST /auth/login)

**Archivo de Test:** `internal/api/auth_test.go`

#### 🔴 Paso 1: Write Tests (RED)

1. **`TestLogin_Success`**:
* Input: Email/Pass correctos.
* Mock: `GetByEmail` retorna User con hash correcto.
* Assert: Status 200. Body contiene `token`.


2. **`TestLogin_UserNotFound`**:
* Mock: `GetByEmail` retorna `ErrNotFound`.
* Assert: Status 401 Unauthorized (¡No digas 404 por seguridad!).


3. **`TestLogin_BadPassword`**:
* Mock: `GetByEmail` retorna User, pero al comparar hash falla.
* Assert: Status 401 Unauthorized.



#### 🟢 Paso 2: Implement (GREEN)

* Implementa método `Login`.
* Busca user -> `bcrypt.Compare` -> `jwt.Generate` -> JSON Response.

### Tarea 3.3: Middleware de Autenticación

**Archivo de Test:** `internal/api/middleware_test.go`

#### 🔴 Paso 1: Write Tests (RED)

Crea un handler "dummy" que solo responda 200 OK si llega a ejecutarse. Envuelve ese handler con tu middleware.

1. **`TestAuthMiddleware_NoHeader`**: Request sin header. Assert: 401. Handler dummy NO ejecutado.
2. **`TestAuthMiddleware_BadFormat`**: Header "Token xyz". Assert: 401.
3. **`TestAuthMiddleware_InvalidToken`**: Header "Bearer <token_invalido>". Assert: 401.
4. **`TestAuthMiddleware_Success`**:
* Setup: Generar token válido.
* Request: Header "Bearer <token>".
* Assert: 200 OK. Handler dummy ejecutado.
* **Crucial:** Verificar que dentro del handler dummy, `r.Context().Value("userID")` no es nil.



#### 🟢 Paso 2: Implement (GREEN)

* Implementa función `func WithAuth(next http.Handler) http.Handler`.
* Parsea Header -> Valida Token -> `context.WithValue` -> `next.ServeHTTP`.

---

## 📝 Fase 4: Notes (CRUD Rápido)

*Ya tienes los patrones. Aplícalos a las notas.*

### Tarea 4.1: Notes Store (Integration)

**Tests:**

* `TestCreateNote`: Verifica que se guarde con el `user_id` correcto.
* `TestListNotes_Isolation`: Crea notas de User A y User B. Pide notas de User A. Asegura que NO vengan las de B.

### Tarea 4.2: Notes Handlers (Unit with Mock)

**Tests:**

* `TestCreateNote_Authorized`: Inyecta el `userID` en el contexto (simulando el middleware) y verifica que el handler lo use para llamar al Store.
* `TestGetNotes_Format`: Verifica que el JSON de respuesta tenga la estructura `data: [], meta: {}`.
