package aplicacion

import (
	"context"
	"sort"

	"Proyecto_web/dominio"
)

type PilotoService struct {
	repo       dominio.PilotoRepository
	escuderias dominio.EscuderiaRepository
}

func NewPilotoService(repo dominio.PilotoRepository, escuderias dominio.EscuderiaRepository) *PilotoService {
	return &PilotoService{repo: repo, escuderias: escuderias}
}

// Listar devuelve los pilotos ordenados por puntos de la temporada.
func (s *PilotoService) Listar(ctx context.Context, orden string) ([]dominio.Piloto, error) {
	pilotos, err := s.repo.Listar(ctx)
	if err != nil {
		return nil, err
	}
	sort.Slice(pilotos, func(i, j int) bool {
		if orden == "asc" {
			return pilotos[i].Puntos < pilotos[j].Puntos
		}
		return pilotos[i].Puntos > pilotos[j].Puntos
	})
	return pilotos, nil
}

func (s *PilotoService) Recuperar(ctx context.Context, id int) (dominio.Piloto, error) {
	return s.repo.Recuperar(ctx, id)
}

func (s *PilotoService) ListarParaApuesta(ctx context.Context) ([]dominio.Piloto, error) {
	return s.repo.ListarParaApuesta(ctx)
}

// Escuderias devuelve las opciones del select de escuderías.
func (s *PilotoService) Escuderias(ctx context.Context) ([]string, error) {
	return s.repo.Escuderias(ctx)
}

// Inscribir aplica las reglas para agregar un piloto a la temporada.
func (s *PilotoService) Inscribir(ctx context.Context, p dominio.Piloto) error {
	if err := p.Validar(); err != nil {
		return err
	}
	enTemporada, err := s.repo.ExisteEnTemporada(ctx, p.ID)
	if err != nil {
		return err
	}
	if enTemporada {
		return dominio.ErrPilotoDuplicado
	}
	if ok, err := s.escuderias.Recuperar(ctx, p.Escuderia); err != nil || ok.ID == "" {
		return dominio.ErrEscuderiaNoExiste
	}
	enHistorico, err := s.repo.ExisteHistorico(ctx, p.ID)
	if err != nil {
		return err
	}
	if !enHistorico && (p.Nombre == "" || p.Pais == "" || p.Nacimiento.IsZero()) {
		return dominio.ErrHistoricoIncompleto
	}
	return s.repo.Inscribir(ctx, p, !enHistorico)
}

func (s *PilotoService) Modificar(ctx context.Context, p dominio.Piloto) error {
	if err := p.Validar(); err != nil {
		return err
	}
	if ok, err := s.escuderias.Recuperar(ctx, p.Escuderia); err != nil || ok.ID == "" {
		return dominio.ErrEscuderiaNoExiste
	}
	return s.repo.Modificar(ctx, p)
}

func (s *PilotoService) Eliminar(ctx context.Context, id int) error {
	if id <= 0 {
		return dominio.ErrValidacion("ID de piloto inválido")
	}
	return s.repo.Eliminar(ctx, id)
}
