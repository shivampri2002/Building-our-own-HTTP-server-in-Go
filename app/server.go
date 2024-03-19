package main

import (
	"bufio"
	"fmt"

	// Uncomment this block to pass the first stage
	"net"
	"os"
	"strings"
)

var statusCodeToString = map[int]string{
	200: "OK",
	404: "Not Found",
}

type Request struct {
	Method    string
	Path      string
	UserAgent string
}

type Response struct {
	StatusCode    int
	Body          string
	ContentType   string
	ContentLength int
}

func (r Response) String() string {
	statusText := statusCodeToString[r.StatusCode]
	r.ContentLength = len(r.Body)

	return fmt.Sprintf("HTTP/1.1 %d %s\r\nContent-Type: %s\r\nContent-Length: %d\r\n\r\n%s", r.StatusCode, statusText, r.ContentType, r.ContentLength, r.Body)
}

func parseRequest(reader *bufio.Reader) (Request, error) {
	first, err := reader.ReadString('\n')
	if err != nil {
		return Request{}, err
	}

	_, _ = reader.ReadString('\n')
	thrd, err := reader.ReadString('\n')

	parts := strings.Split(first, " ")

	if err == nil {
		return Request{Method: parts[0], Path: parts[1], UserAgent: strings.Split(thrd, " ")[1]}, nil
	}

	return Request{Method: parts[0], Path: parts[1]}, nil
}

func handleConnection(conn net.Conn) {
	stream := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	defer func(conn net.Conn) {
		err := conn.Close()

		if err != nil {
			fmt.Println("Error in closing the Connection: ", err.Error())
			os.Exit(1)
		}
	}(conn)

	// _, err := stream.WriteString("HTTP/1.1 200 OK\r\n\r\n")

	request, err := parseRequest(stream.Reader)

	if err != nil {
		fmt.Println("Failed to parse the request: ", err.Error())
		os.Exit(1)
	}

	var response Response

	match := false
	// fmt.Println(request.Path, request.Path[1:5])
	if len(request.Path) > 6 && request.Path[1:5] == "echo" {
		match = true
	}

	if request.Path == "/" || request.Path == "/user-agent" || match {
		if match || request.Path == "/user-agent" {
			body := request.Path[6:]
			if request.UserAgent != "" {
				body = request.UserAgent
			}
			response = Response{StatusCode: 200, ContentType: "text/plain", Body: body}
		} else {
			response = Response{StatusCode: 200}
		}
	} else {
		response = Response{StatusCode: 404}
	}

	_, err = stream.WriteString(response.String())

	if err != nil {
		fmt.Println("Failed to write to the Socket: ", err.Error())
		os.Exit(1)
	}

	err = stream.Flush()

	if err != nil {
		fmt.Println("Failed to flush the Socket: ", err.Error())
		os.Exit(1)
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", ":4221")

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

		go handleConnection(conn)
	}

	// defer l.Close()
	// conn, err := l.Accept()

	// if err != nil {
	// 	fmt.Println("Error accepting connection: ", err.Error())
	// 	os.Exit(1)
	// }

	// buffer := make([]byte, 1024)

	// httpRequest, err := conn.Read(buffer)

	// if err != nil {
	// 	fmt.Println("Error reading the data from the connection: ", err.Error())
	// 	os.Exit(1)
	// }

	// fmt.Println(httpRequest)

	// httpResponse := "HTTP/1.1 200 OK\r\n\r\n"

	// _, err = conn.Write([]byte(httpResponse))

	// if err != nil {
	// 	fmt.Println("Error writing to the connection: ", err.Error())
	// 	os.Exit(1)
	// }

	// defer conn.Close()

}
