package main

import (
	"bufio"
	"fmt"
	"io"

	// Uncomment this block to pass the first stage
	"net"
	"os"
	"strconv"
	"strings"
)

var statusCodeToString = map[int]string{
	200: "OK",
	404: "Not Found",
}

type Request struct {
	Method string
	Path   string
	// UserAgent string
	Body    string
	Headers map[string]string
}

type Response struct {
	StatusCode int
	Body       string
	// ContentType   string
	// ContentLength int
	Headers map[string]string
}

func (r Response) String() string {
	statusText, ok := statusCodeToString[r.StatusCode]
	if !ok {
		statusText = "Unknown"
	}
	// r.Headers["contentLength"] = strconv.Itoa(len(r.Body))

	if r.Headers == nil {
		r.Headers = make(map[string]string)
	}

	if _, ok = r.Headers["Content-Length"]; !ok {
		r.Headers["Content-Length"] = strconv.Itoa(len(r.Body))
	}

	var headerString strings.Builder
	for k, v := range r.Headers {
		headerString.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}

	return fmt.Sprintf("HTTP/1.1 %d %s\r\n%s\r\n%s", r.StatusCode, statusText, headerString.String(), r.Body)
}

func parseRequest(reader *bufio.Reader) (Request, error) {
	request := Request{
		Headers: make(map[string]string),
	}

	firstLine, err := reader.ReadString('\n')
	if err != nil {
		return Request{}, err
	}

	parts := strings.Split(firstLine, " ")
	request.Method = parts[0]
	request.Path = parts[1]

	// _, _ = reader.ReadString('\n')
	// thrd, err := reader.ReadString('\n')

	for {
		curLine, err := reader.ReadString('\n')
		if curLine == "\r\n" {
			break
		}
		if err == io.EOF {
			return request, nil
		} else if err != nil {
			return Request{}, err
		}

		headerParts := strings.SplitN(curLine, ":", 2)
		request.Headers[headerParts[0]] = strings.TrimSpace(headerParts[1])
	}

	// if err == nil {
	// 	return Request{Method: parts[0], Path: parts[1], Headers: map[string]string{
	// 		"userAgent": strings.Split(thrd, " ")[1],
	// 	}}, nil
	// }

	// return Request{Method: parts[0], Path: parts[1]}, nil

	contentLenStr, ok := request.Headers["Content-Length"]
	if !ok {
		return request, nil
	}

	contentLen, _ := strconv.Atoi(contentLenStr)

	buf := make([]byte, contentLen)
	_, err = io.ReadFull(reader, buf)

	if err != nil {
		return request, err
	}

	request.Body = string(buf)

	return request, nil
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

	// var response Response
	response := Response{StatusCode: 200}

	// match := false
	// // fmt.Println(request.Path, request.Path[1:5])
	// if len(request.Path) > 6 && request.Path[1:5] == "echo" {
	// 	match = true
	// }

	// if request.Path == "/" || request.Path == "/user-agent" || match {
	// 	if match || request.Path == "/user-agent" {
	// 		body := request.Path[6:]
	// 		if request.Headers["userAgent"] != "" {
	// 			body = request.Headers["userAgent"]
	// 		}
	// 		response = Response{StatusCode: 200, Headers: map[string]string{
	// 			"contentType": "text/plain",
	// 		}, Body: body}
	// 	} else {
	// 		response = Response{StatusCode: 200}
	// 	}
	// } else {
	// 	response = Response{StatusCode: 404}
	// }
	fmt.Println(request.Path)

	if strings.HasPrefix(request.Path, "/echo") {
		pathParts := strings.SplitN(request.Path, "/echo/", 2)
		response.Body = pathParts[1]
	} else if request.Path == "/user-agent" {
		userAgent := request.Headers["User-Agent"]
		fmt.Printf(userAgent)
		response.Body = userAgent
	} else if request.Path != "/" {
		response.StatusCode = 404
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
