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
	time      int
	processor Processor
	memory    Memory
	ReadyQueue
	completedProcesses []Process
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

const quantum = 2
const biggestPartition = 250

func swapOut(i int) {
	fmt.Println("Se libero la particion numero ", i)
}

func (os *OS) addReady(l *[]Process) {
	copy := *l

	for index := range copy {
		if len(os.queue) == 4 {
			break
		}

		if copy[index].arrivalTime <= os.time {
			os.queue = append(os.queue, copy[index])

		
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

func ReadProcessesFromFile(filename string) ([]Process, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	processes := []Process{}
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
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

		turnaroundTime := -1
		process := Process{pid, size, arrivalTime, turnaroundTime, time, false}
		processes = append(processes, process)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return processes, nil
}

func quicksort2(processes []Process) []Process {
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
	return append(append(quicksort2(less), equal...), quicksort2(greater)...)
}

func filterProcessesBySize(processes []Process, sizeThreshold int) ([]Process, []Process) {
	var filteredProcesses []Process
	var deleted []Process

	for _, p := range processes {
		if p.size <= sizeThreshold {
			filteredProcesses = append(filteredProcesses, p)
		} else {
			deleted = append(deleted, p)
		}
	}

	return filteredProcesses, deleted
}

func mostrarColas(processes []Process) {
	if len(processes) == 0 {
		fmt.Println("-")
		fmt.Println("\n")
		return
	}

	for i, p := range processes {
		if i == 0 {
			fmt.Print(p.pid)
		} else {
			fmt.Print(", ", p.pid)
		}
	}
	fmt.Println("\n")
}

func main() {
	processes, err := ReadProcessesFromFile("ejemplo.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
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
	cola = quicksort2(cola)
	cola, del = filterProcessesBySize(cola, biggestPartition)
	var input string
	fmt.Println("Inicio del Sistema Operativo - Presione ENTER para continuar")
	fmt.Print("\n")
	fmt.Print("•Estos son los procesos ELIMINADOS por exceder el tamaño de memoria: ") 
	mostrarColas(del)

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
					fmt.Println("Se termino de procesar todo - Fin de la Simulación")
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


				for i := range linux.queue {
					
					bestFitLazy(&linux.memory, &linux.queue[i], &linux)
				}
			}
	

			if len(linux.queue) == 0 && len(cola) == 0 && linux.processor.process.time <= 0 { //ver por que no funciona con el igual
				fmt.Println("\n", "Se termino de procesar todo - Fin de la Simulacion", "\n")
				break
			}

			if !linux.processor.process.isEmpty() {
				fmt.Println("\n", "---------------------------------- TIEMPO: ", linux.time, " ----------------------------------", "\n")
				mostrarProcesador(linux.processor.process)
				fmt.Print("\n")
				mostrarDatos2("MEMORIA")
				fmt.Print("\n")
				mostrarDatos(linux.memory, linux.processor.process)
				fmt.Print("\n")
				mostrarDatos2("COLA DE LISTOS")
				fmt.Print("\n")
				mostrarColaListos(linux.queue)
				fmt.Print("\n")

				//fmt.Print("* Esta es la cola de listos: ")
				//	mostrarColas(linux.queue)
				fmt.Print("• Esta es la cola de  procesos NUEVOS: ")
				mostrarColas(cola)
				fmt.Print("• Esta es la cola de procesos FINALIZADOS: ")
				mostrarColas(linux.completedProcesses)

				fmt.Println("\n", "----------------------------------------------------------------------------------")
			}

		} else {
			break
		}
	}
	fmt.Println("\n", "CUADRO ESTADÍSTICO")
	arrancar(linux.completedProcesses, processes)
}
