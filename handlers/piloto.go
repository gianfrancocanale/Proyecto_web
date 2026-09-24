package handlers

import (
	sqlc "Proyecto_web/db/sqlc"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func CrearPiloto(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var piloto sqlc.CrearPilotoTemporadaParams

	err := json.NewDecoder(r.Body).Decode(&piloto)

	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	piloto.PuntosTemporada = 0 // el piloto arranca con 0 puntos

	queries := sqlc.New(db)

	_, err = queries.RecuperarPilotoHistorico(r.Context(), piloto.IDPiloto)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No existe el piloto", http.StatusBadRequest)
		} else {
			http.Error(w, "Error al verificar el piloto", http.StatusInternalServerError)
		}
		return
	}

	_, err = queries.RecuperarEscuderia(r.Context(), piloto.IDEscuderia)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No existe el piloto en esa escudería", http.StatusBadRequest)
		} else {
			http.Error(w, "Error al verificar la escudería", http.StatusInternalServerError)
		}
		return
	}

	devolver, err := queries.CrearPilotoTemporada(r.Context(), piloto)
	if err != nil {
		http.Error(w, "Error al crear el piloto", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(devolver)
}

func ListarPilotos(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	queries := sqlc.New(db)

	pilotos, err := queries.ListarPilotos(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener los pilotos", http.StatusInternalServerError)
		return
	}
	if len(pilotos) == 0 {
		http.Error(w, "No se encontraron pilotos", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pilotos)
}

func ObtenerPiloto(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idPiloto, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "ID de piloto histórico inválido", http.StatusBadRequest)
		return
	}

	queries := sqlc.New(db)

	piloto, err := queries.RecuperarPiloto(r.Context(), int32(idPiloto))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Piloto no encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Error al obtener el piloto", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(piloto)
}

func ActualizarPiloto(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idPiloto, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "ID de piloto inválido", http.StatusBadRequest)
		return
	}

	var piloto sqlc.ModificarPilotoParams

	err = json.NewDecoder(r.Body).Decode(&piloto)
	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	//if piloto.id_piloto == {
	//	http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
	//	return
	//}

	piloto.IDPiloto = int32(idPiloto)

	queries := sqlc.New(db)

	devolver, err := queries.ModificarPiloto(r.Context(), piloto)
	if err != nil {
		http.Error(w, "Error al actualizar el piloto", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devolver)
}

func EliminarPiloto(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idPiloto, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "ID de piloto inválido", http.StatusBadRequest)
		return
	}

	queries := sqlc.New(db)

	err = queries.EliminarPilotoTemporada(r.Context(), int32(idPiloto))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Piloto no encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Error al eliminar el piloto", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
