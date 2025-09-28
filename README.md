# prestadores-api

Una API REST para gestionar prestadores de servicios construida con Go y Gin.

## Comenzar

### Prerrequisitos
- Go 1.25+

### Instalación
Clonar el repositorio e instalar dependencias:
git clone <repository-url>
cd <yourFolderProject>
go mod tidy

### Ejecutar la aplicación
go run cmd/main.go
El servidor se iniciará en el puerto 8080.

## Endpoints de la API

### Verificación de estado
GET /v1/prestadores/ping
Devuelve {"message":"pong"} (healthcheck).

### Afiliados
GET /v1/prestadores/afiliados
Obtiene la lista de afiliados para la tabla.
Respuesta (array):
[
    { "id": 1, "dni": "43521489", "nombre": "María", "apellido": "Candia", "planMedico": "Sancor Salud", "titular": true },
    { "id": 2, "dni": "53521489", "nombre": "Stella", "apellido": "Candia", "planMedico": "Sancor Salud", "titular": false }
]

GET /v1/prestadores/afiliados/:id
Obtiene la información detallada de un afiliado, incluyendo su grupo familiar.
Respuesta ejemplo:
{
    "id": 2,
    "nroAfiliado": "15121231523",
    "dni": "53521489",
    "nombre": "Stella",
    "apellido": "Candia",
    "planMedico": "Sancor Salud",
    "titular": false,
    "grupoFamiliar": [
        {
        "id": 1,
        "dni": "43521489",
        "nombre": "María",
        "apellido": "Candia",
        "planMedico": "Sancor Salud"
        }
    ]
}

GET /v1/prestadores/afiliados/:id/historia-clinica
Devuelve la historia clínica del afiliado (lista de turnos con sus notas).
Query opcional: ?prestadorId=45 → filtra las notas por ese prestador.
Respuesta ejemplo:
{
    "afiliadoId": 1,
    "page": 0,
    "size": 20,
    "total": 2,
    "turnos": [
        {
            "id": 500,
            "fecha": "2025-09-20T10:00:00Z",
            "especialidad": "Clínica",
            "estado": "RESERVADO",
            "notas": [
                { "id": 10, "fecha": "2025-09-20T10:30:00Z", "prestadorId": 45, "texto": "Control general" }
            ]
        },
        {
            "id": 501,
            "fecha": "2025-09-25T15:00:00Z",
            "especialidad": "Kinesiología",
            "estado": "ATENDIDO",
            "notas": [
                { "id": 12, "fecha": "2025-09-25T15:45:00Z", "prestadorId": 55, "texto": "Ejercicios domiciliarios" }
            ]
        }
    ]
}

### Login
POST /v1/prestadores/login
Realiza el inicio de sesión validando el CUIT del usuario.
Body:
{ 
    "cuit": "20304050607" 
}
Respuestas:
Éxito: {"message": "Login success", "cuit": "20304050607"}
Error por CUIT faltante: {"error": "El campo 'cuit' es obligatorio"}
Error por request inválido: {"error": "Formato de request inválido"}

## Desarrollo

### Estructura del proyecto
prestadores-api/
├── cmd/
│   └── main.go                      # Punto de entrada de la aplicación
├── internal/
│   └── handler/
│       ├── login.go                 # Handler de login
│       └── afiliados/               # Módulo Afiliados
│           ├── afiliados.go         # Lista y detalle de afiliados
│           └── historia_clinica.go  # Historia clínica: turnos + notas (mock)
├── go.mod                           # Dependencias del módulo Go
└── README.md                        # Documentación del proyecto

### Agregar nuevos endpoints
Todos los endpoints deben colgar del prefijo /v1/prestadores/ usando grupos de rutas de Gin.