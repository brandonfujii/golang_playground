package main 

import (
  "net"
  "log"
  "os"
  "fmt"
  "strings"
  "bufio"
)

type SimpleServer struct {
  port string
}

type SimpleClient struct {
  connection net.Conn
  Server *SimpleServer 
} 

func main() {
  args := os.Args[1:]
  
  if len(args) < 1 {
    fmt.Println("Usage: go run main.go [port]")
    os.Exit(1)
  }

  port := args[0]

  li, err := net.Listen("tcp", ":" + port)

  if err != nil {
    log.Fatalln(err.Error())
  }

  server := &SimpleServer{port: port}
  server.listen(li)
}

func (c *SimpleClient) Conn() net.Conn {
  return c.connection
}

func (s *SimpleServer) listen(li net.Listener) {
  fmt.Println("Local server listening on :" + s.port)
  defer li.Close()

  for {
    connection, err := li.Accept()
    if err != nil {
      log.Println(err.Error())
      continue
    }
    defer connection.Close()

    go s.handleRequest(connection)
  }
}

func (s *SimpleServer) handleRequest(connection net.Conn) {
  defer connection.Close()
  fmt.Println("Handling Request\n")
  s.request(connection)
}

func (s *SimpleServer) request(connection net.Conn) {
  i := 0
  scanner := bufio.NewScanner(connection)

  for scanner.Scan() {
    ln := scanner.Text()
    fmt.Println(ln)

    if i == 0 {
      method := strings.Fields(ln)[0]
      uri := strings.Fields(ln)[1]
      s.respond(connection, method, uri)
    }

    if ln == "" {
      break
    }

    i++
  }
}

func (s *SimpleServer) respond(connection net.Conn, method string, uri string) {
  fmt.Println("Responding to request:", method, uri)

  if method == "GET" && uri == "/" {
    s.index(connection)
  } else {
    s.resourceNotFound(connection)
  }
}

/* Status Codes */ 
var ok string = "HTTP/1.1 200 OK\r\n"
var notFound string = "HTTP/1.1 404 Not Found\r\n"

func (s *SimpleServer) index(connection net.Conn) {
  body := `<!DOCTYPE html><html lang="en"><head><meta charset="utf-8"><title>Homepage</title></head><body>
  <h1>HOME PAGE</h1> </body></html>`
  fmt.Fprint(connection, ok)
  fmt.Fprintf(connection, "Content-Length: %d\r\n", len(body))
  fmt.Fprint(connection, "Content-Type: text/html\r\n\r\n")
  fmt.Fprint(connection, body)

}

func (s *SimpleServer) resourceNotFound(connection net.Conn) {
  body := `<!DOCTYPE html><html lang="en"><head><meta charset="utf-8"><title>Homepage</title></head><body>
  <h1>404 Not Found Error</h1></body></html>`
  fmt.Fprint(connection, notFound)
  fmt.Fprintf(connection, "Content-Length: %d\r\n", len(body))
  fmt.Fprint(connection, "Content-Type: text/html\r\n\r\n")
  fmt.Fprint(connection, body)
}

