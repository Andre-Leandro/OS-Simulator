package main

import (
	"fmt"
	"math"
	"path/filepath"

	"github.com/sqweek/dialog"
)

const quantum = 2
const biggestPartition = 250

func swapOut(i int) {
	switch i {
	case 1:
		fmt.Println("Se liberó la partición GRANDE (250 kB)")
	case 2:
		fmt.Println("Se liberó la partición MEDIANA (120 kB)")
	case 3:
		fmt.Println("Se liberó la partición PEQUEÑA (60 kB)")
	default:
		fmt.Println("Partición no reconocida")
	}
}

func (os *OS) addReady(l *[]Process) {
	copy := *l

	for index := range copy {
		if len(os.queue) == 4 {
			break
		}

		if copy[index].arrivalTime <= os.time {
			os.queue = append(os.queue, copy[index])
			if os.queue[len(os.queue)-1].turnaroundTime == -1 {
				os.queue[len(os.queue)-1].turnaroundTime = os.time
			}

			if len(*l) > 0 {
				if len(*l) == 1 {
					*l = []Process{}
				} else {
					*l = (*l)[1:]
				}
			}
		}

	}
}

func bestFitLazy(m *Memory, p *Process, os *OS) {
	var internalFragmentation int
	var idPartition int
	idPartition = -1
	internalFragmentation = math.MaxInt

	if p.loaded == false {
		for index := range m.partitions {
			partition := (*m).partitions[index]
			if partition.state && partition.size >= p.size {
				empty := partition.size - p.size
				if empty < internalFragmentation {
					idPartition = index
				}
			}
		}

		if idPartition != -1 {
			(*m).partitions[idPartition].state = false // Ocupado
			(*m).partitions[idPartition].internalFragmentation = (*m).partitions[idPartition].size - (*p).size
			(*p).loaded = true
			if p.turnaroundTime == -1 {
				(*p).turnaroundTime = os.time
			}
			(*m).partitions[idPartition].process = *p
		}

	}

}

func bestFit(m *Memory, p *Process, os *OS) {
	var internalFragmentation int
	var idPartition int
	idPartition = -1
	internalFragmentation = math.MaxInt

	for index := range m.partitions {
		partition := (*m).partitions[index]
		if partition.state && partition.size >= p.size {
			empty := partition.size - p.size
			if empty < internalFragmentation {
				idPartition = index
			}
		}
	}

	if idPartition == -1 {
		idPartition = bestFitSwap(*m, *p, *os)
		swapOut(idPartition + 1)
		currentProcess := os.memory.partitions[idPartition].process
		for i, queueProcess := range os.queue {
			if queueProcess.pid == currentProcess.pid {
				currentProcess = os.queue[i]
				currentProcess.loaded = false
				os.queue[i] = currentProcess
				break
			}
		}
	}
	(*m).partitions[idPartition].state = false // Ocupado
	(*m).partitions[idPartition].internalFragmentation = (*m).partitions[idPartition].size - (*p).size
	(*p).loaded = true
	if p.turnaroundTime == -1 {
		(*p).turnaroundTime = os.time
	}
	(*m).partitions[idPartition].process = *p
}

func bestFitSwap(m Memory, p Process, os OS) int {
	var internalFragmentation int
	var idPartition int
	internalFragmentation = math.MaxInt

	for index := range m.partitions {
		partition := m.partitions[index]
		if partition.size >= p.size {
			empty := partition.size - p.size
			if empty < internalFragmentation {
				idPartition = index
			}
		}
	}
	return idPartition
}

func (p *Process) timeOut(quantum int, queue *[]Process, os *OS, cola *[]Process) {
	if p.time > quantum {
		os.time = os.time + quantum
		os.addReady(cola)
		p.time = p.time - quantum
		*queue = append(*queue, *p)
	} else {
		os.time = os.time + p.time
		os.addReady(cola)
		p.time = 0
		p.turnaroundTime = os.time - p.turnaroundTime
		fmt.Println("\n", "Termino el proceso:", p.pid, "en el instante", os.time, "\n")
		os.completedProcesses = append(os.completedProcesses, *p)

		for index := range os.memory.partitions {
			partition := os.memory.partitions[index]
			if partition.process.pid == p.pid {
				os.memory.partitions[index].process = Process{}
				os.memory.partitions[index].state = true
				os.memory.partitions[index].internalFragmentation = 0
				break
			}
		}
	}
}

func (p Process) isEmpty() bool {
	return p.pid == 0 && p.time == 0 && p.arrivalTime == 0 && p.size == 0 // and loaded
}

func (os *OS) initialize(m Memory) {
	os.memory = m
}

func quicksort(processes []Process) []Process {
	if len(processes) <= 1 {
		return processes
	}

	pivotIndex := len(processes) / 2
	pivot := processes[pivotIndex]
	var less []Process
	var greater []Process
	var equal []Process

	for _, p := range processes {
		if p.arrivalTime < pivot.arrivalTime {
			less = append(less, p)
		} else if p.arrivalTime > pivot.arrivalTime {
			greater = append(greater, p)
		} else {
			if p.pid < pivot.pid {
				less = append(less, p)
			} else if p.pid > pivot.pid {
				greater = append(greater, p)
			} else {
				equal = append(equal, p)
			}
		}
	}
	return append(append(quicksort(less), equal...), quicksort(greater)...)
}

