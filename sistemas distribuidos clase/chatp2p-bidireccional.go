package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var miNombre string

// Escuchar mensajes entrantes
func escucharMensajes(puerto string) {
	listener, err := net.Listen("tcp", ":"+puerto)
	if err != nil {
		log.Fatal("Error al escuchar:", err)
	}
	defer listener.Close()

	fmt.Printf("Escuchando en puerto %s\n", puerto)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error:", err)
			continue
		}
		go recibirMensajes(conn)
	}
}

// Recibir mensajes
func recibirMensajes(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		mensaje := scanner.Text()
		fmt.Printf("\n  %s\n> ", mensaje)
	}
}

// Enviar mensajes
func enviarMensajes(ip string, puerto string) {
	time.Sleep(1 * time.Second) // Esperar que el otro peer esté listo

	direccion := ip + ":" + puerto
	fmt.Printf("🔌 Conectando a %s...\n", direccion)

	conn, err := net.Dial("tcp", direccion)
	if err != nil {
		fmt.Printf(" No se pudo conectar a %s\n", direccion)
		fmt.Println("El otro peer (estudiante)debe estar corriendo primero")
		return
	}
	defer conn.Close()

	fmt.Println("Conectado")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		mensaje := scanner.Text()

		if strings.ToUpper(mensaje) == "SALIR" {
			return
		}

		// Enviar con nombre
		fmt.Fprintf(conn, "%s: %s\n", miNombre, mensaje)
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("USO:")
		fmt.Println("  go run chat-p2p-bidireccional.go [nombre] [mi-puerto] [ip-peer] [puerto-peer]")
		fmt.Println("")
		fmt.Println("EJEMPLO:")
		fmt.Println("  Terminal 1: go run chat-p2p-bidireccional.go Estudiante1 9001 localhost 9002")
		fmt.Println("  Terminal 2: go run chat-p2p-bidireccional.go Estudiante2 9002 localhost 9001")
		return
	}

	miNombre = os.Args[1]
	miPuerto := os.Args[2]

	fmt.Printf("Hola %s\n", miNombre)

	// Escuchar en mi puerto
	go escucharMensajes(miPuerto)

	// Conectar al otro peer si se proporcionó
	if len(os.Args) == 5 {
		ipPeer := os.Args[3]
		puertoPeer := os.Args[4]
		enviarMensajes(ipPeer, puertoPeer)
	} else {
		fmt.Println("Esperando conexiones...")
		select {}
	}
}