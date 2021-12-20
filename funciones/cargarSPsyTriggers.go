package funciones

import (
	"log"
)

func cargarTriggers() {
	_, err = db.Query(
		`CREATE TRIGGER agregar_alerta_de_rechazo
			AFTER INSERT ON rechazo
			FOR EACH ROW
			EXECUTE PROCEDURE agregar_alerta();

		CREATE TRIGGER controlar_diferencias_de_tiempos
			BEFORE INSERT ON compra
			FOR EACH ROW
			EXECUTE PROCEDURE controlar_tiempo_compras();
		`)
	if err != nil {
		log.Fatal(err)
	}
}
func cargarSPs() {
	spAgregarAlertaRechazo()
	spAgregarRechazo()
	spChequearLimiteRechazos()
	spAutorizarCompra()
	spSeguridadTiempoCompras()
	spGenerarResumen()
	spAgregarProbarConsumos()
}
