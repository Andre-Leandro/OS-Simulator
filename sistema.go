package main

import (
    "fmt"
    "math"
)

type OS struct {
    time int
    processor Processor
    memory Memory
    ReadyQueue
}

type Process struct {
    pid int
    size int
    arrivalTime int
    time int
    loaded bool
}

type Processor struct {
    process Process
}

type MemoryPartition struct {
    id int
    size int
    state bool
    process Process
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

func (os *OS) addReady(l []Process) {
    for index:= range l {
        if len(os.queue) == 5 {
            break
        }
        if l[index].arrivalTime < os.time {
            os.queue = append(os.queue, l[index])
        }
        bestFitLazy(os.memory, l[index])
    }
}


func bestFitLazy(m Memory, p Process) {
    var internalFragmentation int
    var idPartition int
    internalFragmentation = math.MaxInt

    for index:= range m.partitions {
        partition:= m.partitions[index]
        if partition.state && partition.size >= p.size {
            empty:= partition.size - p.size
            if empty < internalFragmentation {
                idPartition = index
            }
        }
        if idPartition != 0 {
            selectedPartition:= &m.partitions[idPartition]
            selectedPartition.state = false
            selectedPartition.internalFragmentation = selectedPartition.size - p.size
            selectedPartition.process = p
        }
    }
} 

func bestFit(m Memory, p Process) {
    var internalFragmentation int
    var idPartition int
    internalFragmentation = math.MaxInt

    for index:= range m.partitions {
        partition:= m.partitions[index]
        if partition.state && partition.size >= p.size {
            empty:= partition.size - p.size
            if empty < internalFragmentation {
                idPartition = index
            }
        }
        if idPartition == 0 {
            index:= bestFitSwap(m, p)
            swapOut(index)
        }

        selectedPartition:= &m.partitions[idPartition]
        selectedPartition.state = false
        selectedPartition.internalFragmentation = selectedPartition.size - p.size
        selectedPartition.process = p
        p.loaded = true
    } 
}

func bestFitSwap(m Memory, p Process) int {
    var internalFragmentation int
    var idPartition int
    internalFragmentation = math.MaxInt

    for index:= range m.partitions {
        partition:= m.partitions[index]
        if partition.size >= p.size {
            empty:= partition.size - p.size
            if empty < internalFragmentation {
                idPartition = index
            }
        }
    } 
    return idPartition
}

func (p *Process) timeOut(quantum int, queue *ReadyQueue){
    if p.time >= quantum{
        queue.queue = append(queue.queue[1:], queue.queue[0])
        p.time = p.time - quantum 
    } else {
        queue.queue = append(queue.queue[1:])
        p.time = 0
    }
}

func sort (input ReadyQueue) {
    for i:= range input.queue {
        fmt.Print( "este es el numero ", i)
        for j:= range input.queue[i:] {
            fmt.Print(j+i, " - ")
            
        }
        fmt.Println("")
    }
}

func (p* Process) Escribir(){
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
    cola:= new(ReadyQueue) 
    linux:= new(OS)
    memoria := Memory{
        partitions: [3]MemoryPartition{
            {id: 1, size: 100, state: false},
            {id: 2, size: 75, state: false},
            {id: 3, size: 35, state: false},
        },
    }
    linux.initialize(memoria)

    var input string
    fmt.Print("¿Desea continuar?: ")
    for {
        fmt.Scanln(&input)
        if input == "" {
            if !linux.processor.process.isEmpty() {
                linux.processor.process.timeOut(5, cola)
                linux.processor.process = cola.queue[0]       
            }

            if !linux.processor.process.loaded {
                bestFit(linux.memory, linux.processor.process)
            }

           // linux.addReady()

            fmt.Println("Bien ahi")

        } else {
            break
        }
    }
}

func O() {
    var a [4]Process
    a[2] = Process{pid: 4}
    var p Process
    p = Process{pid: 7, time: 78}
    cola:= new(ReadyQueue)
    cola.queue = append(cola.queue, p, a[1])
    p.Escribir()
    cola.estado()
    p.timeOut(500, cola)
    cola.estado()
    p.Escribir()
}