func main() {
	printLambda()
	var pidMap = make(map[int]bool)
	var opcion int
	var err error

	var restart bool
	var processes []Process

	for {
		pidMap = make(map[int]bool)
		fmt.Println("Seleccione una opción:")
		fmt.Println("1. Ingresar procesos manualmente.")
		fmt.Println("2. Cargar procesos desde un archivo.")

		restart = false

		for {
			_, err = fmt.Scanln(&opcion)

			if err != nil {
				fmt.Println("Error al leer la opción. Por favor, inténtelo nuevamente.")
				// Limpiar el búfer del teclado para evitar problemas con futuras lecturas
				fmt.Scanln()
				continue
			}

			if opcion != 1 && opcion != 2 {
				fmt.Println("Opción no válida. Por favor, ingrese 1 o 2.")
				continue
			}

			break
		}

		switch opcion {
		case 1:
			processes = ingresarProcesosManualmente(pidMap)
		case 2:
			fmt.Println("Seleccione un archivo con procesos para iniciar la simulación.")
			filePath, err := dialog.File().Load()
			if err != nil {
				fmt.Println("Error al abrir el explorador de archivos:", err)
				return
			}

			if filepath.Ext(filePath) != ".txt" {
				fmt.Println("El archivo seleccionado no tiene la extensión .txt. Por favor, seleccione un archivo válido.")
				restart = true
				break
			}

			processes, err = ReadProcessesFromFile(filePath, pidMap)
			if err != nil {
				fmt.Println("Error:", err)
				restart = true
				break
			}
		default:
			fmt.Println("Opción no válida. Saliendo del simulador.")
			return
		}

		if !restart {
			break
		}
	}

	clearScreen()

	var cola []Process

	var del []Process
	var linux OS
	memoria := Memory{
		partitions: [3]MemoryPartition{
			{id: 1, size: 250, state: true},
			{id: 2, size: 120, state: true},
			{id: 3, size: 60, state: true},
		},
	}
	(&linux).initialize(memoria)
	cola = append(cola, processes...)
	cola = quicksort(cola)
	cola, del = filterProcessesBySize(cola, biggestPartition)
	var input string
	fmt.Println("Inicio del Sistema Operativo")
	fmt.Print("\n")
	fmt.Print("• Estos son los procesos ELIMINADOS por exceder el tamaño de memoria: ")
	showQueue(del)
	fmt.Print("Presione ENTER para continuar ")

	for {
		fmt.Scanln(&input)
		if input == "" {
			if !linux.processor.process.isEmpty() {
				linux.processor.process.timeOut(quantum, &linux.queue, &linux, &cola)
				if len(linux.queue) == 0 && len(cola) > 0 && linux.processor.process.time <= 0 { //ver por que no funciona con el igual
					linux.time = cola[0].arrivalTime
					linux.addReady(&cola)
				}
				if len(linux.queue) == 0 && len(cola) == 0 && linux.processor.process.time <= 0 { //ver por que no funciona con el igual
					fmt.Println(" Se termino de procesar todo - Fin de la Simulación")
					break
				}
				if linux.queue[0].loaded == false {
					bestFit(&linux.memory, &linux.queue[0], &linux)
				}
				linux.processor.process = linux.queue[0]
				linux.queue = append(linux.queue[1:])

				for i := range linux.queue {
					bestFitLazy(&linux.memory, &linux.queue[i], &linux)
				}

				linux.addReady(&cola)
			} else {
				//contemplar que es la primera vez y se puede empezar en algo distinto que 0
				linux.time = cola[0].arrivalTime
				linux.addReady(&cola)

				bestFit(&linux.memory, &linux.queue[0], &linux)

				linux.processor.process = linux.queue[0]
				linux.queue = append(linux.queue[1:])
				linux.addReady(&cola)
				for i := range linux.queue {

					bestFitLazy(&linux.memory, &linux.queue[i], &linux)
				}
			}

			if len(linux.queue) == 0 && len(cola) == 0 && linux.processor.process.time <= 0 { //ver por que no funciona con el igual
				fmt.Println("\n", "Se termino de procesar todo - Fin de la Simulacion", "\n")
				break
			}

			if !linux.processor.process.isEmpty() {
				fmt.Println("\n", "---------------------------------------- TIEMPO: ", linux.time, " ---------------------------------------", "\n")
				mostrarProcesador(linux.processor.process)
				fmt.Print("\n")
				fmt.Print("                                             MEMORIA")
				// mostrarDatos2("MEMORIA")
				fmt.Print("\n")
				mostrarDatos(linux.memory, linux.processor.process)
				fmt.Print("\n")
				// mostrarDatos2("COLA DE LISTOS")
				fmt.Print("                                        COLA DE LISTOS")
				fmt.Print("\n")
				mostrarColaListos(linux.queue)
				fmt.Print("\n")
				fmt.Print("                                        COLA DE NUEVOS")
				fmt.Print("\n")
				mostrarColaNuevos(cola)
				fmt.Print("\n")

				//fmt.Print("* Esta es la cola de listos: ")
				//	mostrarColas(linux.queue)
				/* 				fmt.Print("• Esta es la cola de  procesos NUEVOS: ")
				   				mostrarColas(cola) */
				fmt.Print(" • Esta es la cola de procesos FINALIZADOS: ")
				showQueue(linux.completedProcesses)

				fmt.Println(" ---------------------------------------------------------------------------------------------")
			}

		} else {
			break
		}
	}
	fmt.Println("\n", "                     CUADRO ESTADÍSTICO")
	arrancar(linux.completedProcesses, processes)
}
