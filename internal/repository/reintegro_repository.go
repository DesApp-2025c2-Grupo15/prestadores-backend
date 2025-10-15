package repository

import (
	"fmt"
	"prestadores-api/internal/model"
	"sync"
	"time"
)

type ReintegroRepository interface {
	GetAll(estado string, query string, page int, size int, sort string) ([]model.ReintegroListItem, int, error)
	GetByID(id int) (*model.ReintegroDetalle, error)
	Create(req model.CreateReintegroRequest) (*model.ReintegroDetalle, error)
	Update(id int, req model.UpdateReintegroRequest) error
	CambiarEstado(id int, req model.CambioEstadoRequest) (*model.ReintegroDetalle, error)
}

type reintegroRepositoryImpl struct {
	mu         sync.RWMutex
	reintegros map[int]*model.ReintegroDetalle
	nextID     int
}

func NewReintegroRepository() ReintegroRepository {
	repo := &reintegroRepositoryImpl{
		reintegros: make(map[int]*model.ReintegroDetalle),
		nextID:     8801,
	}
	
	repo.initializeDummyData()

	return repo
}

func (r *reintegroRepositoryImpl) initializeDummyData() {
	dummyData := []model.ReintegroDetalle{
		{
			ID:                 8801,
			Tipo:               model.TipoReintegro,
			Estado:             model.EstadoObservado,
			FechaCreacion:      time.Date(2025, 8, 28, 9, 0, 0, 0, time.UTC),
			FechaActualizacion: time.Date(2025, 8, 28, 9, 15, 0, 0, time.UTC),
			Afiliado: model.AfiliadoBasico{
				ID:       45,
				DNI:      "21345633",
				Nombre:   "Daniela",
				Apellido: "Reynoso",
			},
			Prestacion: "Kinesiología",
			Metodo:     "Credito",
			Monto:      40000,
			Historial: []model.HistorialEstado{
				{
					Estado:      model.EstadoRecibido,
					Usuario:     "prestador.202",
					FechaCambio: time.Date(2025, 8, 28, 9, 0, 0, 0, time.UTC),
				},
				{
					Estado:      model.EstadoObservado,
					Usuario:     "prestador.202",
					FechaCambio: time.Date(2025, 8, 28, 9, 15, 0, 0, time.UTC),
					Motivo:      "Ticket ilegible",
				},
			},
		},
		{
			ID:                 8802,
			Tipo:               model.TipoReintegro,
			Estado:             model.EstadoAprobado,
			FechaCreacion:      time.Date(2025, 9, 3, 10, 0, 0, 0, time.UTC),
			FechaActualizacion: time.Date(2025, 9, 3, 12, 0, 0, 0, time.UTC),
			Afiliado: model.AfiliadoBasico{
				ID:       46,
				DNI:      "30123456",
				Nombre:   "Marcos",
				Apellido: "Ledesma",
			},
			Prestacion: "Estudio diagnóstico",
			Metodo:     "Debito",
			Monto:      55000,
			Historial: []model.HistorialEstado{
				{
					Estado:      model.EstadoRecibido,
					Usuario:     "prestador.205",
					FechaCambio: time.Date(2025, 9, 3, 10, 0, 0, 0, time.UTC),
				},
				{
					Estado:      model.EstadoEnAnalisis,
					Usuario:     "prestador.205",
					FechaCambio: time.Date(2025, 9, 3, 10, 30, 0, 0, time.UTC),
				},
				{
					Estado:      model.EstadoAprobado,
					Usuario:     "prestador.205",
					FechaCambio: time.Date(2025, 9, 3, 12, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			ID:                 8803,
			Tipo:               model.TipoReintegro,
			Estado:             model.EstadoRecibido,
			FechaCreacion:      time.Date(2025, 9, 5, 16, 15, 0, 0, time.UTC),
			FechaActualizacion: time.Date(2025, 9, 5, 16, 15, 0, 0, time.UTC),
			Afiliado: model.AfiliadoBasico{
				ID:       47,
				DNI:      "34567890",
				Nombre:   "Lucía",
				Apellido: "Fernández",
			},
			Prestacion: "Consulta clínica",
			Metodo:     "Efectivo",
			Monto:      12000,
			Historial: []model.HistorialEstado{
				{
					Estado:      model.EstadoRecibido,
					Usuario:     "prestador.206",
					FechaCambio: time.Date(2025, 9, 5, 16, 15, 0, 0, time.UTC),
				},
			},
		},
	}

	for _, rgt := range dummyData {
		r.reintegros[rgt.ID] = &rgt
	}
	r.nextID = 8804
}

func (r *reintegroRepositoryImpl) GetAll(estado string, query string, page int, size int, sort string) ([]model.ReintegroListItem, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var items []model.ReintegroListItem
	for _, rgt := range r.reintegros {
		if estado != "" && string(rgt.Estado) != estado {
			continue
		}
		if query != "" {
			// TODO: implementar búsqueda por DNI/nombre/apellido/nro
		}

		item := model.ReintegroListItem{
			ID:                 rgt.ID,
			Tipo:               rgt.Tipo,
			Afiliado:           rgt.Afiliado,
			Estado:             rgt.Estado,
			FechaCreacion:      rgt.FechaCreacion,
			FechaActualizacion: rgt.FechaActualizacion,
			Prestacion:         rgt.Prestacion,
			Metodo:             rgt.Metodo,
			Monto:              rgt.Monto,
		}
		items = append(items, item)
	}

	total := len(items)
	start := page * size
	end := start + size

	if start > total {
		return []model.ReintegroListItem{}, total, nil
	}
	if end > total {
		end = total
	}
	return items[start:end], total, nil
}

func (r *reintegroRepositoryImpl) GetByID(id int) (*model.ReintegroDetalle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rgt, exists := r.reintegros[id]
	if !exists {
		return nil, fmt.Errorf("reintegro no encontrado")
	}
	return rgt, nil
}

func (r *reintegroRepositoryImpl) Create(req model.CreateReintegroRequest) (*model.ReintegroDetalle, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	estadoInicial := req.EstadoInicial
	if estadoInicial == "" {
		estadoInicial = model.EstadoRecibido
	}
	now := time.Now()

	rgt := &model.ReintegroDetalle{
		ID:                 r.nextID,
		Tipo:               model.TipoReintegro,
		Estado:             estadoInicial,
		FechaCreacion:      now,
		FechaActualizacion: now,
		Afiliado: model.AfiliadoBasico{
			ID:       req.AfiliadoID,
			DNI:      "dummy-dni",
			Nombre:   "Dummy",
			Apellido: "Afiliado",
		},
		Prestacion: req.Prestacion,
		Metodo:     req.Metodo,
		Monto:      req.Monto,
		Historial: []model.HistorialEstado{
			{
				Estado:      estadoInicial,
				Usuario:     "sistema",
				FechaCambio: now,
			},
		},
	}

	r.reintegros[r.nextID] = rgt
	r.nextID++

	return rgt, nil
}

func (r *reintegroRepositoryImpl) Update(id int, req model.UpdateReintegroRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	rgt, exists := r.reintegros[id]
	if !exists {
		return fmt.Errorf("reintegro no encontrado")
	}

	if req.Prestacion != "" {
		rgt.Prestacion = req.Prestacion
	}
	if req.Metodo != "" {
		rgt.Metodo = req.Metodo
	}
	if req.Monto != 0 {
		rgt.Monto = req.Monto
	}

	rgt.FechaActualizacion = time.Now()
	return nil
}

func (r *reintegroRepositoryImpl) CambiarEstado(id int, req model.CambioEstadoRequest) (*model.ReintegroDetalle, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rgt, exists := r.reintegros[id]
	if !exists {
		return nil, fmt.Errorf("reintegro no encontrado")
	}

	// Validación mínima de motivo (el service también valida)
	if (req.NuevoEstado == model.EstadoObservado || req.NuevoEstado == model.EstadoRechazado) && req.Motivo == "" {
		return nil, fmt.Errorf("el motivo es obligatorio para estados OBSERVADO y RECHAZADO")
	}

	now := time.Now()
	rgt.Estado = req.NuevoEstado
	rgt.FechaActualizacion = now

	h := model.HistorialEstado{
		Estado:      req.NuevoEstado,
		Usuario:     req.Usuario,
		FechaCambio: now,
		Motivo:      req.Motivo,
	}
	rgt.Historial = append(rgt.Historial, h)

	return rgt, nil
}
