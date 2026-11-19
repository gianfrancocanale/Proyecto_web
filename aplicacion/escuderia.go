// services/escuderia.go
package aplicacion

import (
	"context"
	"sort"

	"Proyecto_web/dominio"
)

type EscuderiaService struct {
	repo dominio.EscuderiaRepository
}

func NewEscuderiaService(repo dominio.EscuderiaRepository) *EscuderiaService {
	return &EscuderiaService{repo: repo}
}

// Listar devuelve las escuderías ordenadas por ID ("asc" o "desc").
func (s *EscuderiaService) Listar(ctx context.Context, orden string) ([]dominio.Escuderia, error) {
	es, err := s.repo.Listar(ctx)
	if err != nil {
		return nil, err
	}
	sort.Slice(es, func(i, j int) bool {
		if orden == "desc" {
			return es[i].ID > es[j].ID
		}
		return es[i].ID < es[j].ID
	})
	return es, nil
}

func (s *EscuderiaService) Recuperar(ctx context.Context, id string) (dominio.Escuderia, error) {
	return s.repo.Recuperar(ctx, id)
}

func (s *EscuderiaService) Crear(ctx context.Context, e dominio.Escuderia) error {
	if err := e.Validar(); err != nil {
		return err
	}
	return s.repo.Crear(ctx, e)
}

func (s *EscuderiaService) Modificar(ctx context.Context, e dominio.Escuderia) error {
	if err := e.Validar(); err != nil {
		return err
	}
	return s.repo.Modificar(ctx, e)
}

func (s *EscuderiaService) Eliminar(ctx context.Context, id string) error {
	if id == "" {
		return dominio.ErrValidacion("id inválido")
	}
	return s.repo.Eliminar(ctx, id)
}
