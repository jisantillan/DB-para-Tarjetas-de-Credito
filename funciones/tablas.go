package funciones

import (
	"log"
)

func crearTablas() {
	//Eliminamos las tablas si existen  y sus objetos que dependen de esta
	_, err = db.Exec(`DROP SCHEMA public CASCADE`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE SCHEMA public`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE cliente (nrocliente int, nombre text, apellido text, domicilio text, telefono char(12));
				CREATE TABLE tarjeta (nrotarjeta char(16), nrocliente int, validadesde char(6), validahasta char(6), 
					codseguridad char(4), limitecompra decimal(8,2), estado char(10));
				CREATE TABLE comercio (nrocomercio int, nombre text, domicilio text, codigopostal text, 
					telefono char(12));
				CREATE TABLE compra (nrooperacion int, nrotarjeta char(16), nrocomercio int, fecha timestamp, 
					monto decimal(7,2), pagado boolean);
				CREATE TABLE rechazo (nrorechazo int, nrotarjeta char(16), nrocomercio int, fecha timestamp, 
					monto decimal(7,2), motivo text);
				CREATE TABLE cierre (anio int, mes int, terminacion int, fechainicio date, fechacierre date, 
					fechavto date);
				CREATE TABLE cabecera (nroresumen int, nombre text, apellido text, domicilio text, nrotarjeta char(16), desde date, hasta date, vence date, total decimal(8,2));
				CREATE TABLE detalle (nroresumen int, nrolinea int, fecha date, nombrecomercio text, monto decimal(7,2));
				CREATE TABLE alerta (nroalerta int, nrotarjeta char(16), fecha timestamp, nrorechazo int, codalerta int, descripcion text);
				CREATE TABLE consumo (nrotarjeta char(16), codseguridad char(4), nrocomercio int, monto decimal(7,2))`)

	if err != nil {
		log.Fatal(err)
	}
}

func crearPK() {
	_, err = db.Exec(`ALTER TABLE cliente ADD CONSTRAINT cliente_pk PRIMARY KEY (nrocliente);
				ALTER TABLE tarjeta ADD CONSTRAINT tarjeta_pk PRIMARY KEY (nrotarjeta); 
				ALTER TABLE comercio ADD CONSTRAINT comercio_pk PRIMARY KEY (nrocomercio);
				ALTER TABLE compra ADD CONSTRAINT compra_pk PRIMARY KEY (nrooperacion); 
				ALTER TABLE rechazo ADD CONSTRAINT rechazo_pk PRIMARY KEY(nrorechazo); 
				ALTER TABLE cierre ADD CONSTRAINT cierre_pk PRIMARY KEY (anio, mes, terminacion); 
				ALTER TABLE cabecera ADD CONSTRAINT cabecera_pk PRIMARY KEY (nroresumen);
				ALTER TABLE detalle ADD CONSTRAINT detalle_pk PRIMARY KEY (nrolinea);
				ALTER TABLE alerta ADD CONSTRAINT alerta_pk PRIMARY KEY (nroalerta);`)

	if err != nil {
		log.Fatal(err)
	}
}

func crearFK() {
	_, err = db.Exec(`ALTER TABLE tarjeta ADD CONSTRAINT tarjeta_fk FOREIGN KEY (nrocliente) REFERENCES cliente(nrocliente);
				ALTER TABLE compra ADD CONSTRAINT compra_fk FOREIGN KEY (nrotarjeta) REFERENCES tarjeta(nrotarjeta);
				ALTER TABLE compra ADD CONSTRAINT compra_fk_1 FOREIGN KEY (nrocomercio) REFERENCES comercio(nrocomercio);
				ALTER TABLE rechazo ADD CONSTRAINT rechazo_fk FOREIGN KEY (nrotarjeta) REFERENCES tarjeta(nrotarjeta);
				ALTER TABLE rechazo ADD CONSTRAINT rechazo_fk_1 FOREIGN KEY (nrocomercio) REFERENCES comercio(nrocomercio);
				ALTER TABLE cabecera ADD CONSTRAINT cabecera_fk FOREIGN KEY (nrotarjeta) REFERENCES tarjeta(nrotarjeta);
				ALTER TABLE alerta ADD CONSTRAINT alerta_fk FOREIGN KEY (nrotarjeta) REFERENCES tarjeta(nrotarjeta);
				ALTER TABLE alerta ADD CONSTRAINT alerta_fk_1 FOREIGN KEY (nrorechazo) REFERENCES rechazo(nrorechazo);
				ALTER TABLE detalle add constraint detalle_fk foreign key(nroresumen) references cabecera(nroresumen);
				`)

	if err != nil {
		log.Fatal(err)
	}
}

