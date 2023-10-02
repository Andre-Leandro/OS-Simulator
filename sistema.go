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

	for index := range *l {
		if len(os.queue) == 5 {
			break
		}

		if copy[index].arrivalTime < os.time {
			os.queue = append(os.queue, copy[index])
            *l = append((*l)[index:])
            fmt.Println("Prueba de puntero ", l)
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
	internalFragmentation = math.MaxInt

	for index := range m.partitions {
		partition := m.partitions[index]
        //fmt.Println("partition", partition)
		if partition.state && partition.size >= p.size {
			//fmt.Println(index, "Entro prof1")
			empty := partition.size - p.size
			if empty < internalFragmentation {
				//fmt.Println("Entre 2")
				idPartition = index
			}
		}
	}

	if idPartition == 0 {
		index := bestFitSwap(*m, *p)
		swapOut(index)
	}

	selectedPartition := &m.partitions[idPartition]

    //fmt.Println(selectedPartition)
	selectedPartition.state = false
    //fmt.Println(selectedPartition.state)
    //fmt.Println(m.partitions[idPartition].state)
	selectedPartition.internalFragmentation = selectedPartition.size - p.size
    p.loaded = true
	selectedPartition.process = *p
    //fmt.Println(p.loaded)
    //fmt.Println(selectedPartition)
    //fmt.Println(m.partitions[idPartition])
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

func (p *Process) timeOut(quantum int, queue []Process, os *OS, cola *[]Process) {
	if p.time >= quantum {
        os.time = os.time + quantum
        os.addReady(cola)


		queue = append(queue, *p)
		p.time = p.time - quantum
	} else {
        os.time = os.time + p.time
        os.addReady(cola)
        fmt.Println("cacique: ", os.queue)


        if len(queue) > 1 {
            queue = append(queue[1:])
            p.time = 0
        } 
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
     // OS definition 
     var cola []Process 
     linux:= new(OS)
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
                 linux.processor.process.timeOut(5, linux.queue, linux, &cola)
                 linux.processor.process = linux.queue[0]
                 linux.queue = append(linux.queue[1:])
                 fmt.Println("Duende")       
             } else {
                linux.processor.process = cola[0]
                cola = append(cola[1:] )
                 bestFit(&linux.memory, &linux.processor.process)
             }
            fmt.Println("TIME: ", linux.time,  "----------------------------------------------------")
            fmt.Println("El proceso que se encuentra en el procesador es: pid", linux.processor.process.pid)
            fmt.Println("Esta es la cola de listos", linux.queue)
            fmt.Println("Esta es la cola de input", cola)

 
         } else {
             break
         }
     }

}