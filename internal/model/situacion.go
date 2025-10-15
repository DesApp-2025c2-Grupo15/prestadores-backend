package model

import "time"

// ===== Enums / Tipos =====

// EstadoSituacion representa el estado de una situación terapéutica
type EstadoSituacion string

const (
	EstadoSituacionActiva EstadoSituacion = "ACTIVA"
	EstadoSituacionBaja   EstadoSituacion = "BAJA"
	EstadoSituacionAlta   EstadoSituacion = "ALTA"
)

// ===== Modelos base =====

// Situacion representa una situación terapéutica de un afiliado o miembro del grupo
type Situacion struct {
	ID                 int             `json:"id"`
	AfiliadoID         int             `json:"afiliadoId"`          // titular del grupo
	MiembroID          *int            `json:"miembroId,omitempty"` // null si es el titular
	Descripcion        string          `json:"descripcion"`
	FechaInicio        string          `json:"fechaInicio"`        // ISO-8601 (yyyy-mm-dd) para simplificar mock
	FechaFin           *string         `json:"fechaFin,omitempty"` // ISO-8601 o null
	Estado             EstadoSituacion `json:"estado"`             // ACTIVA | BAJA | ALTA
	FechaCreacion      time.Time       `json:"fechaCreacion"`
	FechaActualizacion time.Time       `json:"fechaActualizacion"`
}

// IntegranteSituaciones agrupa situaciones por integrante (para vista de grupo familiar)
type IntegranteSituaciones struct {
	MiembroID   int         `json:"miembroId"`
	Nombre      string      `json:"nombre"`
	Parentesco  string      `json:"parentesco"` // "Titular", "Hijo/a", "Cónyuge", etc.
	Situaciones []Situacion `json:"situaciones"`
}

// ===== Respuestas de lectura =====

// SituacionesAfiliadoResponse para GET /afiliados/:afiliadoId/situaciones (scope por defecto = titular)
type SituacionesAfiliadoResponse struct {
	AfiliadoID int         `json:"afiliadoId"`
	Items      []Situacion `json:"items"`
}

// SituacionesGrupoResponse para GET /afiliados/:afiliadoId/situaciones?scope=grupo
type SituacionesGrupoResponse struct {
	AfiliadoID  int                     `json:"afiliadoId"`
	Integrantes []IntegranteSituaciones `json:"integrantes"`
}

// ===== Requests / Responses de escritura =====

// CreateSituacionRequest para POST /afiliados/:afiliadoId/situaciones
// (Nota: AfiliadoID lo setea el handler desde la ruta)
type CreateSituacionRequest struct {
	AfiliadoID  int     `json:"afiliadoId"`          // lo completa el handler
	MiembroID   *int    `json:"miembroId,omitempty"` // opcional (si es del grupo)
	Descripcion string  `json:"descripcion" binding:"required"`
	FechaInicio string  `json:"fechaInicio" binding:"required"` // yyyy-mm-dd
	FechaFin    *string `json:"fechaFin,omitempty"`             // opcional
}

// CreateSituacionResponse al crear una situación
type CreateSituacionResponse struct {
	ID            int             `json:"id"`
	AfiliadoID    int             `json:"afiliadoId"`
	MiembroID     *int            `json:"miembroId,omitempty"`
	Estado        EstadoSituacion `json:"estado"`
	FechaCreacion time.Time       `json:"fechaCreacion"`
}

// PatchSituacionRequest para PATCH /afiliados/:afiliadoId/situaciones/:situacionId
// Usado típicamente para setear/modificar fechaFin u otros campos editables.
type PatchSituacionRequest struct {
	Descripcion *string `json:"descripcion,omitempty"`
	FechaInicio *string `json:"fechaInicio,omitempty"` // yyyy-mm-dd
	FechaFin    *string `json:"fechaFin,omitempty"`    // yyyy-mm-dd (null -> enviar como ausencia de campo)
}

// CambioEstadoSituacionRequest para PATCH /afiliados/:afiliadoId/situaciones/:situacionId/estado
// Permite pasar a BAJA (baja lógica) o re-activar (ACTIVA)
type CambioEstadoSituacionRequest struct {
	Estado  EstadoSituacion `json:"estado" binding:"required"` // ACTIVA | BAJA
	Motivo  string          `json:"motivo,omitempty"`          // opcional, según reglas que definan
	Usuario string          `json:"usuario" binding:"required"`
}

// CambioEstadoSituacionResponse respuesta al cambiar estado
type CambioEstadoSituacionResponse struct {
	ID                 int             `json:"id"`
	Estado             EstadoSituacion `json:"estado"`
	FechaActualizacion time.Time       `json:"fechaActualizacion"`
}
