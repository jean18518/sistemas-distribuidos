package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Println("Servicio Validador iniciado en puerto 9001")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error:", err)
			continue
		}
		go manejarValidacion(conn)
	}
}

func manejarValidacion(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		linea := scanner.Text()
		partes := strings.Fields(linea)

		if len(partes) != 3 || partes[0] != "VALIDAR" {
			fmt.Fprintf(conn, "ERROR: Formato incorrecto\n")
			continue
		}

		nombre := partes[1]
		edadStr := partes[2]

		// Validar nombre
		if nombre == "" || nombre == "_" {
			fmt.Fprintf(conn, "ERROR: El nombre no puede estar vacío\n")
			continue
		}

		// Validar edad
		edad, err := strconv.Atoi(edadStr)
		if err != nil || edad < 16 || edad > 100 {
			fmt.Fprintf(conn, "ERROR: La edad debe estar entre 16 y 100\n")
			continue
		}

		// Todo OK
		fmt.Fprintf(conn, "OK\n")
		fmt.Printf("Validado: %s, %d años\n", nombre, edad)
	}
}