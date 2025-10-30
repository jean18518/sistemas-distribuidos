package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":9003")
	if err != nil {
		log.Fatal("Error al iniciar gateway:", err)
	}
	defer listener.Close()
	fmt.Println("Iniciando el puerto gateway")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error:", err)
			continue
		}
		go manejarCliente(conn)
	}

}

func manejarCliente(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	fmt.Fprintf(conn, "Comandos: REGISTRAR nombre edad | LISTAR\n")
	for scanner.Scan() {
		linea := scanner.Text()
		partes := strings.Fields(linea)
		if len(partes) < 1 {
			continue
		}
		comando := partes[0]
		if comando == "REGISTRAR" {
			if len(partes) != 3 {
				fmt.Fprintf(conn, "ERROR debe usar: REGISTRAR nombre edad\n")
				continue
			}
			nombre := partes[1]
			edad := partes[2]

			// Valida con el servicio validador servicio-validador.go
			validacion := llamarServicioValidador(nombre, edad)
			if strings.HasPrefix(validacion, "ERROR") {
				fmt.Fprintf(conn, "%s \n", validacion)
				continue
			}

			// Guardar con el servicio registro
			resultado := llamarServicioRegistro("GUARDAR", nombre, edad)
			if resultado == "GUARDADO" {
				fmt.Fprintf(conn, "se registro al estudiante %s\n", nombre)
			} else {
				fmt.Fprintf(conn, "error al guardar \n")
			}
		} else if comando == "LISTAR" {
			listarEstudiantes(conn)
		} else {
			fmt.Fprintf(conn, "ERROR comando desconocido \n")
		}
	}
}
func llamarServicioValidador(nombre, edad string) string {
	conn, err := net.Dial("tcp", "localhost:9001")
	if err != nil {
		return "ERROR el servicio validador no esta disponible"
	}
	defer conn.Close()

	fmt.Fprintf(conn, "VALIDAR %s %s \n", nombre, edad)
	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		return scanner.Text()
	}
	return "ERROR no hay respuesta del validador"
}

func llamarServicioRegistro(comando, nombre, edad string) string {
	conn, err := net.Dial("tcp", "localhost:9002")
	if err != nil {
		return "ERROR el servicio registro no esta disponible"
	}
	defer conn.Close()
	fmt.Fprintf(conn, "%s %s %s \n", comando, nombre, edad)
	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		return scanner.Text()
	}
	return "ERROR no hay respuesta del registro"
}

func listarEstudiantes(connCliente net.Conn) {
	conn, err := net.Dial("tcp", "localhost:9002")
	if err != nil {
		fmt.Fprintf(connCliente, "ERROR el servicio registro no esta disponible\n")
		return
	}
	defer conn.Close()
	fmt.Fprintf(conn, "LISTAR\n")
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		linea := scanner.Text()
		if linea == "FIN" {
			break
		}
		fmt.Fprintf(connCliente, "%s \n", linea)
	}
}