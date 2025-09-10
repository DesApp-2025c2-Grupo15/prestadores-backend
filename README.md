# prestadores-api

Una API REST para gestionar prestadores de servicios construida con Go y Gin.

## Comenzar

### Prerrequisitos
- Go 1.25+

### Instalacion
Clonar el repositorio e instalar dependencias:
- git clone <repository-url>
- cd prestadores-api
- go mod tidy

### Ejecutar la Aplicacion
Para ejecutar en modo desarrollo: go run cmd/main.go
Para compilar y ejecutar: go build -o prestadores-api cmd/main.go && ./prestadores-api

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

## Desarrollo

### Estructura del Proyecto
```
prestadores-api/
├── cmd/
│   └── main.go                    # Punto de entrada de la aplicacion
├── internal/
│   └── handler/
│       └── afiliados.go           # Handler para endpoints de afiliados
├── go.mod                         # Dependencias del modulo Go
└── README.md                      # Documentacion del proyecto
```

### Agregar Nuevos Endpoints
Todos los endpoints de la API deben agregarse bajo el prefijo `/v1/prestadores/` usando grupos de rutas de Gin.