package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

var estudiantes []string

func main() {
	listener, err := net.Listen("tcp", ":9002")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Println("Servicio Registro iniciado en puerto 9002")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error:", err)
			continue
		}
		go manejarRegistro(conn)
	}
}

func manejarRegistro(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		linea := scanner.Text()
		partes := strings.Fields(linea)

		if len(partes) < 1 {
			fmt.Fprintf(conn, "ERROR: Comando vacío\n")
			continue
		}

		comando := partes[0]

		if comando == "GUARDAR" {
			if len(partes) != 3 {
				fmt.Fprintf(conn, "ERROR: Formato incorrecto\n")
				continue
			}

			nombre := partes[1]
			edad := partes[2]
			estudiante := fmt.Sprintf("%s - %s años", nombre, edad)

			estudiantes = append(estudiantes, estudiante)

			fmt.Fprintf(conn, "GUARDADO\n")
			fmt.Printf("Guardado: %s\n", estudiante)

		} else if comando == "LISTAR" {
			if len(estudiantes) == 0 {
				fmt.Fprintf(conn, "No hay estudiantes registrados\n")
			} else {
				fmt.Fprintf(conn, "ESTUDIANTES:\n")
				for i, est := range estudiantes {
					fmt.Fprintf(conn, "%d. %s\n", i+1, est)
				}
			}
			fmt.Fprintf(conn, "FIN\n")

		} else {
			fmt.Fprintf(conn, "ERROR: Comando desconocido\n")
		}
	}
}