func eliminarPK() {
	_, err = db.Exec(`ALTER TABLE cliente DROP CONSTRAINT cliente_pk ;
				ALTER TABLE tarjeta DROP CONSTRAINT tarjeta_pk;
				ALTER TABLE comercio DROP CONSTRAINT comercio_pk;
				ALTER TABLE compra DROP CONSTRAINT compra_pk;
				ALTER TABLE rechazo DROP CONSTRAINT rechazo_pk;
				ALTER TABLE cierre DROP CONSTRAINT cierre_pk;
				ALTER TABLE detalle DROP CONSTRAINT detalle_pk;
				ALTER TABLE cabecera DROP CONSTRAINT cabecera_pk;
				ALTER TABLE alerta DROP CONSTRAINT alerta_pk;`)
	if err != nil {
		log.Fatal(err)
	}
}

func eliminarFK() {
	_, err = db.Exec(`ALTER TABLE tarjeta DROP CONSTRAINT  tarjeta_fk;
				ALTER TABLE compra DROP CONSTRAINT compra_fk;
				ALTER TABLE compra DROP CONSTRAINT compra_fk_1;
				ALTER TABLE rechazo DROP CONSTRAINT rechazo_fk;
				ALTER TABLE rechazo DROP CONSTRAINT rechazo_fk_1;
				ALTER TABLE cabecera DROP CONSTRAINT cabecera_fk;
				ALTER TABLE alerta DROP CONSTRAINT alerta_fk;
				ALTER TABLE alerta DROP CONSTRAINT alerta_fk_1;
				ALTER TABLE detalle DROP CONSTRAINT detalle_fk;`)

	if err != nil {
		log.Fatal(err)
	}
}

//cliente (nro cliente int, nombre text, apellido text, domicilio text, telefono char(12));
func cargarClientes() {
	_, err = db.Exec(`INSERT INTO cliente VALUES(0, 'Juan', 'Quintero', 'Av. Pres. Figueroa Alcorta 7597', 110912201831);
				INSERT INTO cliente VALUES(1, 'Fernando', 'Quintero','Diaz Velez 1800', 112315522383);
				INSERT INTO cliente VALUES(2, 'Omar', 'Rodriguez','Hilarion Quintana 3400', 116812121837);
				INSERT INTO cliente VALUES(3, 'Josefina', 'Gonzalez','Debenedetti 2600', 110202254839);
				INSERT INTO cliente VALUES(4, 'Mario', 'Paz','Amador 2100', 110171857171);
				INSERT INTO cliente VALUES(5, 'Ines', 'Rivas','Juan de Garay 2923', 110333301833);
				INSERT INTO cliente VALUES(6, 'Lucas', 'Belli','Ugarte 2800', 113519201831);
				INSERT INTO cliente VALUES(7, 'Oscar', 'Occelli','Jose Ingenieros 2600', 110956206732);
				INSERT INTO cliente VALUES(8, 'Leonel', 'Smith','O Higgins 3500', 110842201990);
				INSERT INTO cliente VALUES(9, 'Cristian', 'Hurtado','Quintana 3000', 110919100835);
				INSERT INTO cliente VALUES(10, 'Sofia', 'Veintemilla','Avellaneda 3600', 110995292817);
				INSERT INTO cliente VALUES(11, 'Carlos', 'Vieytes','Avellaneda 3700', 111975293819);
				INSERT INTO cliente VALUES(12, 'Andres', 'Pertussi','Monteverde 3270',112404691131);
				INSERT INTO cliente VALUES(13, 'Adriana', 'Santangelo','Dorrego 2500',110512201621);
				INSERT INTO cliente VALUES(14, 'Belen', 'Quintana','Hungria 4200', 112633201888);
				INSERT INTO cliente VALUES(15, 'Juan', 'Santos','Parana 1800', 110712225059);
				INSERT INTO cliente VALUES(16, 'Ignacio', 'Fernandez','Salta 2800', 110912201831);
				INSERT INTO cliente VALUES(17, 'Marcelo Daniel', 'Bianchi','Mariano Moreno 1400',110912201831);
				INSERT INTO cliente VALUES(18, 'Roberto', 'Rosales','Bouchard 2700', 112913257738);
				INSERT INTO cliente VALUES(19, 'Malena', 'Martinez','Balcarce 1700', 111952223389);
			`)

	if err != nil {
		log.Fatal(err)
	}
}

