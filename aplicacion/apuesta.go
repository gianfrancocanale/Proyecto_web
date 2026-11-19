package aplicacion

import (
	"context"
	"errors"
	"time"

	"Proyecto_web/dominio"
)

type ApuestaService struct {
	repo       dominio.ApuestaRepository
	carreras   dominio.GranPremioHistoricoRepository
	pilotos    dominio.PilotoRepository
	historicos dominio.PilotoHistoricoRepository
}

func NewApuestaService(
	repo dominio.ApuestaRepository,
	carreras dominio.GranPremioHistoricoRepository,
	pilotos dominio.PilotoRepository,
	historicos dominio.PilotoHistoricoRepository,
) *ApuestaService {
	return &ApuestaService{repo: repo, carreras: carreras, pilotos: pilotos, historicos: historicos}
}

func (s *ApuestaService) ListarPorUsuario(ctx context.Context, usuario string) ([]dominio.Apuesta, error) {
	return s.repo.ListarPorUsuario(ctx, usuario)
}

// CarrerasAbiertas devuelve las ediciones donde todavía se puede apostar.
func (s *ApuestaService) CarrerasAbiertas(ctx context.Context) ([]dominio.GranPremioHistorico, error) {
	return s.carreras.ListarProximos(ctx)
}

// PilotosDisponibles son los pilotos de la temporada que se pueden pronosticar.
func (s *ApuestaService) PilotosDisponibles(ctx context.Context) ([]dominio.Piloto, error) {
	return s.pilotos.ListarParaApuesta(ctx)
}

// NombresPilotos sirve para mostrar apuestas de carreras ya pasadas.
func (s *ApuestaService) NombresPilotos(ctx context.Context) (map[int]string, error) {
	lista, err := s.historicos.Listar(ctx)
	if err != nil {
		return nil, err
	}
	nombres := make(map[int]string, len(lista))
	for _, p := range lista {
		nombres[p.ID] = p.Nombre
	}
	return nombres, nil
}

func (s *ApuestaService) Crear(ctx context.Context, usuario, idGP string, fecha time.Time, prediccion []int) error {
	a := dominio.Apuesta{
		Usuario:      usuario,
		GranPremio:   idGP,
		FechaCarrera: fecha,
		Prediccion:   prediccion,
	}
	if err := s.validar(ctx, a); err != nil {
		return err
	}
	a.FechaApuesta = time.Now()
	return s.repo.Crear(ctx, a)
}

func (s *ApuestaService) Modificar(ctx context.Context, id int, usuario string, prediccion []int) error {
	a, err := s.repo.Recuperar(ctx, id, usuario)
	if errors.Is(err, dominio.ErrApuestaNoEncontrada) {
		return err
	}
	if err != nil {
		return err
	}
	if !s.carreraAbierta(ctx, a) {
		return dominio.ErrCarreraCerrada
	}
	a.Prediccion = prediccion
	if err := s.validar(ctx, a); err != nil {
		return err
	}
	a.FechaApuesta = time.Now()
	return s.repo.Actualizar(ctx, a)
}

func (s *ApuestaService) Eliminar(ctx context.Context, id int, usuario string) error {
	a, err := s.repo.Recuperar(ctx, id, usuario)
	if err != nil {
		return err
	}
	if !s.carreraAbierta(ctx, a) {
		return dominio.ErrCarreraCerrada
	}
	return s.repo.Eliminar(ctx, id, usuario)
}

// Recuperar devuelve la apuesta solo si es del usuario y la carrera sigue abierta.
func (s *ApuestaService) RecuperarEditable(ctx context.Context, id int, usuario string) (dominio.Apuesta, error) {
	a, err := s.repo.Recuperar(ctx, id, usuario)
	if err != nil {
		return dominio.Apuesta{}, err
	}
	if !s.carreraAbierta(ctx, a) {
		return dominio.Apuesta{}, dominio.ErrCarreraCerrada
	}
	return a, nil
}

func (s *ApuestaService) validar(ctx context.Context, a dominio.Apuesta) error {
	if a.Usuario == "" || a.GranPremio == "" || a.FechaCarrera.IsZero() {
		return dominio.ErrValidacion("Elegí un gran premio")
	}
	if !s.carreraAbierta(ctx, a) {
		return dominio.ErrCarreraCerrada
	}
	pilotos, err := s.pilotos.ListarParaApuesta(ctx)
	if err != nil {
		return err
	}
	ids := make([]int, 0, len(pilotos))
	for _, p := range pilotos {
		ids = append(ids, p.ID)
	}
	return a.ValidarPrediccion(ids)
}

func (s *ApuestaService) carreraAbierta(ctx context.Context, a dominio.Apuesta) bool {
	carrera, err := s.carreras.Recuperar(ctx, a.GranPremio, a.FechaCarrera)
	return err == nil && carrera.Abierta()
}
