package repository

import (
	"fmt"
	"prestadores-api/internal/model"
	"sync"
	"time"
)

type SituacionRepository interface {
	GetByAfiliado(afiliadoID int) ([]model.Situacion, error)
	GetByAfiliadoGrupo(afiliadoID int) ([]model.IntegranteSituaciones, error)
	Create(req model.CreateSituacionRequest) (*model.Situacion, error)
	Patch(situacionID int, req model.PatchSituacionRequest) error
	CambiarEstado(situacionID int, req model.CambioEstadoSituacionRequest) (*model.Situacion, error)
	Delete(situacionID int) error
}

type situacionRepositoryImpl struct {
	mu          sync.RWMutex
	situaciones map[int]*model.Situacion
	nextID      int

	// Mapa de grupo familiar (mock): titular -> integrantes
	// Para demo: nombres y parentescos fijos por afiliado.
	grupoFamiliar map[int][]struct {
		MiembroID  int
		Nombre     string
		Parentesco string
	}
}

func NewSituacionRepository() SituacionRepository {
	repo := &situacionRepositoryImpl{
		situaciones: make(map[int]*model.Situacion),
		nextID:      7001,
		grupoFamiliar: map[int][]struct {
			MiembroID  int
			Nombre     string
			Parentesco string
		}{
			// Ejemplo: afiliado 22 tiene 2 integrantes
			22: {
				{MiembroID: 2201, Nombre: "Ana Osorio", Parentesco: "Hija"},
				{MiembroID: 2202, Nombre: "Luis Osorio", Parentesco: "Cónyuge"},
			},
			31: {
				{MiembroID: 3101, Nombre: "Pedro Queen", Parentesco: "Hijo"},
			},
		},
	}
	repo.initializeDummyData()
	return repo
}

func (r *situacionRepositoryImpl) initializeDummyData() {
	now := time.Date(2025, 9, 20, 10, 0, 0, 0, time.UTC)

	// Titular afiliado 22
	r.situaciones[7001] = &model.Situacion{
		ID:                 7001,
		AfiliadoID:         22,
		MiembroID:          nil, // titular
		Descripcion:        "Lumbalgia",
		FechaInicio:        "2025-09-01",
		FechaFin:           nil,
		Estado:             model.EstadoSituacionActiva,
		FechaCreacion:      now.Add(-24 * time.Hour),
		FechaActualizacion: now.Add(-23 * time.Hour),
	}

	// Hija (miembro 2201)
	fin := "2025-09-10"
	r.situaciones[7002] = &model.Situacion{
		ID:                 7002,
		AfiliadoID:         22,
		MiembroID:          intPtr(2201),
		Descripcion:        "Contractura cervical",
		FechaInicio:        "2025-08-15",
		FechaFin:           &fin,
		Estado:             model.EstadoSituacionBaja,
		FechaCreacion:      now.Add(-36 * time.Hour),
		FechaActualizacion: now.Add(-30 * time.Hour),
	}

	// Cónyuge (miembro 2202)
	r.situaciones[7003] = &model.Situacion{
		ID:                 7003,
		AfiliadoID:         22,
		MiembroID:          intPtr(2202),
		Descripcion:        "Tendinitis",
		FechaInicio:        "2025-09-05",
		FechaFin:           nil,
		Estado:             model.EstadoSituacionActiva,
		FechaCreacion:      now.Add(-72 * time.Hour),
		FechaActualizacion: now.Add(-12 * time.Hour),
	}

	// Titular afiliado 31
	r.situaciones[7004] = &model.Situacion{
		ID:                 7004,
		AfiliadoID:         31,
		MiembroID:          nil,
		Descripcion:        "Asma leve",
		FechaInicio:        "2025-08-20",
		FechaFin:           nil,
		Estado:             model.EstadoSituacionActiva,
		FechaCreacion:      now.Add(-100 * time.Hour),
		FechaActualizacion: now.Add(-48 * time.Hour),
	}

	r.nextID = 7005
}

func (r *situacionRepositoryImpl) GetByAfiliado(afiliadoID int) ([]model.Situacion, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]model.Situacion, 0)
	for _, s := range r.situaciones {
		if s.AfiliadoID == afiliadoID {
			out = append(out, *s)
		}
	}
	return out, nil
}

