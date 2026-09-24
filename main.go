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

	// MANEJO DE RUTAS DE ENTIDAD USUARIO:
	http.HandleFunc("/usuarios", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.CrearUsuario(db, w, r)
		case http.MethodGet:
			h.ListarUsuarios(db, w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			fmt.Println("entro al default")
		}
	})

	http.HandleFunc("/usuarios/{id}", func(w http.ResponseWriter, r *http.Request) {
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

	})

	//MANEJO DE RUTAS DE ENTIDAD PILOTO HISTORICO:
	http.HandleFunc("/pilotos_historicos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.CrearPilotoHistorico(db, w, r)
		case http.MethodGet:
			h.ListarPilotosHistoricos(db, w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/pilotos_historicos/{id}", func(w http.ResponseWriter, r *http.Request) {
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
	})

	//MANEJO DE RUTAS DE ENTIDAD ESCUDERIA:
	http.HandleFunc("/escuderias", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.CrearEscuderia(db, w, r)
		case http.MethodGet:
			h.ListarEscuderias(db, w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/escuderias/{id}", func(w http.ResponseWriter, r *http.Request) {
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
	})

	fmt.Println(" ")
	fmt.Println("Servidor escuchando en :8080")
	http.ListenAndServe(":8080", nil)
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
