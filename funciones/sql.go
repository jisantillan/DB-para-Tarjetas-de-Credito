package funciones

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB
var err error

func dbConnection() {
	db, err = sql.Open("postgres", "user=postgres host=localhost dbname=tarjeta sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}
}

func CrearDB() {
	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=postgres sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	_, err = db.Exec(`CREATE DATABASE tarjeta`)

	if err != nil {
		log.Fatal(err)
	}
}

func BorrarDB() {
	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=postgres sslmode=disable")
	defer db.Close()

	_, err = db.Exec(`DROP DATABASE IF EXISTS tarjeta`)

	if err != nil {
		log.Fatal(err)
	}
}

func CrearTablas() {
	dbConnection()
	defer db.Close()
	crearTablas()
}

func CrearPKyFK() {
	dbConnection()
	defer db.Close()
	crearPK()
	crearFK()
}

func EliminarPKyFK() {
	dbConnection()
	defer db.Close()
	eliminarFK()
	eliminarPK()
}

func CargarTablas() {
	dbConnection()
	defer db.Close()
	cargarClientes()
	cargarComercios()
	cargarCierres()
	cargarTarjetas()
	cargarTablaConsumo()
}

func CargarSPs_y_triggers() {
	dbConnection()
	defer db.Close()
	cargarSPs()
	cargarTriggers()
}

func Realizar_Resumen() {
	dbConnection()
	defer db.Close()

	/* probamos de acuerdo a las compras cargadas,
	cliente 3 realizo dos compras
	cliente 0 realizo dos compras
	cliente 4 no realizo compras*/
	_, err = db.Query(

		`	
		SELECT generar_resumen(12,2020,3);
		SELECT generar_resumen(12,2020,0);
		SELECT generar_resumen(12,2020,4);
		`)
	if err != nil {
		log.Fatal(err)
	}
}

func Realizar_compras() {
	dbConnection()
	defer db.Close()
	_, err = db.Query(`SELECT probar_consumos();`)

	if err != nil {
		log.Fatal(err)
	}
}
