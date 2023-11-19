package main

import (
	"fmt"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

const (
	//DataTable
	columnKeySize                  = "size"
	columnKeyState                 = "state"
	columnKeyInternalFragmentation = "internalFragmentation"
	columnKeyProcess               = "process"
	//Processor
	columnKeyProcessor          = "processor"
	columnKeyProcessInProcessor = "processInProcessor"
	//ReadyQueue
	columnKeyLoaded      = "loaded"
	columnKeyArrivalTime = "arrivalTime"
	columnKeyTime        = "time"
)

var (
/* styleSubtle = lipgloss.NewStyle().Foreground(lipgloss.Color("#888")) */
/*
	styleBase = lipgloss.NewStyle().
		Foreground(lipgloss.Color("white")).
		BorderForeground(lipgloss.Color("white")).
		Align(lipgloss.Center) */
)

func NewModelShowData(memoria Memory, proceso Process) Model {

	// Crear columnas con estilo base
	columns := []table.Column{
		table.NewColumn(columnKeySize, "Tamaño (Kb)", 15).WithStyle(styleBase),
		table.NewColumn(columnKeyState, "Estado", 15).WithStyle(styleBase),
		table.NewColumn(columnKeyInternalFragmentation, "Fragm. Interna (Kb)", 25).WithStyle(styleBase),
		table.NewColumn(columnKeyProcess, "Proceso", 20).WithStyle(styleBase),
	}

	// Crear filas con datos de todos los procesos
	var allRows []table.Row

	for i, m := range memoria.partitions {
		partition := memoria.partitions[i]
		colorState := "#f64"
		colorInProcessor := "white"

		state := "Ocupado"
		if partition.state {
			state = "Libre"
			colorState = "#8b8"
		}

		processName := "-"
		if partition.process.pid != 0 {
			processName = fmt.Sprintf("Proceso %d", partition.process.pid)
		}

		if partition.process.pid == proceso.pid {
			colorInProcessor = "#44f"
		}

		coloredState := lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorState)).
			Render(state)

		coloredInProccesor := lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorInProcessor)).
			Render(processName)

		row := table.NewRow(table.RowData{
			columnKeySize:                  m.size,
			columnKeyState:                 coloredState,
			columnKeyInternalFragmentation: m.internalFragmentation,
			columnKeyProcess:               coloredInProccesor,
		})
		allRows = append(allRows, row)

	}

	return Model{
		tabla: table.New(columns).
			WithRows(allRows).
			BorderRounded(),
	}
}

func NewModelShowProcessInProcessor(proceso Process) Model {

	// Crear columnas con estilo base
	columns := []table.Column{table.NewColumn(columnKeyProcessor, "Procesador", 15).WithStyle(styleBase),
		table.NewColumn(columnKeyProcessInProcessor, "NumeroPid", 5).WithStyle(styleBase)}

	// Crear filas con datos de todos los procesos
	var allRows []table.Row
	row := table.NewRow(table.RowData{
		columnKeyProcessor:          "PROCESADOR",
		columnKeyProcessInProcessor: proceso.pid,
	})
	allRows = append(allRows, row)
	return Model{
		tabla: table.New(columns).
			WithRows(allRows).WithHeaderVisibility(false).
			BorderRounded(),
	}
}

func NewModelShowReadyQueue(colaListos []Process) Model {
	// Crear columnas con estilo base
	columns := []table.Column{
		table.NewColumn(columnKeyPid, "Pid", 6).WithStyle(styleBase),
		table.NewColumn(columnKeySize, "Tamaño (Kb)", 15).WithStyle(styleBase),
		table.NewColumn(columnKeyArrivalTime, "Tiempo de Arribo", 25).WithStyle(styleBase),
		table.NewColumn(columnKeyTime, "Tiempo", 14).WithStyle(styleBase),
		table.NewColumn(columnKeyLoaded, "Cargado", 14).WithStyle(styleBase),
	}

	// Crear filas con datos de todos los procesos
	var allRows []table.Row

	if len(colaListos) == 0 {
		row := table.NewRow(table.RowData{
			columnKeyPid:         "-",
			columnKeySize:        "-",
			columnKeyArrivalTime: "-",
			columnKeyTime:        "-",
			columnKeyLoaded:      "-",
		})
		allRows = append(allRows, row)
	} else {
		for i, _ := range colaListos {

			loaded := "Disco"
			if colaListos[i].loaded {
				loaded = "Memoria"
			}

			row := table.NewRow(table.RowData{
				columnKeyPid:         colaListos[i].pid,
				columnKeySize:        colaListos[i].size,
				columnKeyArrivalTime: colaListos[i].arrivalTime,
				columnKeyTime:        colaListos[i].time,
				columnKeyLoaded:      loaded,
			})
			allRows = append(allRows, row)

		}
	}
	return Model{
		tabla: table.New(columns).
			WithRows(allRows).
			BorderRounded(),
	}

}

/* func NewModelAverage(averageTurnaroundTime float64, averageWaitTime float64) Model {
	// Crear columnas sin título
	columns := []table.Column{
		table.NewColumn(columnKeyAverageName, "Promedio", 26).WithStyle(styleBase),
		table.NewColumn(columnKeyAverageValue, "Valor", 20).WithStyle(styleBase),
	}

	// Crear filas con datos de promedio
	var rows []table.Row
	rows = append(rows, table.NewRow(table.RowData{
		columnKeyAverageName:  "P. Tiempo Retorno",
		columnKeyAverageValue: averageTurnaroundTime,
	}))
	rows = append(rows, table.NewRow(table.RowData{
		columnKeyAverageName:  "P. Tiempo Espera",
		columnKeyAverageValue: averageWaitTime,
	}))

	return Model{
		tabla: table.New(columns).
			WithRows(rows).
			BorderRounded(),
	}
} */

func mostrarDatos(memoria Memory, proceso Process) {
	p := tea.NewProgram(NewModelShowData(memoria, proceso))

	go func() {
		time.Sleep(0)                          // Espera 3 segundos (ajusta según sea necesario)
		p.Send(tea.KeyMsg{Type: tea.KeyCtrlC}) // Envía un mensaje de tecla para cerrar la aplicación
	}()

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}

func mostrarProcesador(proceso Process) {
	j := tea.NewProgram(NewModelShowProcessInProcessor(proceso))

	go func() {
		time.Sleep(0)
		j.Send(tea.KeyMsg{Type: tea.KeyCtrlC}) // Envía un mensaje de tecla para cerrar la aplicación
	}()

	if err := j.Start(); err != nil {
		log.Fatal(err)
	}

}

func mostrarColaListos(colaListos []Process) {
	j := tea.NewProgram(NewModelShowReadyQueue(colaListos))

	go func() {
		time.Sleep(0)
		j.Send(tea.KeyMsg{Type: tea.KeyCtrlC}) // Envía un mensaje de tecla para cerrar la aplicación
	}()

	if err := j.Start(); err != nil {
		log.Fatal(err)
	}

}
