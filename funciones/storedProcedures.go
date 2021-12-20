package funciones

import (
	"log"
)

func cargarCierres() {
	spGenerarCierres()
	_, err = db.Query(`SELECT generarCierres(2020);`)
	if err != nil {
		log.Fatal(err)
	}
}

//Fechas de cierre  cada 30 dias/ primero probamos como genera con 29 o 30 dias
func spGenerarCierres() {
	_, err = db.Query(
		`CREATE OR REPLACE FUNCTION generarCierres(anio int) RETURNS VOID AS $$
			DECLARE 
				fdesde date;
				fhasta date;
				fvto date;
			BEGIN
				FOR tarjeta IN 0 .. 9 BY 1 LOOP 
					SELECT INTO fdesde to_date((anio-1)::text || '12' || (SELECT 23 + TRUNC(RANDOM()*4))::text, 'YYYYMMDD');
					SELECT INTO fhasta fdesde::date + CAST((SELECT 29 + TRUNC(RANDOM()*2))::text||' days' AS interval);
					SELECT INTO fvto fhasta::date + CAST('10 days' AS interval);
										
					FOR mes IN 1 .. 12 BY 1 LOOP
						INSERT INTO cierre values(anio, mes, tarjeta, fdesde, fhasta, fvto);
											
						SELECT INTO fdesde fhasta::date + CAST('1 days' AS interval);
						SELECT INTO fhasta fdesde::date + CAST((SELECT 29 + TRUNC(random()*2))::text || ' days' AS interval);
						SELECT INTO fvto fhasta::date + CAST('10 days' AS interval);
					END LOOP;
				END LOOP;
			END;
		$$ LANGUAGE PLPGSQL;`)

	if err != nil {
		log.Fatal(err)
	}
}

func spAgregarRechazo() {
	_, err = db.Query(
		`CREATE OR REPLACE FUNCTION agregarrechazo(_nrotarjeta char(16), _nrocomercio int, _fecha timestamp, _monto decimal(7,2), _motivo text) RETURNS VOID AS $$
				
				DECLARE
					_nrorechazo int;
				BEGIN
				
					PERFORM * FROM rechazo r WHERE r.nrorechazo IS NOT NULL;
					IF (FOUND) THEN
						SELECT MAX(nrorechazo)+1 INTO _nrorechazo FROM rechazo;
					ELSE
						_nrorechazo := 1;
					END IF;
					
					INSERT INTO rechazo(nrorechazo, nrotarjeta, nrocomercio, fecha, monto, motivo) VALUES (_nrorechazo, _nrotarjeta, _nrocomercio, current_timestamp, _monto, _motivo);
							
					PERFORM chequearLimiteRechazos(_nrorechazo);
				END;
			$$ LANGUAGE PLPGSQL;`)
	if err != nil {
		log.Fatal(err)
	}
}

//chequea si una tarjeta recibe 2 rechazos en el mismo dia
func spChequearLimiteRechazos() {
	_, err = db.Query(
		`CREATE OR REPLACE FUNCTION chequearLimiteRechazos(_nroderechazo int) RETURNS VOID AS $$
				DECLARE 
					_nrotarjetarechazada char(16);
					_fecharechazo timestamp;
					_nroalerta int;
				BEGIN 
					SELECT nrotarjeta, fecha INTO _nrotarjetarechazada, _fecharechazo FROM rechazo r WHERE r.nrorechazo = _nroderechazo;
					
					PERFORM * FROM alerta a WHERE a.nroalerta IS NOT NULL;
					IF (FOUND) THEN
						SELECT MAX(nroalerta)+1 INTO _nroalerta FROM alerta;
					ELSE
						_nroalerta := 1;
					END IF;
					
					PERFORM nrotarjeta FROM rechazo
						WHERE
							nrotarjeta = _nrotarjetarechazada 
							AND fecha = _fecharechazo 
							AND motivo = 'Supera limite de tarjeta' 
							GROUP BY nrotarjeta HAVING COUNT(*) > 1;
						
					IF (FOUND) THEN
						INSERT INTO alerta(nroalerta, nrotarjeta, fecha, nrorechazo,codalerta,descripcion)
						VALUES (_nroalerta, _nrotarjetarechazada, _fecharechazo, _nroderechazo,32, 'Tarjeta suspendida');
							 
						 UPDATE tarjeta  SET estado = 'suspendida' WHERE nrotarjeta = _nrotarjetarechazada;
					END IF;
				END;
			$$ LANGUAGE PLPGSQL;`)
	if err != nil {
		log.Fatal(err)
	}
}

