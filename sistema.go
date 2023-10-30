package main

import (
	"fmt"
	"math"
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

	for index := range copy {
		if len(os.queue) == 5 {
			break
		}

		if copy[index].arrivalTime < os.time { // menor o igual tambien 
			os.queue = append(os.queue, copy[index])
			if len(*l) == 1 {
				*l = append([]Process{} )
			} else {
				*l = append((*l)[index+1:])
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

func bestFit(m *Memory, p *Process) {
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
		idPartition = bestFitSwap(*m, *p)
		swapOut(idPartition)
		m.partitions[idPartition].process.loaded = false
	}
	selectedPartition := &m.partitions[idPartition]
	selectedPartition.state = false
	selectedPartition.internalFragmentation = selectedPartition.size - p.size
	p.loaded = true
	selectedPartition.process = *p
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
	if p.time >= quantum {
		os.time = os.time + quantum
		os.addReady(cola)
		p.time = p.time - quantum
		*queue = append(*queue, *p)
	} else {
		os.time = os.time + p.time
		os.addReady(cola)
		p.time = 0
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
	fmt.Printf("hola")
}

func (os *OS) initialize(m Memory) {
	os.memory = m
}

func main() {
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
	cola = append(cola, Process{1, 30, 1, 3, false}, Process{2, 40, 1, 7, false}, Process{3, 100, 2, 10, false})
	var input string
	fmt.Print("Inicio del Sistema Operativo")

	for {
		fmt.Scanln(&input)
		if input == "" {
			if !linux.processor.process.isEmpty() {
				linux.processor.process.timeOut(5, &linux.queue, linux, &cola)
				if linux.queue[0].loaded == false {
					bestFit(&linux.memory, &linux.queue[0])
				}
				linux.processor.process = linux.queue[0]
				linux.queue = append(linux.queue[1:])
				
			} else {
				bestFit(&linux.memory, &cola[0])
				linux.processor.process = cola[0]
				cola = append(cola[1:])
				
			}
			if len(linux.queue) == 0 && linux.processor.process.time == 0 {
				fmt.Println("Se termino de procesar todo - Fin de la Simulacion")
				break
			}
			fmt.Println("TIME: ", linux.time, "----------------------------------------------------")
			fmt.Println("El proceso que se encuentra en el procesador es: pid", linux.processor.process.pid)
			fmt.Println("Esta es la cola de listos", linux.queue)
			fmt.Println("Esta es la cola de input", cola)

		} else {
			break
		}
	}

}