//comercio (nrocomercio int, nombre text, domicilio text, codigopostal text, telefono char(12));
func cargarComercios() {
	_, err = db.Exec(`INSERT INTO comercio VALUES(0, 'Mostaza', 'Av. Lagomarsino 905 ', '1629',023046675071);
				INSERT INTO comercio VALUES(1, 'Grido', 'Rivadavia 875', '1629', 023204666338); 																			
				INSERT INTO comercio VALUES(2, 'Temple Brewery', 'Colectora Este Ramal Pilar Km 40,5 ', '1667',011444931532) ;
				INSERT INTO comercio VALUES(3, 'Santos BMW', 'Ruta Panamericana Ramal Pilar Km 40', '1667', 023204002695) ;
				INSERT INTO comercio VALUES(4, 'Burguer King', 'Colectora Este Ramal Pilar 113', '1629',023044713531) ;
				INSERT INTO comercio VALUES(5, 'Cine Soleil', 'Bernardo de Irigoyen 2647', '1609',080022224632) ;
				INSERT INTO comercio VALUES(6, 'Green Eat', 'Av. Cabildo 1721 ', '1426',011478711992) ;
				INSERT INTO comercio VALUES(7, 'Fullh4rd', 'Concejal Tribulato 194', '1049', 011707999972) ;
				INSERT INTO comercio VALUES(8, 'La fusa', 'Alsina 88', '1704',011465496634) ;
				INSERT INTO comercio VALUES(9, 'Music Store', 'Talcahuano 84', '1013',011521912956) ;
				INSERT INTO comercio VALUES(10, 'Supermercado Dia', 'Valentín Gómez 7040 ', '1629',023204005413);					
				INSERT INTO comercio VALUES(11, 'Hipermercado Carrefour', 'Juan Bautista Alberdi 555 ', '1623',080044484845) ;
				INSERT INTO comercio VALUES(12, 'Supermercado COTO', 'Av. Gral. Belgrano 950', '1619',034846246251) ;
				INSERT INTO comercio VALUES(13, 'El Señor de los Novillos', 'Calle 46', '1900',022142255533) ;
				INSERT INTO comercio VALUES(14, 'Ruta 26 Rock', 'Av. Ing Eduardo Madero 1343, ', '1629',011303400102) ;
				INSERT INTO comercio VALUES(15, 'Hard Rock Cafe', 'Av. Pueyrredón 2501', '1119',011480776255) ;
				INSERT INTO comercio VALUES(16, 'Harrys Digital Rock', 'Juan Francisco Segui 1422', '1615',023204155901) ;
				INSERT INTO comercio VALUES(17, 'Rulo Papa', 'Rivadavia 2687', '1613',011634523692) ;
				INSERT INTO comercio VALUES(18, 'Garbarino', 'Constitución 339', '1646',081044400182) ;
				INSERT INTO comercio VALUES(19, 'Fravega', 'Av. Pres. Juan Domingo Perón 1127', '1646',011445120632) ;
			`)

	if err != nil {
		log.Fatal(err)
	}
}

