package main

import (
	"fmt"
	"net"
	"sync"
	"io"
)

// Canale globale pentru comunicarea inter-goroutine
var messages = make(chan string)
var newConnections = make(chan net.Conn)
var deadConnections = make(chan net.Conn)

// Mutex pentru a proteja harta de clienți, deși vom folosi o goroutine dedicată
// pentru a gestiona harta, care elimină necesitatea de mutex.
// O lăsăm doar pentru ilustrarea conceptului de acces protejat la resurse comune.
var clients = make(map[net.Conn]bool)
var mu sync.Mutex 


func main() {
	fmt.Println("Starting chat server on :8080")
	
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	// 1. Goroutine-ul central care gestionează clienții și difuzează mesajele
	go broadcastManager() 

	fmt.Println("Server is listening on port 8080. Waiting for clients...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		
		// 2. Înregistrează noul client
		newConnections <- conn
		
		// 3. Manipulează conexiunea (citirea mesajelor)
		go handleConnection(conn) 
	}
}

// Goroutine-ul principal de broadcast
func broadcastManager() {
	for {
		select {
		case conn := <-newConnections:
			// Adauga client nou
			mu.Lock()
			clients[conn] = true
			mu.Unlock()
			fmt.Printf("New client connected: %s. Total clients: %d\n", conn.RemoteAddr(), len(clients))

		case conn := <-deadConnections:
			// Sterge clientul deconectat
			mu.Lock()
			delete(clients, conn)
			mu.Unlock()
			conn.Close() // Asigura-te că e închisa
			fmt.Printf("Client disconnected: %s. Total clients: %d\n", conn.RemoteAddr(), len(clients))

		case msg := <-messages:
			// Difuzeaza mesajul tuturor clienților
			fmt.Printf("Broadcasting: %s", msg)
			mu.Lock()
			for conn := range clients {
				_, err := conn.Write([]byte(msg))
				if err != nil {
					// Marchează clientul pentru ștergere în caz de eroare
					deadConnections <- conn 
				}
			}
			mu.Unlock()
		}
	}
}

func handleConnection(conn net.Conn) {
	// Trimite mesaj de deconectare în cazul în care goroutine-ul se termina
	defer func() {
		deadConnections <- conn 
	}()

	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				// Eroare la citire (de obicei deconectare)
				fmt.Println("Error reading:", err) 
			}
			break
		}
	
		// Formatează mesajul înainte de a-l trimite pe canal
		message := fmt.Sprintf("[%s]: %s", conn.RemoteAddr().String(), string(buffer[:n]))
		
		// Trimite mesajul către managerul de broadcast
		messages <- message 
	}
}