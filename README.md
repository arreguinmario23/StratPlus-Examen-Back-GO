# Prueba Técnica StratPlus - Servicio de Autenticación con JWT en Go

## Descripción
Servicio HTTP de autenticación de usuarios implementado en Go como parte de la prueba técnica para StratPlus. El servicio proporciona endpoints para registro y login utilizando tokens JWT, con validaciones exhaustivas de datos de entrada y almacenamiento en memoria de usuarios registrados.

## Requisitos
- Go 1.25.0 o superior
- Módulo JWT: `github.com/golang-jwt/jwt/v5 v5.3.0`

## Instalación

1. Clonar el repositorio
```bash
git clone <url-del-repositorio>
cd StratPlus-Examen-Back-GO-main
```

2. Instalar dependencias
```bash
go mod download
```

3. Ejecutar el servidor
```bash
go run prueba.go
```

El servidor se iniciará en `http://localhost:8080`

## Endpoints

### 1. Registro de Usuario
**POST** `/registro`

Registra un nuevo usuario en el sistema.

#### Request Body
```json
{
  "correo": "usuario@example.com",
  "telefono": "5551234567",
  "password": "Pass123@"
}
```

#### Validaciones
- **Correo**: Formato válido de email (usuario@dominio.extension)
- **Teléfono**: Exactamente 10 dígitos numéricos
- **Contraseña**: 
  - Entre 6 y 12 caracteres
  - Al menos una mayúscula
  - Al menos una minúscula
  - Al menos un número
  - Al menos un carácter especial (@, $, &)

#### Respuestas

**201 Created** - Usuario registrado exitosamente
```json
{
  "mensaje": "Usuario registrado exitosamente"
}
```

**400 Bad Request** - Datos inválidos o faltantes
```json
{
  "error": "Falta el campo correo"
}
```

**409 Conflict** - Usuario duplicado
```json
{
  "error": "El correo ya se encuentra registrado"
}
```

### 2. Login
**POST** `/login`

Autentica un usuario y devuelve un token JWT.

#### Request Body
```json
{
  "correo": "usuario@example.com",
  "password": "Pass123@"
}
```

#### Respuestas

**200 OK** - Login exitoso
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "fecha_inicio": "2025-08-24T17:24:41.190626-06:00"
}
```

**400 Bad Request** - Datos faltantes
```json
{
  "error": "Falta el campo contraseña"
}
```

**401 Unauthorized** - Credenciales incorrectas
```json
{
  "error": "Correo o contraseña incorrectos"
}
```

## Ejemplos de Uso

### Registro exitoso
```bash
curl -X POST http://localhost:8080/registro \
  -H "Content-Type: application/json" \
  -d '{
    "correo": "usuario@example.com",
    "telefono": "5551234567",
    "password": "Pass123@"
  }'
```

### Login exitoso
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "correo": "usuario@example.com",
    "password": "Pass123@"
  }'
```

### Validación de campo faltante
```bash
curl -X POST http://localhost:8080/registro \
  -H "Content-Type: application/json" \
  -d '{
    "correo": "test@test.com",
    "telefono": "5559999999"
  }'
```
Respuesta: `{"error":"Falta el campo contraseña"}`

### Validación de contraseña inválida
```bash
curl -X POST http://localhost:8080/registro \
  -H "Content-Type: application/json" \
  -d '{
    "correo": "otro@test.com",
    "telefono": "5557777777",
    "password": "simple"
  }'
```
Respuesta: `{"error":"Contraseña inválida"}`

## Estructura del Proyecto
```
StratPlus-Examen-Back-GO-main/
├── go.mod          # Configuración del módulo Go
├── go.sum          # Checksums de dependencias
├── prueba.go       # Código fuente principal
└── README.md       # Este archivo
```

## Características Implementadas

- Validación exhaustiva de datos de entrada
- Tokens JWT con expiración de 24 horas
- Validación de formato de correo electrónico
- Requisitos de complejidad de contraseña
- Prevención de usuarios duplicados
- Mensajes de error descriptivos en español

## Token JWT

El token JWT generado contiene:
- **correo**: Email del usuario autenticado
- **exp**: Fecha de expiración (24 horas desde la generación)

Firmado con algoritmo HS256.

## Notas Técnicas

- Base de datos en memoria (slice de Go)
- Puerto: 8080
- Algoritmo JWT: HS256
- Expiración de token: 24 horas

## Requerimientos Cumplidos

✅ **1. Servicio de registro** - Endpoint `/registro` que recibe correo, teléfono y contraseña  
✅ **2. Validación de duplicados** - Verifica correo y teléfono únicos con mensajes de error específicos  
✅ **3. Validación de contraseña** - 6-12 caracteres, mayúscula, minúscula, número y carácter especial (@, $, &)  
✅ **4. Validaciones de campos** - Formato de correo electrónico y teléfono de 10 dígitos  
✅ **5. Servicio de Login** - Endpoint `/login` que retorna token JWT y fecha de inicio  
✅ **6. Validación de campos requeridos** - Mensajes de error específicos para campos faltantes  
