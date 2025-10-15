package service

import (
	"prestadores-api/internal/model"
	"prestadores-api/internal/repository"

	"go.uber.org/zap"
)

type ReintegroService interface {
	GetReintegros(estado string, query string, page int, size int, sort string) (*model.PaginatedReintegrosResponse, error)
	GetReintegroByID(id int) (*model.ReintegroDetalle, error)
	CreateReintegro(req model.CreateReintegroRequest) (*model.CreateReintegroResponse, error)
	UpdateReintegro(id int, req model.UpdateReintegroRequest) error
	CambiarEstadoReintegro(id int, req model.CambioEstadoRequest) (*model.CambioEstadoResponse, error)
}

type reintegroServiceImpl struct {
	repo   repository.ReintegroRepository
	logger *zap.Logger
}

func NewReintegroService(repo repository.ReintegroRepository, logger *zap.Logger) ReintegroService {
	return &reintegroServiceImpl{
		repo:   repo,
		logger: logger,
	}
}

func (s *reintegroServiceImpl) GetReintegros(estado string, query string, page int, size int, sort string) (*model.PaginatedReintegrosResponse, error) {
	s.logger.Info("Obteniendo reintegros",
		zap.String("estado", estado),
		zap.String("query", query),
		zap.Int("page", page),
		zap.Int("size", size),
		zap.String("sort", sort),
	)

	items, total, err := s.repo.GetAll(estado, query, page, size, sort)
	if err != nil {
		s.logger.Error("Error al obtener reintegros", zap.Error(err))
		return nil, err
	}

	response := &model.PaginatedReintegrosResponse{
		Page:  page,
		Size:  size,
		Total: total,
		Items: items,
	}

	return response, nil
}

func (s *reintegroServiceImpl) GetReintegroByID(id int) (*model.ReintegroDetalle, error) {
	s.logger.Info("Obteniendo reintegro por ID", zap.Int("id", id))

	detalle, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("Error al obtener reintegro", zap.Int("id", id), zap.Error(err))
		return nil, err
	}

	return detalle, nil
}

func (s *reintegroServiceImpl) CreateReintegro(req model.CreateReintegroRequest) (*model.CreateReintegroResponse, error) {
	s.logger.Info("Creando reintegro",
		zap.Int("afiliadoId", req.AfiliadoID),
		zap.String("prestacion", req.Prestacion),
		zap.String("metodo", req.Metodo),
		zap.Float64("monto", req.Monto),
	)

	detalle, err := s.repo.Create(req)
	if err != nil {
		s.logger.Error("Error al crear reintegro", zap.Error(err))
		return nil, err
	}

	response := &model.CreateReintegroResponse{
		ID:            detalle.ID,
		Tipo:          detalle.Tipo,
		Estado:        detalle.Estado,
		FechaCreacion: detalle.FechaCreacion,
	}

	return response, nil
}

func (s *reintegroServiceImpl) UpdateReintegro(id int, req model.UpdateReintegroRequest) error {
	s.logger.Info("Actualizando reintegro",
		zap.Int("id", id),
		zap.String("prestacion", req.Prestacion),
		zap.String("metodo", req.Metodo),
		zap.Float64("monto", req.Monto),
	)

	if err := s.repo.Update(id, req); err != nil {
		s.logger.Error("Error al actualizar reintegro", zap.Int("id", id), zap.Error(err))
		return err
	}
	return nil
}

func (s *reintegroServiceImpl) CambiarEstadoReintegro(id int, req model.CambioEstadoRequest) (*model.CambioEstadoResponse, error) {
	s.logger.Info("Cambiando estado de reintegro",
		zap.Int("id", id),
		zap.String("nuevoEstado", string(req.NuevoEstado)),
		zap.String("usuario", req.Usuario),
	)

	// Validación de motivo obligatorio para OBSERVADO y RECHAZADO
	if (req.NuevoEstado == model.EstadoObservado || req.NuevoEstado == model.EstadoRechazado) && req.Motivo == "" {
		s.logger.Warn("Intento de cambio de estado sin motivo",
			zap.Int("id", id),
			zap.String("nuevoEstado", string(req.NuevoEstado)),
		)
		return nil, ErrMotivoRequerido
	}

	// TODO: agregar validaciones adicionales del workflow si corresponde

	detalle, err := s.repo.CambiarEstado(id, req)
	if err != nil {
		s.logger.Error("Error al cambiar estado de reintegro", zap.Int("id", id), zap.Error(err))
		return nil, err
	}

	resp := &model.CambioEstadoResponse{
		ID:                 detalle.ID,
		Tipo:               detalle.Tipo,
		Estado:             detalle.Estado,
		FechaActualizacion: detalle.FechaActualizacion,
	}
	return resp, nil
}
