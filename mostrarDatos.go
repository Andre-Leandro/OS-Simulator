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
	columnKeyDirInicio             = "dirInicio"
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

func NewModelWithTitle(title string) Model {
	columns := []table.Column{
		table.NewColumn(columnKeySize, title, 78).WithStyle(styleBase),
	}

	return Model{
		tabla: table.New(columns).
			BorderRounded(),
	}
}

func NewModelShowData(memoria Memory, proceso Process) Model {

	// Crear columnas con estilo base
	columns := []table.Column{
		table.NewColumn(columnKeyDirInicio, "Dir. Inicio", 15).WithStyle(styleBase),
		table.NewColumn(columnKeySize, "Tamaño (Kb)", 15).WithStyle(styleBase),
		table.NewColumn(columnKeyState, "Estado", 15).WithStyle(styleBase),
		table.NewColumn(columnKeyInternalFragmentation, "Fragm. Interna (Kb)", 22).WithStyle(styleBase),
		table.NewColumn(columnKeyProcess, "Proceso", 19).WithStyle(styleBase),
	}

	// Crear filas con datos de todos los procesos
	var allRows []table.Row

	for i, m := range memoria.partitions {
		partition := memoria.partitions[i]
		colorState := "#f64"
		colorInProcessor := "white"
		var dirInicio string

		switch i {
		case 0:
			dirInicio = "0x1D6" //470
		case 1:
			dirInicio = "0x15E" //350
		case 2:
			dirInicio = "0x64" //100
		default:
			fmt.Println("Partición no reconocida")
		}

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
			columnKeyDirInicio:             dirInicio,
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
		table.NewColumn(columnKeyPid, "PID", 9).WithStyle(styleBase),
		table.NewColumn(columnKeyArrivalTime, "Tiempo de Arribo", 21).WithStyle(styleBase),
		table.NewColumn(columnKeySize, " Tamaño (Kb)", 15).WithStyle(styleBase),
		table.NewColumn(columnKeyTime, "T. Restante", 22).WithStyle(styleBase),
		table.NewColumn(columnKeyLoaded, "Cargado", 19).WithStyle(styleBase),
	}

	// Crear filas con datos de todos los procesos
	var allRows []table.Row

	if len(colaListos) == 0 {
		row := table.NewRow(table.RowData{
			columnKeyPid:         "-",
			columnKeyArrivalTime: "-",
			columnKeySize:        "-",
			columnKeyTime:        "-",
			columnKeyLoaded:      "-",
		})
		allRows = append(allRows, row)
	} else {
		for i, _ := range colaListos {

			loaded := "Disco"
			colorLoaded := "#DC7633"

			if colaListos[i].loaded {
				loaded = "Memoria"
				colorLoaded = "#F4D03F"
			}

			coloredLoaded := lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorLoaded)).
				Render(loaded)

			row := table.NewRow(table.RowData{
				columnKeyPid:         colaListos[i].pid,
				columnKeyArrivalTime: colaListos[i].arrivalTime,
				columnKeySize:        colaListos[i].size,
				columnKeyTime:        colaListos[i].time,
				columnKeyLoaded:      coloredLoaded,
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
func NewModelShowNewQueue(colaNuevos []Process) Model {
	// Crear columnas con estilo base
	columns := []table.Column{
		table.NewColumn(columnKeyPid, "PID", 9).WithStyle(styleBase),
		table.NewColumn(columnKeyArrivalTime, "Tiempo de Irrupción", 21).WithStyle(styleBase),
		table.NewColumn(columnKeySize, " Tamaño (Kb)", 19).WithStyle(styleBase),
		table.NewColumn(columnKeyTime, "Tiempo", 18).WithStyle(styleBase),
	}

	// Crear filas con datos de todos los procesos
	var allRows []table.Row

	if len(colaNuevos) == 0 {
		row := table.NewRow(table.RowData{
			columnKeyPid:         "-",
			columnKeyArrivalTime: "-",
			columnKeySize:        "-",
			columnKeyTime:        "-",
		})
		allRows = append(allRows, row)
	} else {
		for i, _ := range colaNuevos {

			row := table.NewRow(table.RowData{
				columnKeyPid:         colaNuevos[i].pid,
				columnKeyArrivalTime: colaNuevos[i].arrivalTime,
				columnKeySize:        colaNuevos[i].size,
				columnKeyTime:        colaNuevos[i].time,
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

func mostrarDatos2(title string) {
	p := tea.NewProgram(NewModelWithTitle(title))

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

func mostrarColaNuevos(colaNuevos []Process) {
	y := tea.NewProgram(NewModelShowNewQueue(colaNuevos))

	go func() {
		time.Sleep(0)
		y.Send(tea.KeyMsg{Type: tea.KeyCtrlC}) // Envía un mensaje de tecla para cerrar la aplicación
	}()

	if err := y.Start(); err != nil {
		log.Fatal(err)
	}

}
