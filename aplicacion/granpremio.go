package aplicacion

import (
	"context"
	"sort"

	"Proyecto_web/dominio"
)

type GranPremioService struct {
	repo       dominio.GranPremioRepository
	historicos dominio.PilotoHistoricoRepository
}

func NewGranPremioService(repo dominio.GranPremioRepository, historicos dominio.PilotoHistoricoRepository) *GranPremioService {
	return &GranPremioService{repo: repo, historicos: historicos}
}

// Listar devuelve los circuitos ordenados por nombre ("asc" o "desc").
func (s *GranPremioService) Listar(ctx context.Context, orden string) ([]dominio.GranPremio, error) {
	gps, err := s.repo.Listar(ctx)
	if err != nil {
		return nil, err
	}
	sort.Slice(gps, func(i, j int) bool {
		if orden == "desc" {
			return gps[i].ID > gps[j].ID
		}
		return gps[i].ID < gps[j].ID
	})
	return gps, nil
}

func (s *GranPremioService) Recuperar(ctx context.Context, id string) (dominio.GranPremio, error) {
	return s.repo.Recuperar(ctx, id)
}

func (s *GranPremioService) Circuitos(ctx context.Context) ([]string, error) {
	return s.repo.Circuitos(ctx)
}

func (s *GranPremioService) Existe(ctx context.Context, id string) bool {
	gp, err := s.repo.Recuperar(ctx, id)
	return err == nil && gp.ID != ""
}

func (s *GranPremioService) Crear(ctx context.Context, g dominio.GranPremio) error {
	if err := g.Validar(); err != nil {
		return err
	}
	if s.Existe(ctx, g.ID) {
		return dominio.ErrGranPremioDuplicado
	}
	if ok, err := s.historicos.Existe(ctx, g.UltimoGanador); err != nil || !ok {
		return dominio.ErrValidacion("El último ganador no existe")
	}
	return s.repo.Crear(ctx, g)
}

func (s *GranPremioService) Modificar(ctx context.Context, g dominio.GranPremio) error {
	if err := g.Validar(); err != nil {
		return err
	}
	if ok, err := s.historicos.Existe(ctx, g.UltimoGanador); err != nil || !ok {
		return dominio.ErrValidacion("El último ganador no existe")
	}
	return s.repo.Modificar(ctx, g)
}

func (s *GranPremioService) Eliminar(ctx context.Context, id string) error {
	if id == "" {
		return dominio.ErrValidacion("ID de gran premio inválido")
	}
	return s.repo.Eliminar(ctx, id)
}
