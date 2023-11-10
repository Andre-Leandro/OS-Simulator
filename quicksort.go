package main

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