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
	process               *Process
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

func swapOut(p *Process) {
	fmt.Println("Se liberó la partición", p.pid)
	p.loaded = false
}

func (os *OS) addReady(l *[]Process) {
	copy := *l

	fmt.Println(os.queue)
	fmt.Println(len(os.queue))

	for index := range copy {
		if len(os.queue) == 4 {
			break
		}

		if copy[index].arrivalTime <= os.time { // menor o igual también
			os.queue = append(os.queue, copy[index])
			if len(*l) > 0 {
				if len(*l) == 1 {
					*l = []Process{}
				} else {
					*l = (*l)[1:]
				}
			}
			bestFitLazy(&os.memory, &copy[index])
		}
	}
}

func bestFitLazy(m *Memory, p *Process) {
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
		if selectedPartition.process != nil {
			selectedPartition.state = false
			selectedPartition.internalFragmentation = selectedPartition.size - p.size
			selectedPartition.process = p
		}
	}
}

func bestFit(m *Memory, p *Process) {
	var internalFragmentation int
	var idPartition int
	idPartition = -1
	internalFragmentation = math.MaxInt

	for index := range m.partitions {
		partition := &m.partitions[index] // Obtener una referencia a la partición real en la estructura Memory
		if partition.state && partition.size >= p.size {
			empty := partition.size - p.size
			if empty < internalFragmentation {
				idPartition = index
			}
		}
	}
	if idPartition == -1 {
		idPartition = bestFitSwap(m, p)
		swapOut(&m.partitions[idPartition].process)
		fmt.Println("Cacatua", idPartition)
		if m.partitions[idPartition].process != nil {
			m.partitions[idPartition].process.loaded = false
			fmt.Println(m.partitions[idPartition])
		} else {
			fmt.Println("selectedPartition.process es nil")
		}
	}
	selectedPartition := &m.partitions[idPartition]
	if selectedPartition != nil {
		selectedPartition.state = false // Ocupado
		selectedPartition.internalFragmentation = selectedPartition.size - p.size
		p.loaded = true
		selectedPartition.process = p
		fmt.Print("laguna")
		fmt.Print(*selectedPartition)
	} else {
		fmt.Println("selectedPartition es nil")
	}
}

func bestFitSwap(m *Memory, p *Process) int {
	var internalFragmentation int
	var idPartition int
	internalFragmentation = math.MaxInt

	for index := range m.partitions {
		partition := &m.partitions[index]
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
		if p.time > 0 {
			*queue = append(*queue, *p)
		}
	} else {
		os.time = os.time + p.time
		os.addReady(cola)
		p.time = 0
		fmt.Println("Termino el proceso: ", p.pid)
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
	os.queue = []Process{}     // Inicializa la cola
	os.processor = Processor{} // Inicializa el procesador
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
	if partition.process != nil && partition.process.pid != 0 {
		processName = fmt.Sprintf("Process-%d", partition.process.pid)
	}
	fmt.Printf("| %-10d | %-10d | %-10s | %-15d | %-10s |\n",
		partition.id, partition.size, state, partition.internalFragmentation, processName)
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
					bestFit(&linux.memory, &linux.queue[0])
				}
				linux.processor.process = linux.queue[0]
				linux.queue = append(linux.queue[1:])
				linux.addReady(&cola)

			} else {
				//contemplar que es la primera vez y se puede empezar en algo distinto que 0
				linux.time = cola[0].arrivalTime
				bestFit(&linux.memory, &cola[0])
				fmt.Print("leyenda")
				fmt.Print(*&linux.memory.partitions[2])
				linux.processor.process = cola[0]
				cola = append(cola[1:])
			}

			if len(linux.queue) == 0 && len(cola) == 0 && linux.processor.process.time <= 0 { //ver por que no funciona con el igual
				fmt.Println("Se termino de procesar todo - Fin de la Simulacion")
				break
			}
			fmt.Println("TIME: ", linux.time, "----------------------------------------------------")
			fmt.Println("El proceso que se encuentra en el procesador es: pid", linux.processor.process.pid, linux.processor.process.loaded)
			fmt.Println("Esta es la cola de listos", linux.queue)
			fmt.Println("Esta es la cola de input", cola)
			fmt.Println("----------------------------------------------------------------------")
			fmt.Printf("| %-10s | %-10s | %-10s | %-15s | %-10s |\n", "ID", "Size", "State", "Internal Frag.", "Process")
			fmt.Println("------------------------------------------------------------------------")
			printPartitionInfo(linux.memory.partitions[0])
			printPartitionInfo(linux.memory.partitions[1])
			printPartitionInfo(linux.memory.partitions[2])
			fmt.Println("------------------------------------------------------------------------")

		} else {
			break
		}
	}

}
