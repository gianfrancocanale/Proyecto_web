package main

import (
	h "Proyecto_web/handlers"
	"database/sql" //para manejar errores(id no encontrados)
	"fmt"
	"net/http"
	"os" //obtener datos de variables de entorno(bases de datos)

	//convertir string a int

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	db := conectarBase()
	fmt.Printf("La base de datos se inicio")
	defer db.Close()

	// MANEJO DE RUTAS DE ENTIDAD USUARIO:
	http.HandleFunc("/usuarios", handlerUsuario)
	http.HandleFunc("/usuarios/{id}", handlerUsuarioId)

	//MANEJO DE RUTAS DE ENTIDAD PILOTO HISTORICO:
	http.HandleFunc("/pilotos_historicos", handlerPilotoHistorico)
	http.HandleFunc("/pilotos_historicos/{id}", handlerPilotoHistoricoId)

	//MANEJO DE RUTAS DE ENTIDAD ESCUDERIA:
	http.HandleFunc("/escuderias", handlerEscuderia)
	http.HandleFunc("/escuderias/{id}", handlerEscuderiaId)

	//MANEJO DE RUTAS DE ENTIDAD PILOTO:
	http.HandleFunc("/pilotos", handlerPiloto)
	http.HandleFunc("/pilotos/{id}", handlerPilotoId)

	//MANEJO DE RUTAS DE ENTIDAD GRAN PREMIO:
	http.HandleFunc("/gran_premios", h.HandlerGranPremio)
	http.HandleFunc("/gran_premios/{id}", h.HandlerGranPremioId)

	fmt.Println(" ")
	fmt.Println("Servidor escuchando en :8080")
	http.ListenAndServe(":8080", nil)
}

func handlerGranPremio(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CrearGranPremio(db, w, r)
	case http.MethodGet:
		h.ListarGranPremios(db, w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func handlerGranPremioId(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.ObtenerGranPremio(db, w, r)
	case http.MethodPut:
		h.ActualizarGranPremio(db, w, r)
	case http.MethodDelete:
		h.EliminarGranPremio(db, w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func handlerPiloto(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CrearPiloto(db, w, r)
	case http.MethodGet:
		h.ListarPilotos(db, w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func handlerPilotoId(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.ObtenerPiloto(db, w, r)
	case http.MethodPut:
		h.ActualizarPiloto(db, w, r)
	case http.MethodDelete:
		h.EliminarPiloto(db, w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func handlerUsuario(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CrearUsuario(db, w, r)
	case http.MethodGet:
		h.ListarUsuarios(db, w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		fmt.Println("entro al default")
	}
}

func handlerUsuarioId(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.ObtenerUsuario(db, w, r)
	case http.MethodPut:
		h.ActualizarUsuario(db, w, r)
	case http.MethodDelete:
		h.EliminarUsuario(db, w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}

}

func handlerPilotoHistorico(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CrearPilotoHistorico(db, w, r)
	case http.MethodGet:
		h.ListarPilotosHistoricos(db, w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func handlerPilotoHistoricoId(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.ObtenerPilotoHistorico(db, w, r)
	case http.MethodPut:
		h.ActualizarPilotoHistorico(db, w, r)
	case http.MethodDelete:
		h.EliminarPilotoHistorico(db, w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func handlerEscuderia(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CrearEscuderia(db, w, r)
	case http.MethodGet:
		h.ListarEscuderias(db, w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func handlerEscuderiaId(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.ObtenerEscuderia(db, w, r)
	case http.MethodPut:
		h.ActualizarEscuderia(db, w, r)
	case http.MethodDelete:
		h.EliminarEscuderia(db, w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
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
