package main

type OS struct {
	time      int
	processor Processor
	memory    Memory
	ReadyQueue
	completedProcesses []Process
}

type Process struct {
	pid            int
	size           int
	arrivalTime    int
	turnaroundTime int
	time           int
	loaded         bool
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

type NewQueue struct {
	queue []Process
}

type ReadyQueue struct {
	queue []Process
}
