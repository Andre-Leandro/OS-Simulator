package main

import (
	"fmt"
)

func ingresarProcesosManualmente(pidMap map[int]bool) []Process {
	var procesos []Process

	for {
		var cantidadProcesos int
		for {
			fmt.Print("• Ingrese la cantidad de procesos: ")
			_, err := fmt.Scanln(&cantidadProcesos)

			if err != nil || cantidadProcesos <= 0 {
				fmt.Println("Error al leer la cantidad de procesos. Por favor, ingrese un número entero positivo.")
				// Limpiar el búfer del teclado para evitar problemas con futuras lecturas
				continue
			} else {
				break
			}

		}

		for i := 0; i < cantidadProcesos; i++ {
			var pid, size, arrivalTime, time int

			for {
				fmt.Printf("• Ingrese el PID para el proceso %d: ", i+1)
				_, err := fmt.Scanln(&pid)
				if err != nil || pid < 0 {
					fmt.Println("Error al leer el PID. Por favor, ingrese un número entero positivo.")
					// Limpiar el búfer del teclado para evitar problemas con futuras lecturas
					continue
				}
				if !isPIDInUse(pid, pidMap) {
					pidMap[pid] = true
					break
				} else {
					fmt.Println("Error: El PID ya está en uso. Por favor, ingrese un PID único.")
					continue
				}
			}
			for {
				fmt.Printf("• Ingrese el TAMAÑO (kB) para el proceso %d: ", i+1)
				_, err := fmt.Scanln(&size)
				if size > 250 {
					fmt.Println("El tamaño de los procesos no puede se superior a los 250 kB.")
					continue
				}
				if err != nil || size <= 0 {
					fmt.Println("Error al leer el tamaño. Por favor, ingrese un número entero positivo.")
					continue
				} else {
					break
				}
			}
			for {
				fmt.Printf("• Ingrese el TIEMPO DE ARRIBO para el proceso %d: ", i+1)
				_, err := fmt.Scanln(&arrivalTime)
				if err != nil || arrivalTime < 0 {
					fmt.Println("Error al leer el tiempo de arribo. Por favor, ingrese un número entero mayor o igual a 0.")
					continue
				} else {
					break
				}
			}
			for {
				fmt.Printf("• Ingrese el TIEMPO DE IRRUPCIÓN para el proceso %d: ", i+1)
				_, err := fmt.Scanln(&time)

				if err != nil || time <= 0 {
					fmt.Println("Error al leer el tiempo de irrupción. Por favor, ingrese un número entero positivo.")
					continue
				} else {
					break
				}
			}

			proceso := Process{pid, size, arrivalTime, -1, time, false}
			procesos = append(procesos, proceso)
		}

		break
	}

	return procesos
}
