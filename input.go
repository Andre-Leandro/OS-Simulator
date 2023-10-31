package main

import (
    "fmt"
    "os"
)

func main() {
    // Nombre del archivo que deseas abrir
    archivo := "ejemplo.txt"

    // Lee el contenido del archivo
    contenido, err := os.ReadFile(archivo)
    if err != nil {
        fmt.Println("Error al leer el archivo:", err)
        return
    }

    // Convierte el contenido en una cadena y luego imprímelo
    fmt.Println(string(contenido))
}
