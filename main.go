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
}

type Process struct {
	pid         int
	size        int
	arrivalTime int
	time        int
	loaded      bool
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

func swapOut(i int) {
	fmt.Println("Se libero la particion maleta ", i)
}

func (os *OS) addReady(l *[]Process) {
	copy := *l

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

func bestFit(m *Memory, p *Process, os *OS) {
	var internalFragmentation int
	var idPartition int
	idPartition = -1
	internalFragmentation = math.MaxInt

	for index := range m.partitions {
		partition := (*m).partitions[index] // Obtener una referencia a la partición real en la estructura Memory
		if partition.state && partition.size >= p.size {
			empty := partition.size - p.size
			if empty < internalFragmentation {
				idPartition = index
			}
		}
	}
	//fmt.Println("primaria", os.queue)

	if idPartition == -1 {
		idPartition = bestFitSwap(*m, *p)
		swapOut(idPartition + 1)

		//(*os).memory.partitions[idPartition].process.loaded = false\

		currentProcess := os.memory.partitions[idPartition].process
		currentProcess.loaded = false
		currentProcess.time = currentProcess.time - 5

    // Encontrar el índice del proceso en os.queue con el mismo ID
			for i, queueProcess := range os.queue {
				if queueProcess.pid == currentProcess.pid {
					// Reemplazar el proceso en os.queue con el nuevo proceso de la partición
					os.queue[i] = currentProcess
					break
				}
			}

		//fmt.Println((*os).memory.partitions[idPartition].process)
		//fmt.Println("gargante", os.queue)
	}
	// selectedPartition := &m.partitions[idPartition]
	// selectedPartition.state = false // Ocupado
	// selectedPartition.internalFragmentation = selectedPartition.size - p.size
	// selectedPartition.process = *p
	// selectedPartition.process.loaded = newLoaded
	(*m).partitions[idPartition].state = false // Ocupado
	(*m).partitions[idPartition].internalFragmentation = (*m).partitions[idPartition].size - (*p).size
	(*p).loaded = true
	(*m).partitions[idPartition].process = *p
}

func bestFitSwap(m Memory, p Process) int {
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
		fmt.Println("Termino el proceso: ", p.pid, "en el instante", os.time)

		for index := range os.memory.partitions {
			partition := os.memory.partitions[index] 
			if partition.process.pid == p.pid {
				os.memory.partitions[index].process = Process{}
				os.memory.partitions[index].state = true
				os.memory.partitions[index].internalFragmentation = 0
			}
		}
	}
}


func freeMemory () {

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

		process := Process{pid, size, arrivalTime, time, false}
		processes = append(processes, process)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return processes, nil
}

func printPartitionInfo(partition MemoryPartition) {
	state := "Occupied"
	if partition.state {
		state = "Free"
	}
	processName := "N/A"
	if partition.process.pid != 0 {
		processName = fmt.Sprintf("Process-%d", partition.process.pid)
	}
	fmt.Printf("| %-10d | %-10d | %-10s | %-15d | %-10s |\n",
		partition.id, partition.size, state, partition.internalFragmentation, processName)
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

func main() {
	processes, err := ReadProcessesFromFile("ejemplo.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var cola []Process
	var linux OS
	memoria := Memory{
		partitions: [3]MemoryPartition{
			{id: 1, size: 100, state: true},
			{id: 2, size: 75, state: true},
			{id: 3, size: 35, state: true},
		},
	}
	(&linux).initialize(memoria)
	cola = append(cola, processes...)
	cola = quicksort2(cola)
	//fmt.Println(cola)

	var input string
	fmt.Print("Inicio del Sistema Operativo")

	for {
		fmt.Scanln(&input)
		if input == "" {
			if !linux.processor.process.isEmpty() {
				linux.processor.process.timeOut(5, &linux.queue, &linux, &cola)

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
					bestFit(&linux.memory, &linux.queue[0], &linux)
				}
				linux.processor.process = linux.queue[0]
				linux.queue = append(linux.queue[1:])
				linux.addReady(&cola)

			} else {
				//contemplar que es la primera vez y se puede empezar en algo distinto que 0
				linux.time = cola[0].arrivalTime
				linux.addReady(&cola)
				bestFit(&linux.memory, &linux.queue[0], &linux)
				//fmt.Print("leyenda")
				//fmt.Print(*&linux.memory.partitions[2])
				linux.processor.process = linux.queue[0]
				linux.queue = append(linux.queue[1:])
			}

			if len(linux.queue) == 0 && len(cola) == 0 && linux.processor.process.time <= 0 { //ver por que no funciona con el igual
				fmt.Println("Se termino de procesar todo - Fin de la Simulacion")
				break
			}
			fmt.Println("")
			fmt.Println("------------------------------ TIME: ", linux.time, " ------------------------------")
			fmt.Println("")
			fmt.Println("PROCESADOR: Proceso", linux.processor.process.pid)
			fmt.Println("* Esta es la cola de listos: ", linux.queue)
			fmt.Println("* Esta es la cola de input/nuevos: ", cola)
			fmt.Println("")
			fmt.Println("-----------------------------------------------------------------------")
			fmt.Printf("| %-10s | %-10s | %-10s | %-15s | %-10s |\n", "ID", "Size", "State", "Internal Frag.", "Process")
			fmt.Println("-----------------------------------------------------------------------")
			printPartitionInfo(linux.memory.partitions[0])
			printPartitionInfo(linux.memory.partitions[1])
			printPartitionInfo(linux.memory.partitions[2])
			fmt.Println("-----------------------------------------------------------------------")

		} else {
			break
		}
	}

}