package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Dos Cuadrados")

	// Ancho común para ambos cuadrados
	ancho := float32(150)
	memoriaColor := color.RGBA{R: 0, G: 255, B: 55, A: 255}
	procesoColor := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	bordeColor := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	// Altura y color del primer cuadrado
	memoriaGrandeTamano := float32(250)

	// Altura y color del segundo cuadrado
	procesoGrandeTamano := float32(100)

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
	// Contenedor vertical para los cuadrados y el texto
	content := container.NewWithoutLayout(memoriaGrandeContainer, memoriaMedianaContainer, memoriaPequenaContainer)
	content.Resize(fyne.NewSize(250, 370))
	myWindow.SetContent(content)

	myWindow.SetContent(content)

	myWindow.Resize(fyne.NewSize(800, 750))

	myWindow.ShowAndRun()
}
