package dominio

import (
	"context"
	"errors"
	"time"
)

// Errores de dominio de Apuesta.
var (
	ErrApuestaDuplicada    = errors.New("ya apostaste en ese gran premio")
	ErrApuestaNoEncontrada = errors.New("apuesta no encontrada")
)

// Puntaje de las apuestas: regla de negocio pura.
const (
	PuntosPosicionExacta = 3
	PuntosEnTop10        = 1
)

// Apuesta es el top 10 pronosticado por un usuario para una carrera.
type Apuesta struct {
	ID           int       `json:"id_apuesta"`
	Usuario      string    `json:"nombre_usuario"`
	GranPremio   string    `json:"id_gran_premio"`
	FechaCarrera time.Time `json:"fecha_carrera"`
	Prediccion   []int     `json:"prediccion"`
	FechaApuesta time.Time `json:"fecha_apuesta"`
}

// ValidarPrediccion exige 10 pilotos distintos y que todos estén en la grilla.
func (a *Apuesta) ValidarPrediccion(pilotosValidos []int) error {
	if len(a.Prediccion) != 10 {
		return ErrValidacion("Tenés que elegir 10 pilotos")
	}
	validos := make(map[int]bool, len(pilotosValidos))
	for _, id := range pilotosValidos {
		validos[id] = true
	}
	vistos := make(map[int]bool, 10)
	for _, id := range a.Prediccion {
		if !validos[id] {
			return ErrValidacion("Piloto inválido")
		}
		if vistos[id] {
			return ErrValidacion("No podés repetir pilotos")
		}
		vistos[id] = true
	}
	return nil
}

// CalcularPuntos compara el top 10 apostado contra los primeros 10 del resultado.
func CalcularPuntos(prediccion, resultado []int) int {
	top10 := resultado
	if len(top10) > 10 {
		top10 = top10[:10]
	}
	estaEnTop10 := make(map[int]bool, len(top10))
	for _, id := range top10 {
		estaEnTop10[id] = true
	}

	pts := 0
	for i, id := range prediccion {
		switch {
		case i < len(top10) && top10[i] == id:
			pts += PuntosPosicionExacta
		case estaEnTop10[id]:
			pts += PuntosEnTop10
		}
	}
	return pts
}

// ApuestaRepository es el puerto de persistencia de apuestas.
type ApuestaRepository interface {
	ListarPorUsuario(ctx context.Context, usuario string) ([]Apuesta, error)
	ListarPorGranPremio(ctx context.Context, id string, fecha time.Time) ([]Apuesta, error)
	Recuperar(ctx context.Context, id int, usuario string) (Apuesta, error)
	Crear(ctx context.Context, a Apuesta) error
	Actualizar(ctx context.Context, a Apuesta) error
	Eliminar(ctx context.Context, id int, usuario string) error
}
