# prestadores-api

Una API REST para gestionar prestadores de servicios construida con Go y Gin.

## Comenzar

### Prerrequisitos
- Go 1.25+

### Instalacion
Clonar el repositorio e instalar dependencias:
- git clone <repository-url>
- cd *yourFolderProject*
- go mod tidy

### Ejecutar la Aplicacion
- Para ejecutar: go run cmd/main.go

El servidor se iniciara en el puerto 8080.

## Endpoints de la API

### Verificacion de Estado
- **GET** `/v1/prestadores/ping`
  - Devuelve una respuesta simple de pong (es un healthcheck)
  - Respuesta: `{"message": "pong"}`

### Afiliados
- **GET** `/v1/prestadores/afiliados`
  - Obtiene la lista de afiliados
  - Respuesta: `{"data": [...], "count": 5, "message": "Afiliados obtenidos exitosamente"}`

### Login
- **POST** `/v1/prestadores/login`
    - Realiza el inicio de sesión validando el CUIT del usuario.
    - Request body: `{"cuit": "20304050607"}`
    - Respuesta exitosa: `{"message": "Login success", "cuit": "20304050607"}`
    - Error por CUIT faltante: `{"error": "El campo 'cuit' es obligatorio"}`
    - Error por request inválido: `{"error": "Formato de request inválido"}`

## Desarrollo

### Estructura del Proyecto
```
prestadores-api/
├── cmd/
│   └── main.go                    # Punto de entrada de la aplicacion
├── internal/
│   └── handler/
│       └── afiliados.go           # Handler para endpoints de afiliados
│       └── login.go               # Handler para endpoints de login
├── go.mod                         # Dependencias del modulo Go
└── README.md                      # Documentacion del proyecto
```

### Agregar Nuevos Endpoints
Todos los endpoints de la API deben agregarse bajo el prefijo `/v1/prestadores/` usando grupos de rutas de Gin.