func cargarTarjetas() {
	_, err = db.Exec(`INSERT INTO tarjeta VALUES('4456102030401212',0,'201811','202403','7777',100000.00, 'vigente');
				INSERT INTO tarjeta VALUES('4570252530312020',1,'201903','202512','5656',100000.00, 'vigente');		
				INSERT INTO tarjeta VALUES('4301314121381218',2,'201705','202508','6327',100000.00, 'vigente');	
				INSERT INTO tarjeta VALUES('4351223400109919',3,'201901','202510','3141',100000.00,'vigente');
				INSERT INTO tarjeta VALUES('4636116817573601',4,'202002','202612','2138',100000.00,'vigente');
				INSERT INTO tarjeta VALUES('4663204064081011',5,'202005','202701','4305',100000.00,'suspendida');
				INSERT INTO tarjeta VALUES('4714213850503001',6,'202006','202705','2030',100000.00,'vigente');
				INSERT INTO tarjeta VALUES('4670147924096115',7,'202009','202806','1011',100000.00,'vigente');
				INSERT INTO tarjeta VALUES('4343198078796513',8,'201710','202012','1020',50000.00,'vigente');
				INSERT INTO tarjeta VALUES('4450414250909109',9,'201801','202106','9092',50000.00,'vigente');
				INSERT INTO tarjeta VALUES('4490232351236800',10,'20180','202201','4631',50000.00,'vigente');
				INSERT INTO tarjeta VALUES('5455309023212007',11,'202001','202203','1233',20000.00,'vigente');
				INSERT INTO tarjeta VALUES('4678785455121203',12,'202011','202412','4533',50000.00,'vigente');
				INSERT INTO tarjeta VALUES('5400323254578804',13,'201905','202204','1235',100000.00,'vigente');
				INSERT INTO tarjeta VALUES('4610103232659802',14,'202001','202209','3234',50000.00,'suspendida');
				INSERT INTO tarjeta VALUES('4650503031202017',15,'202005','202304','3131',50000.00,'vigente');
				INSERT INTO tarjeta VALUES('4490807522326210',16,'201803','202109','3251',50000.00,'vigente');
				INSERT INTO tarjeta VALUES('4320265751232114',17,'201703','202010','6568',200000.00,'anulada'); 
				INSERT INTO tarjeta VALUES('5421226988783205',18,'201912','202508','3201',200000.00,'vigente');
				INSERT INTO tarjeta VALUES('5658464889992908',18,'202003','202503','9752',200000.00,'vigente');
				INSERT INTO tarjeta VALUES('4421212332323516',19,'201903','202509','9421',200000.00,'vigente');
				INSERT INTO tarjeta VALUES('4421952370323306',19,'201907','202512','5537',200000.00,'vigente');
			`)
	if err != nil {
		log.Fatal(err)
	}
}

/* consumo (nrotarjeta char(16), codseguridad char(4), nrocomercio int, monto decimal(7,2))
1) happy path -> Realiza compra
2) misma tarjeta, mismo cod postal -> salta alerta pero realiza compra
3) happy path 2 > Realiza compra
4) misma tarjeta, diferentes cod postales -> salta alerta pero realiza compra
5) codigo de seguridad erróneo -> genera rechazo y alerta automaticamente
6) tarjeta suspendida -> genera rechazo y alerta automaticamente
7) tarjeta anulada -> genera rechazo y alerta automaticamente
8) supera el límite -> genera rechazo y alerta automaticamente
9) vuelve a superar el límite en el mismo dia -> genera rechazo, alerta automaticamente y alerta por separado
10) tarjeta inexistente -> genera rechazo y alerta automaticamente
*/
func cargarTablaConsumo() {
	_, err = db.Exec(`INSERT INTO consumo VALUES ('4351223400109919', '3141', 10, 1327.5);
				INSERT INTO consumo VALUES ('4351223400109919', '3141', 4, 610.50);	 
				INSERT INTO consumo VALUES ('4456102030401212', '7777', 13, 2100.5); 
				INSERT INTO consumo VALUES ('4456102030401212', '7777', 2, 750.00);
				INSERT INTO consumo VALUES ('4570252530312020', '1234', 3, 750.00);
				INSERT INTO consumo VALUES ('4610103232659802', '3234', 7, 1705.00); 
				INSERT INTO consumo VALUES ('4320265751232114', '6568', 12, 7315.63); 
				INSERT INTO consumo VALUES ('4343198078796513', '1020', 18, 51000.99);
				INSERT INTO consumo VALUES ('4343198078796513', '1020', 19, 54715.00);
				INSERT INTO consumo VALUES ('1234567898765432', '7069', 19, 6245.50); 	
			
				`)

	if err != nil {
		log.Fatal(err)
	}
}
