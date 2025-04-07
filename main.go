package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	var reqs []string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("error reading from connection")
			os.Exit(1)
		}
		if line == "\r\n" {
			break
		}
		reqs = append(reqs, line)
	}
	getRequest := reqs[0]
	getRequestArr := strings.Split(getRequest, " ")
	if getRequestArr[1] == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else if getRequestArr[1][0:6] == "/echo/" {
		contentLength := len(getRequestArr[1][6:])
		content := getRequestArr[1][6:]
		resp := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n%v", contentLength, content)
		conn.Write([]byte(resp))
	} else if getRequestArr[1] == "/user-agent" {
		agentLine := strings.Split(reqs[2], " ")
		contentLength := len(strings.TrimSpace(agentLine[1]))
		agent := agentLine[1]
		resp := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n%v", contentLength, agent)
		conn.Write([]byte(resp))

	} else if getRequestArr[1][0:7] == "/files/" {
		file := getRequestArr[1][7:]
		path := os.Args[2]
		fmt.Printf(path)
		os.Chdir(path)
		contents, err := os.ReadFile(file)
		if err != nil {
			fmt.Println("error reading file")
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
			fmt.Println(file)
		} else {
			contentLength := len(string(contents))
			resp := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %v\r\n\r\n%v", contentLength, string(contents))
			conn.Write([]byte(resp))
		}
	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		defer conn.Close()
		// conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		// reader := bufio.NewReader(conn)
		// var reqs []string
		// for {
		// 	line, err := reader.ReadString('\n')
		// 	if err != nil {
		// 		fmt.Println("error reading from connection")
		// 		os.Exit(1)
		// 	}
		// 	if line == "\r\n" {
		// 		break
		// 	}
		// 	reqs = append(reqs, line)
		// }
		// getRequest := reqs[0]
		// getRequestArr := strings.Split(getRequest, " ")
		// if getRequestArr[1] == "/" {
		// 	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		// } else if getRequestArr[1][0:6] == "/echo/" {
		// 	contentLength := len(getRequestArr[1][6:])
		// 	content := getRequestArr[1][6:]
		// 	resp := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n%v", contentLength, content)
		// 	conn.Write([]byte(resp))
		// } else if getRequestArr[1] == "/user-agent" {
		// 	agentLine := strings.Split(reqs[2], " ")
		// 	contentLength := len(strings.TrimSpace(agentLine[1]))
		// 	agent := agentLine[1]
		// 	resp := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %v\r\n\r\n%v", contentLength, agent)
		// 	conn.Write([]byte(resp))

		// } else {
		// 	conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		// }
		go handleConnection(conn)
	}
}
