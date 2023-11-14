package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type OS struct {
	time               int
	processor          Processor
	memory             Memory
	completedProcesses []Process
	ReadyQueue
}

type Process struct {
	pid            int
	size           int
	arrivalTime    int
	turnaroundTime int
	time           int
	loaded         bool
}

type Processor struct {
	process Process
}

type MemoryPartition struct {
	id                    int
	size                  int
	state                 bool
	process               Process
	internalFragmentation int
}

type Memory struct {
	partitions [3]MemoryPartition
}

type NewQueue struct {
	queue []Process
}

type ReadyQueue struct {
	queue []Process
}

func (q ReadyQueue) estado() {
	fmt.Println(q.queue)
}

func swapOut(i int) {
	fmt.Println("Se libero la particion", i)
}

func (os *OS) addReady(l *[]Process) {
	copy := *l

	fmt.Println(os.queue)
	fmt.Println(len(os.queue))

	for index := range copy {
		if len(os.queue) == 4 {
			break
		}

		if copy[index].arrivalTime <= os.time { // menor o igual tambien
			os.queue = append(os.queue, copy[index])
			if len(*l) > 0 {
				if len(*l) == 1 {
					*l = []Process{}
				} else {
					*l = (*l)[1:]
				}
			}
		}
		bestFitLazy(os.memory, copy[index])
	}
}

func bestFitLazy(m Memory, p Process) {
	var internalFragmentation int
	var idPartition int
	internalFragmentation = math.MaxInt

	for index := range m.partitions {
		partition := m.partitions[index]
		if partition.state && partition.size >= p.size {
			empty := partition.size - p.size
			if empty < internalFragmentation {
				idPartition = index
			}
		}
	}
	if idPartition != 0 {
		selectedPartition := &m.partitions[idPartition]
		selectedPartition.state = false
		selectedPartition.internalFragmentation = selectedPartition.size - p.size
		selectedPartition.process = p
	}
}

func bestFit(os *OS, m *Memory, p *Process) {
	var internalFragmentation int
	var idPartition int
	idPartition = -1
	internalFragmentation = math.MaxInt

	for index := range m.partitions {
		partition := m.partitions[index]
		if partition.state && partition.size >= p.size {
			empty := partition.size - p.size
			if empty < internalFragmentation {
				idPartition = index
			}
		}
	}
	if idPartition == -1 {
		idPartition = bestFitSwap(*os, *m, *p)
		swapOut(idPartition)
		m.partitions[idPartition].process.loaded = false
	}
	selectedPartition := &m.partitions[idPartition]
	selectedPartition.state = false
	selectedPartition.internalFragmentation = selectedPartition.size - p.size
	p.loaded = true
	p.turnaroundTime = os.time
	selectedPartition.process = *p
}

func bestFitSwap(os OS, m Memory, p Process) int {
	var internalFragmentation int
	var idPartition int
	internalFragmentation = math.MaxInt
	var currentProcess Process

	for index := range m.partitions {
		partition := m.partitions[index]
		if partition.size >= p.size {
			empty := partition.size - p.size
			if empty < internalFragmentation {
				idPartition = index
				currentProcess = partition.process
			}
		}
	}
	currentProcess.turnaroundTime = os.time - currentProcess.turnaroundTime
	fmt.Println("SALE EL PROCESO: ", currentProcess.pid, currentProcess.turnaroundTime)
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
		fmt.Println("Termino el proceso: ", p.pid)
		os.completedProcesses = append(os.completedProcesses, *p)
	}
}

func sort(input ReadyQueue) {
	for i := range input.queue {
		fmt.Print("este es el numero ", i)
		for j := range input.queue[i:] {
			fmt.Print(j+i, " - ")

		}
		fmt.Println("")
	}
}

func (p *Process) Escribir() {
	p.pid = 5
	fmt.Println(p.time)
}

var a NewQueue

func (p Process) isEmpty() bool {
	return p.pid == 0 && p.time == 0 && p.arrivalTime == 0 && p.size == 0 // and loaded
}

func (m Memory) idLoaded() {
	fmt.Printf("caragada")
}

func (os *OS) initialize(m Memory) {
	os.memory = m
}

