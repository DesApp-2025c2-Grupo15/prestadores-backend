package repository

import (
	"fmt"
	"prestadores-api/internal/model"
	"sync"
	"time"
)

type RecetaRepository interface {
	GetAll(estado string, query string, page int, size int, sort string) ([]model.RecetaListItem, int, error)
	GetByID(id int) (*model.RecetaDetalle, error)
	Create(req model.CreateRecetaRequest) (*model.RecetaDetalle, error)
	Update(id int, req model.UpdateRecetaRequest) error
	CambiarEstado(id int, req model.CambioEstadoRecetaRequest) (*model.RecetaDetalle, error)
}

type recetaRepositoryImpl struct {
	mu      sync.RWMutex
	recetas map[int]*model.RecetaDetalle
	nextID  int
}

func NewRecetaRepository() RecetaRepository {
	repo := &recetaRepositoryImpl{
		recetas: make(map[int]*model.RecetaDetalle),
		nextID:  9350,
	}

	repo.initializeDummyData()

	return repo
}

func (r *recetaRepositoryImpl) initializeDummyData() {
	dummyData := []model.RecetaDetalle{
		{
			ID:                 9350,
			Tipo:               model.TipoReceta,
			Estado:             model.RecetaEstadoAprobado,
			FechaCreacion:      time.Date(2025, 9, 11, 11, 10, 0, 0, time.UTC),
			FechaActualizacion: time.Date(2025, 9, 11, 12, 0, 0, 0, time.UTC),
			Afiliado: model.AfiliadoBasico{
				ID:       22,
				DNI:      "32654708",
				Nombre:   "Miguel",
				Apellido: "Osorio",
			},
			Medicamento: "Amoxicilina 500mg",
			Dosis:       "1 cap. c/8h x 7d",
			Historial: []model.HistorialEstadoReceta{
				{
					Estado:      model.RecetaEstadoRecibido,
					Usuario:     "prestador.101",
					FechaCambio: time.Date(2025, 9, 11, 11, 10, 0, 0, time.UTC),
				},
				{
					Estado:      model.RecetaEstadoEnAnalisis,
					Usuario:     "prestador.101",
					FechaCambio: time.Date(2025, 9, 11, 11, 30, 0, 0, time.UTC),
				},
				{
					Estado:      model.RecetaEstadoAprobado,
					Usuario:     "prestador.101",
					FechaCambio: time.Date(2025, 9, 11, 12, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			ID:                 9351,
			Tipo:               model.TipoReceta,
			Estado:             model.RecetaEstadoRechazado,
			FechaCreacion:      time.Date(2025, 9, 12, 9, 0, 0, 0, time.UTC),
			FechaActualizacion: time.Date(2025, 9, 12, 10, 30, 0, 0, time.UTC),
			Afiliado: model.AfiliadoBasico{
				ID:       23,
				DNI:      "28456123",
				Nombre:   "Ana",
				Apellido: "Fernández",
			},
			Medicamento: "Ibuprofeno 600mg",
			Dosis:       "1 comp. c/8h",
			Historial: []model.HistorialEstadoReceta{
				{
					Estado:      model.RecetaEstadoRecibido,
					Usuario:     "prestador.102",
					FechaCambio: time.Date(2025, 9, 12, 9, 0, 0, 0, time.UTC),
				},
				{
					Estado:      model.RecetaEstadoRechazado,
					Usuario:     "prestador.102",
					FechaCambio: time.Date(2025, 9, 12, 10, 30, 0, 0, time.UTC),
					Motivo:      "Medicamento no cubierto por el plan",
				},
			},
		},
		{
			ID:                 9352,
			Tipo:               model.TipoReceta,
			Estado:             model.RecetaEstadoEnAnalisis,
			FechaCreacion:      time.Date(2025, 9, 13, 14, 0, 0, 0, time.UTC),
			FechaActualizacion: time.Date(2025, 9, 13, 14, 30, 0, 0, time.UTC),
			Afiliado: model.AfiliadoBasico{
				ID:       24,
				DNI:      "35789456",
				Nombre:   "Roberto",
				Apellido: "Díaz",
			},
			Medicamento: "Omeprazol 20mg",
			Dosis:       "1 cap. c/12h x 30d",
			Historial: []model.HistorialEstadoReceta{
				{
					Estado:      model.RecetaEstadoRecibido,
					Usuario:     "prestador.103",
					FechaCambio: time.Date(2025, 9, 13, 14, 0, 0, 0, time.UTC),
				},
				{
					Estado:      model.RecetaEstadoEnAnalisis,
					Usuario:     "prestador.103",
					FechaCambio: time.Date(2025, 9, 13, 14, 30, 0, 0, time.UTC),
				},
			},
		},
	}

	for _, rec := range dummyData {
		r.recetas[rec.ID] = &rec
	}

	r.nextID = 9353
}

func (r *recetaRepositoryImpl) GetAll(estado string, query string, page int, size int, sort string) ([]model.RecetaListItem, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var items []model.RecetaListItem
	for _, rec := range r.recetas {
		if estado != "" && string(rec.Estado) != estado {
			continue
		}

		if query != "" {
			// TODO: implementar búsqueda
		}

		item := model.RecetaListItem{
			ID:                 rec.ID,
			Tipo:               rec.Tipo,
			Afiliado:           rec.Afiliado,
			Estado:             rec.Estado,
			FechaCreacion:      rec.FechaCreacion,
			FechaActualizacion: rec.FechaActualizacion,
			Medicamento:        rec.Medicamento,
			Dosis:              rec.Dosis,
		}
		items = append(items, item)
	}

	total := len(items)

	start := page * size
	end := start + size

	if start > total {
		return []model.RecetaListItem{}, total, nil
	}

	if end > total {
		end = total
	}

	return items[start:end], total, nil
}

func (r *recetaRepositoryImpl) GetByID(id int) (*model.RecetaDetalle, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rec, exists := r.recetas[id]
	if !exists {
		return nil, fmt.Errorf("receta no encontrada")
	}

	return rec, nil
}

func (r *recetaRepositoryImpl) Create(req model.CreateRecetaRequest) (*model.RecetaDetalle, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	estadoInicial := req.EstadoInicial
	if estadoInicial == "" {
		estadoInicial = model.RecetaEstadoRecibido
	}

	now := time.Now()

	rec := &model.RecetaDetalle{
		ID:                 r.nextID,
		Tipo:               model.TipoReceta,
		Estado:             estadoInicial,
		FechaCreacion:      now,
		FechaActualizacion: now,
		Afiliado: model.AfiliadoBasico{
			ID:       req.AfiliadoID,
			DNI:      "dummy-dni",
			Nombre:   "Dummy",
			Apellido: "Afiliado",
		},
		Medicamento: req.Medicamento,
		Dosis:       req.Dosis,
		Historial: []model.HistorialEstadoReceta{
			{
				Estado:      estadoInicial,
				Usuario:     "sistema",
				FechaCambio: now,
			},
		},
	}

	r.recetas[r.nextID] = rec
	r.nextID++

	return rec, nil
}

func (r *recetaRepositoryImpl) Update(id int, req model.UpdateRecetaRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	rec, exists := r.recetas[id]
	if !exists {
		return fmt.Errorf("receta no encontrada")
	}

	if req.Medicamento != "" {
		rec.Medicamento = req.Medicamento
	}

	if req.Dosis != "" {
		rec.Dosis = req.Dosis
	}

	rec.FechaActualizacion = time.Now()

	return nil
}

func (r *recetaRepositoryImpl) CambiarEstado(id int, req model.CambioEstadoRecetaRequest) (*model.RecetaDetalle, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rec, exists := r.recetas[id]
	if !exists {
		return nil, fmt.Errorf("receta no encontrada")
	}

	if (req.NuevoEstado == model.RecetaEstadoObservado || req.NuevoEstado == model.RecetaEstadoRechazado) && req.Motivo == "" {
		return nil, fmt.Errorf("el motivo es obligatorio para estados OBSERVADO y RECHAZADO")
	}

	now := time.Now()

	rec.Estado = req.NuevoEstado
	rec.FechaActualizacion = now

	historial := model.HistorialEstadoReceta{
		Estado:      req.NuevoEstado,
		Usuario:     req.Usuario,
		FechaCambio: now,
		Motivo:      req.Motivo,
	}

	rec.Historial = append(rec.Historial, historial)

	return rec, nil
}
