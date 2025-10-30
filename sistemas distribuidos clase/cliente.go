package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9003")
	if err != nil {
		log.Fatal("Error: Gateway no disponible. ¿Está corriendo?")
	}
	defer conn.Close()

	// Leer mensaje de bienvenida
	scanner := bufio.NewScanner(conn)
	for i := 0; i < 2; i++ {
		if scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}

	// Leer comandos del usuario
	entrada := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("\n> ")
		if !entrada.Scan() {
			break
		}

		comando := entrada.Text()

		if strings.ToUpper(comando) == "SALIR" {
			fmt.Println("Hasta luego!")
			break
		}

		// Enviar comando al gateway
		fmt.Fprintf(conn, "%s\n", comando)

		// Leer respuesta
		if scanner.Scan() {
			respuesta := scanner.Text()
			fmt.Println(respuesta)

			// Si es LISTAR, leer múltiples líneas
			if strings.HasPrefix(comando, "LISTAR") {
				for scanner.Scan() {
					linea := scanner.Text()
					if strings.HasPrefix(linea, "ERROR") ||
						strings.HasPrefix(linea, "OK") {
						break
					}
					fmt.Println(linea)
				}
			}
		}
	}
}