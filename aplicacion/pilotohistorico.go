package aplicacion

import (
	"context"
	"sort"

	"Proyecto_web/dominio"
)

type PilotoHistoricoService struct {
	repo dominio.PilotoHistoricoRepository
}

func NewPilotoHistoricoService(repo dominio.PilotoHistoricoRepository) *PilotoHistoricoService {
	return &PilotoHistoricoService{repo: repo}
}

func (s *PilotoHistoricoService) Listar(ctx context.Context) ([]dominio.PilotoHistorico, error) {
	return s.repo.Listar(ctx)
}

// Nombres devuelve el nombre de cada piloto histórico, para mostrar apuestas viejas.
func (s *PilotoHistoricoService) Nombres(ctx context.Context) (map[int]string, error) {
	pilotos, err := s.repo.Listar(ctx)
	if err != nil {
		return nil, err
	}
	nombres := make(map[int]string, len(pilotos))
	for _, p := range pilotos {
		nombres[p.ID] = p.Nombre
	}
	return nombres, nil
}

func (s *PilotoHistoricoService) Recuperar(ctx context.Context, id int) (dominio.PilotoHistorico, error) {
	return s.repo.Recuperar(ctx, id)
}

// Crear registra un piloto histórico nuevo, siempre sin títulos ganados.
func (s *PilotoHistoricoService) Crear(ctx context.Context, p dominio.PilotoHistorico) error {
	p.Titulos = 0
	if err := p.Validar(); err != nil {
		return err
	}
	return s.repo.Crear(ctx, p)
}

func (s *PilotoHistoricoService) Modificar(ctx context.Context, p dominio.PilotoHistorico) error {
	if err := p.Validar(); err != nil {
		return err
	}
	return s.repo.Modificar(ctx, p)
}

func (s *PilotoHistoricoService) Eliminar(ctx context.Context, id int) error {
	if id <= 0 {
		return dominio.ErrValidacion("ID de piloto histórico inválido")
	}
	return s.repo.Eliminar(ctx, id)
}

// Opciones devuelve los pilotos históricos ordenados por nombre, para los selects.
func (s *PilotoHistoricoService) Opciones(ctx context.Context) ([]dominio.PilotoHistorico, error) {
	pilotos, err := s.repo.Listar(ctx)
	if err != nil {
		return nil, err
	}
	sort.Slice(pilotos, func(i, j int) bool { return pilotos[i].Nombre < pilotos[j].Nombre })
	return pilotos, nil
}