/* Cuando se realiza una compra, esta función se encarga de validarla o generar rechazos u alertas
según corresponda*/
func spAutorizarCompra() {
	_, err = db.Query(
		`CREATE OR REPLACE FUNCTION autorizarcompra(_nrotarjeta char(16), _codseguridad char(4), _nrocomercio int,
			 	_monto decimal(7,2)) RETURNS BOOL AS $$
				DECLARE
					totalPendiente decimal(8,2);
					fechaVenceTarjeta int;
					fechaVence date;
					limiteCompraTarjeta decimal(8,2);
					_nrooperacion int;
				BEGIN
						
					PERFORM * FROM tarjeta WHERE nrotarjeta = _nrotarjeta;
					IF (NOT FOUND) THEN
						PERFORM agregarrechazo(
							null,
							_nrocomercio,
							CAST(current_timestamp AS timestamp),
							_monto,
							CAST('Tarjeta inexistente' AS text));
						RETURN FALSE;
					END IF;
						
					PERFORM * FROM tarjeta t WHERE t.nrotarjeta = _nrotarjeta AND estado = 'vigente';
					IF (NOT FOUND) THEN 
						PERFORM agregarrechazo(
							_nrotarjeta, 
							_nrocomercio, 
							CAST(current_timestamp AS timestamp),
							_monto, 
							CAST('Tarjeta no valida o no vigente' AS text));
						RETURN FALSE;
					END IF;
					
					PERFORM * FROM tarjeta t WHERE t.nrotarjeta = _nrotarjeta AND estado = 'suspendida';
					IF (FOUND) THEN
						PERFORM agregarrechazo(
							_nrotarjeta, 
							_nrocomercio,
							CAST(current_timestamp AS timestamp),
							_monto, 
							CAST('La tarjeta se encuentra suspendida' AS text));
						RETURN FALSE;
					END IF;
						
					PERFORM * FROM tarjeta t WHERE t.nrotarjeta = _nrotarjeta AND codseguridad = _codseguridad;
					IF (NOT FOUND) THEN
						PERFORM agregarrechazo(
							_nrotarjeta, 
							_nrocomercio,
							CAST(current_timestamp AS timestamp),
							_monto, 
							CAST('Codigo de seguridad invalido' AS text));
							
						RETURN FALSE;
					END IF;
						
					totalPendiente := (SELECT sum(monto) FROM compra WHERE nrotarjeta = _nrotarjeta AND pagado = False);
					limiteCompraTarjeta := (SELECT limitecompra FROM tarjeta WHERE nrotarjeta = _nrotarjeta);
						
					IF (totalPendiente IS NULL AND _monto > limiteCompraTarjeta OR totalPendiente IS NOT NULL 
						AND (totalPendiente + _monto) > limiteCompraTarjeta) THEN
							
						PERFORM agregarrechazo(
							_nrotarjeta,
							_nrocomercio, 
							CAST(current_timestamp AS timestamp),
							_monto, 
							CAST('Supera limite de tarjeta' AS text));
						RETURN FALSE;
					END IF;
						
					SELECT validahasta INTO fechaVenceTarjeta FROM tarjeta t WHERE t.nrotarjeta = _nrotarjeta;
					SELECT INTO fechaVence to_date(fechaVenceTarjeta || '01','YYYYMMDD');
					SELECT INTO fechaVence (fechaVence + interval '1 month')::date;
						
					IF (fechaVence < current_date) THEN
						PERFORM agregarrechazo(
							_nrotarjeta, 
							_nrocomercio, 
							CAST(current_timestamp AS timestamp),
							_monto, 
							CAST('Plazo de vigencia expirado' AS text));
						RETURN FALSE;
					END IF;
					
					PERFORM * FROM compra c WHERE c.nrooperacion IS NOT NULL ;
					IF (FOUND) THEN
						SELECT MAX(nrooperacion)+1 INTO _nrooperacion FROM compra;
					ELSE
						_nrooperacion:= 1;
					END IF;
						
					INSERT INTO compra(nrooperacion,nrotarjeta, nrocomercio, fecha, monto, pagado) 
						VALUES (_nrooperacion,_nrotarjeta, _nrocomercio, current_timestamp, _monto, FALSE);
					RETURN TRUE;	
					END;
					$$ LANGUAGE PLPGSQL;`)
	if err != nil {
		log.Fatal(err)
	}
}

