package dominio

import (
	"context"
	"errors"
	"strings"
	"time"
)

// Errores de dominio de Piloto.
var (
	ErrPilotoDuplicado     = errors.New("ese piloto ya está en la temporada")
	ErrPilotoNoEncontrado  = errors.New("piloto no encontrado")
	ErrEscuderiaNoExiste   = errors.New("no existe la escudería")
	ErrHistoricoIncompleto = errors.New("piloto nuevo: completá nombre, país y fecha de nacimiento")
)

// Piloto es un piloto de la temporada, con los datos históricos ya unidos.
type Piloto struct {
	ID         int       `json:"id_piloto"`
	Nombre     string    `json:"nombre"`
	Pais       string    `json:"pais"`
	Nacimiento time.Time `json:"fecha_nacimiento"`
	Titulos    int       `json:"titulos_ganados"`
	Escuderia  string    `json:"id_escuderia"`
	Puntos     int       `json:"puntos_temporada"`
}

// Validar aplica las reglas de la entidad.
func (p *Piloto) Validar() error {
	p.Nombre = strings.TrimSpace(p.Nombre)
	p.Pais = strings.TrimSpace(p.Pais)
	p.Escuderia = strings.TrimSpace(p.Escuderia)

	if p.ID <= 0 {
		return ErrValidacion("Número de piloto inválido")
	}
	if p.Escuderia == "" {
		return ErrValidacion("Elegí una escudería")
	}
	if p.Puntos < 0 {
		return ErrValidacion("Los puntos tienen que ser un número mayor o igual a 0")
	}
	if p.Titulos < 0 {
		return ErrValidacion("Los títulos ganados tienen que ser un número mayor o igual a 0")
	}
	return nil
}

// PilotoRepository es el puerto de persistencia de pilotos.
type PilotoRepository interface {
	Listar(ctx context.Context) ([]Piloto, error)
	Recuperar(ctx context.Context, id int) (Piloto, error)
	ListarParaApuesta(ctx context.Context) ([]Piloto, error)
	Escuderias(ctx context.Context) ([]string, error)
	ExisteEscuderia(ctx context.Context, id string) (bool, error)
	ExisteEnTemporada(ctx context.Context, id int) (bool, error)
	ExisteHistorico(ctx context.Context, id int) (bool, error)
	// Inscribir crea (si hace falta) el histórico y anota el piloto en la temporada, en una transacción.
	Inscribir(ctx context.Context, p Piloto, crearHistorico bool) error
	Modificar(ctx context.Context, p Piloto) error
	Eliminar(ctx context.Context, id int) error
}
