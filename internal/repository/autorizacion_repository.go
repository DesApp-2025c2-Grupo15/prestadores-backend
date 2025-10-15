package repository

import (
	"fmt"
	"prestadores-api/internal/model"
	"sync"
	"time"
)

type AutorizacionRepository interface {
	GetAll(estado string, query string, page int, size int, sort string) ([]model.AutorizacionListItem, int, error)
	GetByID(id int) (*model.AutorizacionDetalle, error)
	Create(req model.CreateAutorizacionRequest) (*model.AutorizacionDetalle, error)
	Update(id int, req model.UpdateAutorizacionRequest) error
	CambiarEstado(id int, req model.CambioEstadoRequest) (*model.AutorizacionDetalle, error)
}

type autorizacionRepositoryImpl struct {
	mu             sync.RWMutex
	autorizaciones map[int]*model.AutorizacionDetalle
	nextID         int
}

func NewAutorizacionRepository() AutorizacionRepository {
	repo := &autorizacionRepositoryImpl{
		autorizaciones: make(map[int]*model.AutorizacionDetalle),
		nextID:         12001,
	}

	repo.initializeDummyData()

	return repo
}

func (r *autorizacionRepositoryImpl) initializeDummyData() {
	dummyData := []model.AutorizacionDetalle{
		{
			ID:                 12001,
			Tipo:               model.TipoAutorizacion,
			Estado:             model.EstadoRechazado,
			FechaCreacion:      time.Date(2025, 9, 2, 10, 0, 0, 0, time.UTC),
			FechaActualizacion: time.Date(2025, 9, 2, 10, 30, 0, 0, time.UTC),
			Afiliado: model.AfiliadoBasico{
				ID:       31,
				DNI:      "45678089",
				Nombre:   "David",
				Apellido: "Queen",
			},
			Procedimiento: "Consulta de control",
			Especialidad:  "Clínica Médica",
			Historial: []model.HistorialEstado{
				{
					Estado:      model.EstadoRecibido,
					Usuario:     "prestador.201",
					FechaCambio: time.Date(2025, 9, 2, 10, 0, 0, 0, time.UTC),
				},
				{
					Estado:      model.EstadoEnAnalisis,
					Usuario:     "prestador.201",
					FechaCambio: time.Date(2025, 9, 2, 10, 15, 0, 0, time.UTC),
				},
				{
					Estado:      model.EstadoRechazado,
					Usuario:     "prestador.201",
					FechaCambio: time.Date(2025, 9, 2, 10, 30, 0, 0, time.UTC),
					Motivo:      "Falta documentación",
				},
			},
		},
		{
			ID:                 12002,
			Tipo:               model.TipoAutorizacion,
			Estado:             model.EstadoAprobado,
			FechaCreacion:      time.Date(2025, 9, 3, 9, 0, 0, 0, time.UTC),
			FechaActualizacion: time.Date(2025, 9, 3, 11, 0, 0, 0, time.UTC),
			Afiliado: model.AfiliadoBasico{
				ID:       32,
				DNI:      "38567123",
				Nombre:   "Laura",
				Apellido: "García",
			},
			Procedimiento: "Radiografía de tórax",
			Especialidad:  "Diagnóstico por Imágenes",
			Historial: []model.HistorialEstado{
				{
					Estado:      model.EstadoRecibido,
					Usuario:     "prestador.202",
					FechaCambio: time.Date(2025, 9, 3, 9, 0, 0, 0, time.UTC),
				},
				{
					Estado:      model.EstadoAprobado,
					Usuario:     "prestador.202",
					FechaCambio: time.Date(2025, 9, 3, 11, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			ID:                 12003,
			Tipo:               model.TipoAutorizacion,
			Estado:             model.EstadoEnAnalisis,
			FechaCreacion:      time.Date(2025, 9, 4, 14, 30, 0, 0, time.UTC),
			FechaActualizacion: time.Date(2025, 9, 4, 15, 0, 0, 0, time.UTC),
			Afiliado: model.AfiliadoBasico{
				ID:       33,
				DNI:      "42123456",
				Nombre:   "Carlos",
				Apellido: "Martínez",
			},
			Procedimiento: "Consulta cardiológica",
			Especialidad:  "Cardiología",
			Historial: []model.HistorialEstado{
				{
					Estado:      model.EstadoRecibido,
					Usuario:     "prestador.203",
					FechaCambio: time.Date(2025, 9, 4, 14, 30, 0, 0, time.UTC),
				},
				{
					Estado:      model.EstadoEnAnalisis,
					Usuario:     "prestador.203",
					FechaCambio: time.Date(2025, 9, 4, 15, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	for _, aut := range dummyData {
		r.autorizaciones[aut.ID] = &aut
	}

	r.nextID = 12004
}

func (r *autorizacionRepositoryImpl) GetAll(estado string, query string, page int, size int, sort string) ([]model.AutorizacionListItem, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var items []model.AutorizacionListItem
	for _, aut := range r.autorizaciones {
		if estado != "" && string(aut.Estado) != estado {
			continue
		}

		if query != "" {
			// TODO: implementar búsqueda
		}

		item := model.AutorizacionListItem{
			ID:                 aut.ID,
			Tipo:               aut.Tipo,
			Afiliado:           aut.Afiliado,
			Estado:             aut.Estado,
			FechaCreacion:      aut.FechaCreacion,
			FechaActualizacion: aut.FechaActualizacion,
			Procedimiento:      aut.Procedimiento,
			Especialidad:       aut.Especialidad,
		}
		items = append(items, item)
	}

	total := len(items)

	start := page * size
	end := start + size

	if start > total {
		return []model.AutorizacionListItem{}, total, nil
	}

	if end > total {
		end = total
	}

	return items[start:end], total, nil
}

func (r *autorizacionRepositoryImpl) GetByID(id int) (*model.AutorizacionDetalle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	aut, exists := r.autorizaciones[id]
	if !exists {
		return nil, fmt.Errorf("autorización no encontrada")
	}

	return aut, nil
}

func (r *autorizacionRepositoryImpl) Create(req model.CreateAutorizacionRequest) (*model.AutorizacionDetalle, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	estadoInicial := req.EstadoInicial
	if estadoInicial == "" {
		estadoInicial = model.EstadoRecibido
	}

	now := time.Now()

	aut := &model.AutorizacionDetalle{
		ID:                 r.nextID,
		Tipo:               model.TipoAutorizacion,
		Estado:             estadoInicial,
		FechaCreacion:      now,
		FechaActualizacion: now,
		Afiliado: model.AfiliadoBasico{
			ID:       req.AfiliadoID,
			DNI:      "dummy-dni",
			Nombre:   "Dummy",
			Apellido: "Afiliado",
		},
		Procedimiento: req.Procedimiento,
		Especialidad:  req.Especialidad,
		Historial: []model.HistorialEstado{
			{
				Estado:      estadoInicial,
				Usuario:     "sistema",
				FechaCambio: now,
			},
		},
	}

	r.autorizaciones[r.nextID] = aut
	r.nextID++

	return aut, nil
}

func (r *autorizacionRepositoryImpl) Update(id int, req model.UpdateAutorizacionRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	aut, exists := r.autorizaciones[id]
	if !exists {
		return fmt.Errorf("autorización no encontrada")
	}

	if req.Procedimiento != "" {
		aut.Procedimiento = req.Procedimiento
	}

	if req.Especialidad != "" {
		aut.Especialidad = req.Especialidad
	}

	aut.FechaActualizacion = time.Now()

	return nil
}

func (r *autorizacionRepositoryImpl) CambiarEstado(id int, req model.CambioEstadoRequest) (*model.AutorizacionDetalle, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	aut, exists := r.autorizaciones[id]
	if !exists {
		return nil, fmt.Errorf("autorización no encontrada")
	}

	if (req.NuevoEstado == model.EstadoObservado || req.NuevoEstado == model.EstadoRechazado) && req.Motivo == "" {
		return nil, fmt.Errorf("el motivo es obligatorio para estados OBSERVADO y RECHAZADO")
	}

	now := time.Now()

	aut.Estado = req.NuevoEstado
	aut.FechaActualizacion = now

	historial := model.HistorialEstado{
		Estado:      req.NuevoEstado,
		Usuario:     req.Usuario,
		FechaCambio: now,
		Motivo:      req.Motivo,
	}

	aut.Historial = append(aut.Historial, historial)

	return aut, nil
}
