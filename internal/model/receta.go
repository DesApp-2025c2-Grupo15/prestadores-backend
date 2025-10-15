package model

import "time"

type EstadoReceta string

const (
	RecetaEstadoRecibido   EstadoReceta = "RECIBIDO"
	RecetaEstadoEnAnalisis EstadoReceta = "EN_ANALISIS"
	RecetaEstadoAprobado   EstadoReceta = "APROBADO"
	RecetaEstadoRechazado  EstadoReceta = "RECHAZADO"
	RecetaEstadoObservado  EstadoReceta = "OBSERVADO"
)

const (
	TipoReceta TipoSolicitud = "RECETA"
)

type RecetaListItem struct {
	ID                 int          `json:"id"`
	Tipo               TipoSolicitud `json:"tipo"`
	Afiliado           AfiliadoBasico `json:"afiliado"`
	Estado             EstadoReceta `json:"estado"`
	FechaCreacion      time.Time    `json:"fechaCreacion"`
	FechaActualizacion time.Time    `json:"fechaActualizacion"`
	Medicamento        string       `json:"medicamento"`
	Dosis              string       `json:"dosis"`
}

type HistorialEstadoReceta struct {
	Estado      EstadoReceta `json:"estado"`
	Usuario     string       `json:"usuario"`
	FechaCambio time.Time    `json:"fechaCambio"`
	Motivo      string       `json:"motivo,omitempty"`
}

type RecetaDetalle struct {
	ID                 int                     `json:"id"`
	Tipo               TipoSolicitud           `json:"tipo"`
	Estado             EstadoReceta            `json:"estado"`
	FechaCreacion      time.Time               `json:"fechaCreacion"`
	FechaActualizacion time.Time               `json:"fechaActualizacion"`
	Afiliado           AfiliadoBasico          `json:"afiliado"`
	Medicamento        string                  `json:"medicamento"`
	Dosis              string                  `json:"dosis"`
	Historial          []HistorialEstadoReceta `json:"historial"`
}

type CreateRecetaRequest struct {
	AfiliadoID    int          `json:"afiliadoId" binding:"required"`
	Medicamento   string       `json:"medicamento" binding:"required"`
	Dosis         string       `json:"dosis" binding:"required"`
	EstadoInicial EstadoReceta `json:"estadoInicial"`
}

type CreateRecetaResponse struct {
	ID            int           `json:"id"`
	Tipo          TipoSolicitud `json:"tipo"`
	Estado        EstadoReceta  `json:"estado"`
	FechaCreacion time.Time     `json:"fechaCreacion"`
}

type UpdateRecetaRequest struct {
	Medicamento string `json:"medicamento,omitempty"`
	Dosis       string `json:"dosis,omitempty"`
}

type CambioEstadoRecetaRequest struct {
	NuevoEstado EstadoReceta `json:"nuevoEstado" binding:"required"`
	Motivo      string       `json:"motivo,omitempty"`
	Usuario     string       `json:"usuario" binding:"required"`
}

type CambioEstadoRecetaResponse struct {
	ID                 int           `json:"id"`
	Tipo               TipoSolicitud `json:"tipo"`
	Estado             EstadoReceta  `json:"estado"`
	FechaActualizacion time.Time     `json:"fechaActualizacion"`
}

type PaginatedRecetasResponse struct {
	Page  int              `json:"page"`
	Size  int              `json:"size"`
	Total int              `json:"total"`
	Items []RecetaListItem `json:"items"`
}