/*Dado un periodo del año, un año y un cliente, esta función crea una tupla para las tablas detalle y
cabecera. Ambas representan un resumen.*/
func spGenerarResumen() {
	_, err = db.Query(
		`CREATE OR REPLACE FUNCTION generar_resumen(_periodo int,_anyo int, _nrocliente int ) RETURNS VOID AS $$
		DECLARE
				_nombre text;
				_apellido text;
				_direccion text;
				_nrotarjeta char(16); 
				_inicioPeriodo date;
				_finPeriodo date;
				_venc date;
				_terminacion char(4);
				_term int;
				_monto decimal(7,2);
				_total decimal(8,2);
				_nroresumen int;
				_nrolinea int;
				fila compra%rowtype;	
				_fechacompra date;			
				_nombreComercio text;
			
		BEGIN 	
			
			PERFORM * FROM cliente c WHERE c.nrocliente =_nrocliente;
			IF (NOT FOUND) THEN
				RAISE  'El nro de cliente % es invalido',_nrocliente 
				USING HINT = 'Por favor chequear su numero de cliente';
			END IF;
			IF(_periodo < 1 OR _periodo > 12) THEN
				RAISE 'El periodo % es invalido',_periodo
				USING HINT = 'Por favor chequear periodo';
			END IF;

			SELECT nombre INTO _nombre FROM cliente c WHERE c.nrocliente = _nrocliente; 
			SELECT apellido INTO _apellido FROM cliente c WHERE c.nrocliente = _nrocliente; 
			SELECT domicilio INTO _direccion FROM cliente c WHERE c.nrocliente = _nrocliente; 

			SELECT nrotarjeta INTO _nrotarjeta FROM tarjeta t WHERE t.nrocliente = _nrocliente;
			
			SELECT SUBSTRING(_nrotarjeta, 16, 1) INTO _terminacion;		
			SELECT CAST (_terminacion as int) INTO _term;

			SELECT fechainicio INTO _inicioPeriodo FROM cierre c WHERE (mes = _periodo AND anio = _anyo AND _term = terminacion);
			SELECT fechacierre INTO _finPeriodo FROM cierre c WHERE (mes = _periodo AND anio = _anyo AND _term = terminacion);
			SELECT fechavto INTO _venc FROM cierre c WHERE (mes = _periodo AND anio = _anyo AND _term = terminacion);

			_total := 0.00;
			FOR _monto IN SELECT monto FROM compra c WHERE (SELECT EXTRACT(MONTH FROM c.fecha ) = _periodo AND c.nrotarjeta = _nrotarjeta)
			LOOP
				_total := _total + _monto;
			END LOOP;

			PERFORM * FROM cabecera c WHERE c.nroresumen IS NOT NULL ;
			IF(FOUND) THEN
				SELECT MAX(nroresumen)+1 INTO _nroresumen FROM cabecera;
			ELSE
				_nroresumen := 1;
			END IF;
			INSERT INTO cabecera VALUES(_nroresumen, _nombre , _apellido, _direccion , _nrotarjeta , _inicioPeriodo , _finPeriodo ,_venc,_total);
			
			UPDATE compra  SET pagado = TRUE WHERE (SELECT EXTRACT(MONTH FROM compra.fecha ) = _periodo AND compra.nrotarjeta = _nrotarjeta);
	
			FOR fila IN SELECT * FROM compra c WHERE (SELECT EXTRACT(MONTH FROM c.fecha ) = _periodo AND c.nrotarjeta = _nrotarjeta )
			LOOP
				PERFORM * FROM detalle d WHERE d.nrolinea IS NOT NULL  ;
				IF(FOUND) THEN
					SELECT MAX(nrolinea)+1 INTO _nrolinea FROM detalle;
				ELSE
					_nrolinea := 1;
				END IF;


				SELECT fila.fecha ::timestamp::date INTO _fechacompra;
				SELECT nombre FROM comercio c WHERE  fila.nrocomercio = c.nrocomercio INTO _nombreComercio ;
				_monto := fila.monto;

				INSERT INTO detalle VALUES(_nroresumen,_nrolinea,_fechacompra,_nombreComercio,_monto);
			END LOOP;
		END;
		$$ LANGUAGE PLPGSQL;`)
	if err != nil {
		log.Fatal(err)
	}
}

