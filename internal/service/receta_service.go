package service

import (
	"prestadores-api/internal/model"
	"prestadores-api/internal/repository"

	"go.uber.org/zap"
)

type RecetaService interface {
	GetRecetas(estado string, query string, page int, size int, sort string) (*model.PaginatedRecetasResponse, error)
	GetRecetaByID(id int) (*model.RecetaDetalle, error)
	CreateReceta(req model.CreateRecetaRequest) (*model.CreateRecetaResponse, error)
	UpdateReceta(id int, req model.UpdateRecetaRequest) error
	CambiarEstadoReceta(id int, req model.CambioEstadoRecetaRequest) (*model.CambioEstadoRecetaResponse, error)
}

type recetaServiceImpl struct {
	repo   repository.RecetaRepository
	logger *zap.Logger
}

func NewRecetaService(repo repository.RecetaRepository, logger *zap.Logger) RecetaService {
	return &recetaServiceImpl{
		repo:   repo,
		logger: logger,
	}
}

func (s *recetaServiceImpl) GetRecetas(estado string, query string, page int, size int, sort string) (*model.PaginatedRecetasResponse, error) {
	s.logger.Info("Obteniendo recetas",
		zap.String("estado", estado),
		zap.String("query", query),
		zap.Int("page", page),
		zap.Int("size", size),
		zap.String("sort", sort),
	)

	items, total, err := s.repo.GetAll(estado, query, page, size, sort)
	if err != nil {
		s.logger.Error("Error al obtener recetas", zap.Error(err))
		return nil, err
	}

	response := &model.PaginatedRecetasResponse{
		Page:  page,
		Size:  size,
		Total: total,
		Items: items,
	}

	return response, nil
}

func (s *recetaServiceImpl) GetRecetaByID(id int) (*model.RecetaDetalle, error) {
	s.logger.Info("Obteniendo receta por ID", zap.Int("id", id))

	detalle, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("Error al obtener receta", zap.Int("id", id), zap.Error(err))
		return nil, err
	}

	return detalle, nil
}

func (s *recetaServiceImpl) CreateReceta(req model.CreateRecetaRequest) (*model.CreateRecetaResponse, error) {
	s.logger.Info("Creando receta",
		zap.Int("afiliadoId", req.AfiliadoID),
		zap.String("medicamento", req.Medicamento),
		zap.String("dosis", req.Dosis),
	)

	detalle, err := s.repo.Create(req)
	if err != nil {
		s.logger.Error("Error al crear receta", zap.Error(err))
		return nil, err
	}

	response := &model.CreateRecetaResponse{
		ID:            detalle.ID,
		Tipo:          detalle.Tipo,
		Estado:        detalle.Estado,
		FechaCreacion: detalle.FechaCreacion,
	}

	return response, nil
}

func (s *recetaServiceImpl) UpdateReceta(id int, req model.UpdateRecetaRequest) error {
	s.logger.Info("Actualizando receta",
		zap.Int("id", id),
		zap.String("medicamento", req.Medicamento),
		zap.String("dosis", req.Dosis),
	)

	err := s.repo.Update(id, req)
	if err != nil {
		s.logger.Error("Error al actualizar receta", zap.Int("id", id), zap.Error(err))
		return err
	}

	return nil
}

func (s *recetaServiceImpl) CambiarEstadoReceta(id int, req model.CambioEstadoRecetaRequest) (*model.CambioEstadoRecetaResponse, error) {
	s.logger.Info("Cambiando estado de receta",
		zap.Int("id", id),
		zap.String("nuevoEstado", string(req.NuevoEstado)),
		zap.String("usuario", req.Usuario),
	)

	if (req.NuevoEstado == model.RecetaEstadoObservado || req.NuevoEstado == model.RecetaEstadoRechazado) && req.Motivo == "" {
		s.logger.Warn("Intento de cambio de estado sin motivo",
			zap.Int("id", id),
			zap.String("nuevoEstado", string(req.NuevoEstado)),
		)
		return nil, ErrMotivoRequerido
	}

	// aca tenemos q agregar las validaciones adicionales q tengamos

	detalle, err := s.repo.CambiarEstado(id, req)
	if err != nil {
		s.logger.Error("Error al cambiar estado de receta", zap.Int("id", id), zap.Error(err))
		return nil, err
	}

	response := &model.CambioEstadoRecetaResponse{
		ID:                 detalle.ID,
		Tipo:               detalle.Tipo,
		Estado:             detalle.Estado,
		FechaActualizacion: detalle.FechaActualizacion,
	}

	return response, nil
}
