package funciones

import (
	"encoding/json"
	"log"
	"strconv"

	bolt "github.com/boltdb/bolt"
)

var dbbolt *bolt.DB

type Cliente struct {
	Nrocliente int
	Nombre     string
	Apellido   string
	Domicilio  string
	Telefono   string
}

type Tarjeta struct {
	Nrotarjeta   string
	Nrocliente   int
	Validadesde  string
	Validahasta  string
	Codseguridad string
	Limitecompra int
	Estado       string
}

type Comercio struct {
	Nrocomercio  int
	Nombre       string
	Domicilio    string
	Codigopostal string
	Telefono     string
}

type Compra struct {
	Nrooperacion int
	Nrotarjeta   string
	Nrocomercio  int
	Fecha        string
	Monto        int
	Pagado       bool
}

func dbBoltConnection() {
	dbbolt, err = bolt.Open("test.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func cargarClienteBolt(nrocliente int, nombre string, apellido string, domicilio string, telefono string) {
	cliente := Cliente{nrocliente, nombre, apellido, domicilio, telefono}
	data, err := json.Marshal(cliente)
	if err != nil {
		log.Fatal(err)
	}

	CreateUpdate(dbbolt, "Cliente", []byte(strconv.Itoa(cliente.Nrocliente)), data)
}

func cargarTarjetaBolt(nrotarjeta string, nrocliente int, validadesde string, validahasta string, codseguridad string, limitecompra int, estado string) {
	tarjeta := Tarjeta{nrotarjeta, nrocliente, validadesde, validahasta, codseguridad, limitecompra, estado}

	data, err := json.Marshal(tarjeta)
	if err != nil {
		log.Fatal(err)
	}

	CreateUpdate(dbbolt, "Tarjeta", []byte(tarjeta.Nrotarjeta), data)
}

func cargarComercioBolt(nrocomercio int, nombre string, domicilio string, codigopostal string, telefono string) {
	comercio := Comercio{nrocomercio, nombre, domicilio, codigopostal, telefono}

	data, err := json.Marshal(comercio)
	if err != nil {
		log.Fatal(err)
	}

	CreateUpdate(dbbolt, "Comercio", []byte(strconv.Itoa(comercio.Nrocomercio)), data)
}

func cargarCompraBolt(nrooperacion int, nrotarjeta string, nrocomercio int, fecha string, monto int, pagado bool) {
	compra := Compra{nrooperacion, nrotarjeta, nrocomercio, fecha, monto, pagado}

	data, err := json.Marshal(compra)
	if err != nil {
		log.Fatal(err)
	}
	CreateUpdate(dbbolt, "Compra", []byte(strconv.Itoa(compra.Nrooperacion)), data)
}

func CreateUpdate(dbbolt *bolt.DB, bucketName string, key []byte, val []byte) error {
	//abre transaccion de escritura
	tx, err := dbbolt.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	b, _ := tx.CreateBucketIfNotExists([]byte(bucketName))

	err = b.Put(key, val)
	if err != nil {
		return err
	}

	//cierra transaccion
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func ReadUnique(dbbolt *bolt.DB, bucketName string, key []byte) ([]byte, error) {
	var buf []byte

	err := dbbolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		buf = b.Get(key)
		return nil
	})

	return buf, err
}

func CargarDatosBolt() {
	dbBoltConnection()

	defer dbbolt.Close()

	cargarClienteBolt(1, "Jorge", "Martinez", "Av. Presidente Peron 5567", "1156789033")
	cargarClienteBolt(2, "Maria", "Romero", "General Lemos 2233", "1123456789")
	cargarClienteBolt(3, "Juan", "Lazarotti", "Rivadavia 1234", "1145678909")

	cargarTarjetaBolt("1023455667789887", 1, "201506", "202007", "7001", 5000, "vigente")
	cargarTarjetaBolt("1209988776655443", 2, "201709", "202012", "9999", 7000, "vigente")
	cargarTarjetaBolt("1233445566778899", 3, "201410", "202012", "6666", 8000, "vigente")

	cargarComercioBolt(501, "Kevingston", "Av. Ricchieri 987", "1661", "46664566")
	cargarComercioBolt(523, "LaoLao", "Paunero 553", "1612", "47595566")
	cargarComercioBolt(551, "Musimundo", "Belgrano 998", "1677", "44556677")

	cargarCompraBolt(1, "1023455667789887", 501, "2020-04-25 00:00:00", 1500.00, true)
	cargarCompraBolt(2, "1023455667789887", 523, "2020-04-25 00:00:00", 340.00, true)
	cargarCompraBolt(3, "1023455667789887", 551, "2020-04-25 00:00:00", 3000.00, true)
}