//Cada vez que se crea un rechazo, se llama al trigger agregar_alerta()
func spAgregarAlertaRechazo() {
	_, err = db.Query(
		`CREATE OR REPLACE FUNCTION agregar_alerta() RETURNS TRIGGER AS $$
				DECLARE
				_nroalerta int;
				BEGIN
				
					PERFORM * FROM alerta a WHERE a.nroalerta IS NOT NULL;
					IF (FOUND) THEN
						SELECT MAX(nroalerta)+1 INTO _nroalerta FROM alerta;
					ELSE
						_nroalerta := 1;
					END IF;
					
					INSERT INTO alerta(nroalerta,nrotarjeta, fecha, nrorechazo, codalerta, descripcion) 
					VALUES (_nroalerta,new.nrotarjeta, new.fecha, new.nrorechazo, 0, new.motivo);
					RETURN NEW;
				END;
			$$ LANGUAGE PLPGSQL;`)
	if err != nil {
		log.Fatal(err)
	}
}

/*Se encarga de controlar los tiempos entre compras. Si es de menos de un minuto y en el mismo cod postal.
Y si de menos de cinco minutos en diferentes cod postales*/
func spSeguridadTiempoCompras() {
	_, err = db.Query(
		`CREATE OR REPLACE FUNCTION controlar_tiempo_compras() RETURNS TRIGGER AS $$
			DECLARE
				ultimaCompra record;
				diferenciaTiempo decimal;
				codPostalAnterior int;
				codPostalActual int;
				_nroalerta int;
			BEGIN
			
				PERFORM * FROM alerta a WHERE a.nroalerta IS NOT NULL;
				IF (FOUND) THEN
					SELECT MAX(nroalerta)+1 INTO _nroalerta FROM alerta;
				ELSE
					_nroalerta := 1;
				END IF;
			
				SELECT * INTO ultimaCompra FROM compra WHERE nrotarjeta = new.nrotarjeta ORDER BY nrotarjeta DESC LIMIT 1;
				IF (NOT FOUND) THEN
					RETURN NEW;
				END IF;
				
				SELECT INTO diferenciaTiempo EXTRACT(EPOCH FROM new.fecha - ultimaCompra.fecha)/60;
				
				SELECT codigopostal INTO codPostalAnterior FROM comercio WHERE nrocomercio = ultimaCompra.nrocomercio;
				SELECT codigopostal INTO codPostalActual FROM comercio WHERE nrocomercio = new.nrocomercio;
				
				IF(diferenciaTiempo < 1 AND ultimaCompra.nrocomercio != new.nrocomercio AND codPostalAnterior = codPostalActual) THEN
					INSERT INTO alerta(nroalerta, nrotarjeta, fecha, nrorechazo, codalerta, descripcion)
						VALUES (_nroalerta,new.nrotarjeta, new.fecha, NULL, 1,  'Compra en menos de 1 minuto en el mismo comercio');
					RETURN NEW;
				END IF;
				
				IF(diferenciaTiempo < 5 AND ultimaCompra.nrocomercio != new.nrocomercio AND codPostalAnterior != codPostalActual) THEN
					INSERT INTO alerta(nroalerta, nrotarjeta, fecha, nrorechazo, codalerta, descripcion) 
						VALUES (_nroalerta,new.nrotarjeta, new.fecha, NULL, 5, 'Compra en menos de 5 minutos en distintos comercios');
					RETURN NEW;
				END IF;
				
				RETURN NEW;
			END;
			$$ LANGUAGE PLPGSQL;`)
	if err != nil {
		log.Fatal(err)
	}
}

func spAgregarProbarConsumos() {
	_, err = db.Query(`
		CREATE OR REPLACE FUNCTION probar_consumos() RETURNS VOID AS $$
			DECLARE
				v consumo%rowtype;

			BEGIN
				FOR v IN SELECT * FROM consumo LOOP
					PERFORM autorizarcompra(v.nrotarjeta, v.codseguridad, v.nrocomercio, v.monto);
				END LOOP;

			END;
			$$ LANGUAGE PLPGSQL;`)

	if err != nil {
		log.Fatal(err)
	}
}
