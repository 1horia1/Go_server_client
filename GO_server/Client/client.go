package main 

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"io"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to chat server. Start typing your message (press Enter to send).")

	// Goroutine pentru primirea mesajelor de la server (citire din conexiune)
	go receiveMessages(conn) 
	
	// Goroutine pentru trimiterea mesajelor la server (citire de la tastatură)
	sendMessages(conn) 
}

// Asculta constant mesaje de la server
func receiveMessages(conn net.Conn) {
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF { // Importă io
				fmt.Println("\nServer connection lost (read error):", err)
			}
			return
		}
		// Curata linia pentru a nu interfera cu prompt-ul de tastare
		fmt.Printf("\r\033[K%s", string(buffer[:n]))
		fmt.Print("> ") // Reafișează prompt-ul
	}
}

// Citeste de la tastatura și trimite la server
func sendMessages(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	
	
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		
		_, err := conn.Write([]byte(text)) // Trimite mesajul direct

		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
	}
}