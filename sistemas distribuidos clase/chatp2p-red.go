package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

// Esta función se encarga de escuchar los mensajes entrantes (ROL SERVIDOR)
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
			log.Println("Error al aceptar:", err)
			continue
		}

		// Manejar cada conexión entrante
		go recibirMensajes(conn)
	}
}

// Recibir mensajes de otro peer
func recibirMensajes(conn net.Conn) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		mensaje := scanner.Text()
		fmt.Printf("\n Estudiante: %s\n> ", mensaje)
	}
}

// Esta función se encarga de enviar los mensajes (ROL CLIENTE)
func enviarMensajes(ip string, puerto string) {
	direccion := ip + ":" + puerto

	fmt.Printf("Conectando a %s...\n", direccion)

	conn, err := net.Dial("tcp", direccion)
	if err != nil {
		log.Fatal("Error al conectar:", err)
	}
	defer conn.Close()

	fmt.Println("Se logro la conexion!...")
	fmt.Println("Escribe tus mensajes (escribe SALIR para terminar):")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		mensaje := scanner.Text()

		if strings.ToUpper(mensaje) == "SALIR" {
			fmt.Println("Cerrando chat, hasta luego!..")
			return
		}

		// Enviar mensaje al otro peer
		fmt.Fprintf(conn, "%s\n", mensaje)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso:")
		fmt.Println("  Modo ESPERAR: go run chat-p2p.go [mi-puerto]")
		fmt.Println("  Modo CONECTAR: go run chat-p2p.go [mi-puerto] [ip-peer] [puerto-peer]")
		fmt.Println("")
		fmt.Println("Ejemplo:")
		fmt.Println("  Peer 1: go run chat-p2p.go 9001")
		fmt.Println("  Peer 2: go run chat-p2p.go 9002 localhost 9001")
		return
	}

	miPuerto := os.Args[1]

	// Iniciar servidor (escuchar) en goroutine
	go escucharMensajes(miPuerto)

	// Si hay argumentos adicionales, conectar a otro peer
	if len(os.Args) == 4 {
		ipPeer := os.Args[2]
		puertoPeer := os.Args[3]

		// Esperar un poco para que el servidor esté listo
		fmt.Println("Iniciando...")

		enviarMensajes(ipPeer, puertoPeer)
	} else {
		fmt.Printf("Esperando que otro peer (estudiante) se conecte a ti en puerto %s...\n", miPuerto)
		fmt.Println("Cuando se conecten, podrás escribir mensajes aquí.")

		// Mantener el programa corriendo
		select {}
	}
}
