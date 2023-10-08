package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Sistemas Opertivos")

	// Ancho común para ambos cuadrados
	ancho := float32(150)
	memoriaColor := color.RGBA{R: 0, G: 255, B: 55, A: 255}
	procesoColor := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	bordeColor := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	// Altura y color del primer cuadrado
	memoriaGrandeTamano := float32(250)

	// Altura y color del segundo cuadrado
	procesoGrandeTamano := float32(250)

	// Cuadrado 1
	memoriaGrande := canvas.NewRectangle(memoriaColor)
	memoriaGrande.Resize(fyne.NewSize(ancho, memoriaGrandeTamano))
	memoriaGrande.StrokeColor = bordeColor
	memoriaGrande.StrokeWidth = 1

	// Agregar texto centrado en square1
	textoProcesoGrande := widget.NewLabel("P1")
	textoProcesoGrande.Move(fyne.NewPos(0, (procesoGrandeTamano-textoProcesoGrande.MinSize().Height)/2))

	// Cuadrado 2
	procesoGrande := canvas.NewRectangle(procesoColor)
	procesoGrande.Resize(fyne.NewSize(ancho, procesoGrandeTamano))
	procesoGrande.StrokeColor = bordeColor
	procesoGrande.StrokeWidth = 1

	// Contenedor para los cuadrados y el texto
	memoriaGrandeContainer := container.NewWithoutLayout(
		memoriaGrande,
		canvas.NewRectangle(color.Transparent), // Espacio en blanco entre los cuadrados
		procesoGrande,
		textoProcesoGrande,
	)

	// Altura y color del primer cuadrado (Memoria Mediana)
	memoriaMedianaTamano := float32(120)

	// Altura y color del segundo cuadrado (Proceso Mediano)
	procesoMedianoTamano := float32(50)

	// Cuadrado 1 (Memoria Mediana)
	memoriaMediana := canvas.NewRectangle(memoriaColor)
	memoriaMediana.Resize(fyne.NewSize(ancho, memoriaMedianaTamano))
	memoriaMediana.StrokeColor = bordeColor
	memoriaMediana.StrokeWidth = 1

	// Cuadrado 2 (Proceso Mediano)
	procesoMediano := canvas.NewRectangle(procesoColor)
	procesoMediano.Resize(fyne.NewSize(ancho, procesoMedianoTamano))
	procesoMediano.StrokeColor = bordeColor
	procesoMediano.StrokeWidth = 1

	// Agregar texto centrado en procesoMediano
	textoProcesoMediano := widget.NewLabel("P2")
	textoProcesoMediano.Move(fyne.NewPos(0, (procesoMedianoTamano-textoProcesoMediano.MinSize().Height)/2))

	// Contenedor para los cuadrados y el texto
	memoriaMedianaContainer := container.NewWithoutLayout(
		memoriaMediana,
		canvas.NewRectangle(color.Transparent), // Espacio en blanco entre los cuadrados
		procesoMediano,
		textoProcesoMediano,
	)

	memoriaMedianaContainer.Move(fyne.NewPos(0, 250))

	// Altura y color del primer cuadrado (Memoria Pequeña)
	memoriaPequenaTamano := float32(60)

	// Altura y color del segundo cuadrado (Proceso Pequeño)
	procesoPequenoTamano := float32(50)

	// Cuadrado 1 (Memoria Pequeña)
	memoriaPequena := canvas.NewRectangle(memoriaColor)
	memoriaPequena.Resize(fyne.NewSize(ancho, memoriaPequenaTamano))
	memoriaPequena.StrokeColor = bordeColor
	memoriaPequena.StrokeWidth = 1

	// Cuadrado 2 (Proceso Pequeño)
	procesoPequeno := canvas.NewRectangle(procesoColor)
	procesoPequeno.Resize(fyne.NewSize(ancho, procesoPequenoTamano))
	procesoPequeno.StrokeColor = bordeColor
	procesoPequeno.StrokeWidth = 1

	// Agregar texto centrado en procesoPequeno
	textoProcesoPequeno := widget.NewLabel("P3")
	textoProcesoPequeno.Move(fyne.NewPos(0, (procesoPequenoTamano-textoProcesoPequeno.MinSize().Height)/2))

	// Contenedor para los cuadrados y el texto
	memoriaPequenaContainer := container.NewWithoutLayout(
		memoriaPequena,
		canvas.NewRectangle(color.Transparent), // Espacio en blanco entre los cuadrados
		procesoPequeno,
		textoProcesoPequeno,
	)

	memoriaPequenaContainer.Move(fyne.NewPos(0, 370))

	// Parte superior con botón, etiqueta y otro botón centrados
	button1 := widget.NewButton("Restart", func() {
		// Acción del primer botón
	})
	label := widget.NewLabel("Timer")
	button2 := widget.NewButton("Next", func() {
		// Acción del segundo botón
	})

	// Centrar los elementos en la parte superior
	top := container.NewHBox(
		layout.NewSpacer(),
		button1,
		layout.NewSpacer(),
		label,
		layout.NewSpacer(),
		button2,
		layout.NewSpacer(),
	)

	cpu := canvas.NewRectangle(procesoColor)
	cpu.Resize(fyne.NewSize(150, 150))
	cpu.StrokeColor = bordeColor
	cpu.StrokeWidth = 1

	textoCpu := widget.NewLabel("PID")
	textoCpu.Move(fyne.NewPos(0, (150-textoCpu.MinSize().Height)/2))

	cpuContainer := container.NewWithoutLayout(
		cpu,
		textoCpu,
	)

	memory := container.NewVBox(
		widget.NewLabel("Memory"),
		container.NewWithoutLayout(memoriaGrandeContainer, memoriaMedianaContainer, memoriaPequenaContainer),
	)

	process := container.NewVBox(
		widget.NewLabel("CPU"),
		cpuContainer,
	)

	disk := container.NewVBox(
		widget.NewLabel("Disk"),
		widget.NewLabel("Contenido 3"),
	)

	// Contenido debajo de la parte superior
	content := container.NewVBox(
		top,
		container.NewHBox(memory,
			layout.NewSpacer(),
			process,
			layout.NewSpacer(),
			disk),
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(800, 600))

	w.ShowAndRun()
}
