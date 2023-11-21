package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func ReadProcessesFromFile(filename string, pidMap map[int]bool) ([]Process, error) {
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
		if isPIDInUse(pid, pidMap) {
			return nil, fmt.Errorf("Error: El PID %d se repite más de una vez. Por favor, ingrese PIDs únicos.", pid)
		} else {
			pidMap[pid] = true
		}
		size, err := strconv.Atoi(values[1])
		if size <= 0 {
			return nil, fmt.Errorf("Error: El tamaño de los procesos debe ser mayor a 0 y menor a 250 kB.")
		}
		if err != nil {
			return nil, err
		}
		arrivalTime, err := strconv.Atoi(values[2])
		if arrivalTime < 0 {
			return nil, fmt.Errorf("Error: El tiempo de arribo de los procesos debe ser mayor a 0")
		}
		if err != nil {
			return nil, err
		}
		time, err := strconv.Atoi(values[3])
		if time <= 0 {
			return nil, fmt.Errorf("Error: El tiempo de irrupción de los procesos debe ser mayor a 0")
		}
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

func clearScreen() {
	switch runtime.GOOS {
	case "linux", "darwin": // para Linux y macOS
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "windows": // para Windows
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		fmt.Println("No se pudo determinar el sistema operativo para realizar el clear screen.")
	}
}

func isPIDInUse(pid int, pidMap map[int]bool) bool {
	_, exists := pidMap[pid]
	return exists
}
