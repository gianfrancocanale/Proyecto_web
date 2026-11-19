package aplicacion

import (
	"context"
	"time"

	"Proyecto_web/dominio"
)

// AdminService agrupa los casos de uso del panel: programar carreras y cargar resultados.
type AdminService struct {
	usuarios    dominio.UsuarioRepository
	granPremios dominio.GranPremioRepository
	carreras    dominio.GranPremioHistoricoRepository
	apuestas    dominio.ApuestaRepository
	pilotos     dominio.PilotoRepository
}

func NewAdminService(
	usuarios dominio.UsuarioRepository,
	granPremios dominio.GranPremioRepository,
	carreras dominio.GranPremioHistoricoRepository,
	apuestas dominio.ApuestaRepository,
	pilotos dominio.PilotoRepository,
) *AdminService {
	return &AdminService{
		usuarios:    usuarios,
		granPremios: granPremios,
		carreras:    carreras,
		apuestas:    apuestas,
		pilotos:     pilotos,
	}
}

func (s *AdminService) EsAdmin(ctx context.Context, nombre string) bool {
	if nombre == "" {
		return false
	}
	ok, err := s.usuarios.EsAdmin(ctx, nombre)
	return err == nil && ok
}

func (s *AdminService) Circuitos(ctx context.Context) ([]string, error) {
	return s.granPremios.Circuitos(ctx)
}

func (s *AdminService) Pendientes(ctx context.Context) ([]dominio.GranPremioHistorico, error) {
	return s.carreras.ListarPendientes(ctx)
}

func (s *AdminService) Pilotos(ctx context.Context) ([]dominio.Piloto, error) {
	return s.pilotos.ListarParaApuesta(ctx)
}

func (s *AdminService) ProgramarCarrera(ctx context.Context, circuito string, fecha time.Time) error {
	g := dominio.GranPremioHistorico{ID: circuito, Fecha: fecha}
	if err := g.Validar(); err != nil {
		return err
	}
	existe, err := s.granPremios.Existe(ctx, circuito)
	if err != nil {
		return err
	}
	if !existe {
		return dominio.ErrGranPremioNoEncontrado
	}
	return s.carreras.Programar(ctx, g)
}

// CargarResultado valida el orden de llegada, calcula los puntos de cada apuesta
// y lo persiste todo junto en una transacción.
func (s *AdminService) CargarResultado(ctx context.Context, id string, fecha time.Time, orden []int) error {
	g := dominio.GranPremioHistorico{ID: id, Fecha: fecha, Resultado: orden}
	if err := g.Validar(); err != nil {
		return err
	}
	if err := s.validarOrden(ctx, orden); err != nil {
		return err
	}

	apuestas, err := s.apuestas.ListarPorGranPremio(ctx, id, fecha)
	if err != nil {
		return err
	}
	puntos := make([]dominio.PuntosUsuario, 0, len(apuestas))
	for _, a := range apuestas {
		pts := dominio.CalcularPuntos(a.Prediccion, orden)
		if pts > 0 {
			puntos = append(puntos, dominio.PuntosUsuario{Usuario: a.Usuario, Puntos: pts})
		}
	}
	return s.carreras.CargarResultado(ctx, g, puntos)
}

// validarOrden exige los 22 pilotos de la temporada, sin repetir.
func (s *AdminService) validarOrden(ctx context.Context, orden []int) error {
	pilotos, err := s.pilotos.ListarParaApuesta(ctx)
	if err != nil {
		return err
	}
	validos := make(map[int]bool, len(pilotos))
	for _, p := range pilotos {
		validos[p.ID] = true
	}
	if len(orden) != 22 {
		return dominio.ErrValidacion("Tenés que ordenar los 22 pilotos")
	}
	vistos := make(map[int]bool, 22)
	for _, id := range orden {
		if !validos[id] {
			return dominio.ErrValidacion("Piloto inválido")
		}
		if vistos[id] {
			return dominio.ErrValidacion("No podés repetir pilotos")
		}
		vistos[id] = true
	}
	return nil
}
