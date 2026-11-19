package aplicacion

import (
	"context"
	"time"

	"Proyecto_web/dominio"
)

type GranPremioHistoricoService struct {
	repo        dominio.GranPremioHistoricoRepository
	granPremios dominio.GranPremioRepository
	historicos  dominio.PilotoHistoricoRepository
}

func NewGranPremioHistoricoService(
	repo dominio.GranPremioHistoricoRepository,
	granPremios dominio.GranPremioRepository,
	historicos dominio.PilotoHistoricoRepository,
) *GranPremioHistoricoService {
	return &GranPremioHistoricoService{repo: repo, granPremios: granPremios, historicos: historicos}
}

func (s *GranPremioHistoricoService) Listar(ctx context.Context) ([]dominio.GranPremioHistorico, error) {
	return s.repo.Listar(ctx)
}

func (s *GranPremioHistoricoService) ListarProximos(ctx context.Context) ([]dominio.GranPremioHistorico, error) {
	return s.repo.ListarProximos(ctx)
}

func (s *GranPremioHistoricoService) ListarPendientes(ctx context.Context) ([]dominio.GranPremioHistorico, error) {
	return s.repo.ListarPendientes(ctx)
}

func (s *GranPremioHistoricoService) Recuperar(ctx context.Context, id string, fechaTime time.Time) (dominio.GranPremioHistorico, error) {
	return s.repo.Recuperar(ctx, id, fechaTime)
}

// Abierta responde si todavía se puede apostar en esa edición.
func (s *GranPremioHistoricoService) Abierta(ctx context.Context, id string, fechaTime time.Time) bool {
	carrera, err := s.repo.Recuperar(ctx, id, fechaTime)
	return err == nil && carrera.Abierta()
}

// Programar crea una carrera futura para un circuito existente.
func (s *GranPremioHistoricoService) Programar(ctx context.Context, id string, fechaTime time.Time) error {
	g := dominio.GranPremioHistorico{ID: id, Fecha: fechaTime}
	if err := g.Validar(); err != nil {
		return err
	}
	existe, err := s.granPremios.Existe(ctx, id)
	if err != nil {
		return err
	}
	if !existe {
		return dominio.ErrGranPremioNoEncontrado
	}
	return s.repo.Programar(ctx, g)
}

// Crear registra una edición con resultado (uso de la API REST).
func (s *GranPremioHistoricoService) Crear(ctx context.Context, g dominio.GranPremioHistorico) error {
	if err := g.Validar(); err != nil {
		return err
	}
	if len(g.Resultado) == 0 {
		return dominio.ErrResultadoIncompleto
	}
	existe, err := s.granPremios.Existe(ctx, g.ID)
	if err != nil {
		return err
	}
	if !existe {
		return dominio.ErrGranPremioNoEncontrado
	}
	if err := s.validarPilotos(ctx, g.Resultado); err != nil {
		return err
	}
	return s.repo.Crear(ctx, g)
}

// RegistrarResultado completa el resultado de una carrera programada.
func (s *GranPremioHistoricoService) RegistrarResultado(ctx context.Context, g dominio.GranPremioHistorico) error {
	if err := g.Validar(); err != nil {
		return err
	}
	if len(g.Resultado) == 0 {
		return dominio.ErrResultadoIncompleto
	}
	existe, err := s.granPremios.Existe(ctx, g.ID)
	if err != nil {
		return err
	}
	if !existe {
		return dominio.ErrGranPremioNoEncontrado
	}
	if err := s.validarPilotos(ctx, g.Resultado); err != nil {
		return err
	}
	return s.repo.RegistrarResultado(ctx, g)
}

func (s *GranPremioHistoricoService) validarPilotos(ctx context.Context, ids []int) error {
	for _, id := range ids {
		if ok, err := s.historicos.Existe(ctx, id); err != nil || !ok {
			return dominio.ErrValidacion("Uno de los pilotos no existe en la grilla")
		}
	}
	return nil
}

func (s *GranPremioHistoricoService) Eliminar(ctx context.Context, id string, fechaTime time.Time) error {
	if id == "" || fechaTime.IsZero() {
		return dominio.ErrValidacion("El ID del Gran Premio y la fecha de carrera son obligatorios")
	}
	return s.repo.Eliminar(ctx, id, fechaTime)
}
