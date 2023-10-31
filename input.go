package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
    archivo := "ejemplo.txt"
    
    file, err := os.Open(archivo)
    if err != nil {
        fmt.Println("Error al abrir el archivo:", err)
        return
    }
    defer file.Close()


	type Process struct {
		pid         int
		size        int
		arrivalTime int
		time        int
		loaded      bool
	}
	
    processes := []Process{}

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        values := strings.Fields(line)
        if len(values) != 4 {
            fmt.Printf("Invalid input format: %s\n", line)
            continue
        }

        pid, _ := strconv.Atoi(values[0])
        size, _ := strconv.Atoi(values[1])
        arrivalTime, _ := strconv.Atoi(values[2])
        time, _ := strconv.Atoi(values[3])

        process := Process{pid, size, arrivalTime, time, false}
        processes = append(processes, process)
    }

    if err := scanner.Err(); err != nil {
        fmt.Println("Error al leer el archivo:", err)
        return
    }

    // Ahora, processes contiene la lista de procesos desde el archivo
    fmt.Println("Procesos leídos desde el archivo:")
    for _, p := range processes {
        fmt.Printf("PID: %d, Size: %d, Arrival Time: %d, Time: %d\n", p.pid, p.size, p.arrivalTime, p.time)
    }
}