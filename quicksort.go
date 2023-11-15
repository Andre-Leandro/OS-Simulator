package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ReadProcessesFromFile2(filename string) ([]Process, error) {
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

		process := Process{pid, size, arrivalTime, time, false, 0}
		processes = append(processes, process)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return processes, nil
}

func quicksort(processes []Process) []Process {
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
			equal = append(equal, p)
		}
	}
	return append(append(quicksort(less), equal...), quicksort(greater)...)
}

func main2() {
	processes, err := ReadProcessesFromFile("ejemplo.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	pepe := quicksort(processes)
	fmt.Println(pepe)
}
