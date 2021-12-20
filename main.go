package main

import (
	"fmt"
	f "ruiz-sanchez-santillan/funciones"
)

func display_menu() {
	running := true
	var opcion int

	for running {
		mostrarOpciones()
		if ret, _ := fmt.Scanln(&opcion); ret == 1 {
			running = manejarOpciones(opcion)
		}
	}
}

func mostrarOpciones() {
	fmt.Print("\n######## TP BDD ########\n\n")
	fmt.Print("1- Crear base de datos. \n")
	fmt.Print("2- Crear tablas. \n")
	fmt.Print("3- Crear PK's y FK's. \n")
	fmt.Print("4- Eliminar PK's y FK's.\n")
	fmt.Print("5- Cargar tablas.\n")
	fmt.Print("6- Crear stored procedures y Triggers\n")
	fmt.Print("7- Probar consumos.\n")
	fmt.Print("8- Generar Resumen.\n")
	fmt.Print("9- Cargar datos noSQL en BoltDB.\n")
	fmt.Print("0- Salir. \n\n")

	fmt.Print("Elija una opcion: ")
}

func manejarOpciones(opcion int) bool {
	switch {
	case opcion == 0:
		fmt.Println("\033[2J") // esto limpia la terminal
		fmt.Println("Termino!")
		f.BorrarDB()
		return false
	case opcion == 1:
		fmt.Print("\033[H\033[2J")
		f.CrearDB()
		fmt.Println("Base de datos creada correctamente.")
	case opcion == 2:
		fmt.Print("\033[H\033[2J")
		f.CrearTablas()
		fmt.Println("Tablas creadas correctamente.")
	case opcion == 3:
		fmt.Print("\033[H\033[2J")
		f.CrearPKyFK()
		fmt.Println("PK's y FK's creadas correctamente.")
	case opcion == 4:
		fmt.Print("\033[H\033[2J")
		f.EliminarPKyFK()
		fmt.Println("PK's y FK's eliminadas correctamente.")
	case opcion == 5:
		fmt.Print("\033[H\033[2J")
		f.CargarTablas()
		fmt.Println("Datos cargados correctamente.")
	case opcion == 6:
		fmt.Print("\033[H\033[2J")
		f.CargarSPs_y_triggers()
		fmt.Println("Funciones creadas correctamente.")
	case opcion == 7:
		fmt.Print("\033[H\033[2J")
		f.Realizar_compras()
		fmt.Println("Compras realizadas.")
	case opcion == 8:
		fmt.Print("\033[H\033[2J")
		f.Realizar_Resumen()
		fmt.Println("Resumenes generados correctamente.")
	case opcion == 9:
		fmt.Print("\033[H\033[2J")
		f.CargarDatosBolt()
		fmt.Print("\nDatos noSQl cargados correctamente\n.")
	default:
		fmt.Println("Opcion incorrecta. Ingrese un numero valido")
	}
	return true
}

func main() {
	display_menu()
}
