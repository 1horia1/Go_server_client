# 🗨️ Go TCP Chat Server & Client

A small **TCP chat system written in Go**, which allows multiple clients to connect to a server, send messages, and receive them in real time (broadcast).

---

## Features

- Server multi-client based on `net.Listener`  
- Concurrent management with **goroutines**  
- Communication between goroutines using **channels**  
- Automatic broadcast to all connected clients  
- Detection and cleanup of closed connections  
- Interactive client from terminal  

---

## How to start the application

### 1 Start the server

Open a terminal and run:

```bash
go run server.go
```

You should see:
```
Starting chat server on :8080
Server is listening on port 8080. Waiting for clients...
```

### 2 Connect the client
Open another terminal and run:
```
go run client.go
```
You will see the message:
```
Connected to chat server. Start typing your message (press Enter to send).
> 
```
You can open multiple terminals and run `go run client.go` several times.
All clients will receive the messages sent by any of them.

## Example run
Server
```
Starting chat server on :8080
Server is listening on port 8080. Waiting for clients...
New client connected: 127.0.0.1:53412. Total clients: 1
New client connected: 127.0.0.1:53414. Total clients: 2
Broadcasting: [127.0.0.1:53412]: Hello!
Broadcasting: [127.0.0.1:53414]: Hello, everybody!
```
