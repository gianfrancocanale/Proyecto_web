package handlers

import (
	sqlc "Proyecto_web/db/sqlc"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func CrearPilotoHistorico(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var pilotoHistorico sqlc.CrearPilotoHistoricoParams

	err := json.NewDecoder(r.Body).Decode(&pilotoHistorico)

	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	pilotoHistorico.TitulosGanados = 0 //setea 0 titulos ganados al crear un nuevo piloto historico

	if pilotoHistorico.Nombre == "" || pilotoHistorico.Pais == "" || pilotoHistorico.FechaNacimiento.IsZero() {
		http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
		return
	}

	queries := sqlc.New(db)

	devolver, err := queries.CrearPilotoHistorico(r.Context(), pilotoHistorico)
	if err != nil {
		http.Error(w, "Error al crear el piloto histórico", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(devolver)
}

func ListarPilotosHistoricos(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	queries := sqlc.New(db)

	pilotosHistoricos, err := queries.ListarPilotosHistoricos(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener los pilotos históricos", http.StatusInternalServerError)
		return
	}
	if len(pilotosHistoricos) == 0 {
		http.Error(w, "No se encontraron pilotos históricos", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pilotosHistoricos)
}

func ObtenerPilotoHistorico(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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

	pilotoHistorico, err := queries.RecuperarPilotoHistorico(r.Context(), int32(idPiloto))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Piloto histórico no encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Error al obtener el piloto histórico", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pilotoHistorico)
}

func ActualizarPilotoHistorico(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idPiloto, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "ID de piloto histórico inválido", http.StatusBadRequest)
		return
	}

	var pilotoHistorico sqlc.ModificarPilotoHistoricoParams

	err = json.NewDecoder(r.Body).Decode(&pilotoHistorico)
	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	if pilotoHistorico.Nombre == "" || pilotoHistorico.Pais == "" || pilotoHistorico.FechaNacimiento.IsZero() {
		http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
		return
	}

	pilotoHistorico.IDPiloto = int32(idPiloto)

	queries := sqlc.New(db)

	devolver, err := queries.ModificarPilotoHistorico(r.Context(), pilotoHistorico)
	if err != nil {
		http.Error(w, "Error al actualizar el piloto histórico", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devolver)
}

func EliminarPilotoHistorico(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idPiloto, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "ID de piloto histórico inválido", http.StatusBadRequest)
		return
	}

	queries := sqlc.New(db)

	err = queries.EliminarPilotoHistorico(r.Context(), int32(idPiloto))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Piloto histórico no encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Error al eliminar el piloto histórico", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
