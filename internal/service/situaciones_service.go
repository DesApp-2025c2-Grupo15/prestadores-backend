package service

import (
	"fmt"
	"prestadores-api/internal/model"
	"prestadores-api/internal/repository"

	"go.uber.org/zap"
)

type SituacionService interface {
	// scope: "" (titular) | "grupo"
	GetSituaciones(afiliadoID int, scope string) (interface{}, error)
	CreateSituacion(req model.CreateSituacionRequest) (*model.CreateSituacionResponse, error)
	PatchSituacion(situacionID int, req model.PatchSituacionRequest) error
	CambiarEstadoSituacion(situacionID int, req model.CambioEstadoSituacionRequest) (*model.CambioEstadoSituacionResponse, error)
	DeleteSituacion(situacionID int) error
}

type situacionServiceImpl struct {
	repo   repository.SituacionRepository
	logger *zap.Logger
}

func NewSituacionService(repo repository.SituacionRepository, logger *zap.Logger) SituacionService {
	return &situacionServiceImpl{
		repo:   repo,
		logger: logger,
	}
}

// GetSituaciones retorna:
// - *model.SituacionesAfiliadoResponse cuando scope = "" (titular)
// - *model.SituacionesGrupoResponse    cuando scope = "grupo"
func (s *situacionServiceImpl) GetSituaciones(afiliadoID int, scope string) (interface{}, error) {
	s.logger.Info("Obteniendo situaciones terapéuticas",
		zap.Int("afiliadoId", afiliadoID),
		zap.String("scope", scope),
	)

	if scope == "grupo" {
		integrantes, err := s.repo.GetByAfiliadoGrupo(afiliadoID)
		if err != nil {
			s.logger.Error("Error al obtener situaciones de grupo", zap.Error(err))
			return nil, err
		}
		resp := &model.SituacionesGrupoResponse{
			AfiliadoID:  afiliadoID,
			Integrantes: integrantes,
		}
		return resp, nil
	}

	items, err := s.repo.GetByAfiliado(afiliadoID)
	if err != nil {
		s.logger.Error("Error al obtener situaciones del afiliado", zap.Error(err))
		return nil, err
	}
	resp := &model.SituacionesAfiliadoResponse{
		AfiliadoID: afiliadoID,
		Items:      items,
	}
	return resp, nil
}

func (s *situacionServiceImpl) CreateSituacion(req model.CreateSituacionRequest) (*model.CreateSituacionResponse, error) {
	s.logger.Info("Creando situación terapéutica",
		zap.Int("afiliadoId", req.AfiliadoID),
		zap.Any("miembroId", req.MiembroID),
		zap.String("descripcion", req.Descripcion),
		zap.String("fechaInicio", req.FechaInicio),
	)

	// Validaciones mínimas
	if req.Descripcion == "" || req.FechaInicio == "" || req.AfiliadoID <= 0 {
		s.logger.Warn("Request inválido al crear situación")
		return nil, fmt.Errorf("request inválido: afiliadoId, descripcion y fechaInicio son obligatorios")
	}

	detalle, err := s.repo.Create(req)
	if err != nil {
		s.logger.Error("Error al crear situación", zap.Error(err))
		return nil, err
	}

	resp := &model.CreateSituacionResponse{
		ID:            detalle.ID,
		AfiliadoID:    detalle.AfiliadoID,
		MiembroID:     detalle.MiembroID,
		Estado:        detalle.Estado,
		FechaCreacion: detalle.FechaCreacion,
	}
	return resp, nil
}

func (s *situacionServiceImpl) PatchSituacion(situacionID int, req model.PatchSituacionRequest) error {
	s.logger.Info("Actualizando (PATCH) situación terapéutica",
		zap.Int("situacionId", situacionID),
		zap.Any("patch", req),
	)

	if situacionID <= 0 {
		return fmt.Errorf("situacionId inválido")
	}

	if err := s.repo.Patch(situacionID, req); err != nil {
		s.logger.Error("Error al actualizar situación", zap.Int("situacionId", situacionID), zap.Error(err))
		return err
	}
	return nil
}

func (s *situacionServiceImpl) CambiarEstadoSituacion(situacionID int, req model.CambioEstadoSituacionRequest) (*model.CambioEstadoSituacionResponse, error) {
	s.logger.Info("Cambiando estado de situación",
		zap.Int("situacionId", situacionID),
		zap.String("nuevoEstado", string(req.Estado)),
		zap.String("usuario", req.Usuario),
	)

	if situacionID <= 0 {
		return nil, fmt.Errorf("situacionId inválido")
	}
	if req.Usuario == "" {
		return nil, fmt.Errorf("usuario es obligatorio")
	}

	if req.Estado == model.EstadoSituacionBaja && req.Motivo == "" {
		return nil, fmt.Errorf("el motivo es obligatorio para BAJA")
	}

	detalle, err := s.repo.CambiarEstado(situacionID, req)
	if err != nil {
		s.logger.Error("Error al cambiar estado de situación", zap.Int("situacionId", situacionID), zap.Error(err))
		return nil, err
	}

	resp := &model.CambioEstadoSituacionResponse{
		ID:                 detalle.ID,
		Estado:             detalle.Estado,
		FechaActualizacion: detalle.FechaActualizacion,
	}
	return resp, nil
}

func (s *situacionServiceImpl) DeleteSituacion(situacionID int) error {
	s.logger.Info("Eliminando situación terapéutica", zap.Int("situacionId", situacionID))
	if situacionID <= 0 {
		return fmt.Errorf("situacionId inválido")
	}
	if err := s.repo.Delete(situacionID); err != nil {
		s.logger.Error("Error al eliminar situación", zap.Int("situacionId", situacionID), zap.Error(err))
		return err
	}
	return nil
}
