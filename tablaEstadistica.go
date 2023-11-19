package main

import (
	"fmt"
	"log"
	"math"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

type Model struct {
	tabla table.Model
}

const (
	columnKeyPid            = "pid"
	columnKeyTurnaroundTime = "turnaround"
	columnKeyWaitTime       = "wait"
	columnKeyAverageName    = "averageName"
	columnKeyAverageValue   = "averageValue"
)

var (
	/* styleSubtle = lipgloss.NewStyle().Foreground(lipgloss.Color("#888")) */

	styleBase = lipgloss.NewStyle().
		Foreground(lipgloss.Color("white")).
		BorderForeground(lipgloss.Color("white")).
		Align(lipgloss.Center)
)
var averageTurnaroundTime, averageWaitTime float64

func NewModelStats(completedProcesses []Process, allProcesses []Process) Model {
	// Crear columnas con estilo base
	columns := []table.Column{
		table.NewColumn(columnKeyPid, "PID", 5).WithStyle(styleBase),
		table.NewColumn(columnKeyTurnaroundTime, "Tiempo de Retorno", 20).WithStyle(styleBase),
		table.NewColumn(columnKeyWaitTime, "Tiempo de espera", 20).WithStyle(styleBase),
	}

	// Crear filas con datos de todos los procesos
	var allRows []table.Row
	var totalTurnaroundTime, totalWaitTime int
	for _, p := range completedProcesses {
		// Buscar el proceso correspondiente en allProcesses
		var originalProcess Process
		for _, originalP := range allProcesses {
			if originalP.pid == p.pid {
				originalProcess = originalP
				break
			}
		}

		waitTime := p.turnaroundTime - originalProcess.time
		row := table.NewRow(table.RowData{
			columnKeyPid:            p.pid,
			columnKeyTurnaroundTime: p.turnaroundTime,
			columnKeyWaitTime:       waitTime,
		})
		allRows = append(allRows, row)
		// Sumar tiempos para el cálculo promedio
		totalTurnaroundTime += p.turnaroundTime
		totalWaitTime += waitTime
	}
	// Calcular promedios
	averageTurnaroundTime = math.Round((float64(totalTurnaroundTime)/float64(len(completedProcesses)))*100) / 100
	averageWaitTime = math.Round((float64(totalWaitTime)/float64(len(completedProcesses)))*100) / 100

	return Model{
		tabla: table.New(columns).
			WithRows(allRows).
			BorderRounded(),
	}
}

func NewModelAverage(averageTurnaroundTime float64, averageWaitTime float64) Model {
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
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.tabla, cmd = m.tabla.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "enter":
			cmds = append(cmds, tea.Quit)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	view := lipgloss.JoinVertical(
		lipgloss.Left,
		/* 		styleSubtle.Render("Press q or ctrl+c to quit - Sorted by # Conversations"), */
		m.tabla.View(),
	) + "\n"

	return lipgloss.NewStyle().MarginLeft(1).Render(view)
}

func arrancar(completedProcesses []Process, allProcesses []Process) {
	p := tea.NewProgram(NewModelStats(completedProcesses, allProcesses))

	go func() {
		time.Sleep(0)                          // Espera 3 segundos (ajusta según sea necesario)
		p.Send(tea.KeyMsg{Type: tea.KeyCtrlC}) // Envía un mensaje de tecla para cerrar la aplicación
	}()

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
	fmt.Print("\n")

	j := tea.NewProgram(NewModelAverage(averageTurnaroundTime, averageWaitTime))

	if err := j.Start(); err != nil {
		log.Fatal(err)
	}
}
