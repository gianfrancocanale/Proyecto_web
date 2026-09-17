package main

import (
	sqlc "Proyecto_web/db/sqlc"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	db := conectarBase()
	fmt.Printf("La base de datos se inicio")
	defer db.Close()
	http.HandleFunc("POST /usuarios", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		crearUsuario(db, w, r)
	})

	/*http.HandleFunc("GET /usuarios", func(w http.ResponseWriter, r *http.Request) {
		listarUsuarios(db, w, r)
	})

	http.HandleFunc("GET /usuarios/{id}", func(w http.ResponseWriter, r *http.Request) {
		obtenerUsuario(db, w, r)
	})

	http.HandleFunc("PUT /usuarios/{id}", func(w http.ResponseWriter, r *http.Request) {
		actualizarUsuario(db, w, r)
	})

	http.HandleFunc("DELETE /usuarios/{id}", func(w http.ResponseWriter, r *http.Request) {
		eliminarUsuario(db, w, r)
	})*/

	fmt.Println("Servidor escuchando en :8080")
	http.ListenAndServe(":8080", nil)
}

func crearUsuario(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	queries := sqlc.New(db)

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

	devolver, err := queries.CrearUsuario(r.Context(), usuario)

	if err != nil {
		http.Error(w, "Error al crear el usuario", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Usuario creado exitosamente")
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
