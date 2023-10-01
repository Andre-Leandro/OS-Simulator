package main

import (
    "fmt"
)

type Process struct {
    pid int
    size int
    arrivalTime int
    time int
}

type Processor struct {
    process Process
}

type MemoryPartition struct {
    id int
    size int
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

func (p *Process) timeOut(quantum int, queue *ReadyQueue){
    if p.time >= quantum{
        queue.queue = append(queue.queue[1:], queue.queue[0])
        p.time = p.time - quantum 
    } else {
        queue.queue = append(queue.queue[1:])
        p.time = 0
    }
}

func (p* Process) Escribir(){
    p.pid = 5
    fmt.Println(p.time)
}

var a NewQueue

func main() {
    // OS()
}

func OS() {
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