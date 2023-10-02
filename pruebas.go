package main


/*
func prueba() {


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
    cola = append(cola, Process{1, 30, 1, 3, true}, Process{2, 40, 1, 7, true}, Process{3, 100, 2, 10, true})
    var input string
    fmt.Print("Inicio del Sistema Operativo")

    for {
        fmt.Scanln(&input)
        if input == "" {
            if !linux.processor.process.isEmpty() {
                linux.processor.process.timeOut(5, linux.queue)
                linux.processor.process = linux.queue[0]
                fmt.Println("Duende")       
            }

            if !linux.processor.process.loaded {
                bestFit(linux.memory, linux.processor.process)
            }

           linux.addReady(cola)
           fmt.Println("El proceso que se encuentra en el procesador es: pid", linux.processor.process.pid)
           fmt.Println("Esta es la cola de listos", linux.queue)

        } else {
            break
        }
    }
}
*/
