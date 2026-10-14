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

func main() {
	db := conectarBase()
	fmt.Printf("La base de datos se inicio")
	defer db.Close()

	hUsuario := h.NewUsuarioHandler(db)
	hPilotoHistorico := h.NewPilotoHistoricoHandler(db)
	hPiloto := h.NewPilotoHandler(db)
	hEscuderia := h.NewEscuderiaHandler(db)
	hGranPremio := h.NewGranPremioHandler(db)
	hApuesta := h.NewApuestaHandler(db)
	hGranPremioHistorico := h.NewGranPremioHistoricoHandler(db)

	// MANEJO DE RUTAS DE ENTIDAD USUARIO:
	http.HandleFunc("/usuarios", hUsuario.HandlerUsuario)
	http.HandleFunc("/usuarios/{nombre}", hUsuario.HandlerUsuarioId)

	//MANEJO DE RUTAS DE ENTIDAD PILOTO HISTORICO:
	http.HandleFunc("/pilotos_historicos", hPilotoHistorico.HandlerPilotoHistorico)
	http.HandleFunc("/pilotos_historicos/{id}", hPilotoHistorico.HandlerPilotoHistoricoId)

	//MANEJO DE RUTAS DE ENTIDAD ESCUDERIA:
	http.HandleFunc("/escuderias", hEscuderia.HandlerEscuderia)
	http.HandleFunc("/escuderias/{id}", hEscuderia.HandlerEscuderiaId)

	//MANEJO DE RUTAS DE ENTIDAD PILOTO:
	http.HandleFunc("/pilotos", hPiloto.HandlerPiloto)
	http.HandleFunc("/pilotos/{id}", hPiloto.HandlerPilotoId)

	//MANEJO DE RUTAS DE ENTIDAD GRAN PREMIO:
	http.HandleFunc("/gran_premios", hGranPremio.HandlerGranPremio)
	http.HandleFunc("/gran_premios/{id}", hGranPremio.HandlerGranPremioId)

	//MANEJO DE RUTAS DE ENTIDAD APUESTA:
	http.HandleFunc("/apuestas", hApuesta.HandlerApuesta)
	http.HandleFunc("/apuestas/{id}/{usuario}", hApuesta.HandlerApuestaIdUsuario)
	http.HandleFunc("/apuestas/{usuario}", hApuesta.HandlerApuestaUsuario)

	//MANEJO DE RUTAS DE ENTIDAD GRAN PREMIO HISTORICO:
	http.HandleFunc("/gran_premios_historicos", hGranPremioHistorico.HandlerGranPremioHistorico)
	http.HandleFunc("/gran_premios_historicos/{id}/{fecha_carrera}", hGranPremioHistorico.HandlerGranPremioHistoricoId)

	fmt.Println(" ")
	fmt.Println("Servidor escuchando en :8080")
	http.ListenAndServe(":8080", nil)
}

func handlerApuestaUsuario(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.ObtenerApuestaUsuario(db, w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func handlerGranPremioHistorico(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CrearGranPremioHistorico(db, w, r)
	case http.MethodGet:
		h.ListarGrandesPremiosHistoricos(db, w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func handlerGranPremioHistoricoId(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.ObtenerGranPremioHistorico(db, w, r)
	case http.MethodPut:
		h.ActualizarGranPremioHistorico(db, w, r)
	case http.MethodDelete:
		h.EliminarGranPremioHistorico(db, w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
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
