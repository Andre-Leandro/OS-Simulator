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
	columnKeySize                  = "size"
	columnKeyState                 = "state"
	columnKeyInternalFragmentation = "internalFragmentation"
	columnKeyProcess               = "process"
)

var (
/* styleSubtle = lipgloss.NewStyle().Foreground(lipgloss.Color("#888")) */
/*
	styleBase = lipgloss.NewStyle().
		Foreground(lipgloss.Color("white")).
		BorderForeground(lipgloss.Color("white")).
		Align(lipgloss.Center) */
)

func NewModelShowData(memoria Memory) Model {

	// Crear columnas con estilo base
	columns := []table.Column{
		table.NewColumn(columnKeySize, "Tamaño", 20).WithStyle(styleBase),
		table.NewColumn(columnKeyState, "Estado", 20).WithStyle(styleBase),
		table.NewColumn(columnKeyInternalFragmentation, "Fragmentacion Interna", 20).WithStyle(styleBase),
		table.NewColumn(columnKeyProcess, "Proceso", 20).WithStyle(styleBase),
	}

	// Crear filas con datos de todos los procesos
	var allRows []table.Row

	for i, m := range memoria.partitions {
		partition := memoria.partitions[i]
		colorState := "#f64"

		state := "Occupied"
		if partition.state {
			state = "Free"
			colorState = "#8b8"
		}
		processName := "N/A"
		if partition.process.pid != 0 {
			processName = fmt.Sprintf("Process-%d", partition.process.pid)
		}

		coloredState := lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorState)).
			Render(state)

		row := table.NewRow(table.RowData{
			columnKeySize:                  m.size,
			columnKeyState:                 coloredState,
			columnKeyInternalFragmentation: m.internalFragmentation,
			columnKeyProcess:               processName,
		})
		allRows = append(allRows, row)

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

func mostrarDatos(memoria Memory) {
	fmt.Println("hola")

	p := tea.NewProgram(NewModelShowData(memoria))

	go func() {
		time.Sleep(0)                          // Espera 3 segundos (ajusta según sea necesario)
		p.Send(tea.KeyMsg{Type: tea.KeyCtrlC}) // Envía un mensaje de tecla para cerrar la aplicación
	}()

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
