package model

import "time"

// Agregamos el tipo de solicitud REINTEGRO al enum existente
const (
	TipoReintegro TipoSolicitud = "REINTEGRO"
)

// ReintegroListItem representa un item de la lista de reintegros
type ReintegroListItem struct {
	ID                 int                `json:"id"`
	Tipo               TipoSolicitud      `json:"tipo"` // "REINTEGRO"
	Afiliado           AfiliadoBasico     `json:"afiliado"`
	Estado             EstadoAutorizacion `json:"estado"`
	FechaCreacion      time.Time          `json:"fechaCreacion"`
	FechaActualizacion time.Time          `json:"fechaActualizacion"`
	Prestacion         string             `json:"prestacion"`
	Metodo             string             `json:"metodo"` // Efectivo | Debito | Credito (mock)
	Monto              float64            `json:"monto"`
}

// ReintegroDetalle representa el detalle completo de un reintegro
type ReintegroDetalle struct {
	ID                 int                `json:"id"`
	Tipo               TipoSolicitud      `json:"tipo"` // "REINTEGRO"
	Estado             EstadoAutorizacion `json:"estado"`
	FechaCreacion      time.Time          `json:"fechaCreacion"`
	FechaActualizacion time.Time          `json:"fechaActualizacion"`
	Afiliado           AfiliadoBasico     `json:"afiliado"`
	Prestacion         string             `json:"prestacion"`
	Metodo             string             `json:"metodo"`
	Monto              float64            `json:"monto"`
	Historial          []HistorialEstado  `json:"historial"`
}

// CreateReintegroRequest representa el request para crear un reintegro
type CreateReintegroRequest struct {
	AfiliadoID    int                `json:"afiliadoId" binding:"required"`
	Prestacion    string             `json:"prestacion" binding:"required"`
	Metodo        string             `json:"metodo" binding:"required"`
	Monto         float64            `json:"monto" binding:"required"`
	EstadoInicial EstadoAutorizacion `json:"estadoInicial"`
}

// CreateReintegroResponse representa la respuesta al crear un reintegro
type CreateReintegroResponse struct {
	ID            int                `json:"id"`
	Tipo          TipoSolicitud      `json:"tipo"` // "REINTEGRO"
	Estado        EstadoAutorizacion `json:"estado"`
	FechaCreacion time.Time          `json:"fechaCreacion"`
}

// UpdateReintegroRequest representa el request para actualizar datos de un reintegro
type UpdateReintegroRequest struct {
	Prestacion string  `json:"prestacion,omitempty"`
	Metodo     string  `json:"metodo,omitempty"`
	Monto      float64 `json:"monto,omitempty"`
}

// PaginatedReintegrosResponse representa la respuesta paginada de reintegros
type PaginatedReintegrosResponse struct {
	Page  int                 `json:"page"`
	Size  int                 `json:"size"`
	Total int                 `json:"total"`
	Items []ReintegroListItem `json:"items"`
}
