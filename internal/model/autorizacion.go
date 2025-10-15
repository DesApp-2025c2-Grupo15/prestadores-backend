package model

import "time"

// Estado de la autorización
type EstadoAutorizacion string

const (
	EstadoRecibido    EstadoAutorizacion = "RECIBIDO"
	EstadoEnAnalisis  EstadoAutorizacion = "EN_ANALISIS"
	EstadoAprobado    EstadoAutorizacion = "APROBADO"
	EstadoRechazado   EstadoAutorizacion = "RECHAZADO"
	EstadoObservado   EstadoAutorizacion = "OBSERVADO"
)

// TipoSolicitud tipo de solicitud
type TipoSolicitud string

const (
	TipoAutorizacion TipoSolicitud = "AUTORIZACION"
)

// AfiliadoBasico representa los datos básicos del afiliado en la autorización
type AfiliadoBasico struct {
	ID       int    `json:"id"`
	DNI      string `json:"dni"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
}

// AutorizacionListItem representa un item de la lista de autorizaciones
type AutorizacionListItem struct {
	ID                 int                `json:"id"`
	Tipo               TipoSolicitud      `json:"tipo"`
	Afiliado           AfiliadoBasico     `json:"afiliado"`
	Estado             EstadoAutorizacion `json:"estado"`
	FechaCreacion      time.Time          `json:"fechaCreacion"`
	FechaActualizacion time.Time          `json:"fechaActualizacion"`
	Procedimiento      string             `json:"procedimiento"`
	Especialidad       string             `json:"especialidad"`
}

// HistorialEstado representa un cambio de estado en el historial
type HistorialEstado struct {
	Estado      EstadoAutorizacion `json:"estado"`
	Usuario     string             `json:"usuario"`
	FechaCambio time.Time          `json:"fechaCambio"`
	Motivo      string             `json:"motivo,omitempty"`
}

// AutorizacionDetalle representa el detalle completo de una autorización
type AutorizacionDetalle struct {
	ID                 int                `json:"id"`
	Tipo               TipoSolicitud      `json:"tipo"`
	Estado             EstadoAutorizacion `json:"estado"`
	FechaCreacion      time.Time          `json:"fechaCreacion"`
	FechaActualizacion time.Time          `json:"fechaActualizacion"`
	Afiliado           AfiliadoBasico     `json:"afiliado"`
	Procedimiento      string             `json:"procedimiento"`
	Especialidad       string             `json:"especialidad"`
	Historial          []HistorialEstado  `json:"historial"`
}

// CreateAutorizacionRequest representa el request para crear una autorización
type CreateAutorizacionRequest struct {
	AfiliadoID    int                `json:"afiliadoId" binding:"required"`
	Procedimiento string             `json:"procedimiento" binding:"required"`
	Especialidad  string             `json:"especialidad" binding:"required"`
	EstadoInicial EstadoAutorizacion `json:"estadoInicial"`
}

// CreateAutorizacionResponse representa la respuesta al crear una autorización
type CreateAutorizacionResponse struct {
	ID            int                `json:"id"`
	Tipo          TipoSolicitud      `json:"tipo"`
	Estado        EstadoAutorizacion `json:"estado"`
	FechaCreacion time.Time          `json:"fechaCreacion"`
}

// UpdateAutorizacionRequest representa el request para actualizar datos de una autorización
type UpdateAutorizacionRequest struct {
	Procedimiento string `json:"procedimiento,omitempty"`
	Especialidad  string `json:"especialidad,omitempty"`
}

// CambioEstadoRequest representa el request para cambiar el estado de una autorización
type CambioEstadoRequest struct {
	NuevoEstado EstadoAutorizacion `json:"nuevoEstado" binding:"required"`
	Motivo      string             `json:"motivo,omitempty"`
	Usuario     string             `json:"usuario" binding:"required"`
}

// CambioEstadoResponse representa la respuesta al cambiar el estado
type CambioEstadoResponse struct {
	ID                 int                `json:"id"`
	Tipo               TipoSolicitud      `json:"tipo"`
	Estado             EstadoAutorizacion `json:"estado"`
	FechaActualizacion time.Time          `json:"fechaActualizacion"`
}

// PaginatedAutorizacionesResponse representa la respuesta paginada de autorizaciones
type PaginatedAutorizacionesResponse struct {
	Page  int                    `json:"page"`
	Size  int                    `json:"size"`
	Total int                    `json:"total"`
	Items []AutorizacionListItem `json:"items"`
}