func ReadProcessesFromFile(filename string) ([]Process, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	processes := []Process{}

	scanner := bufio.NewScanner(file)
	// Agregar una verificación para ignorar la primera línea
	if scanner.Scan() {
		// Ignorar la primera línea (encabezado o comentario)
	}

	for scanner.Scan() {
		line := scanner.Text()
		values := strings.Fields(line)
		if len(values) != 4 {
			return nil, fmt.Errorf("Formato de entrada invalido: %s", line)
		}
		pid, err := strconv.Atoi(values[0])
		if err != nil {
			return nil, err
		}
		size, err := strconv.Atoi(values[1])
		if err != nil {
			return nil, err
		}
		arrivalTime, err := strconv.Atoi(values[2])
		if err != nil {
			return nil, err
		}
		time, err := strconv.Atoi(values[3])
		if err != nil {
			return nil, err
		}

		turnaroundTime := 0

		process := Process{pid, size, arrivalTime, turnaroundTime, time, false}
		processes = append(processes, process)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return processes, nil
}

func printStatistics(completedProcesses []Process, allProcesses []Process) {
	fmt.Println("+" + strings.Repeat("-", 46) + "+")
	fmt.Printf("| %-4s | %-18s | %-16s |\n", "PID", "Tiempo de Retorno", "Tiempo de Espera")
	fmt.Println("+" + strings.Repeat("-", 46) + "+")
	var totalTurnaroundTime, totalWaitTime int

	for _, p := range completedProcesses {
		// Buscar el proceso correspondiente en allProcesses
		var originalProcess Process
		for _, originalP := range allProcesses {
			if originalP.pid == p.pid {
				originalProcess = originalP
				break
			}
		}

		waitTime := p.turnaroundTime - originalProcess.time
		fmt.Printf("| %-4d | %-18d | %-16d |\n", p.pid, p.turnaroundTime, waitTime)

		// Sumar tiempos para el cálculo promedio
		totalTurnaroundTime += p.turnaroundTime
		totalWaitTime += waitTime
	}
	fmt.Println("+" + strings.Repeat("-", 46) + "+")

	// Calcular promedios
	averageTurnaroundTime := float64(totalTurnaroundTime) / float64(len(completedProcesses))
	averageWaitTime := float64(totalWaitTime) / float64(len(completedProcesses))

	// Imprimir tiempos promedio en un mini cuadro
	fmt.Println("+" + strings.Repeat("-", 46) + "+")
	fmt.Printf("| %-18s | %-23.2f |\n", "Promedio Retorno", averageTurnaroundTime)
	fmt.Printf("| %-18s | %-23.2f |\n", "Promedio Espera", averageWaitTime)
	fmt.Println("+" + strings.Repeat("-", 46) + "+")
}

func main() {
	processes, err := ReadProcessesFromFile("ejemplo.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var cola []Process
	linux := new(OS)
	memoria := Memory{
		partitions: [3]MemoryPartition{
			{id: 1, size: 100, state: true},
			{id: 2, size: 75, state: true},
			{id: 3, size: 35, state: true},
		},
	}
	linux.initialize(memoria)

	cola = append(cola, processes...)
	var input string
	fmt.Print("Inicio del Sistema Operativo")

	for {
		fmt.Scanln(&input)
		if input == "" {
			if !linux.processor.process.isEmpty() {
				linux.processor.process.timeOut(5, &linux.queue, linux, &cola)

				//to not go out of bounds
				if len(linux.queue) == 0 && len(cola) > 0 && linux.processor.process.time <= 0 { //ver por que no funciona con el igual
					linux.time = cola[0].arrivalTime
					linux.addReady(&cola)
				}
				if len(linux.queue) == 0 && len(cola) == 0 && linux.processor.process.time <= 0 { //ver por que no funciona con el igual
					fmt.Println("Se termino de procesar todo - Fin de la Simulacion")
					break
				}
				if linux.queue[0].loaded == false {
					bestFit(linux, &linux.memory, &linux.queue[0])
				}
				linux.processor.process = linux.queue[0]
				linux.queue = append(linux.queue[1:])
				linux.addReady(&cola)

			} else {
				//contemplar que es la primera vez y se puede empezar en algo distinto que 0
				linux.time = cola[0].arrivalTime
				bestFit(linux, &linux.memory, &cola[0])
				linux.processor.process = cola[0]
				cola = append(cola[1:])
			}

			if len(linux.queue) == 0 && len(cola) == 0 && linux.processor.process.time <= 0 { //ver por que no funciona con el igual
				fmt.Println("Se termino de procesar todo - Fin de la Simulacion")
				break
			}
			fmt.Println("TIME: ", linux.time, "----------------------------------------------------")
			fmt.Println("El proceso que se encuentra en el procesador es: pid", linux.processor.process.pid)
			fmt.Println("Tiempo de espera", linux.processor.process.turnaroundTime)
			fmt.Println("Esta es la cola de listos", linux.queue)
			fmt.Println("Esta es la cola de input", cola)
			fmt.Println("Esta es la cola de finalizados:", linux.completedProcesses)
		} else {
			break
		}
	}
	fmt.Println(processes)
	fmt.Println("")
	fmt.Println("CUADRO ESTADISTICO")
	printStatistics(linux.completedProcesses, processes)

}