func (r *situacionRepositoryImpl) GetByAfiliadoGrupo(afiliadoID int) ([]model.IntegranteSituaciones, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 1) Titular
	titular := model.IntegranteSituaciones{
		MiembroID:   afiliadoID,
		Nombre:      mockNombreTitular(afiliadoID),
		Parentesco:  "Titular",
		Situaciones: []model.Situacion{},
	}
	for _, s := range r.situaciones {
		if s.AfiliadoID == afiliadoID && s.MiembroID == nil {
			titular.Situaciones = append(titular.Situaciones, *s)
		}
	}

	// 2) Integrantes
	integrantesMeta := r.grupoFamiliar[afiliadoID]
	integrantes := make([]model.IntegranteSituaciones, 0, len(integrantesMeta)+1)
	integrantes = append(integrantes, titular)

	for _, meta := range integrantesMeta {
		miembro := model.IntegranteSituaciones{
			MiembroID:   meta.MiembroID,
			Nombre:      meta.Nombre,
			Parentesco:  meta.Parentesco,
			Situaciones: []model.Situacion{},
		}
		for _, s := range r.situaciones {
			if s.AfiliadoID == afiliadoID && s.MiembroID != nil && *s.MiembroID == meta.MiembroID {
				miembro.Situaciones = append(miembro.Situaciones, *s)
			}
		}
		integrantes = append(integrantes, miembro)
	}

	return integrantes, nil
}

func (r *situacionRepositoryImpl) Create(req model.CreateSituacionRequest) (*model.Situacion, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UTC()
	s := &model.Situacion{
		ID:                 r.nextID,
		AfiliadoID:         req.AfiliadoID,
		MiembroID:          req.MiembroID, // nil => titular
		Descripcion:        req.Descripcion,
		FechaInicio:        req.FechaInicio,
		FechaFin:           req.FechaFin,
		Estado:             model.EstadoSituacionActiva, // alta -> ACTIVA
		FechaCreacion:      now,
		FechaActualizacion: now,
	}

	r.situaciones[r.nextID] = s
	r.nextID++

	return s, nil
}

func (r *situacionRepositoryImpl) Patch(situacionID int, req model.PatchSituacionRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	s, ok := r.situaciones[situacionID]
	if !ok {
		return fmt.Errorf("situación no encontrada")
	}

	if req.Descripcion != nil {
		s.Descripcion = *req.Descripcion
	}
	if req.FechaInicio != nil {
		s.FechaInicio = *req.FechaInicio
	}
	if req.FechaFin != nil {
		// Si querés permitir “limpiar” la fecha, podés definir una convención.
		// Por ahora, si viene "", la dejamos en nil.
		if *req.FechaFin == "" {
			s.FechaFin = nil
		} else {
			s.FechaFin = req.FechaFin
		}
	}

	s.FechaActualizacion = time.Now().UTC()
	return nil
}

func (r *situacionRepositoryImpl) CambiarEstado(situacionID int, req model.CambioEstadoSituacionRequest) (*model.Situacion, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	s, ok := r.situaciones[situacionID]
	if !ok {
		return nil, fmt.Errorf("situación no encontrada")
	}

	// Reglas mínimas: si pasa a BAJA y no tiene fechaFin, la seteamos a hoy (convenio mock)
	if req.Estado == model.EstadoSituacionBaja && (s.FechaFin == nil || *s.FechaFin == "") {
		hoy := time.Now().UTC().Format("2006-01-02")
		s.FechaFin = &hoy
	}

	s.Estado = req.Estado
	s.FechaActualizacion = time.Now().UTC()
	return s, nil
}

func (r *situacionRepositoryImpl) Delete(situacionID int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.situaciones[situacionID]; !ok {
		return fmt.Errorf("situación no encontrada")
	}
	delete(r.situaciones, situacionID)
	return nil
}

// ===== Helpers =====

func intPtr(v int) *int { return &v }

// Nombre “dummy” para titulares según afiliadoID (mock)
func mockNombreTitular(afiliadoID int) string {
	switch afiliadoID {
	case 22:
		return "Miguel Osorio"
	case 31:
		return "David Queen"
	default:
		return fmt.Sprintf("Titular %d", afiliadoID)
	}
}
