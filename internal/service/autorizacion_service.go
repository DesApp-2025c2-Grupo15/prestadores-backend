package service

import (
	"prestadores-api/internal/model"
	"prestadores-api/internal/repository"

	"go.uber.org/zap"
)

type AutorizacionService interface {
	GetAutorizaciones(estado string, query string, page int, size int, sort string) (*model.PaginatedAutorizacionesResponse, error)
	GetAutorizacionByID(id int) (*model.AutorizacionDetalle, error)
	CreateAutorizacion(req model.CreateAutorizacionRequest) (*model.CreateAutorizacionResponse, error)
	UpdateAutorizacion(id int, req model.UpdateAutorizacionRequest) error
	CambiarEstadoAutorizacion(id int, req model.CambioEstadoRequest) (*model.CambioEstadoResponse, error)
}

type autorizacionServiceImpl struct {
	repo   repository.AutorizacionRepository
	logger *zap.Logger
}

func NewAutorizacionService(repo repository.AutorizacionRepository, logger *zap.Logger) AutorizacionService {
	return &autorizacionServiceImpl{
		repo:   repo,
		logger: logger,
	}
}

func (s *autorizacionServiceImpl) GetAutorizaciones(estado string, query string, page int, size int, sort string) (*model.PaginatedAutorizacionesResponse, error) {
	s.logger.Info("Obteniendo autorizaciones",
		zap.String("estado", estado),
		zap.String("query", query),
		zap.Int("page", page),
		zap.Int("size", size),
		zap.String("sort", sort),
	)

	items, total, err := s.repo.GetAll(estado, query, page, size, sort)
	if err != nil {
		s.logger.Error("Error al obtener autorizaciones", zap.Error(err))
		return nil, err
	}

	response := &model.PaginatedAutorizacionesResponse{
		Page:  page,
		Size:  size,
		Total: total,
		Items: items,
	}

	return response, nil
}

func (s *autorizacionServiceImpl) GetAutorizacionByID(id int) (*model.AutorizacionDetalle, error) {
	s.logger.Info("Obteniendo autorización por ID", zap.Int("id", id))

	detalle, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("Error al obtener autorización", zap.Int("id", id), zap.Error(err))
		return nil, err
	}

	return detalle, nil
}

func (s *autorizacionServiceImpl) CreateAutorizacion(req model.CreateAutorizacionRequest) (*model.CreateAutorizacionResponse, error) {
	s.logger.Info("Creando autorización",
		zap.Int("afiliadoId", req.AfiliadoID),
		zap.String("procedimiento", req.Procedimiento),
		zap.String("especialidad", req.Especialidad),
	)

	detalle, err := s.repo.Create(req)
	if err != nil {
		s.logger.Error("Error al crear autorización", zap.Error(err))
		return nil, err
	}

	response := &model.CreateAutorizacionResponse{
		ID:            detalle.ID,
		Tipo:          detalle.Tipo,
		Estado:        detalle.Estado,
		FechaCreacion: detalle.FechaCreacion,
	}

	return response, nil
}

func (s *autorizacionServiceImpl) UpdateAutorizacion(id int, req model.UpdateAutorizacionRequest) error {
	s.logger.Info("Actualizando autorización",
		zap.Int("id", id),
		zap.String("procedimiento", req.Procedimiento),
		zap.String("especialidad", req.Especialidad),
	)

	err := s.repo.Update(id, req)
	if err != nil {
		s.logger.Error("Error al actualizar autorización", zap.Int("id", id), zap.Error(err))
		return err
	}

	return nil
}

func (s *autorizacionServiceImpl) CambiarEstadoAutorizacion(id int, req model.CambioEstadoRequest) (*model.CambioEstadoResponse, error) {
	s.logger.Info("Cambiando estado de autorización",
		zap.Int("id", id),
		zap.String("nuevoEstado", string(req.NuevoEstado)),
		zap.String("usuario", req.Usuario),
	)

	if (req.NuevoEstado == model.EstadoObservado || req.NuevoEstado == model.EstadoRechazado) && req.Motivo == "" {
		s.logger.Warn("Intento de cambio de estado sin motivo",
			zap.Int("id", id),
			zap.String("nuevoEstado", string(req.NuevoEstado)),
		)
		return nil, ErrMotivoRequerido
	}

	// aca tenemos q agregar las validaciones adicionales q tengamos

	detalle, err := s.repo.CambiarEstado(id, req)
	if err != nil {
		s.logger.Error("Error al cambiar estado de autorización", zap.Int("id", id), zap.Error(err))
		return nil, err
	}

	response := &model.CambioEstadoResponse{
		ID:                 detalle.ID,
		Tipo:               detalle.Tipo,
		Estado:             detalle.Estado,
		FechaActualizacion: detalle.FechaActualizacion,
	}

	return response, nil
}

// Errores personalizados del servicio
var (
	ErrMotivoRequerido = &ServiceError{Message: "El motivo es obligatorio para estados OBSERVADO y RECHAZADO"}
)

// ServiceError representa un error del servicio
type ServiceError struct {
	Message string
}

func (e *ServiceError) Error() string {
	return e.Message
}
