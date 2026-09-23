package main

import (
	sqlc "Proyecto_web/db/sqlc"
	"database/sql"
	"encoding/json"
	"errors" //para manejar errores(id no encontrados)
	"fmt"
	"net/http"
	"os"      //obtener datos de variables de entorno(bases de datos)
	"strconv" //convertir string a int

	// "Proyecto_web/handlers"

	_ "github.com/lib/pq"
)

func main() {
	db := conectarBase()
	fmt.Printf("La base de datos se inicio")
	defer db.Close()

	// MANEJO DE RUTAS DE ENTIDAD USUARIO:
	http.HandleFunc("/usuarios", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			crearUsuario(db, w, r)
		case http.MethodGet:
			listarUsuarios(db, w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}

		crearUsuario(db, w, r)
	})

	http.HandleFunc("/usuarios/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			obtenerUsuario(db, w, r)
		case http.MethodPut:
			actualizarUsuario(db, w, r)
		case http.MethodDelete:
			eliminarUsuario(db, w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	//MANEJO DE RUTAS DE ENTIDAD PILOTO HISTORICO:
	http.HandleFunc("/pilotos_historicos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			crearPilotoHistorico(db, w, r)
		case http.MethodGet:
			listarPilotosHistoricos(db, w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/pilotos_historicos/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			obtenerPilotoHistorico(db, w, r)
		case http.MethodPut:
			actualizarPilotoHistorico(db, w, r)
		case http.MethodDelete:
			eliminarPilotoHistorico(db, w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	//MANEJO DE RUTAS DE ENTIDAD ESCUDERIA:
	http.HandleFunc("/escuderias", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			crearEscuderia(db, w, r)
		case http.MethodGet:
			listarEscuderias(db, w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/escuderias/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			obtenerEscuderia(db, w, r)
		case http.MethodPut:
			actualizarEscuderia(db, w, r)
		case http.MethodDelete:
			eliminarEscuderia(db, w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println(" ")
	fmt.Println("Servidor escuchando en :8080")
	http.ListenAndServe(":8080", nil)
}

// ESCUDERIA
func crearEscuderia(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var escuderia sqlc.CrearEscuderiaParams

	err := json.NewDecoder(r.Body).Decode(&escuderia)

	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	if escuderia.IDEscuderia == "" || escuderia.TeamPrincipal == "" || escuderia.FechaFundacion.IsZero() {
		http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
		return
	}

	escuderia.PuntosTemporada = 0      //setea 0 puntos al crear una nueva escuderia
	escuderia.TitulosConstructores = 0 //setea 0 titulos al crear una nueva escuderia

	queries := sqlc.New(db)

	devolver, err := queries.CrearEscuderia(r.Context(), escuderia)
	if err != nil {
		http.Error(w, "Error al crear la escudería", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(devolver)
}

func listarEscuderias(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	queries := sqlc.New(db)

	escuderias, err := queries.ListarEscuderias(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener las escuderías", http.StatusInternalServerError)
		return
	}
	if len(escuderias) == 0 {
		http.Error(w, "No se encontraron escuderías", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(escuderias)
}

func obtenerEscuderia(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idEscuderia := r.PathValue("id")

	queries := sqlc.New(db)

	escuderia, err := queries.RecuperarEscuderia(r.Context(), idEscuderia)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Escudería no encontrada", http.StatusNotFound)
		} else {
			http.Error(w, "Error al obtener la escudería", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(escuderia)
}

func actualizarEscuderia(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idEscuderia := r.PathValue("id")

	var escuderia sqlc.ModificarEscuderiaParams

	err := json.NewDecoder(r.Body).Decode(&escuderia)
	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	if escuderia.TeamPrincipal == "" || escuderia.FechaFundacion.IsZero() {
		http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
		return
	}

	escuderia.IDEscuderia = idEscuderia

	queries := sqlc.New(db)

	devolver, err := queries.ModificarEscuderia(r.Context(), escuderia)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Escudería no encontrada", http.StatusNotFound)
		} else {
			http.Error(w, "Error al actualizar la escudería", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devolver)
}

func eliminarEscuderia(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idEscuderia := r.PathValue("id")

	queries := sqlc.New(db)

	err := queries.EliminarEscuderia(r.Context(), idEscuderia)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Escudería no encontrada", http.StatusNotFound)
		} else {
			http.Error(w, "Error al eliminar la escudería", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// PILOTO HISTORICO
func crearPilotoHistorico(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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

func listarPilotosHistoricos(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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

func obtenerPilotoHistorico(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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

func actualizarPilotoHistorico(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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

func eliminarPilotoHistorico(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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

// USUARIO
func eliminarUsuario(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idUsuario, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "ID de usuario invalido", http.StatusBadRequest)
		return
	}

	queries := sqlc.New(db)

	err = queries.EliminarUsuario(r.Context(), int32(idUsuario))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Error al eliminar el usuario", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func actualizarUsuario(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idUsuario, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "ID de usuario invalido", http.StatusBadRequest)
		return
	}

	var usuario sqlc.ModificarUsuarioParams

	err = json.NewDecoder(r.Body).Decode(&usuario)
	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	if usuario.NombreUsuario == "" || usuario.ContrasenaHash == "" {
		http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
		return
	}

	if usuario.PuntosGanados < 0 {
		http.Error(w, "Los puntos ganados no pueden ser negativos", http.StatusBadRequest)
		return
	}

	usuario.IDUsuario = int32(idUsuario)

	queries := sqlc.New(db)

	devolver, err := queries.ModificarUsuario(r.Context(), usuario)
	if err != nil {
		http.Error(w, "Error al actualizar el usuario", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devolver)
}

func obtenerUsuario(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	idUsuario, err := strconv.ParseInt(r.PathValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "ID de usuario invalido", http.StatusBadRequest)
		return
	}

	queries := sqlc.New(db)

	usuario, err := queries.RecuperarUsuarioPorId(r.Context(), int32(idUsuario))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		} else {
			http.Error(w, "Error al obtener el usuario", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuario)
}

func listarUsuarios(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	queries := sqlc.New(db)

	usuarios, err := queries.RecuperarUsuarios(r.Context())
	if err != nil {
		http.Error(w, "Error al obtener los usuarios", http.StatusInternalServerError)
		return
	}
	if len(usuarios) == 0 {
		http.Error(w, "No se encontraron usuarios", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuarios)
}

func crearUsuario(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var usuario sqlc.CrearUsuarioParams

	err := json.NewDecoder(r.Body).Decode(&usuario)
	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	if usuario.NombreUsuario == "" || usuario.ContrasenaHash == "" {
		http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
		return
	}

	usuario.PuntosGanados = 0 //setea 0 puntos ganados al crear un nuevo usuario

	queries := sqlc.New(db)

	devolver, err := queries.CrearUsuario(r.Context(), usuario)

	if err != nil {
		http.Error(w, "Error al crear el usuario", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(devolver)
}

func conectarBase() *sql.DB {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Errorf("Error al conectar con la base de datos: %v", err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Errorf("Error al hacer ping a la base de datos: %v", err)
	}

	return db
}
